package yatego

// CallflowLoader interface which defines object to be able to load new callflow
type CallflowLoader interface {
	load(params map[string]string) *Callflow
}
