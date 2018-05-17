package main

import (
	"strings"

	"github.com/rukavina/yatego/pkg/yatego"
)

//ElementFactory is factory method to produce concrete element
type ElementFactory func(c *yatego.CallflowComponent) *Element

//ElementFactories returns map of factories per cf component class
func ElementFactories() map[string]ElementFactory {
	return map[string]ElementFactory{
		"default": func(c *yatego.CallflowComponent) *Element {
			return NewElement(c, "[", "]")
		},
		"switch": func(c *yatego.CallflowComponent) *Element {
			return NewElement(c, "{", "}")
		},
		"menu": func(c *yatego.CallflowComponent) *Element {
			return NewElement(c, ">", "]")
		},
		"player": func(c *yatego.CallflowComponent) *Element {
			return NewElement(c, "(", ")")
		},
	}
}

//Element is base element struct
type Element struct {
	Component *yatego.CallflowComponent
	LTag      string
	RTag      string
}

//NewElement constructs new base element
func NewElement(c *yatego.CallflowComponent, LTag string, RTag string) *Element {
	if LTag == "" {
		LTag = "["
	}
	if LTag == "" {
		LTag = "]"
	}
	return &Element{
		Component: c,
		LTag:      LTag,
		RTag:      RTag,
	}
}

//Render returns string for the element
func (e *Element) Render() string {
	return strings.Join(e.renderConnections(), "\n")
}

func (e *Element) renderConnections() []string {
	res := e.renderConnectionsRegular()
	defConn := e.renderConnectionDefault()
	if defConn != "" {
		res = append(res, defConn)
	}
	return res
}

func (e *Element) renderConnectionsRegular() []string {
	res := []string{}
	conf := e.Component.Config
	target, ok := conf["transfer"]
	if !ok {
		return res
	}
	targets := strings.Split(target.(string), ",")
	for i, t := range targets {
		res = append(res, e.renderConnection(t, e.connectionName(i), false))
	}
	return res
}

func (e *Element) renderConnection(targetName string, connName string, isDefault bool) string {
	connStr := "-->"
	if isDefault {
		connStr = "-.->"
	}
	res := e.renderName() + " " + connStr
	if connName == "" {
		return res + " " + targetName
	}
	return res + "|" + connName + "| " + targetName
}

func (e *Element) renderName() string {
	name := e.Component.Name
	content := name + ": " + e.Component.ClassName
	return name + e.LTag + content + e.RTag
}

func (e *Element) renderConnectionDefault() string {
	conf := e.Component.Config
	target, ok := conf["transfer_default"]
	if !ok {
		return ""
	}
	return e.renderConnection(target.(string), "default", true)
}

func (e *Element) connectionName(connIndex int) string {
	conf := e.Component.Config
	values, ok := conf["values"]
	if !ok {
		return ""
	}
	vals := strings.Split(values.(string), ",")
	if connIndex >= len(vals) {
		return ""
	}
	return vals[connIndex]
}
