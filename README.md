# YateGo
YateGo packes enables communication with YATE core via external protocol http://docs.yate.ro/wiki/Programmer%27s_guide

## Installation

`go get github.com/rukavina/yatego`

## Sample IVR

this simple IVR just answers and plays `tone/congestion`

In order to install it as regalar callflow/ivr in yate, compile, place binary in `share/yate/scripts`.
Then you need to have a route to it, eg:

```sh
^555$=external/nodata/yourbinary
```

```golang
package yatego

import (
	"fmt"
	"os"
)

var (
	engine                                         Engine
	caller, called, billID, partyCallID, ourCallID string
)

func main() {
	engine = Engine{
		In:        os.Stdin,
		Out:       os.Stdout,
		Err:       os.Stderr,
		DebugOn:   true,
		LogPrefix: "[GoYate Test] ",
	}

	engine.Log("start")
	engine.Install("chan.notify", 20)
	for {
		m, err := engine.GetEvent()
		if err != nil || m == nil {
			break
		}
		engine.Log("new message: " + fmt.Sprintf("%+v\n", m))
		if m.Type == TypeIncoming && m.Name == "call.execute" {
			go callExecute(m)
		}
	}
}


func callExecute(m *Message) {
	m.Handled = true
	partyCallID = m.Params["id"]
	billID = m.Params["billid"]
	ourCallID = NewCallID()
	m.Params["targetid"] = ourCallID
	caller = m.Params["caller"]
	called = m.Params["called"]
	_, err := engine.Acknowledge(m)
	if err != nil {
		engine.Log("acknowledge error: " + err.Error())
	} else {
		engine.Log("acknowledge msg: " + fmt.Sprintf("%+v\n", m))
	}

	msgAnswer := NewMessage("call.answered", map[string]string{"id": ourCallID, "targetid": partyCallID})
	_, err = engine.Dispatch(msgAnswer)
	if err != nil {
		engine.Log("dispatched err: " + err.Error())
	} else {
		engine.Log("dispatched msg: " + fmt.Sprintf("%+v\n", msgAnswer))
	}

	msgAttach := NewMessage("chan.attach", map[string]string{"source": "tone/congestion", "notify": ourCallID})
	_, err = engine.Dispatch(msgAttach)
	if err != nil {
		engine.Log("dispatched err: " + err.Error())
	} else {
		engine.Log("dispatched msg: " + fmt.Sprintf("%+v\n", msgAttach))
	}
}
```