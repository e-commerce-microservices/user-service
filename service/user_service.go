package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/e-commerce-microservices/user-service/pb"
	"github.com/e-commerce-microservices/user-service/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type userRepository interface {
	CreateUser(context.Context, repository.CreateUserParams) error
	GetUserByEmail(ctx context.Context, email string) (repository.User, error)
	RegisterSupplier(ctx context.Context, id int64) error
	GetUserByID(ctx context.Context, id int64) (repository.User, error)
}

// UserService implement grpc UserServiceServer
type UserService struct {
	userStore  userRepository
	authClient pb.AuthServiceClient
	pb.UnimplementedUserServiceServer
}

// NewUserService creates a new UserService instance
func NewUserService(userStore userRepository, authClient pb.AuthServiceClient) *UserService {
	service := &UserService{
		userStore:  userStore,
		authClient: authClient,
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GeneralResponse{
		Message: fmt.Sprintf("user %s is created, please check your email(%s) to complete the registration", req.UserName, req.Email),
	}, nil
}

// GetUserByEmail query user in db by input email
func (user *UserService) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.User, error) {
	// get user in db
	tmp, err := user.userStore.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// check password
	if !compareHashedPassword(req.Password, tmp.HashedPassword) {
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}

	return &pb.User{
		Id:           tmp.ID,
		Email:        tmp.Email,
		Role:         pb.UserRole(pb.UserRole_value[string(tmp.Role)]),
		ActiveStatus: tmp.ActiveStatus,
	}, nil
}

// GetUserById ...
func (user *UserService) GetUserById(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.User, error) {
	tmp, err := user.userStore.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.User{
		Id:           tmp.ID,
		Email:        tmp.Email,
		Role:         pb.UserRole(pb.UserRole_value[string(tmp.Role)]),
		ActiveStatus: tmp.ActiveStatus,
	}, nil
}

// GetMe query user in db by id parsed from header
func (user *UserService) GetMe(ctx context.Context, req *empty.Empty) (*pb.User, error) {
	// authenticated
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	claims, err := user.authClient.GetUserClaims(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(claims.Id, 10, 64)
	if err != nil {
		return nil, err
	}

	me, err := user.userStore.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:           id,
		Email:        me.Email,
		Role:         claims.GetUserRole(),
		ActiveStatus: me.ActiveStatus,
	}, nil
}

// SupplierRegister update user role to supplier (if user is customer)
func (user *UserService) SupplierRegister(ctx context.Context, _ *empty.Empty) (*pb.GeneralResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	claims, err := user.authClient.GetUserClaims(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	if claims.UserRole == pb.UserRole_supplier || claims.UserRole == pb.UserRole_admin {
		return nil, errors.New("Đăng kí không thành công, bạn đã đăng kí từ trước đó rồi")
	}

	id, err := strconv.ParseInt(claims.GetId(), 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user claims")
	}

	err = user.userStore.RegisterSupplier(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user claims")
	}

	return &pb.GeneralResponse{
		Message: "Đăng kí thành công, bạn đã trở thành người bán hàng",
	}, nil
}

// Ping pong
func (user *UserService) Ping(context.Context, *empty.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Message: "pong",
	}, nil

}
