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
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccount(int) (*Account, error)
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

func (p *PostgresStore) CreateAccount(*Account) error {
	return nil
}

func (p *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (p *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (p *PostgresStore) GetAccount(id int) (*Account, error) {
	return nil, nil
}

func (p *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    number SERIAL UNIQUE NOT NULL
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
