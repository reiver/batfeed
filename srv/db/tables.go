package dbsrv

import (
	"database/sql"

	. "github.com/reiver/batfeed/srv/log"
)

func initTables(database *sql.DB) error {
	if nil == database {
		Logf("[dbsrv] problem initialzing tables: %s", errNilDB)
		return errNilDB
	}

	var tables []string

	{
		const table string =
		`CREATE TABLE IF NOT EXISTS users` +"\n"+
		`( id INTEGER PRIMARY KEY`         +"\n"+
		`, did TEXT   UNIQUE NOT NULL`     +"\n"+
		`)`
		tables = append(tables, table)
	}

	{
		const table string =
		`CREATE TABLE IF NOT EXISTS feeds`            +"\n"+
		`( id      INTEGER PRIMARY KEY`               +"\n"+
		`, user_id INTEGER NOT NULL`                  +"\n"+
		`, name     TEXT   UNIQUE NOT NULL`           +"\n"+
		`, FOREIGN KEY(user_id) REFERENCES users(id)` +"\n"+
		`)`
		tables = append(tables, table)
	}

	for _, createTable := range tables {
		Logf("TABLE:\n%s", createTable)

		if _, err := database.Exec(createTable); nil != err {
			return err
		}
	}

	return nil
}
