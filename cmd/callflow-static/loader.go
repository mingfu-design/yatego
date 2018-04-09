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
