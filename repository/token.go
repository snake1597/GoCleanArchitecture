package repository

import (
	"GoCleanArchitecture/entities"

	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) entities.TokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) UpdateRefreshToken(id string, refreshToken string) (err error) {
	err = r.db.Table("users").Where("id = ?", id).Update("refresh_token", refreshToken).Error
	return err
}

func (r *tokenRepository) CheckRefreshToken(id string, refreshToken string) (user entities.User, err error) {
	err = r.db.Table("users").Where("id = ? and refresh_token = ?", id, refreshToken).First(&user).Error
	return user, err
}
