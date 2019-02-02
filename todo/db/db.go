package db

import (
	"database/sql"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DB - sql.DB wrapper
type DB struct {
	*sql.DB
}

// Connect - connect to a database and return it
func Connect() (*DB, error) {
	connectionString := fmt.Sprintf("%s:%s@/%s?parseTime=true", "root", "Root1234", "todo")
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// Disconnect - disconnect from specified database
func Disconnect(db *DB) error {
	err := db.Close()
	return err
}
