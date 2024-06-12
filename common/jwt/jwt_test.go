package jwt

import (
	"DisHub/common"
	"DisHub/repository"
	"fmt"
	"testing"
)

func Test_JWT(t *testing.T) {
	// 创建一个用户
	user := repository.User{
		Uid:      common.TEST_UID,
		Username: "john_doe",
	}
	// 生成JWT令牌
	token, err := GenerateJWTToken(user.Username)
	if err != nil {
		t.Fatalf("Error generating JWT token:%v", err)
	}
	fmt.Println("jwt token is ", token)
	// 解析JWT令牌
	claims, err := ParseJWTToken(token)
	if err != nil {
		fmt.Println("Error parsing JWT token:", err)
		return
	}
	if claims.Subject != user.Username {
		t.Fatalf("jwt failed,subject=%s,claims=%+v", user.Username, claims)
	}
}
