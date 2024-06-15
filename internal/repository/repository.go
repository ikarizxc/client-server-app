package repository

import (
	"dcsa-lab/internal/entities"
	"fmt"
	"sync"
)

type Repository struct {
	storage map[int]entities.User
	mutex   sync.Mutex
	lastId  int
}

func NewRepository() *Repository {
	repo := &Repository{
		storage: make(map[int]entities.User),
		lastId:  1,
	}

	repo.storage[1] = entities.User{
		Id:       1,
		Username: "admin",
		Email:    "admin@root.com",
		Password: "admin",
		IsAdmin:  true,
	}

	return repo
}

func (r *Repository) Add(user *entities.User) int {
	r.mutex.TryLock()

	r.lastId++

	user.Id = r.lastId
	r.storage[r.lastId] = *user

	return r.lastId
}

func (r *Repository) GetById(id int) (*entities.User, error) {
	r.mutex.TryLock()

	user, ok := r.storage[id]
	if !ok {
		return nil, fmt.Errorf("no user with id %d", id)
	}

	return &user, nil
}

func (r *Repository) GetByEmail(email string) (int, *entities.User, error) {
	r.mutex.TryLock()
	for k, v := range r.storage {
		if v.Email == email {
			return k, &v, nil
		}
	}

	return -1, nil, fmt.Errorf("no user with email %s", email)
}

func (r *Repository) GetByUsername(username string) (*entities.User, error) {
	r.mutex.TryLock()
	for _, v := range r.storage {
		if v.Username == username {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("no user with username %s", username)
}

func (r *Repository) GetAll() []*entities.User {
	r.mutex.TryLock()
	list := make([]*entities.User, len(r.storage))

	i := 0
	for _, v := range r.storage {
		list[i] = &v
		i++
	}

	return list
}

func (r *Repository) Delete(id int) {
	r.mutex.TryLock()
	delete(r.storage, id)
}

func (r *Repository) Update(id int, user *entities.User) error {
	r.mutex.TryLock()
	if _, ok := r.storage[id]; !ok {
		return fmt.Errorf("no user with id = %d", id)
	}

	userToChange := r.storage[id]

	if user.Username != "" {
		userToChange.Username = user.Username
	}
	if user.Email != "" {
		userToChange.Email = user.Email
	}
	if user.Password != "" {
		userToChange.Password = user.Password
	}

	r.storage[id] = userToChange

	return nil
}
