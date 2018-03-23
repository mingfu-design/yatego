package yatego

// Controller main bot object
type Controller struct {
	callManager          CallManager
	fallbackToController bool
	singleChannelMode    bool
	flowID               string
	staticComponents     []Component
	callflowLoader       CallflowLoader
	logger               Logger
	engine               Engine
}

// Run the IVR system, main loop
func (c *Controller) Run(name string) {
	for {
		msg, err := c.engine.GetEvent()
		if c.singleChannelMode && (err != nil || msg == nil) {
			c.logger.Fatalf("event msg err: %s", err)
			break
		}
		c.logger.Debugf("msg: %+v", msg)

		call := c.getCall(msg)
		result := c.processEvent(msg, call)
		c.logger.Debugf("event process result: %s", result)
		//we need to ack all incoming messages
		if msg.Type == TypeIncoming {
			c.logger.Debugf("ACK incoming event: %s", msg.Name)
			c.engine.Acknowledge(msg)
		}
		if !c.prepareNextEvent(result, call) {
			break
		}
	}
}

func (c *Controller) getCall(msg *Message) *Call {
	return nil
}

func (c *Controller) processIncomingCall(msg *Message) *Call {
	return nil
}

func (c *Controller) getCallChanneID(msg *Message) *Call {
	return nil
}

func (c *Controller) processEvent(msg *Message, call *Call) string {
	return ""
}

func (c *Controller) execCallback(msg *Message, call *Call) string {
	return ""
}

func (c *Controller) activeComponent(call *Call) Component {
	return nil
}

func (c *Controller) prepareNextEvent(processResult string, call *Call) bool {
	return true
}

func (c *Controller) installMessageHandlers(call *Call) {
}

func (c *Controller) installMessageWatches(call *Call) {
}

func (c *Controller) loadComponents(params map[string]string) []Component {
	if c.callflowLoader == nil {
		return c.staticComponents
	}
	callflow := c.callflowLoader.load(params)
	components := make([]Component, 0)
	//build components
	for _, com := range callflow.Components {
		components = append(components, com.Factory(com.ClassName, com.Name, com.Config))
	}
	return components
}

func (c *Controller) addStaticComponent(component Component) {
	c.staticComponents = append(c.staticComponents, component)
}
