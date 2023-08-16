package user_handler

import (
	"errors"
	"github.com/dvdxa/add-to-favorites/internal/domain"
	"github.com/dvdxa/add-to-favorites/internal/services"
	"github.com/dvdxa/add-to-favorites/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

var (
	ErrInvalidUserNameOrPassLength = errors.New("username or password must be at least 5 characters")
	ErrInvalidCharacters           = errors.New("username or password can have only underscores, letters and numbers")
	ErrTooManyUnderscore           = errors.New("username must have 2 underscores maximum")
	ErrInvalidUnderscore           = errors.New("username or password cannot begin or end with underscore")
)

type UserHandler struct {
	log         logger.Logger
	userService services.UserServicePort
}

func NewUserHandler(log logger.Logger, userServicePort services.UserServicePort) *UserHandler {
	return &UserHandler{
		log:         log,
		userService: userServicePort,
	}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var user domain.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.log.Errorf("failed to bind json: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	err = h.ValidateRequest(user)
	if err != nil {
		h.log.Errorf("invalid user credentials: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": err.Error(),
		})
		return
	}
	err = h.userService.CreateUser(user)
	if err != nil {
		h.log.Errorf("failed to create user: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "user created")
}

func (h *UserHandler) SignIn(c *gin.Context) {
	var user domain.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.log.Errorf("failed to bind json: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	err = h.ValidateRequest(user)
	if err != nil {
		h.log.Errorf("invalid user credentials: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": err.Error(),
		})
		return
	}
	tokenStr, err := h.userService.GenerateToken(user)
	if err != nil {
		h.log.Errorf("failed to generate token: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.Writer.Header().Set("Token", tokenStr)
	c.JSON(http.StatusOK, gin.H{
		"message": "access token in header",
	})
}

func (h *UserHandler) ValidateUser(c *gin.Context) {
	tokenStr := c.GetHeader("token")
	userId, err := h.userService.ParseToken(tokenStr)
	if err != nil {
		h.log.Errorf("failed to parse token: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": err.Error(),
		})
		return
	}
	userIdfloat := userId.(float64)
	c.Set("userId", userIdfloat)
	c.Next()
}

func (h *UserHandler) ValidateRequest(user domain.User) error {
	//All checks must be done in api service
	if len(user.Name) < 6 || len(user.Password) < 5 {
		return ErrInvalidUserNameOrPassLength
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(user.Name) {
		return ErrInvalidCharacters
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(user.Password) {
		return ErrInvalidCharacters
	}

	underscoreCount := 0
	for i, char := range user.Name {
		if char == '_' {
			underscoreCount++
			if i == 0 || i == len(user.Name)-1 {
				return ErrInvalidUnderscore
			}
		}
	}
	if underscoreCount > 2 {
		return ErrTooManyUnderscore
	}

	underscoreCountPass := 0
	for i, char := range user.Password {
		if char == '_' {
			underscoreCountPass++
			if i == 0 || i == len(user.Password)-1 {
				return ErrInvalidUnderscore
			}
		}
	}
	if underscoreCountPass > 2 {
		return ErrTooManyUnderscore
	}
	return nil
}
