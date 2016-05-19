package server

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type apiServer struct {
	db  *sqlx.DB
	ros *rosServer
}

func newApiServer(db *sqlx.DB, ros *rosServer) *apiServer {
	s := &apiServer{db, ros}
	s.ros.advertise("snacbot/orders", "snacbot/Order")
	return s
}

func (s *apiServer) login(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !decode(w, r, &b) {
		return
	}
	var uid int64
	var sid sql.NullInt64
	err := s.db.QueryRowx(
		"select id, sid from users where email=? and password=?",
		b.Email, b.Password).
		Scan(&uid, &sid)
	if err != nil {
		if err == sql.ErrNoRows {
			badRequest(w, "Invalid email / password combo.")
			return
		}
		panic(err)
	}
	// Valid user!
	if !sid.Valid {
		// No current login session.
		sid.Int64 = rand.Int63()
		s.db.MustExec("update users set sid=? where id=?", sid.Int64, uid)
	}
	grantCookie(w, sid.Int64)
	w.WriteHeader(http.StatusOK)
}

func (s *apiServer) logout(w http.ResponseWriter, r *http.Request) {
	sid, ok := getSessionID(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	result := s.db.MustExec("update users set sid=NULL where sid=?", sid)
	if n, err := result.RowsAffected(); n == 0 {
		// No rows affected -- invalid session id.
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *apiServer) createAccount(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		OrgName  string `json:"orgName"`
		OrgCode  string `json:"orgCode"`
	}
	if !decode(w, r, &b) {
		return
	}
	if userExists(s.db, "email", b.Email) {
		badRequest(w, "Email already registered.")
		return
	}
	// TODO: Check company table for code.
	var oid int64
	err := s.db.QueryRowx(
		"select id from orgs where name=? and code=?", b.OrgName, b.OrgCode).
		Scan(&oid)
	if err != nil {
		if err == sql.ErrNoRows {
			badRequest(w, "Invalid organization information.")
			return
		}
		panic(err)
	}
	sid := rand.Int63()
	// I know we shouldn't store plaintext passwords, but fuck it.
	s.db.MustExec(`
		insert into users (email, password, oid, sid)
		values (?, ?, ?, ?)`,
		b.Email, b.Password, oid, sid)
	grantCookie(w, sid)
	w.WriteHeader(http.StatusOK)
}

func (s *apiServer) order(w http.ResponseWriter, r *http.Request) {
	var b struct {
		LocationID   int64 `json:"locationID"`
		SaveLocation bool  `json:"saveLocation"`
		Snacks       []struct {
			ID       int64 `json:"id"`
			Quantity int   `json:"quantity"`
		} `json:"snacks"`
	}
	if !decode(w, r, &b) {
		return
	}
	u := getUser(s.db, r)
	if u == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO: Save location if specified.
	log.Println("placing order?")
	s.ros.publish("snacbot/orders", dict{
		"location_id": b.LocationID,
		"snacks":      b.Snacks,
	})
}
