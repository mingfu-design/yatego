package yatego

type componentCommon struct {
	name   string
	config map[string]interface{}
	logger Logger
}

func (c *componentCommon) Name() string {
	return c.name
}

func (c *componentCommon) Config(key string) (interface{}, bool) {
	if key == "" {
		return c.config, true
	}
	val, exists := c.config[key]
	return val, exists
}

func (c *componentCommon) TransferComponent() (string, bool) {
	com, exists := c.Config("transfer")
	if !exists {
		return "", false
	}
	return com.(string), true
}

func (c *componentCommon) TransferCallbackResult() *CallbackResult {
	trCom, exists := c.TransferComponent()
	if !exists {
		return NewCallbackResult(ResStop, "")
	}
	return NewCallbackResult(ResTransfer, trCom)
}

func (c *componentCommon) Logger() Logger {
	return c.logger
}

func (c *componentCommon) CallData(call *Call, key string) (interface{}, bool) {
	return call.Data(c.name, key)
}

func (c *componentCommon) SetCallData(call *Call, key string, value interface{}) {
	call.SetData(c.name, key, value)
}

type componentCallback struct {
	callbacks map[string]Callback
}

func (c *componentCallback) Callback(msgName string) Callback {
	if _, exists := c.callbacks[msgName]; !exists {
		return nil
	}
	return c.callbacks[msgName]
}

func (c *componentCallback) Listen(msgName string, cb Callback) {
	c.callbacks[msgName] = cb
}

func (c *componentCallback) OnEnter(cb Callback) {
	c.callbacks[MsgComponentEnter] = cb
}

func (c *componentCallback) Enter(call *Call) bool {
	cb := c.Callback(MsgComponentEnter)
	if cb == nil {
		return false
	}
	cb(call, nil)
	return true
}

type componentYate struct {
	engine            *Engine
	messagesToWatch   []string
	messagesToInstall map[string]InstallDef
}

func (c *componentYate) SendMessage(msgName string, call *Call, params map[string]string, targetID string) (*Message, error) {
	if targetID == "" {
		targetID = call.PeerID
	}
	if targetID != "" {
		params["targetid"] = targetID
	}
	params["id"] = call.ChannelID
	m := NewMessage(msgName, params)

	_, err := c.engine.Dispatch(m)
	return m, err
}

func (c *componentYate) Answer(call *Call, msg *Message) (*Message, error) {
	msg.Params["targetid"] = call.ChannelID
	msg.Processed = true

	_, err := c.engine.Acknowledge(msg)
	if err != nil {
		return msg, err
	}
	return c.SendMessage(MsgCallAnswered, call, map[string]string{"cdrcreate": "no"}, "")
}

func (c *componentYate) MessagesToWatch() []string {
	return c.messagesToWatch
}

func (c *componentYate) MessagesToInstall() map[string]InstallDef {
	return c.messagesToInstall
}

type componentMedia struct {
	componentYate
}

func (c *componentMedia) PlayWave(wave string, call *Call, params map[string]string) (*Message, error) {
	params["source"] = "wave/play/" + wave
	params["notify"] = call.ChannelID
	return c.SendMessage(MsgChanAttach, call, params, "")
}

func (c *componentMedia) PlayTone(tone string, call *Call, params map[string]string) (*Message, error) {
	params["source"] = "tone/" + tone
	params["notify"] = call.ChannelID
	return c.SendMessage(MsgChanAttach, call, params, "")
}

func (c *componentMedia) Record(file string, recordTime string, call *Call, params map[string]string) (*Message, error) {
	params["consumer"] = "wave/record/" + file
	params["notify"] = call.ChannelID
	params["maxlen"] = recordTime
	params["single"] = "true"
	return c.SendMessage(MsgChanAttach, call, params, "")
}
