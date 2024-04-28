package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgresDB struct {
	instance *gorm.DB
}

func newPostgresDB() (*PostgresDB, error) {
	dsn := "host=localhost user=postgres password=romaroma dbname=banking port=5432 sslmode=disable TimeZone=Etc/GMT+4"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{QueryFields: true})
	if err != nil {
		log.Fatal(err)
	}

	return &PostgresDB{
		instance: db,
	}, nil
}

func (db PostgresDB) getUserById(id string) (*User, error) {
	var user User
	err := db.instance.Model(&User{}).
		Joins("LEFT JOIN accounts ON users.id = accounts.user_id").
		Joins("LEFT JOIN transactions ON users.id = transactions.user_id").
		Where("users.id = ?", id).
		Preload("Accounts").
		Preload("Transactions").
		First(&user).
		Error
	if err != nil {
		fmt.Println(err)
	}
	return &user, nil
}

func (db PostgresDB) insertUser(u *User) {
	db.instance.Create(u)
}
