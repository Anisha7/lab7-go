package trainer

import (
	"golang-starter-pack/model"
)

// Store for trainer
type Store interface {
	GetByID(uint) (*model.Trainer, error)
	GetByEmail(string) (*model.Trainer, error)
	GetByUsername(string) (*model.Trainer, error)
	Create(*model.Trainer) error
	Update(*model.Trainer) error
	AddBadge(trainer *model.Trainer, badgeID uint) error
	RemoveBadge(trainer *model.Trainer, badgeID uint) error
	HasBadge(trainerID, badgeID uint) (bool, error)
}
