package grpcserver

import (
	"context"
	"errors"

	"github.com/junaid9001/zvynt/auth/config"
	"github.com/junaid9001/zvynt/auth/db"
	"github.com/junaid9001/zvynt/auth/utils"
	"github.com/junaid9001/zvynt/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	auth.UnimplementedAuthServiceServer
	store *db.Store
	cfg   *config.Config
}

func NewAuthServiceServer(store *db.Store, cfg *config.Config) *Server {
	return &Server{store: store, cfg: cfg}
}

func (s *Server) Signup(ctx context.Context, req *auth.SignupRequest) (*auth.SignupResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error ")
	}
	user, err := s.store.CreateUser(db.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		if errors.Is(err, db.ErrDuplicateKey) {

			return nil, status.Error(codes.AlreadyExists, "email already exist")
		}

		return nil, status.Error(codes.Internal, "internal server error ")

	}

	return &auth.SignupResponse{UserId: user.ID.String()}, nil
}

func (s *Server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {

	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "email not found")
		}

		return nil, status.Error(codes.Internal, "internal server error ")

	}

	err = utils.CompareHashAndPassword(user.Password, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	accessToken, err := utils.GenerateJwt(s.cfg.JWT_SECRET, user.ID.String(), user.Role)

	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error ")
	}

	return &auth.LoginResponse{UserId: user.ID.String(), AccessToken: accessToken}, nil
}

func (s *Server) Health(ctx context.Context, req *auth.Empty) (*auth.Empty, error) {
	return &auth.Empty{}, nil
}

func (s *Server) GetUser(
	ctx context.Context,
	req *auth.UserRequest,
) (*auth.UserResponse, error) {

	var (
		user *db.User
		err  error
	)

	switch v := req.Identifier.(type) {

	case *auth.UserRequest_Email:
		user, err = s.store.GetUserByEmail(v.Email)

	case *auth.UserRequest_UserId:
		user, err = s.store.GetUserByID(v.UserId)

	default:
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid identifier",
		)
	}

	if err != nil {

		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"user not found",
			)
		}

		return nil, status.Error(
			codes.Internal,
			"internal server error",
		)
	}

	return &auth.UserResponse{
		Email:  user.Email,
		UserId: user.ID.String(),
	}, nil
}
