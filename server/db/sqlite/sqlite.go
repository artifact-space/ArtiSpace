package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	sqlite3 "modernc.org/sqlite"

	"github.com/artifact-space/ArtiSpace/config"
	"github.com/artifact-space/ArtiSpace/db/hooks"
	"github.com/artifact-space/ArtiSpace/log"
)

type sqliteDb struct {
	db *sql.DB
}

func SqliteDB(config *config.SqliteConfig) (*sqliteDb, error) {
	log.Logger().Debug().Msgf("Initiating Sqlite: %s", config.FilePath)

	err := os.MkdirAll(filepath.Dir(config.FilePath), 0777)
	if err != nil {
		log.Logger().Error().Err(err).Msgf("unable to create directories: %s", filepath.Dir(config.FilePath))
		return nil, err
	}

	sql.Register("sqlite-hooked", &hooks.HookedDriver{Driver: &sqlite3.Driver{}})

	db, err := sql.Open("sqlite-hooked", fmt.Sprintf("file:%s?cache=shared&_fk=1", config.FilePath))
	if err != nil {
		log.Logger().Error().Err(err).Msg("error in creating connections to sqlite")
		return nil, err
	}

	// create tables
	sqlBytes, err := os.ReadFile(config.CreateTableScriptPath)
	if err != nil {
		log.Logger().Error().Err(err).Msgf("unable to read queries to create tables: %s", config.CreateTableScriptPath)
		return nil, err
	}

	sqlStmts := string(sqlBytes)
	_, err = db.Exec(sqlStmts)
	if err != nil {
		log.Logger().Error().Err(err).Msgf("error executing scripts from: %s", config.CreateTableScriptPath)
		return nil, err
	}

	return &sqliteDb{
		db,
	}, nil
}

func (s *sqliteDb) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *sqliteDb) Release() {

}
