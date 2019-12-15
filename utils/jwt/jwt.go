package jwt

import (
	"github.com/DualVectorFoil/Zelda/conf"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GetToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(conf.USER_INFO_TTL).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte("webapp"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
