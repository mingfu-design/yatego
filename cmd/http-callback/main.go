package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {
	portPtr := flag.Int("port", 9000, "port number")
	flag.Parse()
	port := strconv.Itoa(*portPtr)

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

	log.Printf("HTTP server up and running on port %s", port)

	err := http.ListenAndServe(":"+port, nil)
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
