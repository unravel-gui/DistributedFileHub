package jwt

import (
	"DisHub/common"
	"DisHub/common/utils"
	"DisHub/repository"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GetJWTKey() []byte {
	return []byte(common.JWT_KEY)

}

func GenerateJWTToken(subject string) (string, error) {
	// 设置token有效期为24小时
	expirationTime := utils.GetNow().Add(24 * time.Hour)

	// 创建JWT的声明
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   subject,
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token并返回字符串
	return token.SignedString(GetJWTKey())
}

func ParseJWTToken(tokenString string) (*jwt.StandardClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTKey(), nil
	})
	if err != nil {
		return nil, err
	}

	// 验证token并返回声明
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func GenerateUserToken(user *repository.User) (string, error) {
	u := &repository.User{
		Uid:      user.Uid,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}
	jsonData, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	jsonString := string(jsonData)
	token, err := GenerateJWTToken(jsonString)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseUserToken(token string) (*repository.User, error) {
	claims, err := ParseJWTToken(token)
	if err != nil {
		return nil, err
	}
	if IsTokenExpired(claims) {
		return nil, nil
	}
	var user repository.User
	// 将 JSON 字符串解析到 User 对象中
	err = json.Unmarshal([]byte(claims.Subject), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsTokenExpired(claims *jwt.StandardClaims) bool {
	now := utils.GetNow().Unix()
	return claims.ExpiresAt < now
}
