package yatego

// Controller main bot object
type Controller struct {
	componentYate
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
	chID := c.getCallChannelID(msg)
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
	call, err := c.callManager.Add(coms, msg.Params, "", "", c.logger)
	if err != nil {
		c.logger.Fatalf("Call not added: %s", err)
		return nil, true
	}
	c.logger.Infof("New call added: %+v", call)
	//install handlers
	c.InstallMessageHandlers(call)
	c.InstallMessageWatches(call)

	return call, true
}

func (c *Controller) getCallChannelID(msg *Message) string {
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
	ac := call.ActiveComponent()
	if ac != nil {
		return ac
	}
	//no new component, activate first one
	startCom := call.Component("")
	c.logger.Debugf("No active component, going to activate first one [%s]", startCom.Name())
	call.ActivateComponent(startCom.Name())
	return call.ActiveComponent()
	//TODO fallback to controller
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
		c.logger.Infof("Entering new active component [%s] in call [%s]", res.transferComponent, call.ChannelID)
		if !call.ActivateComponent(res.transferComponent) {
			c.logger.Errorf("Transfer to component [%s] failed", res.transferComponent)
			return !c.singleChannelMode
		}
		//recursive enter the component
		return c.prepareNextEvent(NewCallbackResult(ResEnter, ""), call)
	case ResEnter:
		com := c.activeComponent(call)
		if com == nil {
			c.logger.Fatalf("Active component not found to enter in call [%s]", call.ChannelID)
			return false
		}
		c.logger.Infof("Entering component [%s] after callback in call [%s]", com.Name(), call.ChannelID)
		enterRes := com.Enter(call)
		//recursive prepare again
		return c.prepareNextEvent(enterRes, call)
	default:
		return true
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
	c.logger.Debug("Building components:", cf.Components)
	//build components
	for _, com := range cf.Components {
		coms = append(coms, com.Factory(com.ClassName, com.Name, com.Config))
	}
	return coms
}

// AddStaticComponent add a component to the controller
func (c *Controller) AddStaticComponent(component Component) {
	c.staticComponents = append(c.staticComponents, component)
}
