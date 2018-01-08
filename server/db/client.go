package db

import (
	"log"

	"gopkg.in/mgo.v2"
)

type Client struct {
	MongoSession *mgo.Session
}

func (c *Client) DbCollection(name string) *mgo.Collection {
	return c.MongoSession.DB(dbName).C(name)
}

func (c *Client) Close() {
	c.MongoSession.Close()
}

var session *mgo.Session

func NewClient() *Client {
	session := getSession().Copy()
	client := &Client{
		MongoSession: session,
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
