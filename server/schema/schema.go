package schema

var Schema = `
	schema {
		query: Query
		mutation: Mutation
	}
	# The query type which represents all entry points into our object graph
	type Query {
		film(id: ID!): Film
		searchFilms(title: String!): [Film]!
		user(id: ID!): User
		searchUsers(name: String!): [User]!
	}
	# The mutation type which represents all updates we can make to our data
	type Mutation {
		createFilm(Film: FilmInput!): Film
	}
	# Film ratings
	enum Rating {
		# General admission
		GA
		# Parental guidance recommended
		PG13
		# Mature audiences
		M
		# Restricted
		R
	}
	# A film
	type Film {
		# The film ID
		id: ID!
		# The title of the film
		title: String!
		# The film description
		description: String!
		# The film rating
		rating: Rating!
		# Any reviews of the film
		reviews: [Review]
		# Viewed by users
		viewedByUsers: [User]
	}
	# Input when creating a new film
	input FilmInput {
		# The title of the film
		title: String!
		# The film description
		description: String!
		# The film rating
		rating: Rating!
	}
	# A film review
	type Review {
		# The review ID
		id: ID!
		# The number of stars out of 5
		stars: Int!
		# Comments on the review
		comments: String!
	}
	# A user
	type User {
		# The member ID
		id: ID!
		# The user's name
		name: String!
		# All the occasions when a user has watched a film
		viewings: [Viewing]
		# The user's film viewings exposed as connections with pagination
		viewingConnection(first: Int, after: ID): ViewingConnection!
	}
	# An occasion when a user watches a film
	type Viewing {
		# The viewing ID
		id: ID!
		# The viewing start time
		startTime: Time!
		# The viewing end time
		endTime: Time!
		# The film being viewed
		film: Film!
	}
	# A connection object for a user's film viewings
	type ViewingConnection {
		# The total number of viewings
		totalCount: Int!
		# The edges for each of the user's viewings
		edges: [ViewingEdge]
		# Pagination info for this connection
		pageInfo: PageInfo!
	}
	# An edge object for a user's film viewings
	type ViewingEdge {
		# A cursor used for pagination
		cursor: ID!
		# The viewing represented by this edge
		node: Viewing
	}
	# Pagination info
	type PageInfo {
		startCursor: ID
		endCursor: ID
		hasNextPage: Boolean!
	}
	# A timestamp
	scalar Time
`
