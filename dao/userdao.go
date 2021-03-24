package dao

import "time"

type UserEntity struct {
	ID       int64
	UserName string
	Email    string
	Password string
	CreateAt time.Time
}

func (UserEntity) TableName() string {
	return "user"
}

type UserDao interface {
	SelectByEmail(email string) (*UserEntity, error)
	Save(*UserEntity) error
}

type UserDaoImpl struct {
}

func (u *UserDaoImpl) SelectByEmail(email string) (*UserEntity, error) {
	user := &UserEntity{}
	err := db.Where("email = ?", email).First(user).Error
	return user, err
}

func (u *UserDaoImpl) Save(user *UserEntity) error {
	return db.Create(user).Error
}
