package pack

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Abc(){
	rtr := mux.NewRouter()

	rtr.HandleFunc("/x/{topic}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Write([]byte( "Topic: " + vars["topic"]))
	})
}