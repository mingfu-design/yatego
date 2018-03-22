package callflow

type CallflowComponent struct {
	Name      string
	ClassName string
	Config    map[string]interface{}
}

type Callflow struct {
	Components []*CallflowComponent
}
