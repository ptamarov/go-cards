package dbrepo

import (
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
)

var SQLTest *SQLRepo

// set up database configuration
var cfg = mysql.Config{
	User:   os.Getenv("DBUSER"),
	Passwd: os.Getenv("DBPASS"),
	Net:    "tcp",
	Addr:   "127.0.0.1:3306",
	DBName: "germancorpora",
}

func TestMain(m *testing.M) {
	SQLTest = MySQLDatabaseFromConfig(cfg)
	os.Exit(m.Run())
}
