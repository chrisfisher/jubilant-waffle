package main

import (
	"github.com/chrisfisher/jubilant-waffle/db"
	"github.com/chrisfisher/jubilant-waffle/models"
	"github.com/chrisfisher/jubilant-waffle/repositories"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var context *db.Context

func main() {
	context = db.NewContext()
	defer context.Close()
	createFilms()
	createUsers()
}

func createFilms() {
	fc := context.DbCollection("films")
	fr := &repositories.FilmRepository{C: fc}

	fr.Create(&models.Film{
		Id:          bson.ObjectIdHex("5a498a8bb518634b9454c1e3"),
		Title:       "Captain America: The Winter Soldier",
		Description: "Steve Rogers teams up with Black Widow to battle an assassin known as the Winter Soldier.",
		Rating:      "PG13",
		Reviews: []models.Review{models.Review{
			Id:       bson.NewObjectId(),
			Stars:    4,
			Comments: "Superior to the first Captain America in every way and the best Marvel stand-alone movie.",
		}},
	})

	fr.Create(&models.Film{
		Id:          bson.ObjectIdHex("5a498a8bb518634b9454c1e4"),
		Title:       "Point Break",
		Description: "An FBI agent goes undercover to catch a gang of surfers who may be bank robbers.",
		Rating:      "R",
		Reviews: []models.Review{models.Review{
			Id:       bson.NewObjectId(),
			Stars:    4,
			Comments: "Point Break is a perfect example of the contemporary B movie.",
		}},
	})
}

func createUsers() {
	uc := context.DbCollection("users")
	ur := &repositories.UserRepository{C: uc}

	ur.Create(&models.User{
		Id:   bson.ObjectIdHex("5a498a8bb518634b9454c1e7"),
		Name: "Chris Fisher",
		Viewings: []models.Viewing{
			models.Viewing{
				Id:        bson.NewObjectId(),
				StartTime: time.Date(2017, time.December, 23, 15, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2017, time.December, 23, 17, 0, 0, 0, time.UTC),
				FilmId:    bson.ObjectIdHex("5a498a8bb518634b9454c1e3"),
			},
			models.Viewing{
				Id:        bson.NewObjectId(),
				StartTime: time.Date(2017, time.December, 24, 15, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2017, time.December, 24, 17, 0, 0, 0, time.UTC),
				FilmId:    bson.ObjectIdHex("5a498a8bb518634b9454c1e4"),
			},
		},
	})

	ur.Create(&models.User{
		Id:   bson.ObjectIdHex("5a498a8bb518634b9454c1e8"),
		Name: "Jacques Cousteau",
		Viewings: []models.Viewing{
			models.Viewing{
				Id:        bson.NewObjectId(),
				StartTime: time.Date(2017, time.December, 11, 11, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2017, time.December, 11, 13, 0, 0, 0, time.UTC),
				FilmId:    bson.ObjectIdHex("5a498a8bb518634b9454c1e3"),
			},
			models.Viewing{
				Id:        bson.NewObjectId(),
				StartTime: time.Date(2017, time.December, 12, 18, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2017, time.December, 12, 20, 0, 0, 0, time.UTC),
				FilmId:    bson.ObjectIdHex("5a498a8bb518634b9454c1e4"),
			},
		},
	})
}
