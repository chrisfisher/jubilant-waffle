package loaders

import (
	"context"

	"github.com/chrisfisher/jubilant-waffle/server/db"
	"github.com/chrisfisher/jubilant-waffle/server/mappers"
	"github.com/chrisfisher/jubilant-waffle/server/models"
	"github.com/chrisfisher/jubilant-waffle/server/repositories"
	"github.com/chrisfisher/jubilant-waffle/server/schema/types"

	"gopkg.in/nicksrandall/dataloader.v4"
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
	m, ok := data.(models.Film)
	if !ok {
		return nil
	}
	return mappers.MapFilm(m)
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
