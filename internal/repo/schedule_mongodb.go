package repo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/internal/model"
)

const scheduleCol = "schedules"

type scheduleMongoDBRepo struct {
	session *mgo.Session
}

func NewScheduleMongoDBRepo(session *mgo.Session) *scheduleMongoDBRepo {
	return &scheduleMongoDBRepo{
		session: session,
	}
}

func (r *scheduleMongoDBRepo) EnsureIndex() error {
	sess := r.session.Clone()
	defer sess.Close()

	err := sess.DB("").C(scheduleCol).EnsureIndex(mgo.Index{
		Key: []string{"reply_to"},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create schedule index")
	}

	return nil
}

func (r *scheduleMongoDBRepo) Create(sched *model.Schedule) (*model.Schedule, error) {
	sess := r.session.Clone()
	defer sess.Close()

	if err := sess.DB("").C(scheduleCol).Insert(sched); err != nil {
		return nil, errors.Wrapf(err, "failed to insert to mongodb")
	}
	return sched, nil
}

func (r *scheduleMongoDBRepo) Get(name, replyTo string) (*model.Schedule, error) {
	sess := r.session.Clone()
	defer sess.Close()

	var sched model.Schedule
	if err := sess.DB("").C(scheduleCol).
		Find(&bson.M{"_id": name, "reply_to": replyTo}).One(&sched); err != nil && err != mgo.ErrNotFound {
		return nil, errors.Wrapf(err, "failed to find schedule with name: %s", name)
	}
	return &sched, nil
}

func (r *scheduleMongoDBRepo) ListAll(replyTo string) ([]*model.Schedule, error) {
	sess := r.session.Clone()
	defer sess.Close()

	var scheds []*model.Schedule
	err := sess.DB("").C(scheduleCol).Find(&bson.M{"reply_to": replyTo}).All(&scheds)
	if err != nil && err != mgo.ErrNotFound {
		return nil, errors.Wrapf(err, "failed to list all scheds")
	}
	return scheds, nil
}

func (r *scheduleMongoDBRepo) Update(sched *model.Schedule) (*model.Schedule, error) {
	sess := r.session.Clone()
	defer sess.Close()

	if err := sess.DB("").C(scheduleCol).Update(&bson.M{"_id": sched.Name, "reply_to": sched.ReplyTo}, sched); err != nil {
		return nil, errors.Wrapf(err, "failed to update schedule: %v", sched)
	}
	return sched, nil
}

func (r *scheduleMongoDBRepo) Delete(name, replyTo string) error {
	sess := r.session.Clone()
	defer sess.Close()

	if err := sess.DB("").C(scheduleCol).Remove(&bson.M{"_id": name, "reply_to": replyTo}); err != nil {
		return errors.Wrapf(err, "failed to delete name %s", name)
	}
	return nil
}
