package schema

import (
	graphql "github.com/neelance/graphql-go"
)

type Viewing struct {
	ID        graphql.ID
	StartTime graphql.Time
	EndTime   graphql.Time
	Film      graphql.ID
}
