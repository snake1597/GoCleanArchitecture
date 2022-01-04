package middlewares

import (
	"GoCleanArchitecture/entities"
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var authAllowList = map[string]string{"/api.User/Register": "auth", "/api.User/Login": "auth", "/api.Token/RefreshAccessToken": "auth"}

type AuthMiddlewares struct {
	tokenUsecase entities.TokenUsecase
}

func NewAuthMiddlewares(tokenUsecase entities.TokenUsecase) *AuthMiddlewares {
	return &AuthMiddlewares{tokenUsecase}
}

func (m *AuthMiddlewares) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var auth string

	if _, ok := authAllowList[info.FullMethod]; !ok {
		md, _ := metadata.FromIncomingContext(ctx)

		if value, ok := md["token"]; ok {
			auth = value[0]
		} else {
			return nil, fmt.Errorf("token is require")
		}

		token := strings.Split(auth, "Bearer ")[1]
		_, err := m.tokenUsecase.VerifyToken(token)
		if err != nil {
			return nil, err
		}
	}

	resp, err := handler(ctx, req)
	return resp, err
}
