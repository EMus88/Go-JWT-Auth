package handler

import (
	"JWT_auth/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

const validRole = "user"

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) Init() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())

	auth := router.Group("/auth")
	{
		auth.POST("/signUp", h.SignUp)
		auth.POST("/signIn", h.SignIn)
		auth.POST("/update", h.TokenRefreshing)
	}

	router.POST("/hello", h.AuthMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Hello": "World"})
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Not allowed request"})
	})
	return router
}
