package dbsrv

import (
	"database/sql"

	"github.com/reiver/go-erorr"
	_ "modernc.org/sqlite"

	"github.com/reiver/batfeed/cfg"
	. "github.com/reiver/batfeed/srv/log"
)

var db *sql.DB

func initDB(filename string) (*sql.DB, error) {

	var database *sql.DB
	{
		Logf("[dbsrv] opening database file %q ...", filename)

		var err error
		database, err = sql.Open("sqlite", filename)
		if nil != err {
			Logf("[dbsrv] FAILURE — could not open database file %q: %w", filename, err)
			return nil, erorr.Errorf("dbsrv: problem opening SQLite database file %q: %w", filename, err)
		}
		if nil == database {
			Logf("[dbsrv] FAILURE — could not open database file %q: nil database", filename)
			return nil, erorr.Errorf("dbsrv: nil database wen opening SQLite database file %q", filename)
		}

		Logf("[dbsrv] success — opened database file %q", filename)
	}

	{
		err := initTables(database)
		if nil != err {
			Logf("[dbsrv] FAILURE — problem initializing tables in database file %q", filename)
			return database, err
		}
	}

	return database, nil
}

func init() {
	const filename string = cfg.SQLiteDBName

	var err error
	db, err = initDB(filename)
	if nil != err {
		panic(err)
	}
	if nil == db {
		panic("nil database")
	}
}
