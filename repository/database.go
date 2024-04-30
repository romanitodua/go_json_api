package repository

import (
	"fmt"
	"go-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgresDB struct {
	instance *gorm.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	dsn := "host=localhost user=postgres password=romaroma dbname=banking port=5432 sslmode=disable TimeZone=Etc/GMT+4"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{QueryFields: true})
	if err != nil {
		log.Fatal(err)
	}

	return &PostgresDB{
		instance: db,
	}, nil
}

func (db PostgresDB) GetUserById(id string) (*User, error) {
	var user User
	err := db.instance.Model(&User{}).
		Joins("LEFT JOIN accounts ON users.id = accounts.user_id").
		Joins("LEFT JOIN transactions ON users.id = transactions.user_id").
		Where("users.id = ?", id).
		Preload("Accounts").
		Preload("Transactions").
		First(&user).
		Error

	user.Password = ""
	if err != nil {
		return &user, err
	}
	return &user, nil
}

func (db PostgresDB) LoginUser(id string, password string) bool {
	var count int64
	err := db.instance.Model(&User{}).Where(&User{ID: id, Password: password}).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return count > 0
}
func (db PostgresDB) InsertUser(u *User) error {
	err := db.instance.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (db PostgresDB) InsertAccount(a *Account) (error, string) {
	var accountNumbers []string
	err := db.instance.Model(&Account{}).Pluck("account_number", &accountNumbers).Error
	if err != nil {
		return err, ""
	}
	var accountNumber string
	isUnique := false

	for !isUnique {
		accountNumber = utils.GenerateAccountNumber()
		isUnique = !utils.Contains(accountNumbers, accountNumber)
	}
	a.AccountNumber = accountNumber
	err = db.instance.Create(a).Error
	if err != nil {
		return err, ""
	}
	return nil, accountNumber
}
