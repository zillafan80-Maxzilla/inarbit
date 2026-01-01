package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务
type AuthService struct {
	db        *Database
	jwtSecret string
}

// NewAuthService 创建认证服务
func NewAuthService(db *Database, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (string, *User, error) {
	// 获取用户
	user, err := s.db.GetUser(username)
	if err != nil {
		log.Printf("登录失败：用户 %s 不存在", username)
		return "", nil, fmt.Errorf("用户名或密码错误")
	}

	// 检查用户是否激活
	if !user.IsActive {
		log.Printf("登录失败：用户 %s 已被禁用", username)
		return "", nil, fmt.Errorf("用户已被禁用")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("登录失败：用户 %s 密码错误", username)
		return "", nil, fmt.Errorf("用户名或密码错误")
	}

	// 生成JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		log.Printf("生成token失败：%v", err)
		return "", nil, fmt.Errorf("生成token失败")
	}

	log.Printf("✓ 用户 %s 登录成功", username)

	// 返回用户信息（不包含密码）
	return token, user, nil
}

// GenerateToken 生成JWT token
func (s *AuthService) GenerateToken(user *User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // token有效期24小时

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      expiresAt.Unix(),
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken 验证JWT token
func (s *AuthService) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("无效的签名方法")
			}
			return []byte(s.jwtSecret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("token解析失败: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("无效的token")
	}

	return claims, nil
}

// HashPassword 哈希密码
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword 验证密码
func (s *AuthService) VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
