package yatego

import (
	"strconv"
)

// Loop component transfers based on DTMF
type Loop struct {
	counter int
	Base
}

// NewLoopComponent generates new Loop component
func NewLoopComponent(base Base) *Loop {
	m := &Loop{
		Base: base,
	}
	m.Init()
	return m
}

// Init pseudo constructor
func (l *Loop) Init() {
	l.logger.Debugf("Loop [%s] init, counter: [%s]", l.Name(), l.counter)

	//on enter make choice
	l.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		l.counter++
		l.logger.Debugf("Loop [%s] on enter", l.Name())
		tr, ok := l.Choice(call)
		if !ok {
			return NewCallbackResult(ResStop, "")
		}
		return NewCallbackResult(ResTransfer, tr)
	})
}

//Choice transfers to next component based on setup
func (l *Loop) Choice(call *Call) (string, bool) {
	trDef, exists := l.ConfigAsString("transfer_default")
	if !exists {
		l.logger.Warningf("Loop [%s] has no default transfer", l.Name())
		return "", false
	}
	tr, exists := l.ConfigAsString("transfer")
	if !exists {
		l.logger.Warningf("Loop [%s] has no transfer, transfering to default", l.Name())
		return trDef, true
	}

	max := l.maxCounter(call)

	eq, ok := l.ConfigAsString("break_on_equal")
	breakOnEq := ok && eq == "true"

	//exit the loop
	if l.counter > max || (breakOnEq && l.counter == max) {
		l.logger.Debugf("Loop [%s] counter [%d] reached max [%d]", l.Name(), l.counter, max)
		return trDef, true
	}
	//looping to tr
	l.logger.Debugf("Loop [%s] counter [%d] still under max [%d]", l.Name(), l.counter, max)
	return tr, true
}

func (l *Loop) maxCounter(call *Call) int {
	k, ok := l.ConfigAsString("key")
	//max from call value
	if ok {
		l.logger.Debugf("Loop [%s] getting max from key [%s]", l.Name(), k)
		smax := l.CallDataNamespace(call, k)
		max, err := strconv.Atoi(smax)
		if err == nil && max > 0 {
			return max
		}
	}
	//max from config
	smax, ok := l.ConfigAsString("max")
	if ok {
		max, err := strconv.Atoi(smax)
		if err == nil && max > 0 {
			return max
		}
	}
	return 1
}
