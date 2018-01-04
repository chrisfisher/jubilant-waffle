package db

import (
	"log"

	"gopkg.in/mgo.v2"
)

var session *mgo.Session

var dbName = "cinema"

type Context struct {
	MongoSession *mgo.Session
}

func (c *Context) DbCollection(name string) *mgo.Collection {
	return c.MongoSession.DB(dbName).C(name)
}

func (c *Context) Close() {
	c.MongoSession.Close()
}

func NewContext() *Context {
	session := getSession().Copy()
	context := &Context{
		MongoSession: session,
	}
	return context
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
