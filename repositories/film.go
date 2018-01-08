package repositories

import (
	"time"

	"github.com/chrisfisher/jubilant-waffle/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type FilmRepository struct {
	C *mgo.Collection
}

func (r *FilmRepository) Create(film *models.Film) error {
	film.CreatedOn = time.Now()
	err := r.C.Insert(&film)
	return err
}

func (r *FilmRepository) GetById(id string) (film models.Film, err error) {
	err = r.C.FindId(bson.ObjectIdHex(id)).One(&film)
	return
}

func (r *FilmRepository) GetByIds(ids []string) []models.Film {
	var films []models.Film
	bsonIds := make([]bson.ObjectId, len(ids))
	for i, id := range ids {
		bsonIds[i] = bson.ObjectIdHex(id)
	}
	iter := r.C.Find(bson.M{"_id": bson.M{"$in": bsonIds}}).Iter()
	result := models.Film{}
	for iter.Next(&result) {
		films = append(films, result)
	}
	return films
}

func (r *FilmRepository) SearchByTitle(title string) []models.Film {
	var films []models.Film
	iter := r.C.Find(bson.M{"title": bson.M{"$regex": bson.RegEx{Pattern: title}}}).Iter()
	result := models.Film{}
	for iter.Next(&result) {
		films = append(films, result)
	}
	return films
}

func (r *FilmRepository) Delete(id string) error {
	err := r.C.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}
