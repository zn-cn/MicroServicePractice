package handler

import (
	userPb "Ethan/MicroServicePractice/interface-center/out/user"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAPIHandler struct {
	userClient userPb.UserServiceClient
}

func GetUserHandler(userClient userPb.UserServiceClient) *UserAPIHandler {
	return &UserAPIHandler{
		userClient: userClient,
	}
}

func (s *UserAPIHandler) Login(c *gin.Context) {
	user := userPb.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := s.userClient.Auth(context.Background(), &user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": resp.Token,
	})
}

func (s *UserAPIHandler) Sign(c *gin.Context) {
	user := userPb.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	_, err := s.userClient.Create(context.Background(), &user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
