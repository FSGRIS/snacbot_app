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

func (p *pages) filename(f string) string {
	return filepath.Join(p.templateDir, f) + ".html"
}

func (p *pages) serveFile(w http.ResponseWriter, f string) {
	b, err := ioutil.ReadFile(p.filename(f))
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
	p.serveFile(w, "login")
}

func (p *pages) createAccount(w http.ResponseWriter, r *http.Request) {
	p.serveFile(w, "create_account")
}

type Snack struct {
	ID    int64
	Name  string
	Image string
}

type OrderPage struct {
	Snacks []Snack
}

func (p *pages) order(w http.ResponseWriter, r *http.Request) {
	u := getUser(p.db, r)
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
	t, err := template.ParseFiles(p.filename("order"))
	if err != nil {
		panic(err)
	}
	if err := t.Execute(w, page); err != nil {
		panic(err)
	}
}
