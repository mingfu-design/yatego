package yatego

import (
	"github.com/rukavina/dicgo"
)

// Factory returns factory object
type Factory struct {
	container dicgo.Container
}

// NewFactory factory constructor
func NewFactory() *Factory {
	return &Factory{
		container: dic(),
	}
}

// Container returns DIC container
func (f *Factory) Container() dicgo.Container {
	if f.container == nil {
		f.container = dic()
	}
	return f.container
}

// Controller get controller service instance
func (f *Factory) Controller(cl CallflowLoader) *Controller {
	c := f.Container().Service("controller").(*Controller)
	if cl != nil {
		c.callflowLoader = cl
	}
	return c
}

// BaseComponent generates base component
func (f *Factory) BaseComponent() Component {
	fac := BaseComponentFactory(f.Container())
	return fac("", "start", map[string]interface{}{})
}

// CallflowLoaderJSON get json loader instance
func (f *Factory) CallflowLoaderJSON() *CallflowLoaderJSON {
	return f.Container().Service("loader_json").(*CallflowLoaderJSON)
}
