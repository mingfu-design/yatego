package yatego

import (
	"strings"
)

// Menu component transfers based on DTMF
type Menu struct {
	Base
}

// NewMenuComponent generates new Menu component
func NewMenuComponent(base Base) *Menu {
	m := &Menu{
		Base: base,
	}
	m.Init()
	return m
}

// Init pseudo constructor
func (m *Menu) Init() {
	m.logger.Debugf("Menu [%s] init", m.Name())
	//install chan.dtml to listen clicks
	m.messagesToInstall[MsgChanDtmf] = InstallDef{Priority: 100}

	//on chan.dtmf
	m.Listen(MsgChanDtmf, func(call *Call, msg *Message) *CallbackResult {
		msg.Processed = true
		text, exists := msg.Params["text"]
		if !exists || text == "" {
			m.logger.Warningf("No text found in [%s]", m.Name())
			return NewCallbackResult(ResStay, "")
		}
		m.logger.Debugf("Chan.dtmf with text [%s] in [%s]", text, m.Name())
		tr, done := m.Pressed(string(text[0]), call)
		if done {
			return NewCallbackResult(ResTransfer, tr)
		}
		return NewCallbackResult(ResStay, "")
	})

	//play prompt if defined
	m.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		prompt, ok := m.ConfigAsString("prompt")
		if !ok {
			return NewCallbackResult(ResStay, "")
		}
		m.logger.Infof("Menu [%s] has prompt [%s] defined, playing it now", m.Name(), prompt)
		m.PlayWave(prompt, call, map[string]string{})
		return NewCallbackResult(ResStay, "")
	})
}

// Pressed returns transfer component if defined
func (m *Menu) Pressed(key string, call *Call) (string, bool) {
	ks, exists := m.Config("keys")
	if !exists {
		m.logger.Warningf("Menu [%s] has no keys defined", m.Name())
		return "", false
	}
	ts, exists := m.Config("transfer")
	if !exists {
		m.logger.Warningf("Menu [%s] has no transfers defined", m.Name())
		return "", false
	}
	keys := strings.Split(ks.(string), ",")
	transfers := strings.Split(ts.(string), ",")
	for i, k := range keys {
		if k != key {
			continue
		}
		m.SetCallData(call, "key", key)

		if i < len(transfers) {
			return transfers[i], true
		}

		if len(transfers) > 0 {
			m.logger.Debugf("Menu [%s] has no enough transfers, using first [%s], for key [%s]", m.Name(), transfers[0], key)
			return transfers[0], true
		}
		m.logger.Warningf("Menu [%s] has no transfer defined for key [%s]", m.Name(), key)
		return "", false
	}
	tr, exists := m.ConfigAsString("transfer_default")
	if !exists {
		m.logger.Warningf("Menu [%s] has no default transfer", m.Name())
		return "", false
	}
	m.logger.Warningf("Menu [%s] has no option defined for key [%s], but transfer to default [%s]", m.Name(), key, tr)
	return tr, true
}
