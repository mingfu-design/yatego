# YateGo
YateGo packes enables communication with YATE core via external protocol http://docs.yate.ro/wiki/Programmer%27s_guide

## Installation

`go get github.com/rukavina/yatego`

## Simple IVR

this simple IVR just answers and plays `tone/congestion`

In order to install it as regalar callflow/ivr in yate, compile, place binary in `share/yate/scripts`.
Then you need to have a route to it, eg:

```sh
^920$=external/nodata/yourbinary
```

```golang
package main

import "github.com/rukavina/yatego/pkg/yatego"

func main() {
	f := yatego.NewFactory()
	c := f.Controller(nil)
	c.Logger().Debug("Starting yatego IVR [inline]")

	com := f.BaseComponent()
	com.OnEnter(func(call *yatego.Call, message *yatego.Message) *yatego.CallbackResult {
		com.(*yatego.Base).PlayTone("congestion", call, map[string]string{})
		return yatego.NewCallbackResult(yatego.ResStay, "")
	})

	c.AddStaticComponent(com)
	c.Run("")
}
```

## Callflow IVR (static)

The callflow definition allows more cleaner and structural setup of your IVRs. IVRs based on callflows are defined as the system of components.
There are different ways to load a callflow, simplest one is by using `CallflowLoaderStatic`

```golang
package main

import "github.com/rukavina/yatego/pkg/yatego"

func main() {
	f := yatego.NewFactory()
	c := f.Controller(nil)
	c.Logger().Debug("Starting yatego IVR [inline]")

	com := f.BaseComponent()
	com.OnEnter(func(call *yatego.Call, message *yatego.Message) *yatego.CallbackResult {
		com.(*yatego.Base).PlayTone("congestion", call, map[string]string{})
		return yatego.NewCallbackResult(yatego.ResStay, "")
	})

	c.AddStaticComponent(com)
	c.Run("")
}
```

loader.go:

```golang
package main

import (
	"os"
	"path/filepath"

	"github.com/rukavina/dicgo"
	"github.com/rukavina/yatego/pkg/yatego"
)

func loader(c dicgo.Container) yatego.CallflowLoader {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(ex)

	return yatego.NewCallflowLoaderStatic(&yatego.Callflow{
		Components: []*yatego.CallflowComponent{
			//player start
			&yatego.CallflowComponent{
				Name:      "start",
				ClassName: "player",
				Config: map[string]interface{}{
					"playlist": dir + "/assets/audio/welcome.sln",
					"transfer": "menu",
				},
				Factory: yatego.PlayerComponentFactory(c),
			},
			//menu
			&yatego.CallflowComponent{
				Name:      "menu",
				ClassName: "menu",
				Config: map[string]interface{}{
					"keys":     "1,2,3",
					"transfer": "playlist1,playlist2,recorder",
				},
				Factory: yatego.MenuComponentFactory(c),
			},
			//playlist1
			&yatego.CallflowComponent{
				Name:      "playlist1",
				ClassName: "player",
				Config: map[string]interface{}{
					"playlist": dir + "/assets/audio/clicked_1.sln",
					"transfer": "goodbye",
				},
				Factory: yatego.PlayerComponentFactory(c),
			},
			//playlist2
			&yatego.CallflowComponent{
				Name:      "playlist2",
				ClassName: "player",
				Config: map[string]interface{}{
					"playlist": dir + "/assets/audio/clicked_2.sln",
					"transfer": "goodbye",
				},
				Factory: yatego.PlayerComponentFactory(c),
			},
			//recorder
			&yatego.CallflowComponent{
				Name:      "recorder",
				ClassName: "recorder",
				Config: map[string]interface{}{
					"file":     dir + "/assets/voicemail/{called}_{caller}_{billingId}.sln",
					"maxlen":   80000,
					"transfer": "goodbye",
				},
				Factory: yatego.RecorderComponentFactory(c),
			},
			//player goodbye
			&yatego.CallflowComponent{
				Name:      "goodbye",
				ClassName: "player",
				Config: map[string]interface{}{
					"playlist": dir + "/assets/audio/goodbye.sln",
				},
				Factory: yatego.PlayerComponentFactory(c),
			},
		},
	})
}
```

## Callflow IVR (json)

More callflow flexibility you gain by using `CallflowLoaderJSON` which allows to define you callflow in an external json file

```golang
package main

import (
	"os"
	"path/filepath"

	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	f := yatego.NewFactory()

	//json loader
	l := f.CallflowLoaderJSON()
	//load json content from external file
	exec, _ := os.Executable()
	dir := filepath.Dir(exec)
	l.SetJSONFile(dir + "/assets/configs/callflow_static.json")

	//controller
	c := f.Controller(l)
	c.Logger().Debug("Starting yatego IVR [callflow-json]")
	c.Run("")
}
```
## Callflow IVR (vars)

