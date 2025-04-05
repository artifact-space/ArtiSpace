package db

import (
	"context"
	"database/sql"

	"github.com/artifact-space/ArtiSpace/config"
	"github.com/artifact-space/ArtiSpace/db/sqlite"
)

type DatabaseProvider interface {

	//Tx
	Begin() (*sql.Tx, error)
	Release()

	//Docker V2
	CreateDockerNamespaceAndRepositoryIfMissing(ctx context.Context,  namespace string, repository string) error
}

var Provider DatabaseProvider

func InitDB(config *config.DatabaseConfig) error {
	//for now, we'll only support sqlite
	sqlite, err := sqlite.SqliteDB(&config.Sqlite)
	if err != nil {
		return err
	}

	Provider = sqlite

	return nil
}

func Release() {
	Provider.Release()
}
