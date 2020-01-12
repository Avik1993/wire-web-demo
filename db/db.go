package db

import (
	"context"
	"database/sql"

	"github.com/fsufitch/testable-web-demo/config"
	_ "github.com/lib/pq" // inject PQ database driver
)

// PostgresDBConn is a database connection to a Postgres DB
type PostgresDBConn *sql.DB

// PreInitPostgresDBConn is a database connection to a Postgres DB which may not have initialized schema
type PreInitPostgresDBConn *sql.DB

// ProvidePreInitPostgresDBConn provides a PostgresDBConn by connecting to a database
func ProvidePreInitPostgresDBConn(dbString config.DatabaseString) (PreInitPostgresDBConn, func(), error) {
	db, err := sql.Open("postgres", string(dbString))
	cleanup := func() { db.Close() }

	if err != nil {
		return nil, cleanup, err
	}
	if err := db.Ping(); err != nil {
		return nil, cleanup, err
	}
	return db, cleanup, nil
}

// ProvidePostgresDBConn performs schema initialization
func ProvidePostgresDBConn(db PreInitPostgresDBConn) (PostgresDBConn, error) {
	tx, err := (*sql.DB)(db).BeginTx(context.Background(), nil)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS counter (value int NOT NULL DEFAULT 0);
		INSERT INTO counter (value)
			SELECT 0 WHERE NOT EXISTS (SELECT * FROM counter);
	`)
	if err != nil {
		return nil, err
	}
	return PostgresDBConn(db), tx.Commit()
}
