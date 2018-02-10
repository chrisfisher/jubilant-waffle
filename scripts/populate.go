package main

import (
	"time"

	"github.com/chrisfisher/jubilant-waffle/server/db"
	"github.com/chrisfisher/jubilant-waffle/server/models"
	"github.com/chrisfisher/jubilant-waffle/server/repositories"

	"gopkg.in/mgo.v2/bson"
)

var client *db.Client

func main() {
	client = db.NewClient()
	defer client.Close()
	createFilms()
	createUsers()
}

func createFilms() {
	fc := client.DbCollection("films")
	fr := &repositories.FilmRepository{C: fc}

	fr.Create(&models.Film{
		Id:          bson.ObjectIdHex("5a498a8bb518634b9454c1e1"),
		Title:       "Captain America: The Winter Soldier",
		Description: "Steve Rogers teams up with Black Widow to battle an assassin known as the Winter Soldier.",
		Rating:      "PG13",
		Reviews: []models.Review{models.Review{
			Id:       bson.NewObjectId(),
			Stars:    4,
			Comments: "Superior to the first Captain America in every way and the best Marvel stand-alone movie.",
		}},
		ViewedByUsers: []bson.ObjectId{
			bson.ObjectIdHex("5a498a8bb518634b9454c1e3"),
			bson.ObjectIdHex("5a498a8bb518634b9454c1e4"),
		},
	})

	fr.Create(&models.Film{
		Id:          bson.ObjectIdHex("5a498a8bb518634b9454c1e2"),
		Title:       "Point Break",
		Description: "An FBI agent goes undercover to catch a gang of surfers who may be bank robbers.",
		Rating:      "R",
		Reviews: []models.Review{models.Review{
			Id:       bson.NewObjectId(),
			Stars:    4,
			Comments: "Point Break is a perfect example of the contemporary B movie.",
		}},
		ViewedByUsers: []bson.ObjectId{
			bson.ObjectIdHex("5a498a8bb518634b9454c1e3"),
			bson.ObjectIdHex("5a498a8bb518634b9454c1e4"),
		},
	})
}

func createUsers() {
	uc := client.DbCollection("users")
	ur := &repositories.UserRepository{C: uc}

	ur.Create(&models.User{
		Id:   bson.ObjectIdHex("5a498a8bb518634b9454c1e3"),
		Name: "Chris Fisher",
		Viewings: []models.Viewing{
			models.Viewing{
				Id:        bson.NewObjectId(),
				StartTime: time.Date(2017, time.December, 23, 15, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2017, time.December, 23, 17, 0, 0, 0, time.UTC),
				FilmId:    bson.ObjectIdHex("5a498a8bb518634b9454c1e1"),
			},
			models.Viewing{
				Id:        bson.NewObjectId(),
				StartTime: time.Date(2017, time.December, 24, 15, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2017, time.December, 24, 17, 0, 0, 0, time.UTC),
				FilmId:    bson.ObjectIdHex("5a498a8bb518634b9454c1e2"),
			},
		},
	})

	ur.Create(&models.User{
		Id:   bson.ObjectIdHex("5a498a8bb518634b9454c1e4"),
		Name: "Jacques Cousteau",
		Viewings: []models.Viewing{
			models.Viewing{
				Id:        bson.NewObjectId(),
				StartTime: time.Date(2017, time.December, 11, 11, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2017, time.December, 11, 13, 0, 0, 0, time.UTC),
				FilmId:    bson.ObjectIdHex("5a498a8bb518634b9454c1e1"),
			},
			models.Viewing{
				Id:        bson.NewObjectId(),
				StartTime: time.Date(2017, time.December, 12, 18, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2017, time.December, 12, 20, 0, 0, 0, time.UTC),
				FilmId:    bson.ObjectIdHex("5a498a8bb518634b9454c1e2"),
			},
		},
	})
}
