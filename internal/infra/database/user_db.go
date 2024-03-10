package database

import (
	"github.com/adrianomota/fullcycle/my-api/internal/entity"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{DB: db}
}

func (u *User) Create(user *entity.User) error {
	return u.DB.Create(user).Error
}

func (p *User) FindById(id string) (*entity.User, error) {
	var user entity.User
	err := p.DB.First(&user, "id = ?", id).Error
	return &user, err
}

func (u *User) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := u.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
