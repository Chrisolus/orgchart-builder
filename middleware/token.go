package middleware

import (
	"crypto/rsa"
	"errors"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Claims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	IsRefresh bool   `json:"isRefresh"`
	jwt.RegisteredClaims
}

func getPrvKey() (*rsa.PrivateKey, error) {
	prvKeyBytes, err := os.ReadFile(viper.GetString("jwt.private"))
	if err != nil {
		return nil, err
	} else {
		if prvKey, err := jwt.ParseRSAPrivateKeyFromPEM(prvKeyBytes); err != nil {
			return nil, err
		} else {
			return prvKey, nil
		}
	}
}

func getPubKey() (*rsa.PublicKey, error) {
	pubKeyBytes, err := os.ReadFile(viper.GetString("jwt.public"))
	if err != nil {
		return nil, err
	} else {
		if pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes); err != nil {
			return nil, err
		} else {
			return pubKey, nil
		}
	}
}

func generateSignedToken(claims *Claims, isRefresh bool, prvKey *rsa.PrivateKey, expTime time.Time) (string, error) {
	if claims == nil {
		return "", errors.New("claims cannot be nil")
	}

	if claims.UserID == 0 {
		return "", errors.New("claims: userID is required")
	}
	claims.IsRefresh = isRefresh
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "org_chart",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(prvKey)

	if err != nil {
		log.Println("GenerateSignedTokenError: ", err.Error())
		return "", errors.New("error while signing the token")
	}

	return tokenString, nil
}

func GenAuthAndRefreshToken(claims *Claims) (gin.H, error) {
	prvKey, err := getPrvKey()
	if err != nil {
		log.Println("GetPrivateKeyError: ", err.Error())
		return nil, errors.New("error fetching private key")
	}
	authString, err := generateSignedToken(claims, false, prvKey, time.Now().Add(viper.GetDuration("jwt.auth_exp")))
	if err != nil {
		return nil, err
	}
	refreshString, err := generateSignedToken(claims, true, prvKey, time.Now().Add(viper.GetDuration("jwt.refresh_exp")))
	if err != nil {
		return nil, err
	}

	return gin.H{"refresh": refreshString, "auth": authString}, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		pubKey, err := getPubKey()
		if err != nil {
			log.Println("GetPublicKeyError: ", err.Error())
			return nil, errors.New("cannot fetch public key")
		} else {
			return pubKey, nil
		}
	})
	if err != nil {
		log.Println("TOKEN: ", tokenString)
		log.Println("ValidateTokenError: ", err.Error())
		return nil, err
	}
	if !token.Valid {
		log.Println("ValidateTokenError: Invalid Token")
		return nil, errors.New("invalid token")
	}
	claims := token.Claims.(*Claims)
	return claims, nil
}
