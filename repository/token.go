package repository

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/project-tilas/svc-auth/domain"
)

type TokenRepository interface {
	Insert(token domain.Token) (*domain.Token, error)
	Remove(token string) error
	FindByUserIDAndToken(userID, token string) (*domain.Token, error)
}

type mongoTokenRepository struct {
	client     *mongoClient
	collection string
}

func NewMongoTokenRespository(client *mongoClient, collection string) (TokenRepository, error) {
	repo := &mongoTokenRepository{
		client:     client,
		collection: collection,
	}

	s := client.session.Copy()
	defer s.Close()
	coll := s.DB("").C(repo.collection)

	err := coll.EnsureIndex(mgo.Index{
		Key:        []string{"userId", "token"},
		Background: false,
	})
	if err != nil {
		return nil, err
	}

	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"token"},
		Unique:      true,
		ExpireAfter: 0,
		Background:  false,
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *mongoTokenRepository) Insert(t domain.Token) (*domain.Token, error) {

	s := repo.client.session.Copy()
	defer s.Close()
	coll := s.DB("").C(repo.collection)

	t.ID = bson.NewObjectId().String()
	t.CreatedAt = time.Now()
	err := coll.Insert(t)

	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (repo *mongoTokenRepository) Remove(token string) error {

	s := repo.client.session.Copy()
	defer s.Close()
	coll := s.DB("").C(repo.collection)

	return coll.Remove(bson.M{"token": token})
}

func (repo *mongoTokenRepository) FindByUserIDAndToken(userID, token string) (*domain.Token, error) {

	s := repo.client.session.Copy()
	defer s.Close()
	coll := s.DB("").C(repo.collection)

	var doc domain.Token
	err := coll.Find(bson.M{
		"userId": userID,
		"token":  token,
	}).One(&doc)

	if err == mgo.ErrNotFound {
		return nil, &domain.NotFoundError{Resource: "Token", ID: token}
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
