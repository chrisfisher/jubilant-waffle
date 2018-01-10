package mappers

import (
	"github.com/chrisfisher/jubilant-waffle/server/models"
	"github.com/chrisfisher/jubilant-waffle/server/schema/types"

	graphql "github.com/neelance/graphql-go"
)

func MapReviews(reviews []models.Review) []schema.Review {
	results := make([]schema.Review, len(reviews))
	for i, review := range reviews {
		results[i] = MapReview(review)
	}
	return results
}

func MapReview(review models.Review) schema.Review {
	return schema.Review{
		ID:       graphql.ID(review.Id.Hex()),
		Stars:    review.Stars,
		Comments: review.Comments,
	}
}
