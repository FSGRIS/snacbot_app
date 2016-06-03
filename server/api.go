package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type apiServer struct {
	db   *sqlx.DB
	ros  *rosServer
	locs map[int64]point
}

func expect(ok bool, t string) {
	if !ok {
		panic("expected " + t)
	}
}

func newApiServer(db *sqlx.DB, ros *rosServer) *apiServer {
	s := &apiServer{
		db:   db,
		ros:  ros,
		locs: make(map[int64]point),
	}
	s.ros.advertise("snacbot/orders", "snacbot/Order")
	r := s.ros.callService("snacbot/locations", make([]interface{}, 0))
	var resp struct {
		Values struct {
			Locs []struct {
				ID int64   `json:"id"`
				X  float64 `json:"x"`
				Y  float64 `json:"y"`
			} `json:"locs"`
		} `json:"values"`
	}
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		panic(err)
	}
	for _, l := range resp.Values.Locs {
		s.locs[l.ID] = point{X: l.X, Y: l.Y}
	}
	log.Printf("%#v\n", s.locs)
	return s
}

func (s *apiServer) getLocations(w http.ResponseWriter, r *http.Request) {
	// Convert locs from int->point to string->point, to satisfy JSON spec.
	locs := make(map[string]point)
	for lid, p := range s.locs {
		locs[strconv.FormatInt(lid, 10)] = p
	}
	writeJSON(w, locs)
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
	snackIDs := make([]int64, len(b.Snacks))
	for i, s := range b.Snacks {
		snackIDs[i] = s.ID
	}
	// TODO: Save location if specified.
	log.Println("placing order?")
	s.ros.publish("snacbot/orders", dict{
		"location_id": b.LocationID,
		"snack_ids":   snackIDs,
	})
}
