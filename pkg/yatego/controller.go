package yatego

// Controller main bot object
type Controller struct {
	callManager          *CallManager
	fallbackToController bool
	singleChannelMode    bool
	flowID               string
	staticComponents     []Component
	callflowLoader       CallflowLoader
	logger               Logger
	engine               *Engine
}

// Run the IVR system, main loop
func (c *Controller) Run(name string) {
	for {
		msg, err := c.engine.GetEvent()
		if err != nil {
			c.logger.Fatalf("event msg err: %s", err)
			if c.singleChannelMode {
				break
			} else {
				continue
			}
		}
		if msg == nil {
			if err == nil {
				c.logger.Warningln("event msg EOF")
			} else {
				c.logger.Fatal("event msg is nil")
			}
			if c.singleChannelMode {
				break
			} else {
				continue
			}
		}
		if msg == nil {
			c.logger.Fatal("event msg is nil")
			if c.singleChannelMode {
				break
			} else {
				continue
			}
		}
		c.logger.Debugf("new msg: %+v", msg)
		call, processed := c.getCall(msg)
		if call == nil && processed {
			c.logger.Fatal("Call not found, we need to exit")
			if c.singleChannelMode {
				break
			}
			continue
		}
		if call == nil {
			c.logger.Debug("Call not for us, ignoring")
			continue
		}
		res := c.processEvent(msg, call)
		c.logger.Debugf("event process result: %v", res)
		//we need to ack all incoming messages
		if msg.Type == TypeIncoming {
			c.logger.Debugf("ACK incoming event: %s", msg.Name)
			c.engine.Acknowledge(msg)
		}
		if !c.prepareNextEvent(res, call) {
			break
		}
	}
}

//Logger get logger
func (c *Controller) Logger() Logger {
	return c.logger
}

func (c *Controller) getCall(msg *Message) (*Call, bool) {
	call, processed := c.processIncomingCall(msg)
	if processed {
		return call, true
	}
	chID := c.getCallChanneID(msg)
	if chID == "" {
		c.logger.Debugf("Channel ID not defined in msg params %+v", msg.Params)
		return nil, false
	}
	call, exists := c.callManager.Call(chID)
	if !exists {
		c.logger.Errorf("Call for channel ID [%s] not found", chID)
		return nil, true
	}
	return call, true
}

func (c *Controller) processIncomingCall(msg *Message) (*Call, bool) {
	if msg.Type != TypeIncoming || msg.Name != MsgCallExecute {
		return nil, false
	}
	if !c.singleChannelMode || msg.Params["flow"] != c.flowID {
		return nil, false
	}
	//new call
	c.logger.Infof("New call received we're going to answer: %+v", msg.Params)
	//load callflow
	coms := c.loadComponents(msg.Params)
	if coms == nil || len(coms) == 0 {
		c.logger.Fatal("No components loaded")
	}
	call, err := c.callManager.Add(coms, msg.Params, "", "")
	if err != nil {
		c.logger.Fatalf("Call not added: %s", err)
		return nil, true
	}
	c.logger.Infof("New call added: %+v", call)
	//install handlers
	c.installMessageHandlers(call)
	c.installMessageWatches(call)

	return call, true
}

func (c *Controller) getCallChanneID(msg *Message) string {
	if msg.Type == TypeAnswer {
		if id, exists := msg.Params["id"]; exists {
			return id
		}
		return ""
	}
	if msg.Type == TypeIncoming {
		if id, exists := msg.Params["targetid"]; exists {
			return id
		}
		return ""
	}
	return ""
}

func (c *Controller) processEvent(msg *Message, call *Call) *CallbackResult {
	//controllerHandled := false
	com := c.activeComponent(call)
	if com == nil {
		c.logger.Fatal("No component found")
		return nil
	}
	cb := com.Callback(msg.Name)
	//TODO fallback to controller
	if cb == nil {
		c.logger.Debugf("No callback found for msg [%s] type [%s] in component [%s]", msg.Name, msg.Type, com.Name())
		return nil
	}

	res := cb(call, msg)
	c.logger.Debugf("For msg [%s] type [%s] component [%s] processed with result: %+v", msg.Name, msg.Type, com.Name(), res)
	if msg.Processed {
		return res
	}
	//TODO fallback to controller
	return res
}

func (c *Controller) activeComponent(call *Call) Component {
	//TODO fallback to controller
	return call.Component(call.ActiveComponentName)
}

func (c *Controller) prepareNextEvent(res *CallbackResult, call *Call) bool {
	if res == nil {
		c.logger.Debug("No process result, returning TRUE")
		return true
	}
	switch res.result {
	case ResStop:
		return !c.singleChannelMode
	case ResTransfer:
		next := call.Component(res.transferComponent)
		if next == nil {
			c.logger.Errorf("Transfer to component [%s] not found", res.transferComponent)
			return !c.singleChannelMode
		}
		call.ActiveComponentName = next.Name()
		c.logger.Infof("Entering new active component [%s] in call [%s]", next.Name(), call.ChannelID)
		next.Enter(call)
		return true
	case ResEnter:
		com := c.activeComponent(call)
		if com == nil {
			c.logger.Fatalf("Active component not found to enter in call [%s]", call.ChannelID)
			return false
		}
		c.logger.Infof("Entering component [%s] after callback in call [%s]", com.Name(), call.ChannelID)
		if !com.Enter(call) {
			c.logger.Infof("Component [%s] Enter callback not defined in call [%s]", com.Name(), call.ChannelID)
		}
		return true
	default:
		return true
	}
}

func (c *Controller) installMessageHandlers(call *Call) {
	msgs := make(map[string]InstallDef)
	coms := call.Components()
	for _, com := range coms {
		for msgName, msgDef := range com.MessagesToInstall() {
			c.logger.Debugf("Analysing msf [%s]: %+v", msgName, msgDef)
			if _, exists := msgs[msgName]; !exists {
				msgs[msgName] = msgDef
				continue
			}
			if msgs[msgName].Priority < msgDef.Priority {
				continue
			}
			delete(msgs, msgName)
			msgs[msgName] = msgDef
		}
	}
	c.logger.Debugf("Going to install [%d] message handlers, from [%d] components", len(msgs), len(coms))
	for msgName, msgDef := range msgs {
		c.engine.InstallFiltered(msgName, msgDef.Priority, msgDef.FilterName, msgDef.FilterValue)
	}
}

func (c *Controller) installMessageWatches(call *Call) {
	msgs := make(map[string]bool)
	coms := call.Components()
	for _, com := range coms {
		for _, msgName := range com.MessagesToWatch() {
			if _, exists := msgs[msgName]; exists {
				continue
			}
			msgs[msgName] = true
		}
	}
	c.logger.Debugf("Going to install [%d] message watchers, from [%d] components", len(msgs), len(coms))
	for msgName := range msgs {
		c.engine.Watch(msgName)
	}
}

func (c *Controller) loadComponents(params map[string]string) []Component {
	if c.callflowLoader == nil {
		return c.staticComponents
	}
	cf, err := c.callflowLoader.Load(params)
	if err != nil {
		c.logger.Fatalf("Error loading callflow: %s", err)
	}
	coms := make([]Component, 0)
	//build components
	for _, com := range cf.Components {
		c.logger.Debugf("Building component: %+v", com)
		coms = append(coms, com.Factory(com.ClassName, com.Name, com.Config))
	}
	return coms
}

// AddStaticComponent add a component to the controller
func (c *Controller) AddStaticComponent(component Component) {
	c.staticComponents = append(c.staticComponents, component)
}
