package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"user/api/token"
	pb "user/genproto/userservice"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body userservice.RegisterUserRequest true "User data"
// @Success 202 {object} userservice.RegisterUserResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /register [post]
func (h Handler) Register(c *gin.Context) {
	req := pb.RegisterUserRequest{}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("bodydan malumotlarni olishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Pasworni hashlashda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Password = string(hashpassword)

	resp, err := h.AuthUser.RegisterUser(context.Background(), &req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Foydalanuvchi malumotlarni bazga yuborishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, resp)
}

// @Summary Login user
// @Description Login user with username and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body userservice.LoginUserRequest true "User data"
// @Success 202 {object} string
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /login [post]
func (h Handler) LoginUser(c *gin.Context) {
	req := pb.LoginUserRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	user, err := h.AuthUser.GetByUser(context.Background(), &req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("GetbyUserda xatolik: %v",err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := token.GenerateJWT(&pb.User{
		Id:       user.Id,
		Username: req.Username,
		Password: req.Password,
		Role:     user.Role,
	})

	_, err = h.AuthUser.StoreRefreshToken(context.Background(),&pb.StoreRefreshTokenReq{
		UserId:    user.Id,
		Token:     token.Refreshtoken,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	if err != nil {
		h.Log.Error(fmt.Sprintf("storefreshtokenda xatolik: %v",err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusAccepted, token)
}

// @Summary Refreshe Token
// @Description This endpoint refreshes the signing key and returns a confirmation message.
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {object} userservice.RefReshTokenRes
// @Failure 500 {object} string
// @Router /refresh-token [post]
func (h Handler) RefReshToken(c *gin.Context){
	req:=pb.RefReshTokenReq{}

	resp,err:=h.AuthUser.RefReshToken(c,&req)
	if err!=nil{
		h.Log.Error(fmt.Sprintf("Api da malumotlarni olishda xatolik: %v",err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":"Api da malumotlarni olishda xatolik: "+err.Error(),
		})
		return
	}

	c.JSON(200,resp)
}