package yatego

// Stop component just halts CF
type Stop struct {
	Base
}

// NewStopComponent generates new Switch component
func NewStopComponent(base Base) *Stop {
	m := &Stop{
		Base: base,
	}
	m.Init()
	return m
}

// Init pseudo constructor
func (s *Stop) Init() {
	s.logger.Debugf("Stop [%s] init", s.Name())

	//on enter just return to stop
	s.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		s.logger.Debugf("Stop [%s] on enter", s.Name())
		return NewCallbackResult(ResStop, "")
	})
}
