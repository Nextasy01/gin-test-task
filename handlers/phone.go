package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Nextasy01/gin-test-task/db"
	"github.com/Nextasy01/gin-test-task/entity"
	"github.com/Nextasy01/gin-test-task/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PhoneInput struct {
	Phone       string `json:"phone" binding:"required,e164,max=12"`
	Description string `json:"description"`
	IsFax       *bool  `json:"is_fax" binding:"required"`
}

type PhoneHandler struct {
	phrepo db.PhoneRepository
}

func NewPhoneInput(phone, desc string, is_fax *bool) PhoneInput {
	return PhoneInput{phone, desc, is_fax}
}

func NewPhoneHandler(pr db.PhoneRepository) PhoneHandler {
	return PhoneHandler{pr}
}

func (ph *PhoneHandler) SavePhone(ctx *gin.Context) {
	uid, err := utils.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no user id"})
		return
	}

	is_fax := new(bool)
	*is_fax = ctx.GetBool("is_fax")

	input := NewPhoneInput(ctx.PostForm("phone"), ctx.PostForm("description"), is_fax)

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("couldn't bind JSON: %v", err)})
		return
	}

	err = ph.phrepo.CheckPhone(input.Phone)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("this phone number '%s' is already registered, %v", input.Phone, err)})
		return
	}

	p := entity.NewPhone()

	if p.ID, err = uuid.NewRandom(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := uuid.Parse(uid)
	p.UserId = id
	p.Number = input.Phone
	p.Description = input.Description
	p.IsFax = *input.IsFax

	err = ph.phrepo.SavePhone(*p)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("couldn't save this phone: %v", err)})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": fmt.Sprintf("successfully saved phone %s!", p.Number)})

}

func (ph PhoneHandler) GetPhone(ctx *gin.Context) {
	number := ctx.Query("q")
	if number == "" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "please provide the phone number to search"})
		return
	}
	if string(number[0]) != "+" {
		number = strings.Replace(number, " ", "", -1)
		number = fmt.Sprintf("+%s", number)
	}
	phones, err := ph.phrepo.GetPhone(number)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("couldn't find the phone: %v", err)})
		return
	}

	if len(phones) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("no phones with number '%s'", number)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"phones": phones})

}

func (ph *PhoneHandler) UpdatePhone(ctx *gin.Context) {

	phone_id := ctx.PostForm("phone_id")
	if phone_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "please provide the id of the phone"})
		return
	}

	is_fax := new(bool)
	*is_fax = ctx.GetBool("is_fax")

	phoneInput := NewPhoneInput(ctx.PostForm("phone"), ctx.PostForm("description"), is_fax)

	input := struct {
		PhoneId     string `json:"phone_id"`
		Phone       string `json:"phone" binding:"e164,max=12"`
		Description string `json:"description"`
		IsFax       bool   `json:"is_fax"`
	}{
		phone_id,
		phoneInput.Phone,
		phoneInput.Description,
		*phoneInput.IsFax,
	}
	err := ctx.Bind(&input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("couldn't bind the phone input: %v", err)})
		return
	}

	err = ph.phrepo.UpdatePhone(input.PhoneId, input.Phone, input.Description, input.IsFax)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("couldn't update the phone: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("successfully update the phone with id: %s", input.PhoneId)})
}

func (ph *PhoneHandler) DeletePhone(ctx *gin.Context) {
	phone_id := ctx.Param("phone_id")
	if phone_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "please specify the phone id you want to delete"})
		return
	}

	err := ph.phrepo.DeletePhone(phone_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("successfully deleted phone with id: '%s'", phone_id)})
}
