package routes

import (
	"github.com/Nextasy01/gin-test-task/db"
	"github.com/Nextasy01/gin-test-task/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var DB db.Database
var userRepo db.UserRepository
var phoneRepo db.PhoneRepository
var registerHandler handlers.RegisterHandler
var loginHandler handlers.LoginHandler
var userHandler handlers.UserHandler
var phoneHandler handlers.PhoneHandler

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	DB = db.NewDatabase()
	userRepo = db.NewUserRepository(&DB)
	phoneRepo = db.NewPhoneRepository(&DB)
	registerHandler = handlers.NewRegisterHandler(userRepo)
	loginHandler = handlers.NewLoginHandler(userRepo)
	userHandler = handlers.NewUserHandler(userRepo)
	phoneHandler = handlers.NewPhoneHandler(phoneRepo)

}

func PublicRoutes(g *gin.RouterGroup) {
	g.POST("/user/register", registerHandler.Register)
	g.POST("/user/auth", loginHandler.Login)
}

func PrivateRoutes(g *gin.RouterGroup) {
	g.GET("/user/:name", userHandler.GetUser)
	g.POST("/user/phone", phoneHandler.SavePhone)
	g.GET("/user/phone", phoneHandler.GetPhone)
	g.PUT("/user/phone", phoneHandler.UpdatePhone)
	g.DELETE("/user/phone/:phone_id", phoneHandler.DeletePhone)
}
