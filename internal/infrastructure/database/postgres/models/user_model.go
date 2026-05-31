package models

import (
	"like-api/internal/domain/entities"
	"time"
)

// UserModel is the GORM representation of a user.
// It is completely separate from the domain entity and is used
// only for database persistence.
type UserModel struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	Username  string    `gorm:"uniqueIndex;not null;type:varchar(50)"`
	Name      string    `gorm:"not null;type:varchar(100)"`
	Email     string    `gorm:"uniqueIndex;not null;type:varchar(255)"`
	Password  string    `gorm:"not null;type:varchar(255)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (UserModel) TableName() string { return "users" }

func (m UserModel) ToEntity() entities.User {
	return entities.User{
		ID:        m.ID,
		Username:  m.Username,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
}

func UserModelFromEntity(u entities.User) UserModel {
	var createdAt time.Time
	if u.CreatedAt != "" {
		createdAt, _ = time.Parse(time.RFC3339, u.CreatedAt)
	}

	return UserModel{
		ID:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: createdAt,
	}
}
