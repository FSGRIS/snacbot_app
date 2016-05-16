package server

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

type pageServer struct {
	templateDir string
	db          *sqlx.DB
}

func (s *pageServer) filename(f string) string {
	return filepath.Join(s.templateDir, f) + ".html"
}

func (s *pageServer) serveFile(w http.ResponseWriter, f string) {
	b, err := ioutil.ReadFile(s.filename(f))
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func (s *pageServer) login(w http.ResponseWriter, r *http.Request) {
	sid, ok := getSessionID(r)
	if ok && userExists(s.db, "sid", sid) {
		http.Redirect(w, r, "/order", http.StatusSeeOther)
		return
	}
	s.serveFile(w, "login")
}

func (s *pageServer) createAccount(w http.ResponseWriter, r *http.Request) {
	s.serveFile(w, "create_account")
}

type Snack struct {
	ID    int64
	Name  string
	Image string
}

type OrderPage struct {
	Snacks []Snack
}

func (s *pageServer) order(w http.ResponseWriter, r *http.Request) {
	u := getUser(s.db, r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	page := OrderPage{
		Snacks: []Snack{
			{ID: 1, Name: "Snickers", Image: "snickers.jpg"},
			{ID: 2, Name: "KitKat", Image: "kitkat.jpg"},
			{ID: 3, Name: "Tim's Potato Chips", Image: "tims.jpg"},
		},
	}
	t, err := template.ParseFiles(s.filename("order"))
	if err != nil {
		panic(err)
	}
	if err := t.Execute(w, page); err != nil {
		panic(err)
	}
}
