package jwtutils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateTokens(userID uint64, exp time.Duration, secret string, refExp time.Duration, refreshSecret string, deviceID string) (signedAccessToken, signedRefreshToken string, err error) {
	// 1. 生成 Access Token
	accessExp := time.Now().Add(exp).Unix()
	accessClaims := jwt.MapClaims{
		"userID": userID,
		"exp":    accessExp, // 过期时间戳
		"type":   "access",
		"device": deviceID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err = accessToken.SignedString([]byte(secret))
	if err != nil {
		return
	}

	// 2. 生成 Refresh Token
	refreshExp := time.Now().Add(refExp).Unix()
	refreshClaims := jwt.MapClaims{
		"userID": userID,
		"exp":    refreshExp, // 过期时间戳
		"type":   "refresh",
		"device": deviceID,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err = refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return
	}

	return
}

func ParseToken(tokenString string, secret string, deviceID string) (userID uint64, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名算法: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	// 处理验证错误（包括过期）
	if err != nil || !token.Valid {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return
	}

	// 验证Token类型（必须是access，refresh token由专门接口处理）
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		err = errors.New("无效的Token")
		return
	}

	// 提取用户ID
	userID, ok = claims["userID"].(uint64)
	if !ok {
		err = errors.New("无效的Token")
		return
	}

	// 验证deviceID
	tokenDeviceID, ok := claims["device"].(string)
	if !ok || tokenDeviceID == "" || tokenDeviceID != deviceID {
		return
	}

	return
}
