package service

//func TestUserRepository_Register(t *testing.T) {
//	// 设置数据库连接
//	db := common.TEST_MYSQL_ADDR
//
//	// 创建 UserRepository
//	us := NewUserRepository(db)
//	testUser := "test" + uuid.New().String()
//	// 注册新用户
//	_, err := us.Register(testUser, "testpassword")
//	if err != nil {
//		t.Errorf("resgister failed, err:%v\n", err)
//	}
//	// 尝试使用相同用户名注册新用户，应该返回错误
//	ok, err := us.Register(testUser, "testpassword")
//	if err != nil {
//		t.Errorf("resgister should be failed, but got unexpected err:%v\n", err)
//	}
//	if ok {
//		t.Errorf("resgister should be failed\n")
//	}
//}

//func TestUserRepository_Login(t *testing.T) {
//	// 设置数据库连接
//	db := common.TEST_MYSQL_ADDR
//	// 创建 UserRepository
//	us := NewUserRepository(db)
//	testUser := "test" + uuid.New().String()
//	pass := "testpassword"
//	// 注册新用户
//	ok, err := us.Register(testUser, pass)
//	if !ok {
//		t.Errorf("resgister failed, err:%v\n", err)
//	}
//
//	// 登录测试
//	isLogin := us.Login(testUser, pass)
//	if !isLogin {
//		t.Errorf("login failed\n")
//	}
//
//	// 使用错误的密码登录，应该返回错误
//	isLogin = us.Login(testUser, "wrongpassword")
//	if isLogin {
//		t.Errorf("login should be failed\n")
//	}
//}
