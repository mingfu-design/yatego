package main

import (
	"os"
	"path/filepath"

	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	f := yatego.NewFactory()

	//set config
	//config := f.Container().Service("config").(map[string]string)
	//config["log_file"] = "./app.log"

	//json loader
	l := f.CallflowLoaderJSON()
	//load json content from external file
	exec, _ := os.Executable()
	dir := filepath.Dir(exec)
	l.SetJSONFile(dir + "/assets/configs/callflow_http_switch.json")

	//controller
	c := f.Controller(l)
	c.Logger().Debug("Starting yatego IVR [callflow-http-switch]")
	c.Run("")
}
