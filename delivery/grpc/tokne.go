package grpc

import (
	"GoCleanArchitecture/entities"
	pbToken "GoCleanArchitecture/proto/token"
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type TokenHandler struct {
	TokenUsecase entities.TokenUsecase
	pbToken.UnimplementedTokenServer
}

func NewTokenHandler(s *grpc.Server, tokenUsecase entities.TokenUsecase) {
	handler := TokenHandler{
		TokenUsecase: tokenUsecase,
	}

	pbToken.RegisterTokenServer(s, &handler)
}

func (h *TokenHandler) RefreshAccessToken(ctx context.Context, in *pbToken.TokenRequest) (response *pbToken.TokenResponse, err error) {
	userId := in.GetUserId()
	refreshToken := in.GetToken()
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	token, err := h.TokenUsecase.RefreshAccessToken(userId, refreshToken)
	if err != nil {
		return nil, err
	}

	response = &pbToken.TokenResponse{
		UserId:       userId,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return response, nil
}
