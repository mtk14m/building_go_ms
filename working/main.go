package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("hello world")
		d, err:= io.ReadAll(r.Body)
		if err!=nil{
			http.Error(w, "Ooops", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, " hello %s\n", d)
	})

	http.HandleFunc("/toto", func(w http.ResponseWriter, r *http.Request) {
		log.Println("hello toto")
	})

	http.ListenAndServe(":9090", nil)
}
