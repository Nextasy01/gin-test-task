package handlers

import (
	"fmt"
	"net/http"

	"github.com/Nextasy01/gin-test-task/db"
	"github.com/Nextasy01/gin-test-task/entity"
	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Login    string `json:"login" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginHandler struct {
	urepo db.UserRepository
}

func NewLoginHandler(ur db.UserRepository) LoginHandler {
	return LoginHandler{ur}
}

func NewLoginInput(login, password string) LoginInput {
	return LoginInput{login, password}
}

func (lh *LoginHandler) Login(ctx *gin.Context) {
	input := NewLoginInput(ctx.PostForm("login"), ctx.PostForm("password"))

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("couldn't bind json: %v", err)})
		return
	}

	u := entity.NewUser()

	u.Email = input.Login
	u.Password = input.Password

	token, userId, err := lh.urepo.LoginCheck(u.Email, u.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("couldn't get token: %v", err)})
		return
	}
	ctx.SetCookie("SESSTOKEN", token, 24*3600, "/", ctx.Request.URL.Hostname(), false, true)
	ctx.Header("Cache-Control", "no-cache, private, max-age=0")
	ctx.JSON(http.StatusAccepted, gin.H{"message": fmt.Sprintf("successfully logged in as %s!", u.Email),
		"user_id": userId,
		"token":   token})
}
