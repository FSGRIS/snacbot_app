package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type user struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	OID      string `db:"oid"`
	SID      int64  `db:"sid"`
}

func userExists(db *sqlx.DB, field string, v interface{}) bool {
	q := fmt.Sprintf("select exists (select 1 from users where %s=?) limit 1", field)
	var n int
	err := db.QueryRowx(q, v).Scan(&n)
	if err != nil {
		panic(err)
	}
	return n == 1
}

func decode(w http.ResponseWriter, r *http.Request, b interface{}) bool {
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&b); err != nil {
		log.Println("[decode]", err)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, res interface{}) {
	b, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		panic(err)
	}
}

func grantCookie(w http.ResponseWriter, sid int64) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    strconv.FormatInt(sid, 10),
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HttpOnly: true,
		Path:     "/",
	})
}

func getSessionID(r *http.Request) (int64, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("no cookie")
			return 0, false
		}
		panic(err)
	}
	sid, err := strconv.ParseInt(c.Value, 10, 64)
	if err != nil {
		log.Println(err)
		return 0, false
	}
	return sid, true
}

func getUser(db *sqlx.DB, r *http.Request) *user {
	sid, ok := getSessionID(r)
	if !ok {
		return nil
	}
	var u user
	err := db.Get(&u, "select * from users where sid=?", sid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}
	return &u
}

func badRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	if _, err := w.Write([]byte(msg)); err != nil {
		panic(err)
	}
}
