package resolvers

import (
	"context"

	"github.com/chrisfisher/jubilant-waffle/loaders"

	graphql "github.com/neelance/graphql-go"
)

type QueryResolver struct {
}

func (r *QueryResolver) Film(ctx context.Context, args struct{ ID graphql.ID }) *filmResolver {
	film := loaders.LoadFilmById(string(args.ID), ctx)
	return &filmResolver{film}
}

func (r *QueryResolver) SearchFilms(ctx context.Context, args struct{ Title string }) []*filmResolver {
	films := loaders.LoadFilmsByTitle(args.Title, ctx)
	var filmResolvers []*filmResolver
	for _, film := range films {
		filmResolvers = append(filmResolvers, &filmResolver{film})
	}
	return filmResolvers
}

func (r *QueryResolver) User(ctx context.Context, args struct{ ID graphql.ID }) *userResolver {
	user := loaders.LoadUserById(string(args.ID), ctx)
	return &userResolver{user}
}

func (r *QueryResolver) SearchUsers(ctx context.Context, args struct{ Name string }) []*userResolver {
	users := loaders.LoadUsersByName(args.Name, ctx)
	var userResolvers []*userResolver
	for _, user := range users {
		userResolvers = append(userResolvers, &userResolver{user})
	}
	return userResolvers
}
