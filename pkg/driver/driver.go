package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib" // drivers
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenConns = 10
const maxIdleConns = 5
const connMaxLifetime = 5 * time.Minute

// ConnectSQL creates a database pool for postgres
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)

	if err != nil {
		panic(err) // crash if connection fails
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	dbConn.SQL = db

	if err := testDB(db); err != nil {
		return nil, err
	} else {
		return dbConn, nil
	}
}

// testDB tries to ping the database
func testDB(db *sql.DB) error {
	return db.Ping()
}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
