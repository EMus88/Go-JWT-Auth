package handler

import (
	"JWT_auth/internal/model"
	"JWT_auth/internal/service"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

const validRole = "user"

var middlewareExpn = []string{"/auth/signIn", "/auth/signUp", "/auth/update"}

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AuthMiddleware(c *gin.Context) {
	for _, val := range middlewareExpn {
		if c.Request.URL.Path == val {
			return
		}
	}
	authHeader := strings.Split(c.GetHeader("Authorization"), " ")
	if len(authHeader) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		c.Abort()
		return
	}
	bearerToken := authHeader[1]
	if _, _, err := h.service.ValidateToken(bearerToken, "access"); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		c.Abort()
		return
	}

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
	if (len(user.Password) < 7) || (len(user.Password) > 20) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password length will be from 7 to 15 simbols"})
		return
	}
	if len(user.Role) == 0 || user.Role != validRole {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not valid role"})
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
	//check user in db
	var user model.User
	user.Email = request.Email
	user.Password = h.service.Auth.HashingPassword(request.Password)
	id, role, err := h.service.Repository.GetUser(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	//create tokens
	t, rt, err := h.service.Auth.GenerateTokenPair(id, role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access token": t, "refresh token": rt})
}

func (h *Handler) TokenRefreshing(c *gin.Context) {
	var request model.UpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}
	id, role, err := h.service.ValidateToken(request.RefreshToken, "refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not valid refresh token"})
		return
	}
	t, rt, err := h.service.Auth.GenerateTokenPair(id, role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access token": t, "refresh token": rt})

}
