package schema

import (
	graphql "github.com/neelance/graphql-go"
)

type Film struct {
	ID            graphql.ID
	Title         string
	Description   string
	Rating        string
	Reviews       []Review
	ViewedByUsers []graphql.ID
}
