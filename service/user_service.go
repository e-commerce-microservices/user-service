package service

import (
	"context"
	"fmt"

	"github.com/e-commerce-microservices/user-service/pb"
	"github.com/e-commerce-microservices/user-service/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userRepository interface {
	CreateUser(context.Context, repository.CreateUserParams) error
}

// UserService implement grpc UserServiceServer
type UserService struct {
	userStore userRepository
	pb.UnimplementedUserServiceServer
}

// NewUserService creates a new UserService instance
func NewUserService(userStore userRepository) *UserService {
	service := &UserService{
		userStore: userStore,
	}

	return service
}

// CreateUser creates new user in db if it isn't exist
func (user *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.GeneralResponse, error) {
	err := user.userStore.CreateUser(ctx, repository.CreateUserParams{
		Email:          req.Email,
		UserName:       req.UserName,
		HashedPassword: createHashedPassword(req.Password),
	})
	if err != nil {
		return &pb.GeneralResponse{
			Message:    fmt.Sprintf("user %s isn't created, please try again", req.UserName),
			StatusCode: 500,
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.GeneralResponse{
		Message:    fmt.Sprintf("user %s is created, please check your email(%s) to complete the registration", req.UserName, req.Email),
		StatusCode: int32(codes.OK),
	}, nil
}
