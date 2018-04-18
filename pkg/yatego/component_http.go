package yatego

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// HTTP component plays list of songs
type HTTP struct {
	Base
	HTTPClient *http.Client
}

// NewHTTPComponent generates new HTTP component
func NewHTTPComponent(base Base, httpClient *http.Client) *HTTP {
	h := &HTTP{
		Base:       base,
		HTTPClient: httpClient,
	}
	h.Init()
	return h
}

// Init pseudo constructor
func (h *HTTP) Init() {
	h.logger.Debugf("HTTP [%s] init", h.Name())

	//on enter play song
	h.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		h.logger.Debugf("HTTP [%s] on enter", h.Name())
		return h.Fetch(call, msg)
	})
}

// Fetch makes http call and fetches and builds new components
func (h *HTTP) Fetch(call *Call, msg *Message) *CallbackResult {
	t, ok := h.ConfigAsString("transfer")
	if !ok {
		h.logger.Errorf("No transfer defined in comp [%s]", h.Name())
		return NewCallbackResult(ResStop, "")
	}
	result, err := h.fetchJSONResponse(call, msg)
	if err != nil {
		h.logger.Errorf("Error fetching in comp [%s]: %s", h.Name(), err)
		return NewCallbackResult(ResStop, "")
	}
	h.storeJSONVals(call, result)
	//transfer to the first of new components
	return NewCallbackResult(ResTransfer, t)
}

func (h *HTTP) fetchJSONResponse(call *Call, msg *Message) (map[string]interface{}, error) {
	url, err := h.url()
	if err != nil {
		return nil, err
	}
	reqParams := h.requestParams(call, msg)
	resp, err := h.HTTPClient.PostForm(url, reqParams)
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

	h.logger.Debugf("Component [%s] got raw json response [%s] at url [%s] and request %+v", h.Name(), string(body), url, reqParams)

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (h *HTTP) url() (string, error) {
	url, exists := h.ConfigAsString("url")
	if !exists || url == "" {
		return "", fmt.Errorf("URL not defined for HTTP component [%s]", h.Name())
	}
	return url, nil
}

func (h *HTTP) storeJSONVals(call *Call, res map[string]interface{}) {
	for key, val := range res {
		h.SetCallData(call, key, val)
	}
}

func (h *HTTP) requestParams(call *Call, msg *Message) url.Values {
	reqData := map[string][]string{}
	reqData["called"] = []string{call.Called}
	reqData["caller"] = []string{call.Caller}
	reqData["ID"] = []string{call.BillingID}

	sFields, ok := h.ConfigAsString("request_fields")
	if !ok {
		return reqData
	}
	sNamespaces, ok := h.ConfigAsString("request_namespaces")
	if !ok {
		return reqData
	}
	sval := ""
	fields := strings.Split(sFields, ",")
	namespaces := strings.Split(sNamespaces, ",")
	for i, f := range fields {
		if i >= len(namespaces) {
			break
		}
		n := strings.Split(namespaces[i], ".")
		if len(n) != 2 {
			continue
		}
		val, ok := call.Data(n[0], n[1])
		if ok {
			sval = val.(string)
		} else {
			sval = ""
		}
		reqData[f] = []string{sval}
	}
	return reqData
}
