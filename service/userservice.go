package service

import (
	"context"
	"errors"
	"time"

	"github.com/dingkegithub/user/dao"
	"github.com/dingkegithub/user/redis"
	"github.com/jinzhu/gorm"
)

// service 给用户
type UserInfoDto struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

// vo 用户view 输入
type RegisterUserVo struct {
	Email    string
	UserName string
	Password string
}

var (
	ErrUserExist  = errors.New("user exist")
	ErrPassword   = errors.New("email or password error")
	ErrRegistring = errors.New("email is registring")
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (*UserInfoDto, error)

	Register(ctx context.Context, vo *RegisterUserVo) (*UserInfoDto, error)
}

type UserServiceImpl struct {
	userDao dao.UserDao
}

func MakeUserService(userDao dao.UserDao) UserService {
	return &UserServiceImpl{
		userDao: userDao,
	}
}

func (u *UserServiceImpl) Login(ctx context.Context, email string, password string) (*UserInfoDto, error) {
	user, err := u.userDao.SelectByEmail(email)
	if err != nil {
		return nil, err
	}

	if user.Password == password {
		return &UserInfoDto{
			ID:       user.ID,
			UserName: user.UserName,
			Email:    user.Email,
		}, nil
	} else {
		return nil, ErrPassword
	}
}

func (u *UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVo) (*UserInfoDto, error) {
	lock := redis.GetLock(vo.Email, time.Duration(5)*time.Second)
	err := lock.Lock()
	if err != nil {
		return nil, ErrRegistring
	}
	defer lock.Unlock()

	existUser, err := u.userDao.SelectByEmail(vo.Email)
	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		newUser := &dao.UserEntity{
			UserName: vo.UserName,
			Email:    vo.Email,
			Password: vo.Password,
		}

		if err := u.userDao.Save(newUser); err == nil {
			return &UserInfoDto{
				ID:       newUser.ID,
				Email:    newUser.Email,
				UserName: newUser.UserName,
			}, nil
		}
	}

	if err == nil {
		return nil, ErrUserExist
	}

	return nil, err
}
