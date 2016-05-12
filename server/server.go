package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	app := negroni.Classic()
	r := mux.NewRouter()
	r.HandleFunc("/", http.FileServer(http.Dir("static")))
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/create_account", CreateAccountHandler).Methods("POST")
	r.HandleFunc("/api/order", OrderHandler).Methods("POST")
	app.UseHandler(r)
	addr := ":3000"
	log.Printf("Listening on %s ...\n", addr)
	panic(http.ListenAndServe(addr, app))
}
