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
	Create(*model.User) (*model.User, error)
	Get(uid string) (*model.User, error)
	ListAll() ([]*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(uid string) error
}
