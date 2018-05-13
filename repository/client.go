package repository

import (
	"github.com/globalsign/mgo"
)

type MongoClientInfo struct {
	Addr     string
	Username string
	Password string
	Database string
}

type mongoClient struct {
	session *mgo.Session
	info    *mgo.DialInfo
}

func NewMongoClient(info *mgo.DialInfo) (*mongoClient, error) {
	s, err := mgo.DialWithInfo(info)
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
