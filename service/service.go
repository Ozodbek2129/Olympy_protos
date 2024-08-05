package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	pb "user/genproto/userservice"
	"user/storage"
	"user/storage/postgres"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	User   storage.IStorage
	Logger *slog.Logger
}

func NewUserService(db *sql.DB, logger *slog.Logger) *UserService {
	return &UserService{
		User:   postgres.NewPostgresStorage(db, logger),
		Logger: logger,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	resp, err := s.User.User().RegisterUser(ctx, req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Servicedan crudlarga malumotni yuborishda xatolik register: %v", err))
		return nil, err
	}
	return resp, nil
}

func (s *UserService) StoreRefreshToken(ctx context.Context, req *pb.StoreRefreshTokenReq) (*pb.StoreRefreshTokenRes,error) {
	_, err := s.User.User().StoreRefreshToken(ctx, req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("servisda storerefresh tokenda xatolik: %v", err))
		return nil,err
	}
	return nil,nil
}

func (s *UserService) GetByUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.GetByUserResponse, error){
	resp,err:=s.User.User().GetByUser(ctx,req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("GetbyUserdan malumotlarni olishda xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *UserService) RefReshToken(ctx context.Context, req *pb.RefReshTokenReq) (*pb.RefReshTokenRes, error){
	resp,err:=s.User.User().RefReshToken(ctx,req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("Signing keyni uzgartirishda xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}