package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/smallbiznis/oauth2-server/internal/pkg/cookie"
	"github.com/smallbiznis/oauth2-server/model"
	"github.com/smallbiznis/oauth2-server/repository"
	"github.com/smallbiznis/oauth2-server/service"
)

type AccountHandler struct {
	v       *validator.Validate
	sess    repository.ISessionRepository
	service service.IAccountService
}

func NewAccountHandler(
	v *validator.Validate,
	sess repository.ISessionRepository,
	service service.IAccountService,
) *AccountHandler {
	return &AccountHandler{
		v,
		sess,
		service,
	}
}

// HandleLookup
func (e *AccountHandler) HandleLookup(c *gin.Context) {
	var req model.RequestLookup
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

// HandleSignUp
func (e *AccountHandler) HandleSignUp(c *gin.Context) {
	ctx := c.Request.Context()

	var body model.RequestSignUp
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	if err := e.v.Struct(&body); err != nil {
		c.Error(err.(validator.ValidationErrors)[0])
		return
	}

	account, err := e.service.HandleSignUp(ctx, body)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}

// HandleSignInWithPassword
func (e *AccountHandler) HandleSignInWithPassword(c *gin.Context) {
	ctx := c.Request.Context()
	tenant := ctx.Value("tenant").(*model.Organization)

	var req model.RequestSignInWithPassword
	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		return
	}

	if err := e.v.Struct(&req); err != nil {
		c.Error(err.(validator.ValidationErrors)[0])
		return
	}

	req.Request = c.Request
	account, err := e.service.HandleSignInWithPassword(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(cookie.SetCookie(tenant.Name, account.SessionID))
	c.JSON(http.StatusOK, account)
}

// HandleSignInWithPhoneNumber
func (e *AccountHandler) HandleSignInWithPhoneNumber(c *gin.Context) {
	ctx := c.Request.Context()

	var req model.RequestSignInWithPhoneNumber
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := e.v.Struct(&req); err != nil {
		c.Error(err.(validator.ValidationErrors)[0])
		return
	}

	account, err := e.service.HandleSignInWithPhoneNumber(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}

// HandleSendVerificationCode
func (e *AccountHandler) HandleSendVerificationCode(c *gin.Context) {
	ctx := c.Request.Context()

	var req model.RequestSendVerificationCode
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := e.v.Struct(&req); err != nil {
		c.Error(err.(validator.ValidationErrors)[0])
		return
	}

	account, err := e.service.HandleSendVerificationCode(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, account)
}
