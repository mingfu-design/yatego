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
package main

import (
	"fmt"
	"os"

	"github.com/rukavina/yatego"
)

var (
	engine                                        yatego.Engine
	caller, called, billID, peerCallID, ourCallID string
)

func main() {
	engine = yatego.Engine{
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
		if m.Type == yatego.TypeIncoming && m.Name == "call.execute" {
			initCallParams(m)
			callExecute(m)
		}
		//we need to ack all incoming messages
		if m.Type == yatego.TypeIncoming {
			engine.Acknowledge(m)
		}
	}
}

func initCallParams(m *yatego.Message) {
	ourCallID = "yatego/" + yatego.NewCallID()

	peerCallID = m.Params["id"]
	billID = m.Params["billid"]
	caller = m.Params["caller"]
	called = m.Params["called"]
}

func callExecute(m *yatego.Message) {
	m.Processed = true
	m.Params["targetid"] = peerCallID
	_, err := engine.Acknowledge(m)
	if err != nil {
		engine.Log("acknowledge error: " + err.Error())
	} else {
		engine.Log("acknowledge msg: " + fmt.Sprintf("%+v\n", m))
	}

	msgAnswer := yatego.NewMessage("call.answered", map[string]string{"id": ourCallID, "targetid": peerCallID})
	_, err = engine.Dispatch(msgAnswer)
	if err != nil {
		engine.Log("dispatched err: " + err.Error())
	} else {
		engine.Log("dispatched msg: " + fmt.Sprintf("%+v\n", msgAnswer))
	}

	msgAttach := yatego.NewMessage("chan.attach", map[string]string{"source": "tone/congestion", "notify": ourCallID})
	_, err = engine.Dispatch(msgAttach)
	if err != nil {
		engine.Log("dispatched err: " + err.Error())
	} else {
		engine.Log("dispatched msg: " + fmt.Sprintf("%+v\n", msgAttach))
	}
}

```