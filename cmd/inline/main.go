package main

import "github.com/rukavina/yatego/pkg/yatego"

func main() {
	f := yatego.NewFactory()
	c := f.Controller(nil)
	c.Logger().Debug("Starting yatego IVR [inline]")

	com := f.BaseComponent()
	com.OnEnter(func(call *yatego.Call, message *yatego.Message) *yatego.CallbackResult {
		com.(*yatego.Base).PlayTone("congestion", call, map[string]string{})
		return yatego.NewCallbackResult(yatego.ResStay, "")
	})

	c.AddStaticComponent(com)
	c.Run("")
}
