package yatego

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

//Engine is communication object with Yate server
type Engine struct {
	In        io.Reader
	Out       io.Writer
	Err       io.Writer
	scanner   *bufio.Scanner
	LogPrefix string
	DebugOn   bool
}

//Dispatch message to yate engine
func (engine *Engine) Dispatch(m *Message) (int, error) {
	if m.Type != TypeOutgoing {
		return 0, errors.New("cannot dispatch non outgoing message type")
	}
	s := m.Encode()
	m.Type = TypeDispatched
	return engine.print(s)
}

//Acknowledge message to yate engine
func (engine *Engine) Acknowledge(m *Message) (int, error) {
	if m.Type != TypeIncoming {
		return 0, errors.New("cannot acknowledge non incoming message type")
	}
	s := m.Encode()
	m.Type = TypeAcknowledged
	return engine.print(s)
}

//GetEvent gets new event message from stdin
func (engine *Engine) GetEvent() (*Message, error) {
	if engine.scanner == nil {
		engine.scanner = bufio.NewScanner(engine.In)
	}
	res := engine.scanner.Scan()
	if !res {
		return nil, engine.scanner.Err()
	}
	s := engine.scanner.Text()
	if s == "" {
		return nil, nil
	}
	engine.Debug("<<< received raw message [" + s + "]")
	return DecodeMessage(s)
}

//Install event handler
func (engine *Engine) Install(event string, priority int) {
	engine.InstallFiltered(event, priority, "", "")
}

//InstallFiltered listens on events filtered
func (engine *Engine) InstallFiltered(event string, priority int, filtname, filtvalue string) {
	var filter string
	if filtname != "" {
		filter = fmt.Sprintf(":%s:%s", filtname, filtvalue)
	}

	msg := "%%" + fmt.Sprintf(">install:%d:%s%s", priority, esc(event), filter)
	go engine.print(msg)
}

//Uninstall event handler
func (engine *Engine) Uninstall(event string) {
	msg := "%%" + fmt.Sprintf(">uninstall:%s", esc(event))
	go engine.print(msg)
}

//Watch particular event
func (engine *Engine) Watch(event string) {
	msg := "%%" + fmt.Sprintf(">watch:%s", esc(event))
	go engine.print(msg)
}

//Unwatch event handler
func (engine *Engine) Unwatch(event string) {
	msg := "%%" + fmt.Sprintf(">unwatch:%s", esc(event))
	go engine.print(msg)
}

//SetLocal variable
func (engine *Engine) SetLocal(name, value string) {
	msg := "%%" + fmt.Sprintf(">setlocal:%s%s", esc(name), esc(value))
	go engine.print(msg)
}

//Log to yate output
func (engine *Engine) Log(s string) (int, error) {
	return fmt.Fprintln(engine.Err, engine.LogPrefix+s)
}

//Debug logs if DebugOn is true
func (engine *Engine) Debug(s string) (int, error) {
	if !engine.DebugOn {
		return 0, nil
	}
	return engine.Log(s)
}

//NewCallID returns new random call id string
func NewCallID() string {
	return randString(10)
}

func (engine *Engine) print(str string) (int, error) {
	engine.Debug(">>> sending message [" + str + "]")
	return fmt.Fprintln(engine.Out, str)
}
