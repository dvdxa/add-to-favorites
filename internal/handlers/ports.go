package handlers

import (
	"github.com/dvdxa/add-to-favorites/internal/handlers/terminal_handler"
	"github.com/dvdxa/add-to-favorites/internal/handlers/user_handler"
	"github.com/dvdxa/add-to-favorites/internal/services"
	"github.com/dvdxa/add-to-favorites/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	user_handler.UserHandler
	terminal_handler.TerminalHandler
}

func NewHandler(log logger.Logger, service services.ServicePort) *Handler {
	return &Handler{
		UserHandler:     *user_handler.NewUserHandler(log, service.UserServicePort),
		TerminalHandler: *terminal_handler.NewTerminalHandler(log, service.TerminalServicePort),
	}
}
func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()
	router.POST("/user/sign-up", h.SignUp)
	router.POST("/user/sign-in", h.SignIn)
	router.GET("/terminals", h.ValidateUser, h.GetTerminalsWithFavorites)
	return router
}
