package resolvers

import (
	"github.com/chrisfisher/jubilant-waffle/server/schema/types"

	graphql "github.com/neelance/graphql-go"
)

type reviewResolver struct {
	review *schema.Review
}

func (r *reviewResolver) ID() graphql.ID {
	return r.review.ID
}

func (r *reviewResolver) Stars() int32 {
	return r.review.Stars
}

func (r *reviewResolver) Comments() string {
	return r.review.Comments
}
