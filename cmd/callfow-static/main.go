package main

import (
	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	f := yatego.NewFactory()
	c := f.Controller(loader(f.Container()))
	c.Logger().Debug("Starting yatego IVR [callflow-static]")
	c.Run("")
}
