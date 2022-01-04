package grpc

import (
	"GoCleanArchitecture/entities"
	pb "GoCleanArchitecture/proto/user"
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	UserUsecase  entities.UserUsecase
	TokenUsecase entities.TokenUsecase
	pb.UnimplementedUserServer
}

func NewUserHandler(s *grpc.Server, userUsecase entities.UserUsecase, tokenUsecase entities.TokenUsecase) {
	handler := UserHandler{
		UserUsecase:  userUsecase,
		TokenUsecase: tokenUsecase,
	}

	pb.RegisterUserServer(s, &handler)
}

func (h *UserHandler) Register(ctx context.Context, in *pb.RegisterRequest) (response *pb.UserResponse, err error) {
	user := entities.User{
		Account:   in.GetAccount(),
		Password:  in.GetPassword(),
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Birthday:  in.GetBirthday(),
	}

	err = h.UserUsecase.Register(user)
	if err != nil {
		return nil, err
	}

	response = &pb.UserResponse{
		Status: "success",
	}

	return response, nil
}

func (h *UserHandler) Login(ctx context.Context, in *pb.LoginRequest) (response *pb.LoginResponse, err error) {
	user := entities.User{
		Account:  in.GetAccount(),
		Password: in.GetPassword(),
	}

	userId, err := h.UserUsecase.Login(user)
	if err != nil {
		return nil, err
	}

	token, err := h.TokenUsecase.CreateToken(userId)
	if err != nil {
		return nil, err
	}

	response = &pb.LoginResponse{
		UserId:       userId,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return response, nil
}

func (h *UserHandler) GetUser(ctx context.Context, in *pb.UserRequest) (response *pb.GetUserResponse, err error) {
	user, err := h.UserUsecase.GetUser(in.GetUserId())
	if err != nil {
		return nil, err
	}

	response = &pb.GetUserResponse{
		UserId:    strconv.Itoa(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Birthday:  user.Birthday,
	}

	return response, nil
}

func (h *UserHandler) GetAllUser(ctx context.Context, in *emptypb.Empty) (response *pb.GetAllUserResponse, err error) {
	userList, err := h.UserUsecase.GetAllUser()
	if err != nil {
		return nil, err
	}

	list := make([]*pb.GetUserResponse, len(userList))
	for _, user := range userList {
		a := &pb.GetUserResponse{
			UserId:    strconv.Itoa(user.ID),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Birthday:  user.Birthday,
		}
		list = append(list, a)
	}

	response = &pb.GetAllUserResponse{
		UserList: list,
	}
	return response, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (response *pb.UserResponse, err error) {
	user := entities.User{
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Birthday:  in.GetBirthday(),
	}

	err = h.UserUsecase.UpdateUser(in.GetUserId(), user)
	if err != nil {
		return nil, err
	}

	response = &pb.UserResponse{
		Status: "success",
	}
	return response, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, in *pb.UserRequest) (response *pb.UserResponse, err error) {

	err = h.UserUsecase.DeleteUser(in.GetUserId())
	if err != nil {
		return nil, err
	}

	response = &pb.UserResponse{
		Status: "success",
	}
	return response, nil
}
