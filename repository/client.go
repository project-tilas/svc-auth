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
	session *mgo.Session
	info    *MongoClientInfo
}

func NewMongoClient(info *MongoClientInfo) (*mongoClient, error) {
	s, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{info.addr},
		Username: info.username,
		Password: info.password,
		Database: info.database,
		Timeout:  time.Second * 8,
	})
	if err != nil {
		return nil, err
	}
	return &mongoClient{
		session: s,
		info:    info,
	}, nil
}

func (c *mongoClient) Close() error {
	if c.session != nil {
		c.session.Close()
		return nil
	}
	return nil
}
