package lib

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

var dbConn *pgx.Conn

const (
	patchStatusTableName = "dbsync_patch_status"
)

type DbConnConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (dbConnCfg DbConnConfig) toDatabaseURL() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v", dbConnCfg.User, dbConnCfg.Password, dbConnCfg.Host, dbConnCfg.Port, dbConnCfg.Database)
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
	initSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v ( id text PRIMARY KEY NOT NULL, applied_at timestamptz NOT NULL DEFAULT NOW());", patchStatusTableName)
	dbConn.Exec(context.Background(), initSQL)

	return nil
}

func CloseDbConn() {
	if !dbConn.IsClosed() {
		dbConn.Close(context.Background())
	}
}

func ApplyPatch(table string, sql string) error {
	err := doSQLtx(table, sql)

	return err
}

func doSQLtx(table string, sqlStmts string, sqlArgs ...interface{}) error {
	logger.Printf("SQL transaction begins\n")
	sqlStatements := strings.Split(sqlStmts, ";")
	sqlStatements = Map(sqlStatements, func(strIn string) string {
		return strIn + ";"
	})
	sqlStatements = sqlStatements[:len(sqlStatements)-1]

	tx, err := dbConn.Begin(context.Background())
	if err != nil {
		logger.Fatalln(err)
		return err
	}

	// uncomment the tx.Rollback line if we want to use pgx library to roll back
	// if the tx commits successfully, this is a no-op
	// defer tx.Rollback(context.Background())

	for _, sqlStatement := range sqlStatements {
		logger.Println("SQL statement: ")
		// logger.Printf("%v\n", sqlStatement)
		logger.Printf("%v\n", strings.Trim(sqlStatement, "\n"))
		tx.Exec(context.Background(), sqlStatement, sqlArgs...)
		// if len(sqlArgs) > 0 {

		// } else {
		// 	tx.Exec(context.Background(), sqlStatement)
		// }

	}

	err = tx.Commit(context.Background())
	if err != nil {
		logger.Fatalln(err)
		return err
	}

	return nil
}

// for the current scope, a prereq can only be in either state:
// in the table AND applied_at = now()
// not in the table
func GetPatchStatus(patchID string) error {
	var patchTime time.Time
	sql := fmt.Sprintf("SELECT applied_at from %v WHERE id = $1 ;", patchStatusTableName)
	err := dbConn.QueryRow(context.Background(), sql, patchID).Scan(&patchTime)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Printf("patch id %v not found in patch status table!\n", patchID)
		} else { // catchall
			logger.Fatalf("error getting patch status. err: %v\n", err)
		}
		return err
	}
	// Patch is applied
	logger.Printf("patch applied at %v\n", patchTime)
	return nil
}

func SetPatchStatus(patchID string) {
	sql := fmt.Sprintf("INSERT INTO %v (id) VALUES ($1);", patchStatusTableName)
	doSQLtx(patchStatusTableName, sql, patchID)
}
