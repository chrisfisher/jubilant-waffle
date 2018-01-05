package resolvers

import (
	"context"
	"github.com/chrisfisher/jubilant-waffle/db"
	"github.com/chrisfisher/jubilant-waffle/models"
	"github.com/chrisfisher/jubilant-waffle/repositories"
	"github.com/chrisfisher/jubilant-waffle/schema/types"

	graphql "github.com/neelance/graphql-go"
)

type userResolver struct {
	user *schema.User
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

func getUserById(id string, client *db.Client) *userResolver {
	repo := &repositories.UserRepository{C: client.DbCollection("users")}
	dbUser, err := repo.GetById(id)
	user := schema.User{
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
		user := schema.User{
			ID:       graphql.ID(dbUser.Id.Hex()),
			Name:     dbUser.Name,
			Viewings: mapViewings(dbUser.Viewings),
		}
		userResolvers = append(userResolvers, &userResolver{&user})
	}
	return userResolvers
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
