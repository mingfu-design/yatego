package yatego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Fetcher component plays list of songs
type Fetcher struct {
	Base
	cfLoader *CallflowLoaderJSON
}

// NewFetcherComponent generates new Fetcher component
func NewFetcherComponent(base Base, cfLoader *CallflowLoaderJSON) *Fetcher {
	f := &Fetcher{
		Base:     base,
		cfLoader: cfLoader,
	}
	f.Init()
	return f
}

// Init pseudo constructor
func (f *Fetcher) Init() {
	f.logger.Debugf("Fetcher [%s] init", f.Name())

	//on enter play song
	f.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		f.logger.Debugf("Fetcher [%s] on enter", f.Name())
		return f.Fetch(call, msg)
	})
}

// Fetch makes http call and fetches and builds new components
func (f *Fetcher) Fetch(call *Call, msg *Message) *CallbackResult {
	jsonStr, err := f.fetchJSONResponse(call, msg)
	if err != nil {
		f.logger.Errorf("Error fetching in comp [%s]: %s", f.Name(), err)
		return NewCallbackResult(ResStop, "")
	}
	f.cfLoader.SetJSON(jsonStr)
	cf, err := f.cfLoader.Load(map[string]string{})
	if err != nil {
		f.logger.Errorf("Error loading fetched CF in comp [%s]: %s", f.Name(), err)
		return NewCallbackResult(ResStop, "")
	}
	f.logger.Debugf("Fetcher [%s] loaded %d new components", f.Name(), len(cf.Components))
	if len(cf.Components) == 0 {
		f.logger.Error("No new components loaded, nothing else to do in comp [%s]: %s", f.Name())
		return NewCallbackResult(ResStop, "")
	}
	//build components
	for _, com := range cf.Components {
		f.logger.Debugf("Building component: %+v", com)
		call.AddComponent(com.Factory(com.ClassName, com.Name, com.Config))
	}
	f.InstallMessageHandlers(call)
	f.InstallMessageWatches(call)
	//transfer to the first of new components
	return NewCallbackResult(ResTransfer, cf.Components[0].Name)
}

func (f *Fetcher) fetchJSONResponse(call *Call, msg *Message) (string, error) {
	url, err := f.url()
	if err != nil {
		return "", err
	}
	reqParams := f.requestParams(call, msg)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqParams)
	resp, err := f.httpClient().PostForm(url, reqParams)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Reponse status not OK/200 but [%d]", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (f *Fetcher) httpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func (f *Fetcher) url() (string, error) {
	url, exists := f.ConfigAsString("url")
	if !exists || url == "" {
		return "", fmt.Errorf("URL not defined for fetcher component [%s]", f.Name())
	}
	return url, nil
}

func (f *Fetcher) requestParams(call *Call, msg *Message) url.Values {
	reqData := map[string][]string{}
	reqData["called"] = []string{call.Called}
	reqData["caller"] = []string{call.Caller}
	reqData["ID"] = []string{call.BillingID}
	for compName, vals := range call.data {
		for key, val := range vals {
			reqData["data."+compName+"."+key] = []string{
				fmt.Sprintf("%v", val),
			}
		}
	}
	return reqData
}
