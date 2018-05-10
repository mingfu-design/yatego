package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	exec, _ := os.Executable()
	dir := filepath.Dir(exec)

	f := yatego.NewFactory()

	//load static config
	config, err := LoadConfiguration(dir + "/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	//correct relating path
	_, prs := config["log_file"]
	if prs && string(config["log_file"][0]) != "/" {
		config["log_file"] = dir + "/" + config["log_file"]
	}
	f.Container().SetValue("config", config)

	//json callflow loader
	l := f.CallflowLoaderJSON()
	//callflow http vars loader

	cfVars := &HTTPCFVarsLoader{
		Config:     config,
		Logger:     f.Container().Service("logger").(yatego.Logger),
		HTTPClient: f.Container().Service("http_client").(*http.Client),
	}
	//custom onload, pull CF vars from json
	l.OnLoad = func(loader *yatego.CallflowLoaderJSON, cf *yatego.Callflow, params map[string]string) error {
		//load vars from http
		vars, err := cfVars.LoadCallflowVars(params)
		if err != nil {
			return err
		}
		//append vars from config
		for k, v := range config {
			_, exists := vars[k]
			if !exists {
				vars[k] = v
			}
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
