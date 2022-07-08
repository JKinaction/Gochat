package model

import (
	"database/sql"
	"fmt"
)

const (
	dbhost     = "localhost"
	dbport     = 5432
	dbuser     = "postgres"
	dbpassword = "123"
	dbname     = "test"
)

var db *sql.DB

func Initsql() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)
	//连接数据库
	db, _ = connectDB(psqlInfo)
}
func connectDB(dbURL string) (*sql.DB, error) {
	Db, err := sql.Open("postgres", dbURL)

	if err != nil {
		return nil, err
	}

	err = Db.Ping()

	if err != nil {
		return nil, err
	}

	_, err = Db.Exec(`
    CREATE TABLE IF NOT EXISTS messages (
      username VARCHAR(64),
	  msg VARCHAR(100),
	  color VARCHAR(10),
	  time VARCHAR(10),
      CHECK (CHAR_LENGTH(TRIM(user)) > 0)
    );
  `)

	if err != nil {
		return nil, err
	}
	_, err = Db.Exec(`
    CREATE TABLE IF NOT EXISTS primess (
      username VARCHAR(64),
	  msg VARCHAR(100),
      touser varchar(64),  
	  time VARCHAR(10),      
      CHECK (CHAR_LENGTH(TRIM(user)) > 0)
    );
  `)

	if err != nil {
		return nil, err
	}
	return Db, nil
}
