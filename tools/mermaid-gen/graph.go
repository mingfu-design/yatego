package main

import (
	"strings"

	"github.com/rukavina/yatego/pkg/yatego"
)

//Graph top struct describing one graph
type Graph struct {
	Orientation string
	factories   map[string]ElementFactory
}

//NewGraph constructor
func NewGraph(orientation string) *Graph {
	return &Graph{
		Orientation: orientation,
		factories:   ElementFactories(),
	}
}

//NewElement generate new element in the graph
func (g *Graph) NewElement(c *yatego.CallflowComponent) *Element {
	fac, ok := g.factories[c.ClassName]
	if !ok {
		fac = g.factories["default"]
	}
	return fac(c)
}

//Elements generates graph elements for CF
func (g *Graph) Elements(cf *yatego.Callflow) []*Element {
	elements := []*Element{}
	for _, c := range cf.Components {
		elements = append(elements, g.NewElement(c))
	}
	return elements
}

//Render renders full graph mmd syntax for defined CF
func (g *Graph) Render(c *yatego.Callflow) string {
	elements := g.Elements(c)
	lines := []string{"graph " + g.Orientation, ""}
	for _, e := range elements {
		lines = append(lines, e.Render(), "")
	}
	return strings.Join(lines, "\n")
}
