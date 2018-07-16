package yatego

import (
	"strconv"
	"strings"
)

type componentCommon struct {
	name      string
	config    map[string]interface{}
	configTpl map[string]interface{}
	logger    Logger
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

func (c *componentCommon) ConfigTpl(key string) (interface{}, bool) {
	if key == "" {
		return c.configTpl, true
	}
	val, exists := c.configTpl[key]
	return val, exists
}

func (c *componentCommon) ConfigKeys() []string {
	keys := make([]string, 0, len(c.config))
	for k := range c.config {
		keys = append(keys, k)
	}
	return keys
}

func (c *componentCommon) SetConfig(key string, val interface{}) {
	c.config[key] = val
}

func (c *componentCommon) ConfigAsString(key string) (string, bool) {
	val, exists := c.Config(key)
	if !exists {
		return "", false
	}
	switch t := val.(type) {
	default:
		c.logger.Errorf("Unable to convert component [%s] config key [%s] val : %+v to string", c.Name(), key, val)
		return "", false
	case string:
		return t, true
	case bool:
		return strconv.FormatBool(t), true
	case int:
		return strconv.Itoa(t), true
	case uint64:
		return strconv.FormatUint(t, 10), true
	case float64:
		return strconv.FormatFloat(t, 'f', 6, 64), true
	}
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

func (c *componentCommon) CallDataNamespace(call *Call, namespace string) string {
	ns := strings.Split(namespace, ".")
	if len(ns) != 2 {
		return ""
	}
	data, ok := call.Data(ns[0], ns[1])
	if !ok {
		return ""
	}
	return data.(string)
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

func (c *componentCallback) Enter(call *Call) *CallbackResult {
	cb := c.Callback(MsgComponentEnter)
	if cb == nil {
		return NewCallbackResult(ResStay, "")
	}
	return cb(call, nil)
}

type componentYate struct {
	componentCommon
	engine            *Engine
	messagesToWatch   []string
	messagesToInstall map[string]InstallDef
}

func (c *componentYate) Engine() *Engine {
	return c.engine
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

//CallProgress makes early media rining/progressing
func (c *componentYate) CallProgress(call *Call, msg *Message) error {
	msg.Params["targetid"] = call.ChannelID
	msg.Processed = true

	_, err := c.Engine().Acknowledge(msg)
	if err != nil {
		return err
	}
	c.SendMessage(MsgCallProgress, call, map[string]string{"earlymedia": "true"}, "")
	return nil
}

func (c *componentYate) MessagesToWatch() []string {
	return c.messagesToWatch
}

func (c *componentYate) MessagesToInstall() map[string]InstallDef {
	return c.messagesToInstall
}

func (c *componentYate) InstallMessageHandlers(call *Call) {
	msgs := make(map[string]InstallDef)
	coms := call.Components()
	for _, com := range coms {
		for msgName, msgDef := range com.MessagesToInstall() {
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
	c.logger.Debugf("Going to install [%+v] message handlers, from [%d] components", msgs, len(coms))
	for msgName, msgDef := range msgs {
		filtValue := msgDef.FilterValue
		//replace template variable
		if filtValue == "{channelID}" {
			filtValue = call.ChannelID
		}
		c.engine.InstallFiltered(msgName, msgDef.Priority, msgDef.FilterName, filtValue)
	}
}

func (c *componentYate) InstallMessageWatches(call *Call) {
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
