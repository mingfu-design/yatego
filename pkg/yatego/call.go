package yatego

import (
	"fmt"
	"strconv"
	"strings"
)

// Call is object persisting all single call related info, including callflow components
type Call struct {
	ChannelID           string
	PeerID              string
	Caller              string
	CallerName          string
	Called              string
	BillingID           string
	ActiveComponentName string
	data                map[string]map[string]interface{}
	components          []Component
}

// CallData returns global call data, like Caller, Called etc.
func (call *Call) CallData() map[string]interface{} {
	return map[string]interface{}{
		"channelId":  call.ChannelID,
		"peerId":     call.PeerID,
		"caller":     call.Caller,
		"callerName": call.CallerName,
		"called":     call.Called,
		"billingId":  call.BillingID,
	}
}

// Data returns the component's data. If key is present, returns data subkey
func (call *Call) Data(componentName string, key string) (interface{}, bool) {
	var (
		data   interface{}
		exists bool
	)
	data, exists = call.data[componentName]
	if !exists {
		return data, exists
	}
	if key == "" {
		return data, exists
	}
	data, exists = call.data[componentName][key]
	return data, exists
}

// DataAll returns all components' data
func (call *Call) DataAll() map[string]map[string]interface{} {
	return call.data
}

// SetData sets the component's data key value
func (call *Call) SetData(componentName string, key string, value interface{}) {
	_, exists := call.data[componentName]
	if !exists {
		call.data[componentName] = make(map[string]interface{})
	}
	call.data[componentName][key] = value
}

// ParseConfig updates a component's config tpl variables {component.variable}
// per stored call data values
func (call *Call) ParseConfig(c Component) {
	keys := c.ConfigKeys()
	s := ""
	rep := call.dataStrReplacer()
	//loop all config keys
	for _, key := range keys {
		v, exists := c.Config(key)
		if !exists {
			continue
		}
		//mind only string config vals
		switch v.(type) {
		default:
			continue
		case string:
			s = v.(string)
		}
		s := rep.Replace(s)
		c.SetConfig(key, s)
	}
}

func (call *Call) dataStrReplacer() *strings.Replacer {
	data := call.DataAll()
	pairs := []string{}
	curr := ""
	for comp, vals := range data {
		for key, val := range vals {
			switch t := val.(type) {
			default:
				continue
			case string:
				curr = t
			case bool:
				curr = strconv.FormatBool(t)
			case int:
				curr = strconv.Itoa(t)
			case uint64:
				curr = strconv.FormatUint(t, 10)
			case float64:
				curr = strconv.FormatFloat(t, 'f', 6, 64)
			}

			pairs = append(pairs, "{"+comp+"."+key+"}", curr)
		}
	}
	return strings.NewReplacer(pairs...)
}

// Components returns all components define
func (call *Call) Components() []Component {
	return call.components
}

// Component return named or first component if name is ""
func (call *Call) Component(name string) Component {
	for _, component := range call.components {
		if name == "" || component.Name() == name {
			return component
		}
	}
	return nil
}

// AddComponent appends new component
func (call *Call) AddComponent(component Component) {
	call.components = append(call.components, component)
}

// ActivateComponent transfer to new active component
func (call *Call) ActivateComponent(newComponent string) bool {
	next := call.Component(newComponent)
	if next == nil {
		return false
	}
	call.ActiveComponentName = next.Name()
	//prepare dynamic config variables, eg. all config values in the form of
	//{component.key} and replaced with actual values from call data
	call.ParseConfig(call.ActiveComponent())
	return true
}

// ActiveComponent returns active compoment in the call
func (call *Call) ActiveComponent() Component {
	if call.ActiveComponentName == "" {
		return nil
	}
	return call.Component(call.ActiveComponentName)
}

// CallManager is the repository for calls
type CallManager struct {
	calls map[string]*Call
}

// Calls returns all calls in a map channelId => *Call
func (cm *CallManager) Calls() map[string]*Call {
	return cm.calls
}

// Call returns the call in the channel
func (cm *CallManager) Call(channelID string) (*Call, bool) {
	call, exists := cm.calls[channelID]
	return call, exists
}

// Remove deletes the Call for the channel
func (cm *CallManager) Remove(channelID string) {
	delete(cm.calls, channelID)
}

// Add new Call to be tracked
func (cm *CallManager) Add(
	components []Component,
	params map[string]string,
	channelID string,
	activeComponentName string) (*Call, error) {
	if channelID == "" {
		channelID = "yatego/" + NewCallID()
	}
	if _, exists := cm.calls[channelID]; exists {
		return nil, fmt.Errorf("Channel [%s] already exists", channelID)
	}
	call := &Call{
		ChannelID:           channelID,
		ActiveComponentName: activeComponentName,
		data:                make(map[string]map[string]interface{}),
		components:          components,
	}
	if _, exists := params["id"]; exists {
		call.PeerID = params["id"]
	}
	if _, exists := params["billid"]; exists {
		call.BillingID = params["billid"]
	}
	if _, exists := params["caller"]; exists {
		call.Caller = params["caller"]
	}
	if _, exists := params["callername"]; exists {
		call.CallerName = params["callername"]
	}
	if _, exists := params["called"]; exists {
		call.Called = params["called"]
	}
	//store call data in fake component key "call"
	call.data["call"] = call.CallData()
	cm.calls[call.ChannelID] = call
	return call, nil
}
