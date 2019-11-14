package pokemon

import (
	"golang-starter-pack/model"
)

// Store for pokemon
type Store interface {
	GetBySlug(string) (*model.Pokemon, error)
	GetTrainerPokemonBySlug(userID uint, slug string) (*model.Pokemon, error)
	CreatePokemon(*model.Pokemon) error
	UpdatePokemon(*model.Pokemon, []string) error
	DeletePokemon(*model.Pokemon) error
	List(offset, limit int) ([]model.Pokemon, int, error)
	ListByTag(tag string, offset, limit int) ([]model.Pokemon, int, error)
	ListByTrainer(username string, offset, limit int) ([]model.Pokemon, int, error)
	ListByWhoFavorited(username string, offset, limit int) ([]model.Pokemon, int, error)

	AddPower(*model.Pokemon, *model.Power) error
	GetPowersBySlug(string) ([]model.Power, error)
	GetPowerByID(uint) (*model.Power, error)
	DeletePower(*model.Power) error

	AddFavorite(*model.Pokemon, uint) error
	RemoveFavorite(*model.Pokemon, uint) error
	ListPokeTags() ([]model.PokeTag, error)
}
