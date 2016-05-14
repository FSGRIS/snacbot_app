package main

import (
	"flag"
	"log"
	"os"
	"snacbot_app/server"
)

func isDir(d string) bool {
	s, err := os.Stat(d)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	return s.IsDir()
}

func main() {
	dbFile := flag.String("db", "", "sqlite3 database file")
	staticDir := flag.String("sd", "./static/", "static assets directory")
	templateDir := flag.String("td", "./templates/", "template directory")
	flag.Parse()
	if *dbFile == "" {
		log.Fatalln("error: no database file specified")
	}
	if !isDir(*staticDir) {
		log.Fatalf("error: cannot find static dir (%s)\n", *staticDir)
	}
	if !isDir(*templateDir) {
		log.Fatalf("error: cannot find template dir (%s)\n", *templateDir)
	}
	server.Run(*dbFile, *staticDir, *templateDir)
}
