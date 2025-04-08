package repository

import (
	"database/sql"
	"time"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (repository *AccountRepository) Save(account *domain.Account) error {
	stmt, err := repository.db.Prepare(`
		INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		account.ID,
		account.Name,
		account.Email,
		account.APIKey,
		account.Balance,
		account.CreatedAt,
		account.UpdateAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repository *AccountRepository) FindByAPIKey(apiKey string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, updatedAt time.Time

	err := repository.db.QueryRow(`
		SELECT id, name, email, api_key, balance, created_at, update_at
		FROM account
		WHERE api_key = $1
	`, apiKey).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdateAt = updatedAt
	return &account, nil
}

func (repository *AccountRepository) FindByID(id string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, updatedAt time.Time

	err := repository.db.QueryRow(`
		SELECT id, name, email, api_key, balance, created_at, update_at
		FROM account
		WHERE ID = $1
	`, id).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdateAt = updatedAt
	return &account, nil
}

func (repository *AccountRepository) UpdateBalance(account *domain.Account) error {
	transaction, err := repository.db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	var currentBalance float64
	err = transaction.QueryRow(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`,
		account.ID).Scan(&currentBalance)

	if err == sql.ErrNoRows {
		return domain.ErrAccountNotFound
	}
	if err != nil {
		return err
	}
	_, err = transaction.Exec(`
		UPDATE accounts
		SET balance = $q, updated_at = $2
		WHERE id = $3
	`, account.Balance, time.Now(), account.ID)
	if err != nil {
		return err
	}
	return transaction.Commit()
}
