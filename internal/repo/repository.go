package repo

import "github.com/lnquy/line-catalyst-server/internal/model"

type MessageRepository interface {
	EnsureIndex() error
	Create(message *model.Message, isUserMessage bool) (*model.Message, error)
	Get(mid string, isUserMessage bool) (*model.Message, error)
	ListLastMessages(id string, limit int, isUserMessage bool) ([]*model.Message, error)
	Delete(mid string, isUserMessage bool) error
}

type UserRepository interface {
	EnsureIndex() error
	Create(*model.User) (*model.User, error)
	Get(uid string) (*model.User, error)
	ListAll() ([]*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(uid string) error
}

type ScheduleRepository interface {
	EnsureIndex() error
	Create(sched *model.Schedule) (*model.Schedule, error)
	Get(name, replyTo string) (*model.Schedule, error)
	ListAll(replyTo string) ([]*model.Schedule, error)
	Update(sched *model.Schedule) (*model.Schedule, error)
	Delete(name, replyTo string) error
}
