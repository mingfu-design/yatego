package yatego

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
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
	CallflowVarsParser
	data       []byte
	factoryMap map[string]ComponentFactory
	OnLoad     CallflowLoaderJSONOnLoad
}

//CallflowLoaderJSONOnLoad is hook called from CallflowLoaderJSON for custom loading
type CallflowLoaderJSONOnLoad func(loader *CallflowLoaderJSON, cf *Callflow, params map[string]string) error

// NewCallflowLoaderJSON generates new CallflowLoaderJSON
func NewCallflowLoaderJSON(strJSON string, factoryMap map[string]ComponentFactory) *CallflowLoaderJSON {
	return &CallflowLoaderJSON{
		CallflowVarsParser: CallflowVarsParser{
			vars: map[string]string{},
		},
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
	//hook
	if cl.OnLoad != nil {
		err := cl.OnLoad(cl, cf, params)
		if err != nil {
			return cf, err
		}
	}
	cl.parseCallflow(cf)
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

// CallflowVarsParser is CF loader mixing to dynamically parse cf template by vars
type CallflowVarsParser struct {
	vars map[string]string
}

// SetVars defines vars map to be used to parse CF template
func (cp *CallflowVarsParser) SetVars(vars map[string]string) {
	cp.vars = vars
}

// parseCallflow based on internal vars map
func (cp *CallflowVarsParser) parseCallflow(cf *Callflow) {
	if len(cp.vars) == 0 {
		return
	}
	vals := []string{}
	for key, val := range cp.vars {
		vals = append(vals, "{"+key+"}", val)
	}
	r := strings.NewReplacer(vals...)
	for _, com := range cf.Components {
		for key, val := range com.Config {
			switch val.(type) {
			default:
				continue
			case string:
				com.Config[key] = r.Replace(val.(string))
			}
		}
	}
}
