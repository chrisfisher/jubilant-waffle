package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	Film struct {
		Id            bson.ObjectId   `bson:"_id,omitempty" json:"id"`
		Title         string          `json:"title"`
		Description   string          `json:"description"`
		Rating        string          `json:"rating"`
		CreatedOn     time.Time       `json:"created_on,omitempty"`
		Reviews       []Review        `json:"review"`
		ViewedByUsers []bson.ObjectId `json:"watched_by"`
	}

	Review struct {
		Id       bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Stars    int32         `json:"stars"`
		Comments string        `json:"comments"`
	}

	User struct {
		Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Name      string        `json:"name"`
		CreatedOn time.Time     `json:"created_on,omitempty"`
		Viewings  []Viewing     `json:"viewings"`
	}

	Viewing struct {
		Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
		StartTime time.Time     `json:"start_time"`
		EndTime   time.Time     `json:"end_time"`
		FilmId    bson.ObjectId `json:"film_id"`
	}
)
