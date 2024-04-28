package main

import "time"

const (
	SENT = iota
	RECEIVED
)
const (
	DEBIT = iota
	CREDIT
	SAVING
)
const (
	ACTIVE = iota
	DISABLED
)

type User struct {
	Name             string        `json:"name"`
	Surname          string        `json:"surname"`
	ID               string        `json:"id" gorm:"primaryKey"`
	Transactions     []Transaction `json:"transactions" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Accounts         []Account     `json:"accounts" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	RegistrationDate time.Time     `json:"registration_date"`
}

type Transaction struct {
	OriginAccountNumber      string    `json:"origin_account_number"`
	DestinationAccountNumber string    `json:"destination_account_number"`
	Type                     int       `json:"type"`
	Date                     time.Time `json:"date"`
	Amount                   int       `json:"amount"`
	Description              string    `json:"description"`
	TransactionID            string    `json:"transaction_id" gorm:"primaryKey"`
	UserID                   string    `json:"user_id" `
}

type Account struct {
	AccountNumber string    `json:"account_number" gorm:"primaryKey"`
	AccountType   int       `json:"account_type"`
	Balance       int       `json:"balance"`
	OpeningDate   time.Time `json:"opening_date"`
	Status        int       `json:"status"`
	UserID        string    `json:"user_id"`
}
