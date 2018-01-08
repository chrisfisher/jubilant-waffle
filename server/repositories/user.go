package repositories

import (
	"time"

	"github.com/chrisfisher/jubilant-waffle/server/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
	C *mgo.Collection
}

func (r *UserRepository) Create(user *models.User) error {
	user.CreatedOn = time.Now()
	err := r.C.Insert(&user)
	return err
}

func (r *UserRepository) Delete(id string) error {
	err := r.C.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}

func (r *UserRepository) GetById(id string) (user models.User, err error) {
	err = r.C.FindId(bson.ObjectIdHex(id)).One(&user)
	return
}

func (r *UserRepository) GetByIds(ids []string) []models.User {
	var users []models.User
	bsonIds := make([]bson.ObjectId, len(ids))
	for i, id := range ids {
		bsonIds[i] = bson.ObjectIdHex(id)
	}
	iter := r.C.Find(bson.M{"_id": bson.M{"$in": bsonIds}}).Iter()
	result := models.User{}
	for iter.Next(&result) {
		users = append(users, result)
	}
	return users
}

func (r *UserRepository) SearchByName(name string) []models.User {
	var users []models.User
	iter := r.C.Find(bson.M{"name": bson.M{"$regex": bson.RegEx{Pattern: name}}}).Iter()
	result := models.User{}
	for iter.Next(&result) {
		users = append(users, result)
	}
	return users
}
