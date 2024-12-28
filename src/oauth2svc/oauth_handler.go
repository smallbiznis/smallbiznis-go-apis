package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/smallbiznis/oauth2-server/model"
	"github.com/smallbiznis/oauth2-server/service"
)

type OAuthHandler struct {
	validator    *validator.Validate
	oauthService service.IOAuthService
}

func NewOAuthHandler(
	validator *validator.Validate,
	oauthService service.IOAuthService,
) *OAuthHandler {
	return &OAuthHandler{
		validator,
		oauthService,
	}
}

func (e *OAuthHandler) HandleUserInfo(c *gin.Context) {
	e.oauthService.HandleUserInfo(c)
}

func (e *OAuthHandler) HandleIntrospect(c *gin.Context) {
	e.oauthService.HandleIntrospect(c)
}

func (e *OAuthHandler) HandleRevoke(c *gin.Context) {
	e.oauthService.HandleRevoke(c)
}

func (e *OAuthHandler) HandleRequestAuthorization(c *gin.Context) {
	e.oauthService.HandleRequestAuthorization(c)
}

func (e *OAuthHandler) HandleRequestToken(c *gin.Context) {
	ctx := c.Request.Context()
	var req model.TokenRequest
	if err := c.Bind(&req); err != nil {
		c.Error(err)
		return
	}

	clientId, clientSecret, ok := c.Request.BasicAuth()
	if ok {
		req.ClientID = clientId
		req.ClientSecret = clientSecret
	}

	resp, err := e.oauthService.HandleRequestToken(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (e *OAuthHandler) HandleGetKeys(c *gin.Context) {
	keys, err := e.oauthService.HandleGetKeys(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"keys": keys,
	})
}

func (e *OAuthHandler) HandleOpenIDConfiguration(c *gin.Context) {
	e.oauthService.HandleOpenIDConfiguration(c)
}
