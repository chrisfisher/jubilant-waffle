package cinema

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	graphql "github.com/neelance/graphql-go"
)

var Schema = `
	schema {
		query: Query
	}
	# The query type which represents all entry points into our object graph
	type Query {
		cinema(id: ID!): Cinema
		film(id: ID!): Film
		searchFilms(text: String!): [Film]!
		member(id: ID!): Member
		searchMembers(text: String!): [Member]!
	}
	# A cinema
	type Cinema {
		# The cinema id
		id: ID!
		# The name of the cinema
		name: String!
		# The cinema location latitude
		latitude: Int!
		# The cinema location longitude
		longitude: Int!
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
	# A session when a film is playing at a cinema at a given time
	type Session {
		# The session ID
		id: ID!
		# The session start time
		start: Time!
		# The cinema where the session is playing
		cinema: Cinema!
		# The film being played
		film: Film!
	}
	# A timestamp
	scalar Time
	# A member of a cinema loyalty program
	type Member {
		# The member ID
		id: ID!
		# The member's past sessions
		sessionHistory: [Session]
		# The member's past sessions exposed as connections with pagination
		sessionHistoryConnection(first: Int, after: ID): SessionHistoryConnection!
	}
	# A connection object for a member's past sessions
	type SessionHistoryConnection {
		# The total number of sessions
		totalCount: Int!
		# The edges for each of the member's past sessions
		edges: [SessionHistoryEdge]
		# Pagination info for this connection
		pageInfo: PageInfo!
	}
	# An edge object for a member's past sessions
	type SessionHistoryEdge {
		# A cursor used for pagination
		cursor: ID!
		# The session represented by this edge
		node: Session
	}
	# Pagination info
	type PageInfo {
		startCursor: ID
		endCursor: ID
		hasNextPage: Boolean!
	}
`

type cinema struct {
	ID        graphql.ID
	Name      string
	Latitude  int32
	Longitude int32
}

type film struct {
	ID          graphql.ID
	Title       string
	Description string
	Rating      string
	Reviews     []graphql.ID
}

type review struct {
	ID       graphql.ID
	Stars    int32
	Comments string
}

type session struct {
	ID     graphql.ID
	Start  graphql.Time
	Cinema graphql.ID
	Film   graphql.ID
}

type member struct {
	ID             graphql.ID
	Name           string
	SessionHistory []graphql.ID
}

var cinemas = []*cinema{
	{
		ID:        "1001",
		Name:      "ABC Queen St",
		Latitude:  0,
		Longitude: 0,
	},
	{
		ID:        "1002",
		Name:      "ABC Newmarket",
		Latitude:  0,
		Longitude: 0,
	},
}

var films = []*film{
	{
		ID:          "2001",
		Title:       "Captain America: The Winter Soldier",
		Description: "Steve Rogers teams up with Black Widow to battle an assassin known as the Winter Soldier.",
		Rating:      "PG13",
		Reviews:     []graphql.ID{"3001"},
	},
	{
		ID:          "2002",
		Title:       "Point Break",
		Description: "An FBI agent goes undercover to catch a gang of surfers who may be bank robbers.",
		Rating:      "R",
		Reviews:     []graphql.ID{"3002"},
	},
}

var reviews = []*review{
	{
		ID:       "3001",
		Stars:    4,
		Comments: "Superior to the first Captain America in every way and the best Marvel stand-alone movie.",
	},
	{
		ID:       "3002",
		Stars:    4,
		Comments: "Point Break is a perfect example of the contemporary B movie.",
	},
}

var sessions = []*session{
	{
		ID:     "4001",
		Start:  graphql.Time{Time: time.Date(2017, time.December, 23, 15, 0, 0, 0, time.UTC)},
		Cinema: "1001",
		Film:   "2001",
	},
	{
		ID:     "4002",
		Start:  graphql.Time{Time: time.Date(2017, time.December, 23, 17, 0, 0, 0, time.UTC)},
		Cinema: "1002",
		Film:   "2002",
	},
}

var members = []*member{
	{
		ID:             "5001",
		Name:           "Chris Fisher",
		SessionHistory: []graphql.ID{"4001", "4002"},
	},
	{
		ID:             "5002",
		Name:           "Jacques Cousteau",
		SessionHistory: []graphql.ID{"4001", "4002"},
	},
}

var cinemaData = make(map[graphql.ID]*cinema)

var filmData = make(map[graphql.ID]*film)

var reviewData = make(map[graphql.ID]*review)

var sessionData = make(map[graphql.ID]*session)

var memberData = make(map[graphql.ID]*member)

func init() {
	for _, f := range films {
		filmData[f.ID] = f
	}
	for _, r := range reviews {
		reviewData[r.ID] = r
	}
	for _, c := range cinemas {
		cinemaData[c.ID] = c
	}
	for _, s := range sessions {
		sessionData[s.ID] = s
	}
	for _, m := range members {
		memberData[m.ID] = m
	}
}

type Resolver struct{}

type searchResultResolver struct {
	result interface{}
}

type cinemaResolver struct {
	c *cinema
}

type filmResolver struct {
	f *film
}

type reviewResolver struct {
	r *review
}

type sessionResolver struct {
	s *session
}

type memberResolver struct {
	m *member
}

type sessionHistoryConnectionResolver struct {
	ids  []graphql.ID
	from int
	to   int
}

type sessionHistoryEdgeResolver struct {
	cursor graphql.ID
	id     graphql.ID
}

