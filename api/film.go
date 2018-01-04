package api

import (
	"github.com/chrisfisher/jubilant-waffle/db"
	"github.com/chrisfisher/jubilant-waffle/models"
	"github.com/chrisfisher/jubilant-waffle/repositories"

	graphql "github.com/neelance/graphql-go"
)

type Film struct {
	ID          graphql.ID
	Title       string
	Description string
	Rating      string
	Reviews     []Review
}

type Review struct {
	ID       graphql.ID
	Stars    int32
	Comments string
}

type filmResolver struct {
	film      *Film
	dbContext *db.Context
}

type reviewResolver struct {
	review    *Review
	dbContext *db.Context
}

func (r *Resolver) Film(args struct{ ID graphql.ID }) *filmResolver {
	context := db.NewContext()
	return getFilmById(string(args.ID), context)
}

func (r *Resolver) SearchFilms(args struct{ Title string }) []*filmResolver {
	context := db.NewContext()
	return searchFilmsByTitle(args.Title, context)
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
		l[i] = &reviewResolver{&review, r.dbContext}
	}
	return &l
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

func getFilmById(id string, context *db.Context) *filmResolver {
	filmRepo := &repositories.FilmRepository{C: context.DbCollection("films")}
	dbFilm, err := filmRepo.GetById(id)
	if err != nil {
		return nil
	}
	film := Film{
		ID:          graphql.ID(dbFilm.Id.Hex()),
		Title:       dbFilm.Title,
		Description: dbFilm.Description,
		Rating:      dbFilm.Rating,
		Reviews:     mapReviews(dbFilm.Reviews),
	}
	return &filmResolver{&film, context}
}

func searchFilmsByTitle(text string, context *db.Context) []*filmResolver {
	repo := &repositories.FilmRepository{C: context.DbCollection("films")}
	dbFilms := repo.SearchByTitle(text)
	var filmResolvers []*filmResolver
	for _, dbFilm := range dbFilms {
		film := Film{
			ID:          graphql.ID(dbFilm.Id.Hex()),
			Title:       dbFilm.Title,
			Description: dbFilm.Description,
			Rating:      dbFilm.Rating,
			Reviews:     mapReviews(dbFilm.Reviews),
		}
		filmResolvers = append(filmResolvers, &filmResolver{&film, context})
	}
	return filmResolvers
}

func mapReviews(dbReviews []models.Review) []Review {
	reviews := make([]Review, len(dbReviews))
	for i, dbReview := range dbReviews {
		reviews[i] = Review{
			ID:       graphql.ID(dbReview.Id.Hex()),
			Stars:    dbReview.Stars,
			Comments: dbReview.Comments,
		}
	}
	return reviews
}
