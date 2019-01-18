package usermap

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/internal/repo"
)

type UserMap struct {
	um       map[string]string
	mux      *sync.RWMutex
	userRepo repo.UserRepository
}

var glb *UserMap

func New(userRepo repo.UserRepository) (*UserMap, error) {
	if glb != nil {
		return glb, nil
	}

	glb = &UserMap{
		um:       make(map[string]string),
		mux:      &sync.RWMutex{},
		userRepo: userRepo,
	}

	// Build in-mem cache from user data in database
	users, err := glb.userRepo.ListAll()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to init user map")
	}
	glb.mux.Lock()
	for _, user := range users {
		glb.um[user.ID] = user.Name
	}
	glb.mux.Unlock()

	return glb, nil
}

func (u *UserMap) Get(uid string) string {
	glb.mux.RLock()
	username := u.um[uid]
	glb.mux.RUnlock()
	return username
}

func (u *UserMap) Set(uid, name string) {
	glb.mux.Lock()
	u.um[uid] = name
	glb.mux.Unlock()
}
