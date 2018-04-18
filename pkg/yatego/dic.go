package yatego

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rukavina/dicgo"
	"github.com/sirupsen/logrus"
)

func dic() dicgo.Container {
	c := dicgo.NewContainer()

	// classname => factory map
	c.SetValue("component_factories", map[string]ComponentFactory{
		"base":     BaseComponentFactory(c),
		"player":   PlayerComponentFactory(c),
		"recorder": RecorderComponentFactory(c),
		"menu":     MenuComponentFactory(c),
		"fetcher":  FetcherComponentFactory(c),
		"switch":   SwitchComponentFactory(c),
		"http":     HTTPComponentFactory(c),
		"loop":     LoopComponentFactory(c),
	})

	c.SetValue("stderr", os.Stderr)

	c.SetValue("stdout", os.Stdout)

	c.SetValue("stdin", os.Stdin)

	c.SetSingleton("logger", func(cont dicgo.Container) interface{} {
		return &logrus.Logger{
			Out:       cont.Service("stderr").(io.Writer),
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		}
	})

	c.SetSingleton("http_client", func(cont dicgo.Container) interface{} {
		return &http.Client{
			Timeout: time.Second * 10,
		}
	})

	c.SetSingleton("call_manager", func(cont dicgo.Container) interface{} {
		return &CallManager{
			calls: make(map[string]*Call),
		}
	})

	c.SetSingleton("engine", func(cont dicgo.Container) interface{} {
		return &Engine{
			In:     cont.Service("stdin").(io.Reader),
			Out:    cont.Service("stdout").(io.Writer),
			Logger: cont.Service("logger").(Logger),
		}
	})

	c.SetSingleton("controller", func(cont dicgo.Container) interface{} {
		return &Controller{
			componentYate: componentYate{
				componentCommon: componentCommon{
					name:   "controller",
					logger: cont.Service("logger").(Logger),
					config: map[string]interface{}{},
				},
				engine: cont.Service("engine").(*Engine),
			},
			callManager:       cont.Service("call_manager").(*CallManager),
			logger:            cont.Service("logger").(Logger),
			engine:            cont.Service("engine").(*Engine),
			singleChannelMode: true,
			staticComponents:  make([]Component, 0),
		}
	})

	c.SetSingleton("loader_json", func(cont dicgo.Container) interface{} {
		return NewCallflowLoaderJSON("", cont.Service("component_factories").(map[string]ComponentFactory))
	})

	return c
}
