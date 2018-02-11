package db

import (
	"log"
	"os"

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
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		var err error
		session, err = mgo.Dial(host + ":27017")
		if err != nil {
			log.Fatalf("[getSession]: %s\n", err)
		}
	}
	return session
}
