package handlers

import (
	"fmt"
	"net/http"

	"github.com/Nextasy01/gin-test-task/db"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	urepo db.UserRepository
}

func NewUserHandler(ur db.UserRepository) UserHandler {
	return UserHandler{ur}
}

func (uh *UserHandler) GetUser(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "please specify the user name in url address"})
		return
	}

	users, err := uh.urepo.GetUsersByName(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if len(users) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("no users with name '%s'", name)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})

}
