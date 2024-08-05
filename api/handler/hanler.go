package handler

import (
	"log/slog"
	pb "user/genproto/userservice"
)

type Handler struct {
	AuthUser pb.UserServiceClient
	Log      *slog.Logger
}