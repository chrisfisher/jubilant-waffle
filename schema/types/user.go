package schema

import (
	graphql "github.com/neelance/graphql-go"
)

type User struct {
	ID       graphql.ID
	Name     string
	Viewings []Viewing
}
