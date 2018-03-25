package yatego

import "github.com/rukavina/minidic"

// Factory returns factory object
type Factory struct {
	container minidic.Container
}

// NewFactory factory constructor
func NewFactory() *Factory {
	return &Factory{
		container: dic(),
	}
}

// Container returns DIC container
func (f *Factory) Container() minidic.Container {
	if f.container == nil {
		f.container = dic()
	}
	return f.container
}

// Controller get controller service instance
func (f *Factory) Controller() *Controller {
	return f.Container().Get("controller").(*Controller)
}

// BaseComponent generates base component
func (f *Factory) BaseComponent() Component {
	fac := BaseComponentFactory(f.Container())
	return fac("", "start", map[string]interface{}{})
}

// BaseComponentFactory is base component factory
func BaseComponentFactory(c minidic.Container) ComponentFactory {
	return func(class string, name string, config map[string]interface{}) Component {
		return NewBaseComponent(name, c.Get("engine").(*Engine), c.Get("logger").(Logger), config)
	}
}
