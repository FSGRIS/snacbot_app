package server

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// One week in seconds.
const cookieMaxAge = 60 * 60 * 24 * 7

const cookieName = "sid"

func Run(dbFile, staticDir, templateDir string) {
	db, err := sqlx.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	tryCreateTables(db)
	a := &api{db}
	p := &pages{templateDir: templateDir, db: db}
	n := negroni.Classic()
	r := mux.NewRouter()
	r.HandleFunc("/login", p.login).Methods("GET")
	r.HandleFunc("/create_account", p.createAccount).Methods("GET")
	r.HandleFunc("/order", p.order).Methods("GET")
	r.HandleFunc("/api/login", a.login).Methods("POST")
	r.HandleFunc("/api/logout", a.logout).Methods("POST")
	r.HandleFunc("/api/create_account", a.createAccount).Methods("POST")
	r.HandleFunc("/api/order", a.order).Methods("POST")
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	n.UseHandler(r)
	addr := ":3000"
	log.Printf("Listening on %s ...\n", addr)
	panic(http.ListenAndServe(addr, n))
}

func tryCreateTables(db *sqlx.DB) {
	qs := []string{
		`create table if not exists orgs (
			id integer primary key,
			name string,
			code string
		)`,
		`create table if not exists users (
			id integer primary key,
			username text,
			email text,
			password text,
			oid integer,
			sid integer
		)`,
		`create table if not exists snacks (
			id integer primary key,
			name text,
			quantity integer,
			image text
		)`,
	}
	for _, q := range qs {
		if _, err := db.Exec(q); err != nil {
			panic(err)
		}
	}
}