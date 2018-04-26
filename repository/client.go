package repository

import (
	"time"

	"github.com/globalsign/mgo"
)

type MongoClientInfo struct {
	addr     string
	username string
	password string
	database string
}

type mongoClient struct {
	s *mgo.Session
	i *MongoClientInfo
}

func NewMongoClient(i *MongoClientInfo) (*mongoClient, error) {
	s, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{i.addr},
		Username: i.username,
		Password: i.password,
		Database: i.database,
		Timeout:  time.Second * 8,
	})
	if err != nil {
		return nil, err
	}
	return &mongoClient{
		s: s,
		i: i,
	}, nil
}

func (c *mongoClient) Close() error {
	if c.s != nil {
		c.s.Close()
		return nil
	}
	return nil
}

func (c *mongoClient) DB() (*mgo.Database, *mgo.Session) {
	s := c.s.Copy()
	return s.DB(c.i.database), s
}

func (c *mongoClient) C(col string) (*mgo.Collection, *mgo.Session) {
	db, s := c.DB()
	return db.C(col), s
}
