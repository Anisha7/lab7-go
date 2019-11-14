package store

import (
	"golang-starter-pack/model"

	"github.com/jinzhu/gorm"
)

// TrainerStore ..
type TrainerStore struct {
	db *gorm.DB
}

// NewTrainerStore ..
func NewTrainerStore(db *gorm.DB) *TrainerStore {
	return &TrainerStore{
		db: db,
	}
}

// GetByID ..
func (us *TrainerStore) GetByID(id uint) (*model.Trainer, error) {
	var m model.Trainer
	if err := us.db.First(&m, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// GetByEmail ..
func (us *TrainerStore) GetByEmail(e string) (*model.Trainer, error) {
	var m model.Trainer
	if err := us.db.Where(&model.Trainer{Email: e}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// GetByUsername ..
func (us *TrainerStore) GetByUsername(username string) (*model.Trainer, error) {
	var m model.Trainer
	if err := us.db.Where(&model.Trainer{Username: username}).Preload("Followers").First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// Create ..
func (us *TrainerStore) Create(u *model.Trainer) (err error) {
	return us.db.Create(u).Error
}

// Update ..
func (us *TrainerStore) Update(u *model.Trainer) error {
	return us.db.Model(u).Update(u).Error
}

// AddBadge ..
func (us *TrainerStore) AddBadge(u *model.Trainer, trainerID uint) error {
	return us.db.Model(u).Association("Badge").Append(&model.Badge{Level: 1, Title: "Badge", TrainerId: trainerID, BadgeID: u.ID}).Error
}

// RemoveBadge ..
func (us *TrainerStore) RemoveBadge(u *model.Trainer, trainerID uint) error {
	f := model.Badge{
		TrainerId: trainerID,
		BadgeID:   u.ID,
	}
	if err := us.db.Model(u).Association("Followers").Find(&f).Error; err != nil {
		return err
	}
	if err := us.db.Delete(f).Error; err != nil {
		return err
	}
	return nil
}

// HasBadge ..
func (us *TrainerStore) HasBadge(trainerID, badgeID uint) (bool, error) {
	var f model.Follow
	if err := us.db.Where("following_id = ? AND follower_id = ?", trainerID, badgeID).Find(&f).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
