package handler

import (
	"JWT_auth/internal/model"
	"JWT_auth/internal/service"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	if err := h.service.Auth.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) SignUp(c *gin.Context) {
	var request model.SignUpRequest
	//parse request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//do hesh password
	bytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	//check user in db
	var user model.User
	user.Email = request.Email
	user.Password = string(bytes)
	id, err := h.service.Repository.GetUser(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	//create tokens
	t, rt, err := h.service.Auth.GenerateTokenPair(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"token": t, "refresh token": rt})
}
