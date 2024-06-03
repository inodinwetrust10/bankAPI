package main

import (
	"math/rand"
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAccountParams struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func NewAccount(FirstName, LastName string) *Account {
	return &Account{
		FirstName: FirstName,
		LastName:  LastName,
		Number:    int64(rand.Intn(1000000)),
		Balance:   0,
	}
}

type TransferParams struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}
