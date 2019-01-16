package repo

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/internal/config"
	"github.com/lnquy/line-catalyst-server/internal/model"
)

const (
	userMessageCol  = "user_messages"
	groupMessageCol = "group_messages"

	expirationDuration = time.Duration(15*24)*time.Hour // 15 days
)

type messageMongoDBRepo struct {
	session *mgo.Session
}

func NewMessageMongoDBRepo(conf config.MongoDB) (MessageRepository, error) {
	session, err := mgo.DialWithTimeout(conf.URI, 30*time.Second)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial mongodb")
	}
	return &messageMongoDBRepo{
		session: session,
	}, nil
}

func (r *messageMongoDBRepo) EnsureIndex() error {
	sess := r.session.Clone()
	defer sess.Close()

	err := sess.DB("").C(userMessageCol).EnsureIndex(mgo.Index{
		ExpireAfter: 1 * time.Second,
		Key:         []string{"expiration_date"},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create %s TTL index", userMessageCol)
	}

	err = sess.DB("").C(groupMessageCol).EnsureIndex(mgo.Index{
		ExpireAfter: 1 * time.Second,
		Key:         []string{"expiration_date"},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create %s TTL index", groupMessageCol)
	}

	return nil
}

func (r *messageMongoDBRepo) Create(message *model.Message, isUserMessage bool) (*model.Message, error) {
	sess := r.session.Clone()
	defer sess.Close()

	col := getCollection(isUserMessage)
	message.ID = uuid.New().String()
	message.ExpirationDate = time.Now().Add(expirationDuration)
	if err := sess.DB("").C(col).Insert(message); err != nil {
		return nil, errors.Wrapf(err, "failed to insert message to mongodb at %s collection", col)
	}
	return message, nil
	// Should get from db after inserted but omitted here to prevent db overload.
	// return r.Get(message.ID, isUserMessage)
}

func (r *messageMongoDBRepo) Get(mid string, isUserMessage bool) (*model.Message, error) {
	sess := r.session.Clone()
	defer sess.Close()

	var message model.Message
	if err := sess.DB("").C(getCollection(isUserMessage)).
		Find(&bson.M{"_id": mid}).One(&message); err != nil && err != mgo.ErrNotFound {
		return nil, errors.Wrapf(err, "failed to find message with id: %s", mid)
	}
	return &message, nil
}

func (r *messageMongoDBRepo) List(uid string, limit int, isUserMessage bool) ([]*model.Message, error) {
	sess := r.session.Clone()
	defer sess.Close()

	col := getCollection(isUserMessage)
	idKey := "group_id"
	if isUserMessage {
		idKey = "user_id"
	}
	var messages []*model.Message
	if err := sess.DB("").C(col).Find(&bson.M{idKey: uid}).Sort("timestamp").
		Limit(limit).All(&messages); err != nil && err != mgo.ErrNotFound {
		return nil, errors.Wrapf(err, "failed to list messages from %s collection", col)
	}
	return messages, nil
}

func (r *messageMongoDBRepo) Delete(mid string, isUserMessage bool) error {
	sess := r.session.Clone()
	defer sess.Close()

	col := getCollection(isUserMessage)
	if err := sess.DB("").C(col).Remove(&bson.M{"_id": mid}); err != nil {
		return errors.Wrapf(err, "failed to delete document %s from %s collection", mid, col)
	}
	return nil
}

func getCollection(isUserMessage bool) string {
	if isUserMessage {
		return userMessageCol
	}
	return groupMessageCol
}
