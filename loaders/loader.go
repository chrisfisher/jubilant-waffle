package loaders

import (
	"context"
	"fmt"

	"github.com/chrisfisher/jubilant-waffle/db"

	"gopkg.in/nicksrandall/dataloader.v4"
)

type key string

const filmLoaderKey key = "film"
const userLoaderKey key = "user"

var loadersByKey map[key]*dataloader.Loader = make(map[key]*dataloader.Loader)

func (k key) String() string {
	return string(k)
}

func NewContext(ctx context.Context, client *db.Client) context.Context {
	loadersByKey[filmLoaderKey] = dataloader.NewBatchedLoader(newFilmLoader(client).loadBatch)
	loadersByKey[userLoaderKey] = dataloader.NewBatchedLoader(newUserLoader(client).loadBatch)

	for k, loader := range loadersByKey {
		ctx = context.WithValue(ctx, k, loader)
	}
	return ctx
}

func fromContext(ctx context.Context, k key) (*dataloader.Loader, error) {
	ldr, ok := ctx.Value(k).(*dataloader.Loader)
	if !ok {
		return nil, fmt.Errorf("cannot extract %s loader from context", k)
	}
	return ldr, nil
}

func convert(keys []interface{}) []string {
	var ids = make([]string, len(keys))
	for i, key := range keys {
		ids[i] = key.(string)
	}
	return ids
}
