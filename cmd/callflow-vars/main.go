package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	//get current dir
	exec, _ := os.Executable()
	dir := filepath.Dir(exec)

	f := yatego.NewFactory()
	//json loader
	l := f.CallflowLoaderJSON()
	//custom onload, pull vars from json
	l.OnLoad = func(loader *yatego.CallflowLoaderJSON, cf *yatego.Callflow, params map[string]string) error {
		//load vars from json
		data, err := ioutil.ReadFile(dir + "/assets/configs/callflow_vars.json")
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
	//load json CF template from external file
	l.SetJSONFile(dir + "/assets/configs/callflow_tpl.json")

	//controller
	c := f.Controller(l)
	c.Logger().Debug("Starting yatego IVR [callflow-vars]")
	c.Run("")
}
