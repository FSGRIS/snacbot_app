package server

import (
	"database/sql"
	"math/rand"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type api struct {
	db *sqlx.DB
}

func (a *api) login(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if !decode(w, r, &b) {
		return
	}
	var uid int64
	var sid sql.NullInt64
	err := a.db.QueryRowx(
		"select id, sid from users where username=? and password=?",
		b.Username, b.Password).
		Scan(&uid, &sid)
	if err != nil {
		if err == sql.ErrNoRows {
			badRequest(w, "invalid username / password")
			return
		}
		panic(err)
	}
	// Valid user!
	if !sid.Valid {
		// No current login session.
		sid.Int64 = rand.Int63()
		a.db.MustExec("update users set sid=? where id=?", sid.Int64, uid)
	}
	grantCookie(w, sid.Int64)
	w.WriteHeader(http.StatusOK)
}

func (a *api) logout(w http.ResponseWriter, r *http.Request) {
	sid, ok := getSessionID(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := a.db.MustExec("update users set sid=NULL where sid=?", sid)
	if n, err := result.RowsAffected(); n == 0 {
		// No rows affected -- invalid session id.
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func (a *api) createAccount(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		OrgName  string `json:"orgName"`
		OrgCode  string `json:"orgCode"`
	}
	if !decode(w, r, &b) {
		return
	}
	if userExists(a.db, "username", b.Username) {
		badRequest(w, "username taken")
		return
	}
	if userExists(a.db, "email", b.Email) {
		badRequest(w, "email taken")
		return
	}
	// TODO: Check company table for code.
	var oid int64
	err := a.db.QueryRowx(
		"select id from orgs where name=? and code=?", b.OrgName, b.OrgCode).
		Scan(&oid)
	if err != nil {
		if err == sql.ErrNoRows {
			badRequest(w, "invalid org")
			return
		}
		panic(err)
	}
	sid := rand.Int63()
	// I know we shouldn't store plaintext passwords, but fuck it.
	a.db.MustExec(`
		insert into users (username, email, password, oid, sid)
		values (?, ?, ?, ?, ?)`,
		b.Username, b.Email, b.Password, oid, sid)
	grantCookie(w, sid)
	w.WriteHeader(http.StatusOK)
}

func (a *api) order(w http.ResponseWriter, r *http.Request) {
	// TODO
}
