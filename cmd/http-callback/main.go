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

	http.HandleFunc("/winner", handleWinner)

	log.Println("HTTP server up and running on port 9000 and serving file [assets/configs/callflow_static.json]")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var counter int

func handleWinner(w http.ResponseWriter, r *http.Request) {
	counter++
	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %s", err)
	} else {
		log.Println("Provided values on request:")
		for key := range r.PostForm {
			log.Printf("\"%s\":\"%s\"", key, r.PostFormValue(key))
		}
	}

	res := "winner"

	//pot 2 always looses, pot 1 each other wins
	if r.PostFormValue("pot") == "2" || counter%2 == 0 {
		res = "looser"
	}

	http.ServeFile(w, r, "assets/response/"+res+".json")
}
