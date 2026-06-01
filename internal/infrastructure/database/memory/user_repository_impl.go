package memory

import (
	"errors"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"sync"
)

type userRepositoryImpl struct {
	users map[string]entities.User
	mu    sync.RWMutex
}

func NewUserRepository() repositories.UserRepository {
	return &userRepositoryImpl{
		users: make(map[string]entities.User),
	}
}

func (r *userRepositoryImpl) Create(user entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
	return nil
}

func (r *userRepositoryImpl) FindAll() ([]entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]entities.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepositoryImpl) FindByID(id string) (entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return entities.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return entities.User{}, errors.New("user not found")
}

func (r *userRepositoryImpl) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.users[id]
	return exists
}

func (r *userRepositoryImpl) ExistsByEmail(email string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return true
		}
	}
	return false
}
