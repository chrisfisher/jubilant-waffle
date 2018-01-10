package mappers

import (
	"github.com/chrisfisher/jubilant-waffle/server/models"
	"github.com/chrisfisher/jubilant-waffle/server/schema/types"

	graphql "github.com/neelance/graphql-go"
)

func MapUser(user models.User) *schema.User {
	return &schema.User{
		ID:       graphql.ID(user.Id.Hex()),
		Name:     user.Name,
		Viewings: MapViewings(user.Viewings),
	}
}
