package yatego

import (
	"strings"
)

// Store component stores literals as call data vals
type Store struct {
	Base
}

// NewStoreComponent generates new Store component
func NewStoreComponent(base Base) *Store {
	m := &Store{
		Base: base,
	}
	m.Init()
	return m
}

// Init pseudo constructor
func (s *Store) Init() {
	s.logger.Debugf("Store [%s] init", s.Name())

	//on enter make choice
	s.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		s.logger.Debugf("Store [%s] on enter", s.Name())
		tr, ok := s.Stores(call)
		if !ok {
			return NewCallbackResult(ResStop, "")
		}
		return NewCallbackResult(ResTransfer, tr)
	})
}

//Stores stores data and transfers to next component based on setup
func (s *Store) Stores(call *Call) (string, bool) {
	tr, exists := s.ConfigAsString("transfer")
	if !exists {
		s.logger.Warningf("Store [%s] has no transfers defined", s.Name())
		return "", false
	}
	ks, exists := s.ConfigAsString("keys")
	if !exists {
		s.logger.Warningf("Store [%s] has no keys defined", s.Name())
		return "", false
	}
	vs, exists := s.ConfigAsString("values")
	if !exists {
		s.logger.Warningf("Store [%s] has no values defined", s.Name())
		return "", false
	}
	keys := strings.Split(ks, ",")
	vals := strings.Split(vs, ",")

	for i, key := range keys {
		if i >= len(vals) {
			s.logger.Warningf("Store [%s] has no val defined for key [%s]", s.Name(), key)
			break
		}
		s.SetCallData(call, key, vals[i])
	}
	s.logger.Debugf("Store [%s] transfer default [%s]", s.Name(), tr)
	return tr, true
}
