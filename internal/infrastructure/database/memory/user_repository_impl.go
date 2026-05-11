package memory

import (
	"errors"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"sync"
)

type UserRepositoryImpl struct {
	users map[string]entities.User
	mu    sync.RWMutex
}

func NewUserRepository() repositories.UserRepository {
	return &UserRepositoryImpl{
		users: make(map[string]entities.User),
	}
}

func (r *UserRepositoryImpl) Create(user entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
	return nil
}

func (r *UserRepositoryImpl) FindAll() ([]entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]entities.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepositoryImpl) FindByID(id string) (entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return entities.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepositoryImpl) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.users[id]
	return exists
}
