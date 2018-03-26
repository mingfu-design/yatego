package yatego

import "github.com/rukavina/minidic"

// BaseComponentFactory is base component factory
func BaseComponentFactory(c minidic.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		com := baseComponent(c, name, config)
		com.Init()
		return &com
	}
}

// PlayerComponentFactory is Player component factory
func PlayerComponentFactory(c minidic.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		com := Player{
			currSong: 0,
			Base:     baseComponent(c, name, config),
		}
		com.Init()
		return &com
	}
}

// RecorderComponentFactory is Recorder component factory
func RecorderComponentFactory(c minidic.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		com := Recorder{
			status: stPrompt,
			Base:   baseComponent(c, name, config),
		}
		com.Init()
		return &com
	}
}

// MenuComponentFactory is Menu component factory
func MenuComponentFactory(c minidic.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		com := Menu{
			Base: baseComponent(c, name, config),
		}
		com.Init()
		return &com
	}
}

// baseComponent helper function get base object by value
func baseComponent(c minidic.Container, name string, config map[string]interface{}) Base {
	base := NewBaseComponent(name, c.Get("engine").(*Engine), c.Get("logger").(Logger), config)
	return *base
}
