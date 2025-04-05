package consts

// config related constants
const (
	DbTypeSqlite   = "sqlite3"
	DbTypePostgres = "postgres"
	DbTypeMySQL    = "mysql"
)

// error codes
// for now, We'll define error codes as we need
// later, it has to restructured
const (
	ErrLoadingConfig      = 1001
	ErrDatabaseSaveFailed = 1002
	ErrBadRequest         = 1003
)
