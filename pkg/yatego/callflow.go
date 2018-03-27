package yatego

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// CallflowComponent is the definition of a single callflow component
type CallflowComponent struct {
	Name      string                 `json:"name"`
	ClassName string                 `json:"class"`
	Config    map[string]interface{} `json:"config"`
	Factory   ComponentFactory
}

// Callflow is the definition of a IVR callflow and components def. repos
type Callflow struct {
	Components []*CallflowComponent `json:"components"`
}

// CallflowLoader interface which defines object to be able to load new callflow
type CallflowLoader interface {
	Load(params map[string]string) (*Callflow, error)
}

// CallflowLoaderStatic is simplest CallflowLoader implementation
type CallflowLoaderStatic struct {
	callflow *Callflow
}

// Load callflow
func (cl *CallflowLoaderStatic) Load(params map[string]string) (*Callflow, error) {
	return cl.callflow, nil
}

// NewCallflowLoaderStatic generates new CallflowLoaderStatic
func NewCallflowLoaderStatic(c *Callflow) *CallflowLoaderStatic {
	return &CallflowLoaderStatic{
		callflow: c,
	}
}

// CallflowLoaderJSON is CF loader from json string
type CallflowLoaderJSON struct {
	data       []byte
	factoryMap map[string]ComponentFactory
}

// NewCallflowLoaderJSON generates new CallflowLoaderJSON
func NewCallflowLoaderJSON(strJSON string, factoryMap map[string]ComponentFactory) *CallflowLoaderJSON {
	return &CallflowLoaderJSON{
		data:       []byte(strJSON),
		factoryMap: factoryMap,
	}
}

// Load callflow
func (cl *CallflowLoaderJSON) Load(params map[string]string) (*Callflow, error) {
	cf := &Callflow{}
	err := json.Unmarshal(cl.data, cf)
	if err != nil {
		return nil, err
	}
	err = CallflowLoaderPopulateFactories(cf, cl.factoryMap)
	if err != nil {
		return cf, err
	}
	return cf, nil
}

// SetJSON sets json string for loader
func (cl *CallflowLoaderJSON) SetJSON(strJSON string) {
	cl.data = []byte(strJSON)
}

// SetJSONFile loads json from file
func (cl *CallflowLoaderJSON) SetJSONFile(fileJSON string) error {
	data, err := ioutil.ReadFile(fileJSON)
	if err != nil {
		return err
	}
	cl.data = data
	return nil
}

// CallflowLoaderPopulateFactories based on a component's className, sets its factory method
func CallflowLoaderPopulateFactories(cf *Callflow, factoryMap map[string]ComponentFactory) error {
	for _, com := range cf.Components {
		fac, exists := factoryMap[com.ClassName]
		if !exists {
			return fmt.Errorf("Factory for class [%s] not found", com.ClassName)
		}
		com.Factory = fac
	}
	return nil
}
