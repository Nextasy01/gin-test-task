package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const schemaQuery = "test.sql"

type Database struct {
	conn *sql.DB
}

func NewDatabase() Database {
	return Database{
		conn: Setup(),
	}
}

func Setup() *sql.DB {
	var dbfile *os.File
	var info fs.FileInfo
	var err error

	if info, err = os.Stat(os.Getenv("DBSourceName")); errors.Is(err, os.ErrNotExist) {
		dbfile, err = os.Create(os.Getenv("DBSourceName"))
		if err != nil {
			panic(err)
		}
		defer dbfile.Close()
	}

	db, err := sql.Open(os.Getenv("DBdriver"), os.Getenv("DBSourceName"))
	if err != nil {
		panic(err)
	}

	if info == nil {
		err = CreateDBScheme(db)
		if err != nil {
			panic(err)
		}
	}

	return db
}

func CreateDBScheme(db *sql.DB) error {
	query, err := os.ReadFile(schemaQuery)
	if err != nil {
		return fmt.Errorf("couldn't find or open '%s' file: %v", schemaQuery, err)
	}

	sql := string(query)
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("couldn't execute query: %v", err)
	}

	return nil
}
