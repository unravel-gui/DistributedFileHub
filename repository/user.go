package repository

type User struct {
	Uid      int    `gorm:"primaryKey;column:uid" json:"uid"`
	Username string `gorm:"unique;column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
	Nickname string `gorm:"column:nickname" json:"nickname"`
	Email    string `gorm:"column:email" json:"email"`
	Avatar   string `gorm:"column:avatar" json:"avatar"`
	IsAdmin  int    `gorm:"column:isAdmin;default:0" json:"isAdmin"`
}

func (User) TableName() string {
	return "oss_user" // 设置表名
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: password,
	}
}

func (u *User) UpdateUserInfo(newUser *User) {
	u.Avatar = newUser.Avatar
	u.Nickname = newUser.Nickname
	u.Email = newUser.Email
}
