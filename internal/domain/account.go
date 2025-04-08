package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	Name      string
	Email     string
	APIKey    string
	Balance   float64
	mu        sync.RWMutex
	CreatedAt time.Time
	UpdateAt  time.Time
}

func generateAPIKey() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func NewAccount(name, email string) *Account {
	account := &Account{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Balance:   0,
		APIKey:    generateAPIKey(),
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}

	return account
}

/**
*	Add Balance to the account
 */
func (account *Account) AddBalance(amount float64) {
	account.mu.Lock()         //Lock the balance alteration while the value is being changed
	defer account.mu.Unlock() //Unlock the var at the end of alteration

	account.Balance += amount
	account.UpdateAt = time.Now()
}