type pageInfoResolver struct {
	startCursor graphql.ID
	endCursor   graphql.ID
	hasNextPage bool
}

func (r *Resolver) Cinema(args struct{ ID graphql.ID }) *cinemaResolver {
	if c := cinemaData[args.ID]; c != nil {
		return &cinemaResolver{c}
	}
	return nil
}

func (r *Resolver) Film(args struct{ ID graphql.ID }) *filmResolver {
	if f := filmData[args.ID]; f != nil {
		return &filmResolver{f}
	}
	return nil
}

func (r *Resolver) SearchFilms(args struct{ Text string }) []*filmResolver {
	var l []*filmResolver
	for _, f := range films {
		if strings.Contains(f.Title, args.Text) {
			l = append(l, &filmResolver{f})
		}
	}
	return l
}

func (r *Resolver) Review(args struct{ ID graphql.ID }) *reviewResolver {
	if r := reviewData[args.ID]; r != nil {
		return &reviewResolver{r}
	}
	return nil
}

func (r *Resolver) Session(args struct{ ID graphql.ID }) *sessionResolver {
	if s := sessionData[args.ID]; s != nil {
		return &sessionResolver{s}
	}
	return nil
}

func (r *Resolver) Member(args struct{ ID graphql.ID }) *memberResolver {
	if m := memberData[args.ID]; m != nil {
		return &memberResolver{m}
	}
	return nil
}

func (r *Resolver) SearchMembers(args struct{ Text string }) []*memberResolver {
	var l []*memberResolver
	for _, m := range members {
		if strings.Contains(m.Name, args.Text) {
			l = append(l, &memberResolver{m})
		}
	}
	return l
}

func (r *cinemaResolver) ID() graphql.ID {
	return r.c.ID
}

func (r *cinemaResolver) Name() string {
	return r.c.Name
}

func (r *cinemaResolver) Latitude() int32 {
	return r.c.Latitude
}

func (r *cinemaResolver) Longitude() int32 {
	return r.c.Longitude
}

func (r *filmResolver) ID() graphql.ID {
	return r.f.ID
}

func (r *filmResolver) Title() string {
	return r.f.Title
}

func (r *filmResolver) Description() string {
	return r.f.Description
}

func (r *filmResolver) Rating() string {
	return r.f.Rating
}

func (r *filmResolver) Reviews() *[]*reviewResolver {
	l := make([]*reviewResolver, len(r.f.Reviews))
	for i, id := range r.f.Reviews {
		l[i] = &reviewResolver{reviewData[id]}
	}
	return &l
}

func (r *reviewResolver) ID() graphql.ID {
	return r.r.ID
}

func (r *reviewResolver) Stars() int32 {
	return r.r.Stars
}

func (r *reviewResolver) Comments() string {
	return r.r.Comments
}

func (r *sessionResolver) ID() graphql.ID {
	return r.s.ID
}

func (r *sessionResolver) Start() graphql.Time {
	return r.s.Start
}

func (r *sessionResolver) Film() *filmResolver {
	if f := filmData[r.s.Film]; f != nil {
		return &filmResolver{f}
	}
	return nil
}

func (r *sessionResolver) Cinema() *cinemaResolver {
	if c := cinemaData[r.s.Cinema]; c != nil {
		return &cinemaResolver{c}
	}
	return nil
}

func (r *memberResolver) ID() graphql.ID {
	return r.m.ID
}

func (r *memberResolver) Name() string {
	return r.m.Name
}

func (r *memberResolver) SessionHistory() *[]*sessionResolver {
	l := make([]*sessionResolver, len(r.m.SessionHistory))
	for i, id := range r.m.SessionHistory {
		l[i] = &sessionResolver{sessionData[id]}
	}
	return &l
}

type sessionHistoryConnectionArgs struct {
	First *int32
	After *graphql.ID
}

func (r *memberResolver) SessionHistoryConnection(args sessionHistoryConnectionArgs) (*sessionHistoryConnectionResolver, error) {
	return newSessionHistoryConnectionResolver(r.m.SessionHistory, args)
}

func newSessionHistoryConnectionResolver(ids []graphql.ID, args sessionHistoryConnectionArgs) (*sessionHistoryConnectionResolver, error) {
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

	to := len(ids)
	if args.First != nil {
		to = from + int(*args.First)
		if to > len(ids) {
			to = len(ids)
		}
	}

	return &sessionHistoryConnectionResolver{
		ids:  ids,
		from: from,
		to:   to,
	}, nil
}

func (r *sessionHistoryConnectionResolver) TotalCount() int32 {
	return int32(len(r.ids))
}

func (r *sessionHistoryConnectionResolver) Edges() *[]*sessionHistoryEdgeResolver {
	l := make([]*sessionHistoryEdgeResolver, r.to-r.from)
	for i := range l {
		l[i] = &sessionHistoryEdgeResolver{
			cursor: encodeCursor(r.from + i),
			id:     r.ids[r.from+i],
		}
	}
	return &l
}

func (r *sessionHistoryConnectionResolver) PageInfo() *pageInfoResolver {
	return &pageInfoResolver{
		startCursor: encodeCursor(r.from),
		endCursor:   encodeCursor(r.to - 1),
		hasNextPage: r.to < len(r.ids),
	}
}

func encodeCursor(i int) graphql.ID {
	return graphql.ID(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i+1))))
}

func (r *sessionHistoryEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *sessionHistoryEdgeResolver) Node() *sessionResolver {
	if s := sessionData[r.id]; s != nil {
		return &sessionResolver{s}
	}
	return nil
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
