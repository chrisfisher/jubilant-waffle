package resolvers

import (
	"context"

	"github.com/chrisfisher/jubilant-waffle/server/db"
	"github.com/chrisfisher/jubilant-waffle/server/loaders"
	"github.com/chrisfisher/jubilant-waffle/server/mappers"
	"github.com/chrisfisher/jubilant-waffle/server/models"
	"github.com/chrisfisher/jubilant-waffle/server/repositories"

	graphql "github.com/neelance/graphql-go"
)

type SchemaResolver struct {
}

func (r *SchemaResolver) Film(ctx context.Context, args struct{ ID graphql.ID }) *filmResolver {
	film := loaders.LoadFilmById(string(args.ID), ctx)
	return &filmResolver{film}
}

func (r *SchemaResolver) SearchFilms(ctx context.Context, args struct{ Title string }) []*filmResolver {
	client := db.GetFromContext(ctx)
	fr := repositories.NewFilmRepository(client)
	films := fr.SearchByTitle(args.Title)
	var filmResolvers []*filmResolver
	for _, f := range films {
		filmResolvers = append(filmResolvers, &filmResolver{mappers.MapFilm(f)})
	}
	return filmResolvers
}

func (r *SchemaResolver) User(ctx context.Context, args struct{ ID graphql.ID }) *userResolver {
	user := loaders.LoadUserById(string(args.ID), ctx)
	return &userResolver{user}
}

func (r *SchemaResolver) SearchUsers(ctx context.Context, args struct{ Name string }) []*userResolver {
	client := db.GetFromContext(ctx)
	ur := repositories.NewUserRepository(client)
	users := ur.SearchByName(args.Name)
	var userResolvers []*userResolver
	for _, u := range users {
		userResolvers = append(userResolvers, &userResolver{mappers.MapUser(u)})
	}
	return userResolvers
}

func (r *SchemaResolver) CreateFilm(ctx context.Context, args *struct{ Film *filmInput }) *filmResolver {
	client := db.GetFromContext(ctx)
	fr := repositories.NewFilmRepository(client)
	film, err := fr.Create(&models.Film{
		Title:       args.Film.Title,
		Description: args.Film.Description,
		Rating:      args.Film.Rating,
	})
	if err != nil {
		return nil
	}
	return &filmResolver{mappers.MapFilm(*film)}
}

type filmInput struct {
	Title       string
	Description string
	Rating      string
}
