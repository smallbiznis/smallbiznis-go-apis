package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/smallbiznis/oauth2-server/model"
	"github.com/smallbiznis/oauth2-server/service"
)

type ApplicationHandler struct {
	v       *validator.Validate
	service service.IApplicationService
}

func NewApplicationHandler(
	v *validator.Validate,
	service service.IApplicationService,
) *ApplicationHandler {
	return &ApplicationHandler{
		v,
		service,
	}
}

func (a *ApplicationHandler) HandleList(c *gin.Context) {
	ctx := c.Request.Context()
	tenant := ctx.Value("tenant").(*model.Organization)

	var req model.RequestListApplication
	if err := c.Bind(&req); err != nil {
		c.Error(err)
		return
	}

	req.OrganizationID = tenant.ID
	apps, count, err := a.service.HandleList(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, model.ResponseListApplication{
		TotalData: count,
		Data:      apps,
	})
}

func (a *ApplicationHandler) HandleCreate(c *gin.Context) {
	ctx := c.Request.Context()
	tenant := ctx.Value("tenant").(*model.Organization)

	var req model.RequestCreateApplication
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := a.v.Struct(&req); err != nil {
		c.Error(err.(validator.ValidationErrors)[0])
		return
	}

	req.OrganizationID = tenant.ID
	app, err := a.service.HandleCreate(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, app)
}

func (a *ApplicationHandler) HandleGet(c *gin.Context) {
	ctx := c.Request.Context()
	tenant := ctx.Value("tenant").(*model.Organization)

	app, err := a.service.HandleGet(ctx, &model.Application{
		OrganizationID: tenant.ID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, app)
}

func (a *ApplicationHandler) HandleUpdate(c *gin.Context) {
	c.Status(http.StatusServiceUnavailable)
}

func (a *ApplicationHandler) HandleDelete(c *gin.Context) {
	c.Status(http.StatusServiceUnavailable)
}
