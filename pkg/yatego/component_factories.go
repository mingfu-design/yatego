package yatego

import "github.com/rukavina/minidic"

// BaseComponentFactory is base component factory
func BaseComponentFactory(c minidic.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		com := baseComponent(c, name, config)
		return &com
	}
}

// PlayerComponentFactory is Player component factory
func PlayerComponentFactory(c minidic.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		return NewPlayerComponent(baseComponent(c, name, config))
	}
}

// RecorderComponentFactory is Recorder component factory
func RecorderComponentFactory(c minidic.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		return NewRecorderComponent(baseComponent(c, name, config))
	}
}

// MenuComponentFactory is Menu component factory
func MenuComponentFactory(c minidic.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		return NewMenuComponent(baseComponent(c, name, config))
	}
}

// baseComponent helper function get base object by value
func baseComponent(c minidic.Container, name string, config map[string]interface{}) Base {
	base := NewBaseComponent(name, c.Get("engine").(*Engine), c.Get("logger").(Logger), config)
	return *base
}
