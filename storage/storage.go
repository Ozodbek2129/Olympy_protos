package storage

import (
	"context"
	pb "user/genproto/userservice"
)

type IStorage interface {
	User() IUserStorage
	Close()
}

type IUserStorage interface {
	RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error)
	StoreRefreshToken(ctx context.Context, req *pb.StoreRefreshTokenReq) (*pb.StoreRefreshTokenRes, error)
	GetByUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.GetByUserResponse, error)
	RefReshToken(ctx context.Context, req *pb.RefReshTokenReq) (*pb.RefReshTokenRes, error)
}
