package main

import (
	"JWT_auth/configs"
	"JWT_auth/internal/handler"
	"JWT_auth/internal/repository"
	"JWT_auth/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	//init configs
	if err := configs.InitConfig(); err != nil {
		log.Fatal(err)
	}
	//db connection
	db, err := repository.NewDB(context.Background())
	if err != nil {
		log.Fatal("No database connection ")
	}
	//migration
	if err := repository.AutoMigration(viper.GetBool("db.migration.isAllowed")); err != nil {
		log.Fatal(err)
	}

	//init main components
	r := repository.NewRepository(db)
	s := service.NewService(r)
	h := handler.NewHandler(s)

	//init routes
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(h.AuthMiddleware)

	router.POST("/auth/signIn", h.SignIn)
	router.POST("/auth/signUp", h.SignUp)
	router.POST("/auth/update", h.TokenRefreshing)
	router.POST("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Hello": "World"})
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Not allowed request"})
	})

	//run server
	router.Run(fmt.Sprint(viper.GetString("host"), ":",
		viper.GetString("port")))

}
