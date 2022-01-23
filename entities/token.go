package entities

// swagger:model
type Token struct {
	// The API access token, made by JWT
	AccessToken string `json:"access_token" `

	// When access token was expired, you have to use this to refresh access token
	// max length: 30
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
