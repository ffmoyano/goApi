package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

var dbPool *sql.DB

// Open is used at the beginning of app.go to initialize the database
func Open() {
	var err error
	dbPool, err = sql.Open(os.Getenv("dbDriver"), os.Getenv("dburl"))
	if err != nil {
		log.Fatal(err)
	}
	dbPool.SetMaxOpenConns(20)
	dbPool.SetMaxIdleConns(20)
	dbPool.SetConnMaxIdleTime(60 * time.Second)
	dbPool.SetConnMaxLifetime(30 * time.Minute)
}

// Get is called by service methods to make available a database connection
func Get() *sql.DB {
	return dbPool
}

// Close is deferred after Open() in app.go to close the db
func Close() {
	dbPool.Close()
}
