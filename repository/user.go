package repository

import (
	"time"

	"github.com/labstack/gommon/log"

	"golang.org/x/crypto/bcrypt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/project-tilas/svc-auth/domain"
)

type UserRepository interface {
	Insert(u domain.User) (domain.User, error)
	Update(u domain.User) (domain.User, error)
	Save(u domain.User) (domain.User, error)
	FindByUsername(u string) (domain.User, error)
}

type mongoUserRepository struct {
	client *mongoClient
}

func NewMongoUserRespository(m *mongoClient) UserRepository {
	repo := &mongoUserRepository{
		client: m,
	}
	coll, s := repo.collection()
	defer s.Close()
	coll.EnsureIndex(mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		Background: false,
	})
	return repo
}

func (repo *mongoUserRepository) collection() (*mgo.Collection, *mgo.Session) {
	return repo.client.C("user")
}

func (repo *mongoUserRepository) Insert(u domain.User) (domain.User, error) {
	u.ID = bson.NewObjectId().String()

	coll, s := repo.collection()
	defer s.Close()
	err := coll.Insert(bson.M{
		"_id":       u.ID,
		"username":  u.Username,
		"password":  u.Password,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	})

	if err != nil {
		return u, err
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
		changes["password"] = encrypt(u.Password)
	}

	coll, s := repo.collection()
	defer s.Close()
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

func (repo *mongoUserRepository) FindByUsername(u string) (domain.User, error) {
	coll, s := repo.collection()
	defer s.Close()

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

func encrypt(s string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
	}
	return string(b)
}
