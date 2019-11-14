package store

import (
	"golang-starter-pack/model"

	"github.com/jinzhu/gorm"
)

// PokemonStore ..
type PokemonStore struct {
	db *gorm.DB
}

// NewPokemonStore ..
func NewPokemonStore(db *gorm.DB) *PokemonStore {
	return &PokemonStore{
		db: db,
	}
}

// GetBySlug ..
func (as *PokemonStore) GetBySlug(s string) (*model.Pokemon, error) {
	var m model.Pokemon
	err := as.db.Where(&model.Pokemon{Slug: s}).Preload("Favorites").Preload("PokeTags").Preload("Trainer").Find(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

// GetTrainerPokemonBySlug ..
func (as *PokemonStore) GetTrainerPokemonBySlug(trainerID uint, slug string) (*model.Pokemon, error) {
	var m model.Pokemon
	err := as.db.Where(&model.Pokemon{Slug: slug, OwnerID: trainerID}).Find(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

// CreatePokemon ..
func (as *PokemonStore) CreatePokemon(a *model.Pokemon) error {
	tags := a.Tags
	tx := as.db.Begin()
	if err := tx.Create(&a).Error; err != nil {
		return err
	}
	for _, t := range a.Tags {
		err := tx.Where(&model.PokeTag{Tag: t.Tag}).First(&t).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&a).Association("PokeTags").Append(t).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where(a.ID).Preload("Favorites").Preload("PokeTags").Preload("Trainer").Find(&a).Error; err != nil {
		tx.Rollback()
		return err
	}
	a.Tags = tags
	return tx.Commit().Error
}

// UpdatePokemon ..
func (as *PokemonStore) UpdatePokemon(a *model.Pokemon, tagList []string) error {
	tx := as.db.Begin()
	if err := tx.Model(a).Update(a).Error; err != nil {
		return err
	}
	tags := make([]model.PokeTag, 0)
	for _, t := range tagList {
		tag := model.PokeTag{Tag: t}
		err := tx.Where(&tag).First(&tag).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			return err
		}
		tags = append(tags, tag)
	}
	if err := tx.Model(a).Association("Tags").Replace(tags).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(a.ID).Preload("Favorites").Preload("PokeTags").Preload("Trainer").Find(a).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// DeleteArticle ..
func (as *PokemonStore) DeleteArticle(a *model.Pokemon) error {
	return as.db.Delete(a).Error
}

func (as *PokemonStore) List(offset, limit int) ([]model.Pokemon, int, error) {
	var (
		pokemons []model.Pokemon
		count    int
	)
	as.db.Model(&pokemons).Count(&count)
	as.db.Preload("Favorites").Preload("PokeTags").Preload("Trainer").Offset(offset).Limit(limit).Order("created_at desc").Find(&pokemons)
	return pokemons, count, nil
}

// ListByTag ..
func (as *PokemonStore) ListByTag(tag string, offset, limit int) ([]model.Pokemon, int, error) {
	var (
		t        model.PokeTag
		pokemons []model.Pokemon
		count    int
	)
	err := as.db.Where(&model.Tag{Tag: tag}).First(&t).Error
	if err != nil {
		return nil, 0, err
	}
	as.db.Model(&t).Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Association("Articles").Find(&pokemons)
	count = as.db.Model(&t).Association("Articles").Count()
	return pokemons, count, nil
}

// ListByAuthor ..
func (as *PokemonStore) ListByAuthor(username string, offset, limit int) ([]model.Pokemon, int, error) {
	var (
		u        model.Trainer
		pokemons []model.Pokemon
		count    int
	)
	err := as.db.Where(&model.User{Username: username}).First(&u).Error
	if err != nil {
		return nil, 0, err
	}
	as.db.Where(&model.Article{AuthorID: u.ID}).Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Find(&pokemons)
	as.db.Where(&model.Article{AuthorID: u.ID}).Model(&model.Article{}).Count(&count)

	return pokemons, count, nil
}

// ListByWhoFavorited ..
func (as *PokemonStore) ListByWhoFavorited(username string, offset, limit int) ([]model.Pokemon, int, error) {
	var (
		u        model.Trainer
		pokemons []model.Pokemon
		count    int
	)
	err := as.db.Where(&model.Trainer{Username: username}).First(&u).Error
	if err != nil {
		return nil, 0, err
	}
	as.db.Model(&u).Preload("Favorites").Preload("PokeTags").Preload("Trainer").Offset(offset).Limit(limit).Order("created_at desc").Association("Favorites").Find(&pokemons)
	count = as.db.Model(&u).Association("Favorites").Count()
	return pokemons, count, nil
}

// AddPower ..
func (as *PokemonStore) AddPower(a *model.Pokemon, c *model.Power) error {
	err := as.db.Model(a).Association("Powers").Append(c).Error
	if err != nil {
		return err
	}
	return as.db.Where(c.ID).Preload("Pokemon").First(c).Error
}

// GetPowersBySlug ..
func (as *PokemonStore) GetPowersBySlug(slug string) ([]model.Power, error) {
	var m model.Pokemon
	if err := as.db.Where(&model.Pokemon{Slug: slug}).Preload("Powers").Preload("Powers.Trainer").First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return m.Powers, nil
}

// GetPowerByID ..
func (as *PokemonStore) GetPowerByID(id uint) (*model.Power, error) {
	var m model.Power
	if err := as.db.Where(id).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// DeletePower ..
func (as *PokemonStore) DeletePower(c *model.Power) error {
	return as.db.Delete(c).Error
}

// AddFavorite ..
func (as *PokemonStore) AddFavorite(a *model.Pokemon, trainerID uint) error {
	usr := model.Trainer{}
	usr.ID = trainerID
	return as.db.Model(a).Association("Favorites").Append(&usr).Error
}

// RemoveFavorite ..
func (as *PokemonStore) RemoveFavorite(a *model.Pokemon, trainerID uint) error {
	usr := model.Trainer{}
	usr.ID = trainerID
	return as.db.Model(a).Association("Favorites").Delete(&usr).Error
}

// ListPokeTags ..
func (as *PokemonStore) ListPokeTags() ([]model.PokeTag, error) {
	var tags []model.PokeTag
	if err := as.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
