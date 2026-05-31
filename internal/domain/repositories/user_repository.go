package repositories

import "like-api/internal/domain/entities"

type UserRepository interface {
	Create(user entities.User) error
	FindAll() ([]entities.User, error)
	FindByID(id string) (entities.User, error)
	FindByEmail(email string) (entities.User, error)
	Exists(id string) bool
	ExistsByEmail(email string) bool
}
