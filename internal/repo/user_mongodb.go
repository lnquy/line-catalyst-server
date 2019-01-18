package repo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/internal/model"
)

const userCol = "users"

type userMongoDBRepo struct {
	session *mgo.Session
}

func NewUserMongoDBRepo(session *mgo.Session) UserRepository {
	return &userMongoDBRepo{
		session: session,
	}
}

func (r *userMongoDBRepo) Create(user *model.User) (*model.User, error) {
	sess := r.session.Clone()
	defer sess.Close()

	if err := sess.DB("").C(userCol).Insert(user); err != nil {
		return nil, errors.Wrapf(err, "failed to insert user to mongodb")
	}
	return user, nil
}

func (r *userMongoDBRepo) Get(uid string) (*model.User, error) {
	sess := r.session.Clone()
	defer sess.Close()

	var user model.User
	if err := sess.DB("").C(userCol).
		Find(&bson.M{"_id": uid}).One(&user); err != nil && err != mgo.ErrNotFound {
		return nil, errors.Wrapf(err, "failed to find user with id: %s", uid)
	}
	return &user, nil
}

func (r *userMongoDBRepo) ListAll() ([]*model.User, error) {
	sess := r.session.Clone()
	defer sess.Close()

	var users []*model.User
	err := sess.DB("").C(userCol).Find(&bson.M{}).All(&users)
	if err != nil && err != mgo.ErrNotFound {
		return nil, errors.Wrapf(err, "failed to list all users")
	}
	return users, nil
}

func (r *userMongoDBRepo) Update(user *model.User) (*model.User, error) {
	sess := r.session.Clone()
	defer sess.Close()

	if err := sess.DB("").C(userCol).Update(&bson.M{"_id": user.ID}, user); err != nil {
		return nil, errors.Wrapf(err, "failed to update user: %v", user)
	}
	return user, nil
}

func (r *userMongoDBRepo) Delete(uid string) error {
	sess := r.session.Clone()
	defer sess.Close()

	if err := sess.DB("").C(userCol).Remove(&bson.M{"_id": uid}); err != nil {
		return errors.Wrapf(err, "failed to delete user %s", uid)
	}
	return nil
}
