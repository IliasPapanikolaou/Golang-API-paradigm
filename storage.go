package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	UpdateAccount(*Account) error
	DeleteAccount(int) error
}

type PostgresStore struct {
	db *sql.DB
}

func newPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		number BIGINT NOT NULL,
		balance DOUBLE PRECISION NOT NULL,
		created_at TIMESTAMP
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatalf("Unable to create table: %v", err)
	}
	// fmt.Println("Table 'account' ensured to exist!")

	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	command := `INSERT INTO account (first_name, last_name, number, balance, created_at) VALUES ($1, $2, $3, $4, $5)`
	resp, err := s.db.Exec(command, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM account ORDER BY id ASC`
	rows, err := s.db.Query(query)
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

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `SELECT * FROM account WHERE ID = $1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	command := `DELETE FROM account WHERE id = $1`

	result, err := s.db.Exec(command, id)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to retrieve rows affected: %w", err)
	}
	fmt.Printf("Deleted %d row(s)\n", rowsAffected)

	return nil
}
