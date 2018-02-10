package resolvers

import (
	"github.com/chrisfisher/jubilant-waffle/server/schema/types"

	graphql "github.com/neelance/graphql-go"
)

type userResolver struct {
	user *schema.User
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
		viewingResolvers[i] = &viewingResolver{viewing}
	}
	return &viewingResolvers
}

func (r *userResolver) ViewingConnection(args viewingConnectionArgs) (*viewingConnectionResolver, error) {
	return newViewingConnectionResolver(r.user.Viewings, args)
}
