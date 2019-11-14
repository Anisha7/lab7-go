package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Trainer struct {
	gorm.Model
	Username string `gorm:"unique_index;not null"`
	Email    string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Bio      *string
	Image    *string
	Badges   []Badge   `gorm:"foreignkey:FollowingID"`
	Pokemons []Pokemon `gorm:"foreignkey:FollowerID"`
}

type Badge struct {
	TrainerId uint
	BadgeID   uint
	Level     int
	Title     string
}

func (u *Trainer) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

func (u *Trainer) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}
