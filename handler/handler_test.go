package handler

import (
	"log"
	"os"
	"testing"

	"encoding/json"

	"golang-starter-pack/db"
	"golang-starter-pack/model"
	"golang-starter-pack/pokemon"
	"golang-starter-pack/router"
	"golang-starter-pack/store"
	"golang-starter-pack/trainer"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
)

var (
	d  *gorm.DB
	us trainer.Store
	as pokemon.Store
	h  *Handler
	e  *echo.Echo
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func authHeader(token string) string {
	return "Token " + token
}

func setup() {
	d = db.TestDB()
	db.AutoMigrate(d)
	us = store.NewTrainerStore(d)
	as = store.NewPokemonStore(d)
	h = NewHandler(us, as)
	e = router.New()
	loadFixtures()
}

func tearDown() {
	_ = d.Close()
	if err := db.DropTestDB(); err != nil {
		log.Fatal(err)
	}
}

func responseMap(b []byte, key string) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m[key].(map[string]interface{})
}

func loadFixtures() error {
	u1bio := "user1 bio"
	u1image := "http://realworld.io/user1.jpg"
	u1 := model.Trainer{
		Username: "user1",
		Email:    "user1@realworld.io",
		Bio:      &u1bio,
		Image:    &u1image,
	}
	u1.Password, _ = u1.HashPassword("secret")
	if err := us.Create(&u1); err != nil {
		return err
	}

	u2bio := "user2 bio"
	u2image := "http://realworld.io/user2.jpg"
	u2 := model.Trainer{
		Username: "user2",
		Email:    "user2@realworld.io",
		Bio:      &u2bio,
		Image:    &u2image,
	}
	u2.Password, _ = u2.HashPassword("secret")
	if err := us.Create(&u2); err != nil {
		return err
	}
	us.AddBadge(&u2, u1.ID)

	a := model.Pokemon{
		Slug:        "article1-slug",
		Name:        "pokemon1 title",
		Description: "article1 description",
		Level:       1,
		OwnerID:     1,
		Tags: []model.PokeTag{
			{
				Tag: "tag1",
			},
			{
				Tag: "tag2",
			},
		},
	}
	// TO-FIX
	as.CreatePokemon(&a)
	// as.AddPower(&a, &model.Power{
	// 	name:  "power",
	// 	power: 1,
	// })

	// a2 := model.Pokemon{
	// 	Slug:        "article2-slug",
	// 	Name:        "article2 title",
	// 	Description: "article2 description",
	// 	Level:       1,
	// 	OwnerID:     2,
	// 	Favorites: []model.Pokemon{
	// 		u1,
	// 	},
	// 	Tags: []model.Tag{
	// 		{
	// 			Tag: "tag1",
	// 		},
	// 	},
	// }
	// as.CreatePokemon(&a2)
	// as.AddPower(&a2, &model.Power{
	// 	Body:      "article2 comment1 by user1",
	// 	ArticleID: 2,
	// 	UserID:    1,
	// })
	// as.AddFavorite(&a2, 1)

	return nil
}
