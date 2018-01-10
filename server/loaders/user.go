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

func LoadUserById(id string, ctx context.Context) *schema.User {
	ldr := GetFromContext(ctx, userLoaderKey)
	if ldr == nil {
		return nil
	}
	data, err := ldr.Load(ctx, id)()
	if err != nil {
		return nil
	}
	user, ok := data.(models.User)
	if !ok {
		return nil
	}
	return mappers.MapUser(user)
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
