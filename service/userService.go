package service

import (
	"DisHub/common/utils"
	"DisHub/config"
	"DisHub/repository"
)

type UserService struct {
	r *repository.UserRepository
}

var G_User UserService

func NewUserService(mysqlAddr string) *UserService {
	if mysqlAddr == "" {
		mysqlAddr = config.GetMySQLAddr()
	}
	r := repository.NewUserRepository(mysqlAddr)
	return &UserService{
		r: r,
	}
}

func (us *UserService) Load() {
	mysqlAddr := config.GetMySQLAddr()
	r := repository.NewUserRepository(mysqlAddr)
	us.r = r
}

func (us *UserService) Close() {
	us.r.Close()
}

func (us *UserService) Register(user *repository.User) (bool, error) {
	// 检查用户名是否已经存在
	if us.r.IsExisted(user.Username) {
		return false, nil
	}
	// 创建新用户
	user.Password = utils.CalculateStringHash(user.Password)
	err := us.r.Insert(user)
	if err != nil {
		return false, err
	}
	user.Password = ""
	return true, nil
}

func (us *UserService) Remove(user *repository.User) error {
	return us.r.Remove(user)
}

func (us *UserService) Login(params *repository.User) (*repository.User, error) {
	pass := utils.CalculateStringHash(params.Password)
	user, err := us.r.GetUserByUsername(params.Username)
	if err != nil {
		return nil, err
	}
	if user.Password != pass {
		return nil, nil
	}
	return user, nil
}

func (us *UserService) GetUserInfo(uid int) (*repository.User, error) {
	user, err := us.r.GetUserByUid(uid)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (us *UserService) UpdateUserInfo(uid int, userReq *repository.User) (*repository.User, error) {
	user, err := us.r.GetUserByUid(uid)
	if err != nil {
		return nil, err
	}
	user.UpdateUserInfo(userReq)
	us.r.UpdateUserInfo(user)
	user.Password = ""
	return user, nil
}

func (us *UserService) UpdateUserPassword(uid int, oldPass, NewPass string) (bool, error) {
	pass := utils.CalculateStringHash(oldPass)
	user, err := us.r.GetUserByUid(uid)
	if err != nil {
		return false, err
	}
	if user.Password != pass {
		return false, nil
	}
	user.Password = utils.CalculateStringHash(NewPass)
	us.r.UpdateUserInfo(user)
	user.Password = ""
	return true, nil
}
