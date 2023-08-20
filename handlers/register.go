package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Nextasy01/gin-test-task/db"
	"github.com/Nextasy01/gin-test-task/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterInput struct {
	Login    string `json:"login" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Age      int    `json:"age" binding:"required"`
}

type RegisterHandler struct {
	urepo db.UserRepository
}

func NewRegisterHandler(ur db.UserRepository) RegisterHandler {
	return RegisterHandler{ur}
}

func NewRegisterInput(email, username, password string, age int) RegisterInput {
	return RegisterInput{email, password, username, age}
}

func (r *RegisterHandler) Register(ctx *gin.Context) {
	ageForm, _ := strconv.Atoi(ctx.PostForm("age"))
	log.Println(ageForm)
	input := NewRegisterInput(ctx.PostForm("login"), ctx.PostForm("name"), ctx.PostForm("password"), ageForm)

	err := ctx.ShouldBind(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := entity.NewUser()

	if u.ID, err = uuid.NewRandom(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u.Email = input.Login
	u.Name = input.Name
	u.Password = input.Password
	u.Age = ageForm

	err = r.urepo.SaveUser(*u)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("cannot register the user: %s", err.Error())})
		return
	}

	token, _, err := r.urepo.LoginCheck(u.Email, u.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.SetCookie("SESSTOKEN", token, 24*3600, "/", ctx.Request.URL.Hostname(), false, true)
	ctx.Header("Cache-Control", "no-cache, private, max-age=0")

	ctx.JSON(http.StatusAccepted, gin.H{"token": token})

}
