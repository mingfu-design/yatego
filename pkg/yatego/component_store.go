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
	ks, exists := s.ConfigAsString("to_keys")
	if !exists {
		s.logger.Warningf("Store [%s] has no keys defined", s.Name())
		return "", false
	}
	vs, exists := s.ConfigAsString("from_values")
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
		s.storeKeyValue(call, key, vals[i])
	}
	s.logger.Debugf("Store [%s] transfer default [%s]", s.Name(), tr)
	return tr, true
}

func (s *Store) storeKeyValue(call *Call, key string, val string) {
	keyNs := strings.Split(key, ".")
	//save under component's key
	if len(keyNs) == 1 {
		s.SetCallData(call, key, val)
		return
	}
	if len(keyNs) == 2 {
		call.SetData(keyNs[0], keyNs[1], val)
		return
	}
	s.logger.Errorf("Store [%s] cannot save invalid key defined [%s]", s.Name(), key)
	return
}
