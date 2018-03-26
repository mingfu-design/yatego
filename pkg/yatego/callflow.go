package yatego

// CallflowComponent is the definition of a single callflow component
type CallflowComponent struct {
	Name      string
	ClassName string
	Config    map[string]interface{}
	Factory   ComponentFactory
}

// Callflow is the definition of a IVR callflow and components def. repos
type Callflow struct {
	Components []*CallflowComponent
}

// CallflowLoader interface which defines object to be able to load new callflow
type CallflowLoader interface {
	Load(params map[string]string) *Callflow
}

// CallflowLoaderStatic is simplest CallflowLoader implementation
type CallflowLoaderStatic struct {
	callflow *Callflow
}

// Load callflow
func (cl *CallflowLoaderStatic) Load(params map[string]string) *Callflow {
	return cl.callflow
}

// NewCallflowLoaderStatic generates new CallflowLoaderStatic
func NewCallflowLoaderStatic(c *Callflow) *CallflowLoaderStatic {
	return &CallflowLoaderStatic{
		callflow: c,
	}
}
