package terminal_handler

import (
	"github.com/dvdxa/add-to-favorites/internal/services"
	"github.com/dvdxa/add-to-favorites/pkg/logger"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TerminalHandler struct {
	log                 logger.Logger
	terminalServicePort services.TerminalServicePort
}

func NewTerminalHandler(log logger.Logger, terminalServicePort services.TerminalServicePort) *TerminalHandler {
	return &TerminalHandler{
		log:                 log,
		terminalServicePort: terminalServicePort,
	}
}

type Request struct {
	TerminalID int    `json:"terminal_id"`
	IsFavorite string `json:"is_favorite" binding:"required"`
}

func (h *TerminalHandler) GetTerminalsWithFavorites(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok || userId == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "failed to get user_id",
		})
		return
	}
	userIdfloat := userId.(float64)
	userIdInt := int(userIdfloat)

	var body Request
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	if body.TerminalID == 0 && body.IsFavorite == "nil" {
		userTerminalsIDS, err := h.terminalServicePort.GetFavoriteTerminalIds(userIdInt)
		if err != nil {
			h.log.Errorf("failed to get user terminal ids: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"failed to get user terminal ids": err.Error(),
			})
			return
		}
		sortedTerminals, err := h.terminalServicePort.SortTerminals(userTerminalsIDS)
		if err != nil {
			h.log.Errorf("failed to sort terminals: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"failed to sort terminals": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, sortedTerminals)
		return
	}
	if body.IsFavorite == "true" {
		err = h.terminalServicePort.AddToFavorite(body.TerminalID, userIdInt)
		if err != nil {
			h.log.Errorf("failed to add to favorites: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"failed to add to favorites": err.Error(),
			})
			return
		}
		userTerminalsIDS, err := h.terminalServicePort.GetFavoriteTerminalIds(userIdInt)
		if err != nil {
			h.log.Errorf("failed to get user terminal ids: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"failed to get user terminal ids": err.Error(),
			})
			return
		}
		log.Println(userTerminalsIDS)
		sortedTerminals, err := h.terminalServicePort.SortTerminals(userTerminalsIDS)
		if err != nil {
			h.log.Errorf("failed to sort terminals: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"failed to sort terminals": err.Error(),
			})
			return
		}
		log.Println(sortedTerminals)
		c.JSON(http.StatusOK, sortedTerminals)
		return
	}
	if body.IsFavorite == "false" {
		err = h.terminalServicePort.RemoveFromFavoriteTerminal(body.TerminalID, userIdInt)
		if err != nil {
			h.log.Errorf("failed to remove terminal from favorites: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"failed to remove terminal from favorites": err.Error(),
			})
			return
		}
		userTerminalsIDS, err := h.terminalServicePort.GetFavoriteTerminalIds(userIdInt)
		if err != nil {
			h.log.Errorf("failed to get user terminal ids: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"failed to get user terminal ids": err.Error(),
			})
			return
		}
		sortedTerminals, err := h.terminalServicePort.SortTerminals(userTerminalsIDS)
		if err != nil {
			h.log.Errorf("failed to sort terminals: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"failed to sort terminals": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, sortedTerminals)
		return
	}
}
