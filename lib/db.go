package lib

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// TODO: support socket

type DbConnection struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (dbConn DbConnection) toDatabaseURL() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v", dbConn.User, dbConn.Password, dbConn.Host, dbConn.Port, dbConn.Database)
}

func (dbConn DbConnection) connect() {
	pgx.Connect(context.Background(), dbConn.toDatabaseURL())
}

// func (dbConn DbConnection) disconnect() {
// 	pgx.
// }

func (dbConn DbConnection) test() {

}
