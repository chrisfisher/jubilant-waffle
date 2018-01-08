package resolvers

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/chrisfisher/jubilant-waffle/loaders"
	"github.com/chrisfisher/jubilant-waffle/schema/types"

	graphql "github.com/neelance/graphql-go"
)

type viewingResolver struct {
	viewing *schema.Viewing
}

type viewingConnectionResolver struct {
	viewings []schema.Viewing
	from     int
	to       int
}

type viewingConnectionArgs struct {
	First *int32
	After *graphql.ID
}

type viewingEdgeResolver struct {
	cursor  graphql.ID
	viewing schema.Viewing
}

type pageInfoResolver struct {
	startCursor graphql.ID
	endCursor   graphql.ID
	hasNextPage bool
}

func (r *viewingResolver) ID() graphql.ID {
	return r.viewing.ID
}

func (r *viewingResolver) StartTime() graphql.Time {
	return r.viewing.StartTime
}

func (r *viewingResolver) EndTime() graphql.Time {
	return r.viewing.EndTime
}

func (r *viewingResolver) Film(ctx context.Context) *filmResolver {
	film := loaders.LoadFilmById(string(r.viewing.Film), ctx)
	return &filmResolver{film}
}

func (r *viewingConnectionResolver) TotalCount() int32 {
	return int32(len(r.viewings))
}

func (r *viewingConnectionResolver) Edges() *[]*viewingEdgeResolver {
	l := make([]*viewingEdgeResolver, r.to-r.from)
	for i := range l {
		l[i] = &viewingEdgeResolver{
			cursor:  encodeCursor(r.from + i),
			viewing: r.viewings[r.from+i],
		}
	}
	return &l
}

func (r *viewingConnectionResolver) PageInfo() *pageInfoResolver {
	return &pageInfoResolver{
		startCursor: encodeCursor(r.from),
		endCursor:   encodeCursor(r.to - 1),
		hasNextPage: r.to < len(r.viewings),
	}
}

func encodeCursor(i int) graphql.ID {
	return graphql.ID(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i+1))))
}

func (r *viewingEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *viewingEdgeResolver) Node() *viewingResolver {
	return &viewingResolver{&r.viewing}
}

func (r *pageInfoResolver) StartCursor() *graphql.ID {
	return &r.startCursor
}

func (r *pageInfoResolver) EndCursor() *graphql.ID {
	return &r.endCursor
}

func (r *pageInfoResolver) HasNextPage() bool {
	return r.hasNextPage
}

func newViewingConnectionResolver(viewings []schema.Viewing, args viewingConnectionArgs) (*viewingConnectionResolver, error) {
	from := 0
	if args.After != nil {
		b, err := base64.StdEncoding.DecodeString(string(*args.After))
		if err != nil {
			return nil, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return nil, err
		}
		from = i
	}

	to := len(viewings)
	if args.First != nil {
		to = from + int(*args.First)
		if to > len(viewings) {
			to = len(viewings)
		}
	}

	return &viewingConnectionResolver{
		viewings: viewings,
		from:     from,
		to:       to,
	}, nil
}
