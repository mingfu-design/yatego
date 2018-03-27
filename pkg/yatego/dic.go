package yatego

import (
	"io"
	"os"

	"github.com/rukavina/minidic"
	"github.com/sirupsen/logrus"
)

func dic() minidic.Container {
	c := minidic.NewContainer()

	// classname => factory map
	c.Add(minidic.NewInjection("component_factories", func(cont minidic.Container) map[string]ComponentFactory {
		return map[string]ComponentFactory{
			"base":     BaseComponentFactory(cont),
			"player":   PlayerComponentFactory(cont),
			"recorder": RecorderComponentFactory(cont),
			"menu":     MenuComponentFactory(cont),
		}
	}))

	c.Add(minidic.NewInjection("stderr", func(cont minidic.Container) io.Writer {
		return os.Stderr
	}))

	c.Add(minidic.NewInjection("stdout", func(cont minidic.Container) io.Writer {
		return os.Stdout
	}))

	c.Add(minidic.NewInjection("stdin", func(cont minidic.Container) io.Reader {
		return os.Stdin
	}))

	c.Add(minidic.NewInjection("logger", func(cont minidic.Container) Logger {
		return &logrus.Logger{
			Out:       cont.Get("stderr").(io.Writer),
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		}
	}))

	c.Add(minidic.NewInjection("call_manager", func(cont minidic.Container) *CallManager {
		return &CallManager{
			calls: make(map[string]*Call),
		}
	}))

	c.Add(minidic.NewInjection("engine", func(cont minidic.Container) *Engine {
		return &Engine{
			In:     cont.Get("stdin").(io.Reader),
			Out:    cont.Get("stdout").(io.Writer),
			Logger: cont.Get("logger").(Logger),
		}
	}))

	c.Add(minidic.NewInjection("controller", func(cont minidic.Container) *Controller {
		return &Controller{
			callManager:       cont.Get("call_manager").(*CallManager),
			logger:            cont.Get("logger").(Logger),
			engine:            cont.Get("engine").(*Engine),
			singleChannelMode: true,
			staticComponents:  make([]Component, 0),
		}
	}))

	c.Add(minidic.NewInjection("loader_json", func(cont minidic.Container) *CallflowLoaderJSON {
		return NewCallflowLoaderJSON("", cont.Get("component_factories").(map[string]ComponentFactory))
	}))

	return c
}
