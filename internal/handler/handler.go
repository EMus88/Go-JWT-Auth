package handler

import (
	"JWT_auth/internal/model"
	"JWT_auth/internal/service"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AuthMiddleware(c *gin.Context) {
	c.Next()
}

//Registration
func (h *Handler) SignIn(c *gin.Context) {
	var user model.User
	//parse request
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//validation request
	if ok, _ := govalidator.ValidateStruct(user); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}
	//save in db
	if err := h.service.SignIn(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
