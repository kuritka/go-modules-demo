package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gomodules/pack"
	"log"
	"net/http"
	"rsc.io/quote"
)

func main() {

	fmt.Println(quote.Hello())
	//fmt.Println(quotev2.Hello())
	//fmt.Println(quote2.Hello())


	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte( "Hello from ABC again!"))
	//})
	//
	//err := http.ListenAndServe(":3001", nil)
	//if err != nil {
	//	log.Fatal(err)
	//}

	rtr := mux.NewRouter()

	rtr.HandleFunc("/{topic}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Write([]byte( "Topic: " + vars["topic"]))
	})

	pack.Abc()

	http.Handle("/", rtr)

	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
