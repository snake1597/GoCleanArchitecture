package usecase

import (
	"GoCleanArchitecture/entities"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenUsecase struct {
	authKey   string
	tokenRepo entities.TokenRepository
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewTokenUsecase(jwtKey string, tokenRepo entities.TokenRepository) entities.TokenUsecase {
	return &TokenUsecase{
		jwtKey,
		tokenRepo,
	}
}

func (u *TokenUsecase) CreateToken(id string) (authToken entities.Token, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["user_id"] = id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token.Claims = claims
	authToken.AccessToken, err = token.SignedString([]byte(u.authKey))
	if err != nil {
		return entities.Token{}, fmt.Errorf("generate access token failed")
	}
	authToken.RefreshToken = randSeq(30)

	err = u.tokenRepo.UpdateRefreshToken(id, authToken.RefreshToken)
	if err != nil {
		return entities.Token{}, fmt.Errorf("update token failed")
	}

	return authToken, nil
}

func (u *TokenUsecase) VerifyToken(token string) (jwtClaim map[string]interface{}, err error) {
	authToken, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(u.authKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token error, " + err.Error())
	}

	if claim, ok := authToken.Claims.(jwt.MapClaims); ok {
		jwtClaim = map[string]interface{}{
			"user_id": claim["user_id"],
		}
		return jwtClaim, nil
	}

	return nil, fmt.Errorf("token is not valid")
}

func (u *TokenUsecase) RefreshAccessToken(id string, refreshToken string) (token entities.Token, err error) {
	_, err = u.tokenRepo.CheckRefreshToken(id, refreshToken)
	if err != nil {
		return entities.Token{}, fmt.Errorf("refresh token is not exist")
	}

	token, err = u.CreateToken(id)
	if err != nil {
		return entities.Token{}, err
	}

	return token, nil
}

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
