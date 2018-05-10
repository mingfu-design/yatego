package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rukavina/yatego/pkg/yatego"
)

//HTTPCFVarsLoader is http client, getting config from http server
type HTTPCFVarsLoader struct {
	Config     map[string]string
	HTTPClient *http.Client
	Logger     yatego.Logger
}

//LoadCallflowVars loads external json file into map
func (c *HTTPCFVarsLoader) LoadCallflowVars(params map[string]string) (map[string]string, error) {

	url, ok := c.Config["config_api_endpoint"]
	if !ok {
		return nil, fmt.Errorf("Config URL not defined in the config")
	}
	paramsData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling params: %s", err)
	}
	c.Logger.Debugf("CF HTTP vars request: %s", string(paramsData))
	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(paramsData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Reponse status not OK/200 but [%d]", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.Logger.Debugf("CF HTTP vars response: %s", string(body))
	var res map[string]string
	err = json.Unmarshal(body, &res)
	if err != nil {
		c.Logger.Errorf("Error decoding JSON data from CF vars: %s", err)
		return nil, err
	}
	return res, nil
}
