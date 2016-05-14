package server

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

type pages struct {
	templateDir string
	db          *sqlx.DB
}

func (p *pages) filename(r *http.Request) string {
	return filepath.Join(p.templateDir, r.URL.Path) + ".html"
}

func (p *pages) serveFile(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile(p.filename(r))
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func (p *pages) login(w http.ResponseWriter, r *http.Request) {
	sid, ok := getSessionID(r)
	if ok && userExists(p.db, "sid", sid) {
		http.Redirect(w, r, "/order", http.StatusSeeOther)
		return
	}
	p.serveFile(w, r)
}

func (p *pages) createAccount(w http.ResponseWriter, r *http.Request) {
	p.serveFile(w, r)
}

type Snack struct {
	Name  string
	Image string
}

type OrderPage struct {
	Snacks []Snack
}

func (p *pages) order(w http.ResponseWriter, r *http.Request) {
	page := OrderPage{
		Snacks: []Snack{{Name: "Snickers", Image: "/static/img/snickers.jpg"}, {Name: "KitKat"}, {Name: "M&Ms"}},
	}
	t, err := template.ParseFiles(p.filename(r))
	if err != nil {
		panic(err)
	}
	if err := t.Execute(w, page); err != nil {
		panic(err)
	}
}
