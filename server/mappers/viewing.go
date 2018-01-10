package mappers

import (
	"github.com/chrisfisher/jubilant-waffle/server/models"
	"github.com/chrisfisher/jubilant-waffle/server/schema/types"

	graphql "github.com/neelance/graphql-go"
)

func MapViewings(viewings []models.Viewing) []schema.Viewing {
	results := make([]schema.Viewing, len(viewings))
	for i, viewing := range viewings {
		results[i] = MapViewing(viewing)
	}
	return results
}

func MapViewing(viewing models.Viewing) schema.Viewing {
	return schema.Viewing{
		ID:        graphql.ID(viewing.Id.Hex()),
		StartTime: graphql.Time{Time: viewing.StartTime},
		EndTime:   graphql.Time{Time: viewing.EndTime},
		Film:      graphql.ID(viewing.FilmId.Hex()),
	}
}
