package server

import (
	"database/sql"
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
	Name  string
	Image string
}

type OrderPage struct {
	Snacks []Snack
}

func (p *pages) order(w http.ResponseWriter, r *http.Request) {
	sid, ok := getSessionID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var u user
	err := p.db.Get(&u, "select * from users where sid=?", sid)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		panic(err)
	}
	page := OrderPage{
		Snacks: []Snack{
			{Name: "Snickers", Image: "snickers.jpg"},
			{Name: "KitKat", Image: "kitkat.jpg"},
			{Name: "Tim's Potato Chips", Image: "tims.jpg"},
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
