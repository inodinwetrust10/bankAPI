package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	godotenv.Load()
	conn, err := sql.Open("postgres", os.Getenv("DB_STRING"))
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	connection := &PostgresStore{db: conn}
	return connection, err
}

func (p *PostgresStore) CreateAccount(Account *Account) error {
	_, err := p.db.Exec(`INSERT INTO account (
    first_name,
    last_name,
    balance,
    number)
    VALUES ($1,$2,$3,$4)
    `, Account.FirstName, Account.LastName, Account.Balance, Account.Number)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (p *PostgresStore) DeleteAccount(id int) error {
	_, err := p.db.Exec(`DELETE FROM account WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := p.db.Query(
		`SELECT * FROM account WHERE id=$1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (p *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    number SERIAL UNIQUE NOT NULL,
    balance INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
`
	_, err := p.db.Exec(query)
	fmt.Println("Table created successfully")
	return err
}

func (p *PostgresStore) Init() error {
	err := p.CreateAccountTable()
	return err
}

func (p *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := p.db.Query(`SELECT * from account`)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// Copying row into account struct
// /////////////////////////////////////////////////////
func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	if err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt); err != nil {
		return nil, err
	}
	return account, nil
}
