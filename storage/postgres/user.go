package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"
	pb "user/genproto/userservice"
	"user/logger"
	"user/storage"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type UserRepository struct {
	Db  *sql.DB
	Log *slog.Logger
}

func NewUserRepository(db *sql.DB) storage.IUserStorage {
	return &UserRepository{Db: db, Log: logger.NewLogger()}
}

func (s *UserRepository) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	query := `	INSERT INTO users (
						id, username, password, role, created_at, updated_at
						) VALUES (
						 	$1, $2, $3, $4, $5, $6)`

	id := uuid.NewString()
	role := "user"
	newtime := time.Now().Format("2006-01-02 15:04:05")

	_, err := s.Db.ExecContext(ctx, query, id, req.Username, req.Password, role, newtime, newtime)
	if err != nil {
		s.Log.Error(fmt.Sprintf("Register qilishda xatolik: %v", err))
		return nil, err
	}

	user := pb.User{
		Id:        id,
		Username:  req.Username,
		Password:  req.Password,
		Role:      role,
		CreatedAt: newtime,
		UpdatedAt: newtime,
	}
	return &pb.RegisterUserResponse{User: &user}, nil
}

func (s *UserRepository) StoreRefreshToken(ctx context.Context, req *pb.StoreRefreshTokenReq) (*pb.StoreRefreshTokenRes, error) {
	query := `	insert into refresh_token(
					id, user_id, token, expires_at
				)values(
					$1, $2, $3, $4)`

	id := uuid.NewString()
	_, err := s.Db.ExecContext(ctx, query, id, req.UserId, req.Token, req.ExpiresAt)
	if err != nil {
		s.Log.Error(fmt.Sprintf("Tokendi bazaga qushishda xatolik: %v", err))
		return nil, err
	}
	return nil, nil
}

func (s *UserRepository) GetByUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.GetByUserResponse, error) {
	query := `
		SELECT 
			id, username, password, role 
		FROM 
			users 
		WHERE 
			username = $1 AND deleted_at is null
	`
	var user pb.GetByUserResponse
	err := s.Db.QueryRowContext(ctx, query, req.Username).Scan(&user.Id, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Log.Error(fmt.Sprintf("user not found: %v", err))
			return nil, fmt.Errorf("user not found")
		}
		s.Log.Error(fmt.Sprintf("Malumotlarni olishda xatolik: %v", err))
		return nil, err
	}
	return &user, nil
}

func (s *UserRepository) RefReshToken(ctx context.Context, req *pb.RefReshTokenReq) (*pb.RefReshTokenRes, error) {
	err := godotenv.Load(".env")
	if err != nil {
		s.Log.Error(fmt.Sprintf(".env faylini yuklashda xatolik yuz berdi: %v", err))
		return nil, fmt.Errorf(".env faylini yuklashda xatolik yuz berdi")
	}

	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	newSigningKey := stringWithCharset(6, charset, seededRand)

	err = os.Setenv("SIGNING_KEY", newSigningKey)
	if err != nil {
		s.Log.Error(fmt.Sprintf("Yangi SIGNING_KEY sozlamasida xatolik yuz berdi: %v", err))
		return nil, fmt.Errorf("yangi SIGNING_KEY sozlamasida xatolik yuz berdi")
	}

	err = godotenv.Write(map[string]string{
		"USER_SERVICE": os.Getenv("USER_SERVICE"),
		"USER_ROUTER":  os.Getenv("USER_ROUTER"),
		"DB_USER":      os.Getenv("DB_USER"),
		"DB_HOST":      os.Getenv("DB_HOST"),
		"DB_NAME":      os.Getenv("DB_NAME"),
		"DB_PASSWORD":  os.Getenv("DB_PASSWORD"),
		"DB_PORT":      os.Getenv("DB_PORT"),
		"SIGNING_KEY":  newSigningKey,
	}, ".env")

	if err != nil {
		s.Log.Error(fmt.Sprintf(".env fayliga yozishda xatolik yuz berdi: %v", err))
		return nil, fmt.Errorf(".env fayliga yozishda xatolik yuz berdi")
	}

	s.Log.Info("SIGNING_KEY muvaffaqiyatli yangilandi")

	return &pb.RefReshTokenRes{
		Message: "SIGNING_KEY muvaffaqiyatli yangilandi",
	}, nil
}

func stringWithCharset(length int, charset string, seededRand *rand.Rand) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
