package db

import (
	"context"
)

type key string

func (k key) String() string {
	return string(k)
}

const clientKey key = "client"

const dbName = "cinema"

func AttachToContext(ctx context.Context, client *Client) context.Context {
	ctx = context.WithValue(ctx, clientKey, client)
	return ctx
}

func GetFromContext(ctx context.Context) *Client {
	client, ok := ctx.Value(clientKey).(*Client)
	if !ok {
		return nil
	}
	return client
}
