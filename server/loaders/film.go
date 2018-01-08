package loaders

import (
	"context"

	"github.com/chrisfisher/jubilant-waffle/server/db"
	"github.com/chrisfisher/jubilant-waffle/server/models"
	"github.com/chrisfisher/jubilant-waffle/server/repositories"
	"github.com/chrisfisher/jubilant-waffle/server/schema/types"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/nicksrandall/dataloader.v4"

	graphql "github.com/neelance/graphql-go"
)

func LoadFilmById(id string, ctx context.Context) *schema.Film {
	ldr := GetFromContext(ctx, filmLoaderKey)
	if ldr == nil {
		return nil
	}
	data, err := ldr.Load(ctx, id)()
	if err != nil {
		return nil
	}
	dbFilm, ok := data.(models.Film)
	if !ok {
		return nil
	}
	film := schema.Film{
		ID:            graphql.ID(dbFilm.Id.Hex()),
		Title:         dbFilm.Title,
		Description:   dbFilm.Description,
		Rating:        dbFilm.Rating,
		Reviews:       mapReviews(dbFilm.Reviews),
		ViewedByUsers: mapIds(dbFilm.ViewedByUsers),
	}
	return &film
}

func LoadFilmsByTitle(title string, ctx context.Context) []*schema.Film {
	client := db.GetFromContext(ctx)
	ldr := newFilmLoader(client)
	dbFilms := ldr.r.SearchByTitle(title)
	var films []*schema.Film
	for _, dbFilm := range dbFilms {
		film := schema.Film{
			ID:            graphql.ID(dbFilm.Id.Hex()),
			Title:         dbFilm.Title,
			Description:   dbFilm.Description,
			Rating:        dbFilm.Rating,
			Reviews:       mapReviews(dbFilm.Reviews),
			ViewedByUsers: mapIds(dbFilm.ViewedByUsers),
		}
		films = append(films, &film)
	}
	return films
}

type filmLoader struct {
	r repositories.FilmRepository
}

func newFilmLoader(client *db.Client) filmLoader {
	c := client.DbCollection("films")
	r := repositories.FilmRepository{C: c}
	return filmLoader{r}
}

func (ldr filmLoader) loadBatch(ctx context.Context, keys []interface{}) []*dataloader.Result {
	ids := convert(keys)
	results := make([]*dataloader.Result, len(ids))
	if len(ids) == 1 {
		film, err := ldr.r.GetById(ids[0])
		results[0] = &dataloader.Result{Data: film, Error: err}
	} else {
		films := ldr.r.GetByIds(ids)
		for i, film := range films {
			results[i] = &dataloader.Result{Data: film, Error: nil}
		}
	}
	return results
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

func mapIds(objectIds []bson.ObjectId) []graphql.ID {
	ids := make([]graphql.ID, len(objectIds))
	for i, objectId := range objectIds {
		ids[i] = graphql.ID(objectId.Hex())
	}
	return ids
}
