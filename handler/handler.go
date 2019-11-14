package handler

import (
	"golang-starter-pack/pokemon"
	"golang-starter-pack/trainer"
)

// Handler ..
type Handler struct {
	trainerStore trainer.Store
	pokemonStore pokemon.Store
}

// NewHandler ..
func NewHandler(us trainer.Store, as pokemon.Store) *Handler {
	return &Handler{
		trainerStore: us,
		pokemonStore: as,
	}
}
