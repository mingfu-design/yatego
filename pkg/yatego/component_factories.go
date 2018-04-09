package yatego

import (
	"github.com/rukavina/dicgo"
)

// BaseComponentFactory is base component factory
func BaseComponentFactory(c dicgo.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		com := baseComponent(c, name, config)
		return &com
	}
}

// PlayerComponentFactory is Player component factory
func PlayerComponentFactory(c dicgo.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		return NewPlayerComponent(baseComponent(c, name, config))
	}
}

// RecorderComponentFactory is Recorder component factory
func RecorderComponentFactory(c dicgo.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		return NewRecorderComponent(baseComponent(c, name, config))
	}
}

// MenuComponentFactory is Menu component factory
func MenuComponentFactory(c dicgo.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		return NewMenuComponent(baseComponent(c, name, config))
	}
}

// FetcherComponentFactory is Fetcher component factory
func FetcherComponentFactory(c dicgo.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		return NewFetcherComponent(baseComponent(c, name, config), c.Service("loader_json").(*CallflowLoaderJSON))
	}
}

// baseComponent helper function get base object by value
func baseComponent(c dicgo.Container, name string, config map[string]interface{}) Base {
	base := NewBaseComponent(name, c.Service("engine").(*Engine), c.Service("logger").(Logger), config)
	return *base
}
