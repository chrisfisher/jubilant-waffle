package api

import (
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
	user      *User
	dbContext *db.Context
}

type viewingResolver struct {
	viewing   *Viewing
	dbContext *db.Context
}

type viewingConnectionResolver struct {
	viewings  []Viewing
	from      int
	to        int
	dbContext *db.Context
}

type viewingConnectionArgs struct {
	First *int32
	After *graphql.ID
}

type viewingEdgeResolver struct {
	cursor    graphql.ID
	viewing   Viewing
	dbContext *db.Context
}

type pageInfoResolver struct {
	startCursor graphql.ID
	endCursor   graphql.ID
	hasNextPage bool
}

func (r *Resolver) User(args struct{ ID graphql.ID }) *userResolver {
	context := db.NewContext()
	return getUserById(string(args.ID), context)
}

func (r *Resolver) SearchUsers(args struct{ Name string }) []*userResolver {
	context := db.NewContext()
	return searchUsersByName(args.Name, context)
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
		viewingResolvers[i] = &viewingResolver{&viewing, r.dbContext}
	}
	return &viewingResolvers
}

func (r *userResolver) ViewingConnection(args viewingConnectionArgs) (*viewingConnectionResolver, error) {
	return newViewingConnectionResolver(r.user.Viewings, args, r.dbContext)
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

func (r *viewingResolver) Film() *filmResolver {
	return getFilmById(string(r.viewing.Film), r.dbContext)
}

func (r *viewingConnectionResolver) TotalCount() int32 {
	return int32(len(r.viewings))
}

func (r *viewingConnectionResolver) Edges() *[]*viewingEdgeResolver {
	l := make([]*viewingEdgeResolver, r.to-r.from)
	for i := range l {
		l[i] = &viewingEdgeResolver{
			cursor:    encodeCursor(r.from + i),
			viewing:   r.viewings[r.from+i],
			dbContext: r.dbContext,
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
	return &viewingResolver{&r.viewing, r.dbContext}
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

func getUserById(id string, context *db.Context) *userResolver {
	repo := &repositories.UserRepository{C: context.DbCollection("users")}
	dbUser, err := repo.GetById(id)
	user := User{
		ID:       graphql.ID(dbUser.Id.Hex()),
		Name:     dbUser.Name,
		Viewings: mapViewings(dbUser.Viewings),
	}
	if err == nil {
		return &userResolver{&user, context}
	}
	return nil
}

func searchUsersByName(name string, context *db.Context) []*userResolver {
	repo := &repositories.UserRepository{C: context.DbCollection("users")}
	dbUsers := repo.SearchByName(name)
	var userResolvers []*userResolver
	for _, dbUser := range dbUsers {
		user := User{
			ID:       graphql.ID(dbUser.Id.Hex()),
			Name:     dbUser.Name,
			Viewings: mapViewings(dbUser.Viewings),
		}
		userResolvers = append(userResolvers, &userResolver{&user, context})
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

func newViewingConnectionResolver(viewings []Viewing, args viewingConnectionArgs, context *db.Context) (*viewingConnectionResolver, error) {
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
		viewings:  viewings,
		from:      from,
		to:        to,
		dbContext: context,
	}, nil
}
