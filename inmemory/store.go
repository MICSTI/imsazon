/*
	This package implements the in-memory repository stores for all important data
 */

package inmemory

import (
	"sync"
	"github.com/MICSTI/imsazon/user"
)

type userRepository struct {
	mtx		sync.RWMutex
	users	map[user.UserId]*user.User
}

func (r *userRepository) Add(u *user.User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.users[u.Id] = u
	return nil
}

func (r *userRepository) Find(id user.UserId) (*user.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.users[id]; ok {
		return val, nil
	}
	return nil, user.ErrUnknown
}

func (r *userRepository) FindAll() []*user.User {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	u := make([]*user.User, 0, len(r.users))
	for _, val := range r.users {
		u = append(u, val)
	}
	return u
}

func (r *userRepository) CheckLogin(username string, password string) (*user.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	for _, val := range r.users {
		// check if the username and password match with a user in the store
		// -------------------- IMPORTANT NOTICE: -----------------
		// obviously it is unacceptable to compare plaintext passwords in a real-world application
		// however, since this is a prototype application we can overlook this flaw
		// admittedly though, it would be quite simpleto add a hashing method so only the hashes would be stored and compared
		if username == val.Username && password == val.Password {
			return val, nil
		}
	}
	return nil, user.ErrUnknown
}

func NewUserRepository() user.Repository {
	r := &userRepository{
		users: make(map[user.UserId]*user.User),
	}

	return r
}