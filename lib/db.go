package lib

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

var dbConn *pgx.Conn

const patch_status_table_name = "dbsync_patch_status"

type DbConnConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (dbConn DbConnConfig) toDatabaseURL() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v", dbConn.User, dbConn.Password, dbConn.Host, dbConn.Port, dbConn.Database)
}

func (dbConnCfg DbConnConfig) getInstance() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), dbConnCfg.toDatabaseURL())
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v\n", err)
		panic("getInstance")
	}
	return conn
}

func InitDbConn(dbConnCfg DbConnConfig) error {
	dbConn = dbConnCfg.getInstance()
	// dbsync_patch_status
	initSql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v ( id text PRIMARY KEY NOT NULL, applied boolean NOT NULL);", patch_status_table_name)
	dbConn.Exec(context.Background(), initSql)

	return nil
}

func CloseDbConn() {
	if !dbConn.IsClosed() {
		dbConn.Close(context.Background())
	}
}

// for the current scope, a prereq can only be in either state:
// the table AND applied=true
// not in the table
func CheckPrereqStatus(prereq_id string) bool {
	var prereq_status bool
	err := dbConn.QueryRow(context.Background(), "SELECT applied from $1 WHERE id = $2 ;", patch_status_table_name, prereq_id).Scan(&prereq_status)
	if err != nil {
		logger.Fatalf("prerequisite patch %v not found in patch status table! aborted.\n", prereq_id)
		panic("missing prereq")
	}
	logger.Printf("%v\n", prereq_status)
	return prereq_status
}

func ApplyPatchTx(table string, sql string) {
	logger.Printf("SQL transactiono begins\n")
	tx, err := dbConn.Begin(context.Background())
	if err != nil {
		logger.Fatalln(err)
	}

	// if the tx commits successfully, this is a no-op
	// uncomment this line if we want to use alternative way to roll back
	// defer tx.Rollback(context.Background())
	logger.Println(sql)
	tx.Exec(context.Background(), sql)
	err = tx.Commit(context.Background())
	if err != nil {
		logger.Fatalln(err)
	}
}
