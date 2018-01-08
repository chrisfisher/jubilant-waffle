package resolvers

import (
	"context"

	"github.com/chrisfisher/jubilant-waffle/loaders"
	"github.com/chrisfisher/jubilant-waffle/schema/types"

	graphql "github.com/neelance/graphql-go"
)

type filmResolver struct {
	film *schema.Film
}

func (r *filmResolver) ID() graphql.ID {
	return r.film.ID
}

func (r *filmResolver) Title() string {
	return r.film.Title
}

func (r *filmResolver) Description() string {
	return r.film.Description
}

func (r *filmResolver) Rating() string {
	return r.film.Rating
}

func (r *filmResolver) Reviews() *[]*reviewResolver {
	l := make([]*reviewResolver, len(r.film.Reviews))
	for i, review := range r.film.Reviews {
		l[i] = &reviewResolver{&review}
	}
	return &l
}

func (r *filmResolver) ViewedByUsers(ctx context.Context) *[]*userResolver {
	l := make([]*userResolver, len(r.film.ViewedByUsers))
	for i, id := range r.film.ViewedByUsers {
		user := loaders.LoadUserById(string(id), ctx)
		l[i] = &userResolver{user}
	}
	return &l
}
