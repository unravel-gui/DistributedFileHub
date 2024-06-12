package repository

import (
	"DisHub/common"
	"fmt"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestUserRepository_Insert(t *testing.T) {
	// 设置数据库连接
	db := common.TEST_MYSQL_ADDR

	// 创建 UserRepository
	us := NewUserRepository(db)
	testUser := "test" + uuid.New().String()

	// 注册新用户
	err := us.Insert(&User{Username: testUser, Password: "testpassword"})
	if err != nil {
		t.Errorf("resgister failed, err:%v\n", err)
	}
	// 尝试使用相同用户名注册新用户，应该返回错误
	err = us.Insert(&User{Username: testUser, Password: "testpassword"})
	if err != nil {
		t.Errorf("resgister should be failed, but got unexpected err:%v\n", err)
	}
}

func Test_time(t *testing.T) {
	now := time.Now()
	fmt.Println(now)
}
