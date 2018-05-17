package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rukavina/yatego/pkg/yatego"
)

func main() {
	jsonPtr := flag.String("cf", "", "callflow json file to load")
	outPtr := flag.String("out", "", "output mmd file path")
	flag.Parse()
	jsonFile := *jsonPtr
	mmdFile := *outPtr
	if jsonFile == "" {
		log.Fatalln("CF json file not defined")
	}
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		log.Fatalf("File [%s] does not exist", jsonFile)
	}

	loader := yatego.NewCallflowLoaderJSON("", map[string]yatego.ComponentFactory{})
	loader.SetJSONFile(jsonFile)
	cf, err := loader.Load(map[string]string{})
	if cf == nil && err != nil {
		log.Fatalf("Error loading json callflow: %s", err)
	}

	g := NewGraph("TD")
	mmd := g.Render(cf)

	//print to std out
	if mmdFile == "" {
		fmt.Printf("\n%s\n", mmd)
		return
	}

	//write to file
	err = ioutil.WriteFile(mmdFile, []byte(mmd), 0644)
	if err != nil {
		log.Fatalf("Error saving mmd file %s: %s", mmdFile, err)
	}
	fmt.Printf("\nCallflow json file [%s] generated mmd output in [%s]\n", jsonFile, mmdFile)

}
