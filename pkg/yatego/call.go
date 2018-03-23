package yatego

import "fmt"

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

// SetData sets the component's data key value
func (call *Call) SetData(componentName string, key string, value interface{}) {
	_, exists := call.data[componentName]
	if !exists {
		call.data[componentName] = make(map[string]interface{})
	}
	call.data[componentName][key] = value
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
	cm.calls[call.ChannelID] = call
	return call, nil
}
