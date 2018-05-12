package repository

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/project-tilas/svc-auth/domain"
)

type UserRepository interface {
	Insert(u domain.User) (domain.User, error)
	Update(u domain.User) (domain.User, error)
	Save(u domain.User) (domain.User, error)
	FindByID(id string) (domain.User, error)
	FindByUsername(u string) (domain.User, error)
}

type mongoUserRepository struct {
	client     *mongoClient
	collection string
}

func NewMongoUserRespository(client *mongoClient, collection string) UserRepository {
	repo := &mongoUserRepository{
		client:     client,
		collection: collection,
	}

	s := repo.client.session.Copy()
	defer s.Close()
	coll := s.DB("").C(repo.collection)

	coll.EnsureIndex(mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		Background: false,
	})
	return repo
}

func (repo *mongoUserRepository) Insert(u domain.User) (domain.User, error) {

	s := repo.client.session.Copy()
	defer s.Close()
	coll := s.DB("").C(repo.collection)

	u.ID = bson.NewObjectId().String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	err := coll.Insert(u)

	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (repo *mongoUserRepository) Update(u domain.User) (domain.User, error) {
	selector := bson.M{"_id": u.ID}
	changes := bson.M{
		"username":  u.Username,
		"updatedAt": time.Now(),
	}
	if u.Password != "" {
		changes["password"] = u.Password
	}

	s := repo.client.session.Copy()
	defer s.Close()
	coll := s.DB("").C(repo.collection)

	err := coll.Update(selector, changes)

	if err != nil {
		return u, err
	}
	return u, nil
}

func (repo *mongoUserRepository) Save(u domain.User) (domain.User, error) {
	if u.ID == "" {
		return repo.Insert(u)
	}
	return repo.Update(u)
}

func (repo *mongoUserRepository) FindByID(id string) (domain.User, error) {

	s := repo.client.session.Copy()
	defer s.Close()
	coll := s.DB("").C(repo.collection)

	var doc domain.User
	err := coll.FindId(id).One(&doc)

	if err == mgo.ErrNotFound {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return doc, nil
}

func (repo *mongoUserRepository) FindByUsername(u string) (domain.User, error) {

	s := repo.client.session.Copy()
	defer s.Close()
	coll := s.DB("").C(repo.collection)

	var doc domain.User
	err := coll.Find(bson.M{"username": u}).One(&doc)

	if err == mgo.ErrNotFound {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return doc, nil
}
