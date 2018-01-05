package resolvers

import (
	"context"
	"github.com/chrisfisher/jubilant-waffle/db"
	"github.com/chrisfisher/jubilant-waffle/models"
	"github.com/chrisfisher/jubilant-waffle/repositories"
	"github.com/chrisfisher/jubilant-waffle/schema/types"

	graphql "github.com/neelance/graphql-go"
)

type filmResolver struct {
	film *schema.Film
}

func (r *Resolver) Film(ctx context.Context, args struct{ ID graphql.ID }) *filmResolver {
	client := db.FromContext(ctx)
	return getFilmById(string(args.ID), client)
}

func (r *Resolver) SearchFilms(ctx context.Context, args struct{ Title string }) []*filmResolver {
	client := db.FromContext(ctx)
	return searchFilmsByTitle(args.Title, client)
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

func getFilmById(id string, client *db.Client) *filmResolver {
	filmRepo := &repositories.FilmRepository{C: client.DbCollection("films")}
	dbFilm, err := filmRepo.GetById(id)
	if err != nil {
		return nil
	}
	film := schema.Film{
		ID:          graphql.ID(dbFilm.Id.Hex()),
		Title:       dbFilm.Title,
		Description: dbFilm.Description,
		Rating:      dbFilm.Rating,
		Reviews:     mapReviews(dbFilm.Reviews),
	}
	return &filmResolver{&film}
}

func searchFilmsByTitle(title string, client *db.Client) []*filmResolver {
	repo := &repositories.FilmRepository{C: client.DbCollection("films")}
	dbFilms := repo.SearchByTitle(title)
	var filmResolvers []*filmResolver
	for _, dbFilm := range dbFilms {
		film := schema.Film{
			ID:          graphql.ID(dbFilm.Id.Hex()),
			Title:       dbFilm.Title,
			Description: dbFilm.Description,
			Rating:      dbFilm.Rating,
			Reviews:     mapReviews(dbFilm.Reviews),
		}
		filmResolvers = append(filmResolvers, &filmResolver{&film})
	}
	return filmResolvers
}

func mapReviews(dbReviews []models.Review) []schema.Review {
	reviews := make([]schema.Review, len(dbReviews))
	for i, dbReview := range dbReviews {
		reviews[i] = schema.Review{
			ID:       graphql.ID(dbReview.Id.Hex()),
			Stars:    dbReview.Stars,
			Comments: dbReview.Comments,
		}
	}
	return reviews
}
