package main

import (
	"encoding/json"
	"io/ioutil"
)

//LoadConfiguration loads external json file into map
func LoadConfiguration(path string) (map[string]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config map[string]string
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
