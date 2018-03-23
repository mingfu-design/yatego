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
