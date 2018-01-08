package loaders

import (
	"context"

	"github.com/chrisfisher/jubilant-waffle/db"

	"gopkg.in/nicksrandall/dataloader.v4"
)

type key string

func (k key) String() string {
	return string(k)
}

const filmLoaderKey key = "film"
const userLoaderKey key = "user"

var loadersByKey map[key]*dataloader.Loader = make(map[key]*dataloader.Loader)

func AttachToContext(ctx context.Context, client *db.Client) context.Context {
	loadersByKey[filmLoaderKey] = dataloader.NewBatchedLoader(newFilmLoader(client).loadBatch)
	loadersByKey[userLoaderKey] = dataloader.NewBatchedLoader(newUserLoader(client).loadBatch)

	for k, loader := range loadersByKey {
		ctx = context.WithValue(ctx, k, loader)
	}
	return ctx
}

func GetFromContext(ctx context.Context, k key) *dataloader.Loader {
	ldr, ok := ctx.Value(k).(*dataloader.Loader)
	if !ok {
		return nil
	}
	return ldr
}

func convert(keys []interface{}) []string {
	var ids = make([]string, len(keys))
	for i, key := range keys {
		ids[i] = key.(string)
	}
	return ids
}
