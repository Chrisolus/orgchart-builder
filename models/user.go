package models

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique;not null" json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSummary struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type JWT struct {
	Token string `json:"token" binding:"required"`
}

func (u User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return errors.New("cannot hash the string")
	}
	u.Password = string(hash)
	return nil
}

func (u *User) ValidateHash(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
