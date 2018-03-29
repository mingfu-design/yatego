package main

import (
	"os"
	"path/filepath"

	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	f := yatego.NewFactory()

	//json loader
	l := f.CallflowLoaderJSON()
	//load json content from external file
	exec, _ := os.Executable()
	dir := filepath.Dir(exec)
	l.SetJSONFile(dir + "/assets/configs/callflow_dynamic.json")

	//controller
	c := f.Controller(l)
	c.Logger().Debug("Starting yatego IVR [callflow-dynamic]")
	c.Run("")
}
