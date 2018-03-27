package main

import (
	"os"
	"path/filepath"

	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	f := yatego.NewFactory()

	//json loader
	l := yatego.NewCallflowLoaderJSON("", f.Container().Get("component_factories").(map[string]yatego.ComponentFactory))
	//load json content from external file
	exec, _ := os.Executable()
	dir := filepath.Dir(exec)
	l.SetJSONFile(dir + "/assets/configs/callflow_static.json")

	//controller
	c := f.Controller(l)
	c.Logger().Debug("Starting yatego IVR [callflow-json]")
	c.Run("")
}
