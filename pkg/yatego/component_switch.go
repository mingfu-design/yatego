package yatego

import (
	"strings"
)

// Switch component transfers based on DTMF
type Switch struct {
	Base
}

// NewSwitchComponent generates new Switch component
func NewSwitchComponent(base Base) *Switch {
	m := &Switch{
		Base: base,
	}
	m.Init()
	return m
}

// Init pseudo constructor
func (s *Switch) Init() {
	s.logger.Debugf("Switch [%s] init", s.Name())

	//on enter make choice
	s.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		s.logger.Debugf("Switch [%s] on enter", s.Name())
		tr, ok := s.Choice(call)
		if !ok {
			return NewCallbackResult(ResStop, "")
		}
		return NewCallbackResult(ResTransfer, tr)
	})
}

//Choice transfers to next component based on setup
func (s *Switch) Choice(call *Call) (string, bool) {
	ks, exists := s.ConfigAsString("keys")
	if !exists {
		s.logger.Warningf("Switch [%s] has no keys defined", s.Name())
		return "", false
	}
	ts, exists := s.ConfigAsString("transfer")
	if !exists {
		s.logger.Warningf("Switch [%s] has no transfers defined", s.Name())
		return "", false
	}
	vs, exists := s.ConfigAsString("values")
	if !exists {
		s.logger.Warningf("Switch [%s] has no values defined", s.Name())
		return "", false
	}
	keys := strings.Split(ks, ",")
	transfers := strings.Split(ts, ",")
	vals := strings.Split(vs, ",")

	s.logger.Debugf("[%s] Making choice for keys %v, values %v against existing call data: %+v", s.Name(), keys, vals, call.DataAll())
	for i, k := range keys {
		if i >= len(vals) || i >= len(transfers) {
			s.logger.Warningf("Switch [%s] has no transfer or value defined for key [%s]", s.Name(), k)
			break
		}
		if vals[i] != s.CallDataNamespace(call, k) {
			continue
		}
		s.logger.Debugf("Switch [%s] found val [%s] on key [%s] which leads to transfer [%s]", s.Name(), vals[i], k, transfers[i])
		s.SetCallData(call, "key", k)
		s.SetCallData(call, "val", vals[i])
		return transfers[i], true
	}
	tr, exists := s.ConfigAsString("transfer_default")
	if !exists {
		s.logger.Warningf("Switch [%s] has no default transfer", s.Name())
		return "", false
	}
	s.logger.Warningf("Switch [%s] transfer default [%s]", s.Name(), tr)
	return tr, true
}
