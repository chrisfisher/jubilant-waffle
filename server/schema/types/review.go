package schema

import (
	graphql "github.com/neelance/graphql-go"
)

type Review struct {
	ID       graphql.ID
	Stars    int32
	Comments string
}
