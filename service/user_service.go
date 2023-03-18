package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"sync"

	"github.com/e-commerce-microservices/user-service/pb"
	"github.com/e-commerce-microservices/user-service/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UserService implement grpc UserServiceServer
type UserService struct {
	userStore  *repository.Queries
	authClient pb.AuthServiceClient
	pb.UnimplementedUserServiceServer
}

// NewUserService creates a new UserService instance
func NewUserService(userStore *repository.Queries, authClient pb.AuthServiceClient) *UserService {
	service := &UserService{
		userStore:  userStore,
		authClient: authClient,
	}

	return service
}

// GetListUser ...
func (user *UserService) GetListUser(ctx context.Context, req *pb.GetListUserRequest) (*pb.GetListUserResponse, error) {
	result := make([]*pb.User, 0)
	wg := sync.WaitGroup{}
	wg.Add(len(req.ListUserId))
	for _, id := range req.ListUserId {
		go func(id int64) {
			defer wg.Done()
			user, err := user.userStore.GetUserByID(ctx, id)
			log.Println("user: ", user, err)
			if err == nil {
				result = append(result, &pb.User{
					Id:    id,
					Email: user.Email,
					Profile: &pb.UserProfile{
						UserName: user.UserName,
						Phone:    "",
						Avatar:   "",
					},
					Address: []*pb.UserAddress{},
				})
			}
		}(id)
	}

	wg.Wait()
	return &pb.GetListUserResponse{
		ListUser: result,
	}, nil
}
func isVietnamesePhoneNumber(number string) bool {
	// return ``.test(number)
	ok, _ := regexp.MatchString("/(03|05|07|08|09|01[2|6|8|9])+([0-9]{8})\b/", number)
	return ok
}

func (user *UserService) UpdateProfile(ctx context.Context, req *pb.UserProfile) (*pb.GeneralResponse, error) {
	// if len(req.GetPhone()) > 0 && !isVietnamesePhoneNumber(req.GetPhone()) {
	// 	return nil, errors.New("Định dạng số điện thoại không đúng")
	// }
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

	err = user.userStore.UpdateUserName(ctx, repository.UpdateUserNameParams{
		UserName: req.GetUserName(),
		Phone: sql.NullString{
			String: req.GetPhone(),
			Valid:  true,
		},
		ID: id,
		Address: sql.NullString{
			String: req.GetAddress(),
			Valid:  true,
		},
		Note: sql.NullString{
			String: req.GetNote(),
			Valid:  true,
		},
	})
	if err != nil {
		return nil, errors.New("Cập nhật không thành công")
	}

	return &pb.GeneralResponse{
		Message:    "Cập nhật thông tin thành công",
		StatusCode: 0,
	}, nil
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
		Profile: &pb.UserProfile{
			UserName: tmp.UserName,
			Phone:    tmp.Phone.String,
			Avatar:   "",
		},
		Address: []*pb.UserAddress{{
			Address: tmp.Address.String,
			Note:    tmp.Note.String,
		}},
		Gender: tmp.Gender.String,
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

	tmp, err := user.userStore.GetUserByID(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.User{
		Id:           tmp.ID,
		Email:        tmp.Email,
		Role:         pb.UserRole(pb.UserRole_value[string(tmp.Role)]),
		ActiveStatus: tmp.ActiveStatus,
		Profile: &pb.UserProfile{
			UserName: tmp.UserName,
			Phone:    tmp.Phone.String,
			Avatar:   "",
		},
		Address: []*pb.UserAddress{{
			Address: tmp.Address.String,
			Note:    tmp.Note.String,
		}},
		Gender: tmp.Gender.String,
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
