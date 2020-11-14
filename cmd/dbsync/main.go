package main

import (
	"flag"
	"log"

	"example.org/york/dbsync/lib"
)

var (
	patchDir         = flag.String("path_dir", ".", "path to the directory that contains patches")
	dbConnConfigFile = flag.String("db_conn", "psql_conn.json", "path to the json file that contains the info to connect to the db")

	logger = lib.GetLogger()
)

func main() {
	flag.Parse()
	var dbConnCfg lib.DbConnConfig
	err := lib.ParseFromYamlFile(*dbConnConfigFile, &dbConnCfg)
	if err != nil {
		logger.Fatalln("failed to parse DB connection config file!")
	}
	lib.InitDbConn(dbConnCfg)
	defer lib.CloseDbConn()
	patches, err := lib.ParsePatches(*patchDir)
	lib.ProcessPatches(patches)
	if err != nil {
		log.Fatalln("failed to parse Patches!")
	}
	logger.Println(patches)
}
