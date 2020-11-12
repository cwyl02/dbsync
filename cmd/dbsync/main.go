package main

import (
	"flag"
	"log"

	"example.org/york/dbsync/lib"
	// "dbsync/lib"
)

var (
	patchDir     = flag.String("path_dir", ".", "path to the directory that contains patches")
	dbConnConfig = flag.String("db_conn", "psql_conn.json", "path to the json file that contains the info to connect to the db")
)

// TODO: set up logging
func main() {
	flag.Parse()
	var dbConn lib.DbConnection
	// parse db connection
	err := lib.ParseFromJsonFile(*dbConnConfig, dbConn)
	if err != nil {
		log.Fatalln("failed to parse DB connection config file!")
	}
	// test db connection
	// err := dbConn.test()
	p, err := lib.ParsePatches(*patchDir)
	if err != nil {
		log.Fatalln("failed to parse Patches!")
	}
	log.Println(p)
	// create patch status table if not exist - exit if no create table access   ## table: status
	// parse patch - inactives will be skipped
	// build the graph - for unapplied patch          ## table: status
	// walk the graph - for each node prompt user     ## table: status, @patch.table
}
