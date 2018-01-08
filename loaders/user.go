package loaders

import (
	"context"

	"github.com/chrisfisher/jubilant-waffle/db"
	"github.com/chrisfisher/jubilant-waffle/models"
	"github.com/chrisfisher/jubilant-waffle/repositories"
	"github.com/chrisfisher/jubilant-waffle/schema/types"

	"gopkg.in/nicksrandall/dataloader.v4"

	graphql "github.com/neelance/graphql-go"
)

func LoadUserById(id string, ctx context.Context) *schema.User {
	ldr := GetFromContext(ctx, userLoaderKey)
	if ldr == nil {
		return nil
	}
	data, err := ldr.Load(ctx, id)()
	if err != nil {
		return nil
	}
	dbUser, ok := data.(models.User)
	if !ok {
		return nil
	}
	user := schema.User{
		ID:       graphql.ID(dbUser.Id.Hex()),
		Name:     dbUser.Name,
		Viewings: mapViewings(dbUser.Viewings),
	}
	return &user
}

func LoadUsersByName(name string, ctx context.Context) []*schema.User {
	client := db.GetFromContext(ctx)
	ldr := newUserLoader(client)
	dbUsers := ldr.r.SearchByName(name)
	var users []*schema.User
	for _, dbUser := range dbUsers {
		user := schema.User{
			ID:       graphql.ID(dbUser.Id.Hex()),
			Name:     dbUser.Name,
			Viewings: mapViewings(dbUser.Viewings),
		}
		users = append(users, &user)
	}
	return users
}

type userLoader struct {
	r repositories.UserRepository
}

func newUserLoader(client *db.Client) userLoader {
	c := client.DbCollection("users")
	r := repositories.UserRepository{C: c}
	return userLoader{r}
}

func (ldr userLoader) loadBatch(ctx context.Context, keys []interface{}) []*dataloader.Result {
	ids := convert(keys)
	results := make([]*dataloader.Result, len(ids))
	if len(ids) == 1 {
		user, err := ldr.r.GetById(ids[0])
		results[0] = &dataloader.Result{Data: user, Error: err}
	} else {
		users := ldr.r.GetByIds(ids)
		for i, user := range users {
			results[i] = &dataloader.Result{Data: user, Error: nil}
		}
	}
	return results
}

func mapViewings(dbViewings []models.Viewing) []schema.Viewing {
	viewings := make([]schema.Viewing, len(dbViewings))
	for i, dbViewing := range dbViewings {
		viewings[i] = schema.Viewing{
			ID:        graphql.ID(dbViewing.Id.Hex()),
			StartTime: graphql.Time{Time: dbViewing.StartTime},
			EndTime:   graphql.Time{Time: dbViewing.EndTime},
			Film:      graphql.ID(dbViewing.FilmId.Hex()),
		}
	}
	return viewings
}
