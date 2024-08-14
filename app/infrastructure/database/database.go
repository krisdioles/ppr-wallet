package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/krisdioles/ppr-wallet/app/infrastructure/database/sqlite3"
)

func Init() *sqlx.DB {
	return sqlite3.Init()
}
