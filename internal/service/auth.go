// Package service 认证服务。JWT 令牌生成/验证 + Redis Token 存储 + 密码哈希。
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// AuthService JWT 令牌与密码操作。
type AuthService struct {
	secret      []byte
	redisClient *redis.Client
}

// NewAuthService 创建认证服务实例。
func NewAuthService(secret string, rdb *redis.Client) *AuthService {
	return &AuthService{secret: []byte(secret), redisClient: rdb}
}

// GenerateToken 生成 JWT 令牌，有效期 24 小时。
func (s *AuthService) GenerateToken(userID int, username string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ParseToken 解析并验证 JWT 令牌，返回 claims。
func (s *AuthService) ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("无效令牌")
}

// ======================== Redis  Token 存储 ========================

// StoreToken 将令牌存入 Redis，TTL 与 JWT 过期时间一致。
func (s *AuthService) StoreToken(ctx context.Context, userID int, tokenStr string, ttl time.Duration) error {
	if s.redisClient == nil {
		return nil
	}
	return s.redisClient.Set(ctx, "token:"+tokenStr, userID, ttl).Err()
}

// ValidateToken 检查令牌是否存在于 Redis（未被撤销）。
func (s *AuthService) ValidateToken(ctx context.Context, tokenStr string) (bool, error) {
	if s.redisClient == nil {
		return true, nil // Redis 未配置，跳过校验
	}
	n, err := s.redisClient.Exists(ctx, "token:"+tokenStr).Result()
	return n > 0, err
}

// RevokeToken 从 Redis 删除令牌，使其立即失效。
func (s *AuthService) RevokeToken(ctx context.Context, tokenStr string) error {
	if s.redisClient == nil {
		return nil
	}
	return s.redisClient.Del(ctx, "token:"+tokenStr).Err()
}

// ======================== 密码哈希 ========================

// HashPassword 对明文密码进行 bcrypt 哈希。
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证明文密码与哈希是否匹配。
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
