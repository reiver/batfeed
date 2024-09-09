package dbsrv

import (
	"database/sql"

	"github.com/reiver/go-erorr"
	_ "modernc.org/sqlite"
)

func Feeds(domain string) ([]string, error) {
	var rows *sql.Rows
	{
		if nil == db {
			return nil, errNilDB
		}

		const code string =
		`SELECT name`         +"\n"+
		`FROM feeds`       +"\n"+
		`WHERE domain = ?`

		var err error
		rows, err = db.Query(code, domain)
		if nil != err {
			return nil, erorr.Errorf("dbsrv: problem querying 'feeds' database table: %w", err)
		}
		if nil == rows {
			return nil, erorr.Error("dbsrv: problem querying 'feeds' database table: nil rows")
		}
		defer rows.Close()
	}

	var feeds []string
	{
		for rows.Next() {
			var feed string

			err := rows.Scan(&feed)
			if nil != err {
				return nil, erorr.Errorf("dbsrv: decoding problem querying 'feeds' database table: %w", err)
			}

			feeds = append(feeds, feed)
		}
		if err := rows.Err(); nil != err {
			return nil, erorr.Errorf("dbsrv: post-loop problem querying 'feeds' database table: %w", err)
		}
	}

	return feeds, nil
}
