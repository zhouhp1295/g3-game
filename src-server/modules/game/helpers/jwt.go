// Copyright (c) 554949297@qq.com . 2022-2022. All rights reserved

package helpers

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type wsJwtClaims struct {
	jwt.StandardClaims
	Uid int64
}

func NewWsJwtToken(secret []byte, uid int64) (string, error) {
	expiredTime := time.Now().Add(time.Hour * 24 * 365)
	claims := wsJwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
		Uid: uid,
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(secret)

	return token, err
}

func ParseWsJwtToken(token string, secret []byte) (*wsJwtClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, new(wsJwtClaims), func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*wsJwtClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
