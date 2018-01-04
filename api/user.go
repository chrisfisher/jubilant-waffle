package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/chrisfisher/jubilant-waffle/db"
	"github.com/chrisfisher/jubilant-waffle/models"
	"github.com/chrisfisher/jubilant-waffle/repositories"
	"strconv"
	"strings"

	graphql "github.com/neelance/graphql-go"
)

type User struct {
	ID       graphql.ID
	Name     string
	Viewings []Viewing
}

type Viewing struct {
	ID        graphql.ID
	StartTime graphql.Time
	EndTime   graphql.Time
	Film      graphql.ID
}

type userResolver struct {
	user *User
}

type viewingResolver struct {
	viewing *Viewing
}

type viewingConnectionResolver struct {
	viewings []Viewing
	from     int
	to       int
}

type viewingConnectionArgs struct {
	First *int32
	After *graphql.ID
}

type viewingEdgeResolver struct {
	cursor  graphql.ID
	viewing Viewing
}

type pageInfoResolver struct {
	startCursor graphql.ID
	endCursor   graphql.ID
	hasNextPage bool
}

func (r *Resolver) User(ctx context.Context, args struct{ ID graphql.ID }) *userResolver {
	client := db.FromContext(ctx)
	return getUserById(string(args.ID), client)
}

func (r *Resolver) SearchUsers(ctx context.Context, args struct{ Name string }) []*userResolver {
	client := db.FromContext(ctx)
	return searchUsersByName(args.Name, client)
}

func (r *userResolver) ID() graphql.ID {
	return r.user.ID
}

func (r *userResolver) Name() string {
	return r.user.Name
}

func (r *userResolver) Viewings() *[]*viewingResolver {
	viewingResolvers := make([]*viewingResolver, len(r.user.Viewings))
	for i, viewing := range r.user.Viewings {
		viewingResolvers[i] = &viewingResolver{&viewing}
	}
	return &viewingResolvers
}

func (r *userResolver) ViewingConnection(args viewingConnectionArgs) (*viewingConnectionResolver, error) {
	return newViewingConnectionResolver(r.user.Viewings, args)
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
	client := db.FromContext(ctx)
	return getFilmById(string(r.viewing.Film), client)
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

func getUserById(id string, client *db.Client) *userResolver {
	repo := &repositories.UserRepository{C: client.DbCollection("users")}
	dbUser, err := repo.GetById(id)
	user := User{
		ID:       graphql.ID(dbUser.Id.Hex()),
		Name:     dbUser.Name,
		Viewings: mapViewings(dbUser.Viewings),
	}
	if err == nil {
		return &userResolver{&user}
	}
	return nil
}

func searchUsersByName(name string, client *db.Client) []*userResolver {
	repo := &repositories.UserRepository{C: client.DbCollection("users")}
	dbUsers := repo.SearchByName(name)
	var userResolvers []*userResolver
	for _, dbUser := range dbUsers {
		user := User{
			ID:       graphql.ID(dbUser.Id.Hex()),
			Name:     dbUser.Name,
			Viewings: mapViewings(dbUser.Viewings),
		}
		userResolvers = append(userResolvers, &userResolver{&user})
	}
	return userResolvers
}

func mapViewings(dbViewings []models.Viewing) []Viewing {
	viewings := make([]Viewing, len(dbViewings))
	for i, dbViewing := range dbViewings {
		viewings[i] = Viewing{
			ID:        graphql.ID(dbViewing.Id.Hex()),
			StartTime: graphql.Time{Time: dbViewing.StartTime},
			EndTime:   graphql.Time{Time: dbViewing.EndTime},
			Film:      graphql.ID(dbViewing.FilmId.Hex()),
		}
	}
	return viewings
}

func newViewingConnectionResolver(viewings []Viewing, args viewingConnectionArgs) (*viewingConnectionResolver, error) {
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
