package config

import (
	"os"
	"path/filepath"

	"github.com/artifact-space/ArtiSpace/consts"
	"github.com/artifact-space/ArtiSpace/log"
)

var config *AppConfig

type AppConfig struct {
	Storage  StorageConfig
	Database DatabaseConfig
	Server   ServerConfig
}

type ServerConfig struct {
	Port             int
	DockerV2Port     int
	HostnameOrNodeIP string
}

type StorageConfig struct {
	StorageClaass string
	Properties    map[string]string
}

type DatabaseConfig struct {
	Type   string
	Sqlite SqliteConfig
}

type SqliteConfig struct {
	FilePath          string
	CreateTableScriptPath string
}

func defaultConfig() *AppConfig {

	cwd, err := os.Getwd()
	if err != nil {
		log.Logger().Error().Err(err).Msg("unable to get current working directory")
		return nil
	}

	fsStoragePath := cwd + "/temp/artispace"

	return &AppConfig{
		Storage: StorageConfig{
			StorageClaass: "fs",
			Properties: map[string]string{
				"fs.storage.path": fsStoragePath,
			},
		},
		Database: DatabaseConfig{
			Type: consts.DbTypeSqlite,
			Sqlite: SqliteConfig{
				FilePath: filepath.Join(cwd, "temp", "sqlite", "app.db"),
				CreateTableScriptPath: "sql-scripts/sqlite.sql",
			},
		},
		Server: ServerConfig{
			Port:             8000,
			DockerV2Port:     7000,
			HostnameOrNodeIP: "localhost",
		},
	}
}

func Config() (*AppConfig, error) {
	//for now, We'll just return default config
	// later, we have to parse configuration from file and return
	if config == nil {
		config = defaultConfig()
	}
	return config, nil
}
