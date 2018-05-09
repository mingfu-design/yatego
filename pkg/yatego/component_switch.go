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
	k, keyOk := s.ConfigAsString("compare_key")
	v, valOk := s.ConfigAsString("compare_val")
	if !keyOk && !valOk {
		s.logger.Warningf("Switch [%s] has no compare key nor compare val defined", s.Name())
		return "", false
	}
	cmp := ""
	if keyOk {
		cmp = s.CallDataNamespace(call, k)
	} else {
		cmp = v
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

	transfers := strings.Split(ts, ",")
	vals := strings.Split(vs, ",")

	s.logger.Debugf("[%s] Making choice comparing value [%s], to values %v against existing call data: %+v", s.Name(), cmp, vals, call.DataAll())
	for i, val := range vals {
		if i >= len(transfers) {
			s.logger.Warningf("Switch [%s] has no transfer for value [%s]", s.Name(), val)
			break
		}
		if val != cmp {
			continue
		}
		s.logger.Debugf("Switch [%s] found val [%s] which leads to transfer [%s]", s.Name(), val, transfers[i])
		s.SetCallData(call, "val", val)
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
