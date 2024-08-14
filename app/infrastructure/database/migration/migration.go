package migration

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/krisdioles/ppr-wallet/app/domain"
)

func CreateUserBalancesTable(db *sqlx.DB) {
	createUserBalancesTableDDL := `CREATE TABLE IF NOT EXISTS user_balances (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) NOT NULL,
		balance INTEGER,
		bank_code VARCHAR(50),
		account_no VARCHAR(50),
		account_name VARCHAR(100),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	log.Println("Create user_balances table...")
	result := db.MustExec(createUserBalancesTableDDL)
	log.Println("user_balances table created, result:", result)
}

func InsertUserBalancesRecord(db *sqlx.DB, userBalance domain.UserBalance) {
	insertUserBalanceDML := `INSERT INTO user_balances(username, balance, bank_code, account_no, account_name) VALUES(:username, :balance, :bank_code, :account_no, :account_name)`

	log.Println("Insert user_balances...")
	result, err := db.NamedExec(insertUserBalanceDML, userBalance)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Insert success. Result:", result)
}

func CreateJournalEntriesTable(db *sqlx.DB) {
	createJournalEntriesTableDDL := `CREATE TABLE IF NOT EXISTS journal_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transaction_name VARCHAR(100),
		account_id INTEGER NOT NULL,
		debit_amount INTEGER,
		credit_amount INTEGER,
		folio VARCHAR(50)
	);`

	log.Println("Create journal_entries table...")
	db.MustExec(createJournalEntriesTableDDL)
	log.Println("journal_entries table created.")
}
