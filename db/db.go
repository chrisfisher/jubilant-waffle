package db

import (
	"context"
	"log"

	"gopkg.in/mgo.v2"
)

var session *mgo.Session

const dbName = "cinema"

const ClientKey = "dbClient"

type Client struct {
	MongoSession *mgo.Session
}

func (c *Client) DbCollection(name string) *mgo.Collection {
	return c.MongoSession.DB(dbName).C(name)
}

func (c *Client) Close() {
	c.MongoSession.Close()
}

func NewClient() *Client {
	session := getSession().Copy()
	client := &Client{
		MongoSession: session,
	}
	return client
}

func FromContext(ctx context.Context) *Client {
	client, ok := ctx.Value(ClientKey).(*Client)
	if !ok {
		return nil
	}
	return client
}

func getSession() *mgo.Session {
	if session == nil {
		var err error
		session, err = mgo.Dial("127.0.0.1:27017")

		if err != nil {
			log.Fatalf("[getSession]: %s\n", err)
		}
	}
	return session
}
