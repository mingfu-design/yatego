package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	exec, _ := os.Executable()
	dir := filepath.Dir(exec)

	f := yatego.NewFactory()

	//set config
	config := f.Container().Service("config").(map[string]string)
	config["log_file"] = "./app.log"

	//json loader
	l := f.CallflowLoaderJSON()
	//custom onload, pull vars from json
	l.OnLoad = func(loader *yatego.CallflowLoaderJSON, cf *yatego.Callflow, params map[string]string) error {
		//load vars from json
		data, err := ioutil.ReadFile(dir + "/assets/configs/callflow_http_switch_vars.json")
		if err != nil {
			return err
		}
		var vars map[string]string
		if err := json.Unmarshal(data, &vars); err != nil {
			return err
		}
		loader.SetVars(vars)
		return nil
	}
	//load json content from external file
	l.SetJSONFile(dir + "/assets/configs/callflow_http_switch.json")

	//controller
	c := f.Controller(l)
	c.Logger().Debug("Starting yatego IVR [callflow-http-switch]")
	c.Run("")
}
