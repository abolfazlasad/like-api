package postgres

import (
	"errors"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"like-api/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	writeDB *gorm.DB
	readDB  *gorm.DB
}

func newUserRepository(db *DB) repositories.UserRepository {
	return &userRepositoryImpl{
		writeDB: db.WriteDB,
		readDB:  db.ReadDB,
	}
}

func (r *userRepositoryImpl) Create(user entities.User) error {
	model := models.UserModelFromEntity(user)
	return r.writeDB.Create(&model).Error
}

func (r *userRepositoryImpl) FindAll() ([]entities.User, error) {
	var userModels []models.UserModel
	if err := r.readDB.Find(&userModels).Error; err != nil {
		return nil, err
	}

	users := make([]entities.User, len(userModels))
	for i, m := range userModels {
		users[i] = m.ToEntity()
	}
	return users, nil
}

func (r *userRepositoryImpl) FindByID(id string) (entities.User, error) {
	var model models.UserModel
	err := r.readDB.Where("id = ?", id).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.User{}, errors.New("user not found")
	}
	if err != nil {
		return entities.User{}, err
	}
	return model.ToEntity(), nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (entities.User, error) {
	var model models.UserModel
	err := r.readDB.Where("email = ?", email).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.User{}, errors.New("user not found")
	}
	if err != nil {
		return entities.User{}, err
	}
	return model.ToEntity(), nil
}

func (r *userRepositoryImpl) Exists(id string) bool {
	var count int64
	r.readDB.Model(&models.UserModel{}).Where("id = ?", id).Count(&count)
	return count > 0
}

func (r *userRepositoryImpl) ExistsByEmail(email string) bool {
	var count int64
	r.readDB.Model(&models.UserModel{}).Where("email = ?", email).Count(&count)
	return count > 0
}
