package models

type RefreshToken struct {
	UserId    string
	Token     string
	ExpiresAt int64
}

type UserInfo struct {
	Id       string
	Username string
	Password string
	Role     string
}