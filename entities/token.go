package entities

type Token struct {
	AccessToken  string `json:"access_token" `
	RefreshToken string `json:"refresh_token" `
}

type TokenRepository interface {
	UpdateRefreshToken(id string, refreshToken string) (err error)
	CheckRefreshToken(id string, refreshToken string) (user User, err error)
}

type TokenUsecase interface {
	CreateToken(id string) (token Token, err error)
	VerifyToken(token string) (data map[string]interface{}, err error)
	RefreshAccessToken(id string, refreshToken string) (token Token, err error)
}
