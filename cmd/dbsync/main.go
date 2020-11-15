package main

import (
	"flag"
	"log"

	"example.org/york/dbsync/lib"
)

var (
	patchDir         = flag.String("path_dir", ".", "path to the directory that contains patches")
	dbConnConfigFile = flag.String("db_conn", "db_connconfig.yaml", "path to the yaml file that contains the info to connect to the db")
	logFile          = flag.String("log_file", "", "path to the log file output")
)

func main() {
	flag.Parse()
	logger := lib.GetMainLogger(*logFile)
	defer lib.CloseLogStream()
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
	logger.Println("Bye")
}
