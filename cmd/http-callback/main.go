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
