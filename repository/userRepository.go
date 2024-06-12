package repository

import (
	"DisHub/common/db"
	"DisHub/config"
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	db *gorm.DB
}

var G_User UserRepository

func NewUserRepository(mysqlAddr string) *UserRepository {
	if mysqlAddr == "" {
		mysqlAddr = config.GetLocalAddr()
	}
	d, err := db.NewConnect(mysqlAddr)
	if err != nil {
		log.Fatalln("connect to mysql err: ", err)
	}
	d.AutoMigrate(&User{})
	return &UserRepository{
		db: d,
	}
}

func (us *UserRepository) Close() {
	sqlDB, err := us.db.DB()
	if err != nil {
		log.Printf("close Mysql connect err:%v\n", err)
		return
	}
	sqlDB.Close()
}

func (us *UserRepository) Insert(user *User) error {
	result := us.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (us *UserRepository) Remove(user *User) error {
	result := us.db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// IsExisted 检查用户是否存在
func (us *UserRepository) IsExisted(username string) bool {
	var count int64
	us.db.Model(&User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

func (us *UserRepository) GetUserByUsername(username string) (*User, error) {
	var user User
	if err := us.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserRepository) GetUserByUid(uid int) (*User, error) {
	var user User
	if err := us.db.First(&user, uid).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserRepository) UpdateUserInfo(user *User) (*User, error) {
	if err := us.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
