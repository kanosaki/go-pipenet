package pipenet

import (
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"database/sql"
)

func CreateDbMap(dbtype, dbparam string, dialect gorp.Dialect) (*gorp.DbMap, error) {
	db, err := sql.Open(dbtype, dbparam)
	if err != nil {
		return nil, err
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: dialect}
	return dbmap, nil
}

// returns "uninitialized (without any table) dbmap"
func CreateSqliteDbMap(dbfile string) (*gorp.DbMap, error) {
	return CreateDbMap("sqlite3", dbfile, gorp.SqliteDialect{})
}

func CreatePostgresDbMap(connstring string) (*gorp.DbMap, error) {
	return CreateDbMap("postgres", connstring, gorp.PostgresDialect{})
}