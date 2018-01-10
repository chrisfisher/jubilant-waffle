package mappers

import (
	"github.com/chrisfisher/jubilant-waffle/server/models"
	"github.com/chrisfisher/jubilant-waffle/server/schema/types"

	"gopkg.in/mgo.v2/bson"

	graphql "github.com/neelance/graphql-go"
)

func MapFilm(film models.Film) *schema.Film {
	return &schema.Film{
		ID:            graphql.ID(film.Id.Hex()),
		Title:         film.Title,
		Description:   film.Description,
		Rating:        film.Rating,
		Reviews:       MapReviews(film.Reviews),
		ViewedByUsers: mapIds(film.ViewedByUsers),
	}
}

func mapIds(objectIds []bson.ObjectId) []graphql.ID {
	ids := make([]graphql.ID, len(objectIds))
	for i, objectId := range objectIds {
		ids[i] = graphql.ID(objectId.Hex())
	}
	return ids
}
