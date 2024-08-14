package sqlite3

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/krisdioles/ppr-wallet/app/domain"
	"github.com/krisdioles/ppr-wallet/app/infrastructure/database/migration"
	_ "github.com/mattn/go-sqlite3"
)

func Init() *sqlx.DB {
	os.Remove("database.db")

	log.Println("Creating sqlite3 database.db...")
	file, err := os.Create("database.db")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	log.Println("sqlite3 database.db created.")

	sqliteDb := sqlx.MustConnect("sqlite3", "./database.db")

	migration.CreateUserBalancesTable(sqliteDb)
	migration.CreateJournalEntriesTable(sqliteDb)

	// development purpose, should be deleted when ready to be pushed to prod
	migration.InsertUserBalancesRecord(sqliteDb, domain.UserBalance{
		Username:    "andy123",
		Balance:     10000,
		BankCode:    "arthagraha",
		AccountNo:   "083012322138",
		AccountName: "Andy Garcia",
	})
	migration.InsertUserBalancesRecord(sqliteDb, domain.UserBalance{
		Username:    "brandy345",
		Balance:     15000,
		BankCode:    "bca",
		AccountNo:   "0810123456878",
		AccountName: "Brandy Joe",
	})
	migration.InsertUserBalancesRecord(sqliteDb, domain.UserBalance{
		Username:    "cindy789",
		Balance:     8000,
		BankCode:    "cempakabank",
		AccountNo:   "11298800345",
		AccountName: "Cindy Kat",
	})

	return sqliteDb
}
