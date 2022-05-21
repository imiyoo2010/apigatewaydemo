package data

import (
	"database/sql"
	//_ "github.com/go-sql-driver/mysql" //mysql driver
	_ "github.com/mattn/go-sqlite3" 	// sqlite3 dirver
)

type DB struct {

	ReadDb 	*sql.DB

	WriteDb *sql.DB

}

func NewDB() *DB {

	d := new(DB)

	ReadDb, err := sql.Open("sqlite3","abc.db")

	if err != nil {
		return nil
	}

	WriteDb, err := sql.Open("sqlite3","abc.db")

	if err != nil {
		return nil
	}

	d.ReadDb = ReadDb
	d.WriteDb = WriteDb

	return d
}



