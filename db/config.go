package db

import "database/sql"

const (
	Host     = "localhost"
	Port     = 5432
	UserPq   = "postgres"
	Password = "1234"
	Dbname   = "finalproject"
)

var (
	Db  *sql.DB
	Err error
)
