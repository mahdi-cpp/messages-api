package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/application"
)

type ChatHandler struct {
	appManager *application.Manager
}

func NewChatHandler(appManager *application.Manager) *ChatHandler {
	return &ChatHandler{
		appManager: appManager,
	}
}

func (h *ChatHandler) Create(c *gin.Context) {

}

func (h *ChatHandler) GetFilter(c *gin.Context) {

}

func (h *ChatHandler) Update(c *gin.Context) {

}

func (h *ChatHandler) Delete(c *gin.Context) {

}
