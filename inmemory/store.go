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

// adds a user to the repository store
func (r *userRepository) Add(u *user.User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.users[u.Id] = u
	return nil
}

// attempts to find the user by UserId inside the repository store
func (r *userRepository) Find(id user.UserId) (*user.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.users[id]; ok {
		return val, nil
	}
	return nil, user.ErrUnknown
}

// returns alls users in an array
func (r *userRepository) FindAll() []*user.User {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	u := make([]*user.User, 0, len(r.users))
	for _, val := range r.users {
		u = append(u, val)
	}
	return u
}

// checks the login credentials of a user
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

// returns an instance of a user repository
func NewUserRepository() user.Repository {
	r := &userRepository{
		users: make(map[user.UserId]*user.User),
	}

	r.users[user.U0001] = user.Rey
	r.users[user.U0002] = user.Kylo
	r.users[user.U0003] = user.Luke

	return r
}