The next level of flexibility is achieved by using json file as callflow template, which means that it can contain variables instead of hard-coded values for component configuration. At the runtime, it's possible to obtain values to be used to parse template variables

```golang
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	//get current dir
	exec, _ := os.Executable()
	dir := filepath.Dir(exec)

	f := yatego.NewFactory()
	//json loader
	l := f.CallflowLoaderJSON()
	//custom onload, pull vars from json
	l.OnLoad = func(loader *yatego.CallflowLoaderJSON, cf *yatego.Callflow, params map[string]string) error {
		//load vars from json
		data, err := ioutil.ReadFile(dir + "/assets/configs/callflow_vars.json")
		if err != nil {
			return err
		}
		var vars map[string]string
		if err := json.Unmarshal(data, &vars); err != nil {
			return err
		}
		loader.SetVars(vars)
		return nil
	}
	//load json CF template from external file
	l.SetJSONFile(dir + "/assets/configs/callflow_tpl.json")

	//controller
	c := f.Controller(l)
	c.Logger().Debug("Starting yatego IVR [callflow-vars]")
	c.Run("")
}
```

template json file:

```json
{
    "components":[
        {
            "name": "start",
            "class": "player",
            "config": {
                "playlist": "/vagrant/assets/audio/{prompt_file}",
                "transfer": "menu"
            }
        },
        {
            "name": "menu",
            "class": "menu",
            "config": {
                "keys": "1,2,3",
                "transfer": "playlist1,playlist2,recorder"
            } 
        },
        {
            "name": "playlist1",
            "class": "player",
            "config": {
                "playlist": "/vagrant/assets/audio/clicked_1.sln",
                "transfer": "goodbye"
            }
        }, 
        {
            "name": "playlist2",
            "class": "player",
            "config": {
                "playlist": "/vagrant/assets/audio/clicked_2.sln",
                "transfer": "goodbye"
            }
        },
        {
            "name": "recorder",
            "class": "recorder",
            "config": {
                "file": "/vagrant/assets/voicemail/{called}_{caller}_{billingId}.sln",
                "maxlen": "{rec_maxlen}",
                "transfer": "goodbye"
            }
        },
        {
            "name": "goodbye",
            "class": "player",
            "config": {
                "playlist": "/vagrant/assets/audio/goodbye.sln"
            }
        }                                       
    ]
}
```

values :

```json
{
    "prompt_file":  "welcome.sln",
    "rec_maxlen":   "160000" 
}
```

## Callflow IVR (dynamic)

The most flexible callflow you can get by using component `fetcher`. It has config `url`. When fetcher enters execution it will make http POST request to defined url and expects json result. The json response should be new callflow. New components are generated and appended to the existing call components. The execution is trasfered to the very next component among new ones.

Your http server/handler receives the following as a form's post data:

```
ID=1522323299-11
called=924
caller=41587000201
```

Thus you can generate json callflow to return dynamically based on provided params.

```golang
package main

import (
	"os"
	"path/filepath"

	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	f := yatego.NewFactory()

	//json loader
	l := f.CallflowLoaderJSON()
	//load json content from external file
	exec, _ := os.Executable()
	dir := filepath.Dir(exec)
	l.SetJSONFile(dir + "/assets/configs/callflow_dynamic.json")

	//controller
	c := f.Controller(l)
	c.Logger().Debug("Starting yatego IVR [callflow-dynamic]")
	c.Run("")
}
```

demo url handler:

```golang
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("error parsing form: %s", err)
		} else {
			log.Println("Provided values on request:")
			for key := range r.PostForm {
				log.Printf("\"%s\":\"%s\"", key, r.PostFormValue(key))
			}
		}
		http.ServeFile(w, r, "assets/configs/callflow_static.json")
	})

	log.Println("HTTP server up and running on port 9000 and serving file [assets/configs/callflow_static.json]")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```

## Running examples in vagrant

* First build all `cmd` executables with `go build` in each subfolder
* From project's root folder run `vagrant up`
* connect sip softphone (eg. zoiper) to sip account `41587000201@172.28.128.3`, pass: `milan`
* call example destinations: `920`, `921`, `922`, `923`, `924`, `925`

NOTE: due to vagrant issues sometimes mock http srv is not started in `vagrant up`. To check and start do:

```bash
vagrant ssh
sudo -i
ps aux|grep http
# if not started
cd /vagrant/cmd/http-callback
./http-callback &
```

## Yate management

SSH to vagrant:

```sh
vagrant ssh
sudo -i
```

Log:

```sh
tail -f /var/log/yate/messages
```

Stop:

```sh
pkill yate
```

Start:

```sh
/opt/yate/startyate.sh
```