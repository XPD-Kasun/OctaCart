package rest

import (
	"net/http"
	"octacart/internal/auth"

	"github.com/gin-gonic/gin"
)

func FormLogin(ctx *gin.Context) {

	email := ctx.Request.PostFormValue("email")
	password := ctx.Request.PostFormValue("password")

	s := auth.FormAuthSvc{}

	user, err := s.Authenticate(ctx, email, password)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, user.Name)

}
