package main

import "database/sql"

func errorCheck(err error) bool {
	if err != nil {
		panic(err.Error())
		return false
	}
	return true
}

func pingDb(db *sql.DB) error {
	err := db.Ping()
	return err
}

func initDb() *sql.DB {
	db, e := sql.Open("mysql", "test:PassworD12312312?@tcp(127.0.0.1)/pymnt_db")
	errorCheck(e)
	return db
}
