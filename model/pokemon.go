package model

import (
	"github.com/jinzhu/gorm"
)

type Pokemon struct {
	gorm.Model
	Slug        string `gorm:"unique_index;not null"`
	Name        string `gorm:"not null"`
	Description string
	Level       int
	Owner       Trainer
	OwnerID     uint
	Powers      []Power
	Favorites   []Pokemon `gorm:"many2many:favorites;"`
	Tags        []PokeTag `gorm:"many2many:article_tags;association_autocreate:false"`
}

type Power struct {
	gorm.Model
	name  string
	power int
}

type PokeTag struct {
	gorm.Model
	Tag    string    `gorm:"unique_index"`
	Owners []Trainer `gorm:"many2many:article_tags;"`
}
