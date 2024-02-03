package _interface

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"test_ina_bank/internal/common/errs"
	"test_ina_bank/internal/domain/users"
	"test_ina_bank/internal/usecase/app"
	"test_ina_bank/pkg/telemetry"
	"test_ina_bank/pkg/utils"
)

type userServer struct {
	tracer *telemetry.OtelSdk
	app    app.UserApplication
}

func NewUserServer(tracer *telemetry.OtelSdk, app app.UserApplication) *userServer {
	return &userServer{tracer: tracer, app: app}
}

func (u *userServer) InsertUser(c *gin.Context) {
	ctx, span := u.tracer.Tracer.Start(c.Request.Context(), "http.insert_user")
	defer span.End()

	request := new(users.User)
	if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, utils.ParseMessage(err))
		return
	}

	err := u.app.Commands.InsertUserHandler.Handle(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusCreated, "success", nil))
}

func (u *userServer) ListUser(c *gin.Context) {
	ctx, span := u.tracer.Tracer.Start(c.Request.Context(), "http.list_user")
	defer span.End()

	out, err := u.app.Queries.ListUserHandler.Handle(ctx, &users.User{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", out))
}

func (u *userServer) GetUser(c *gin.Context) {
	ctx, span := u.tracer.Tracer.Start(c.Request.Context(), "http.get_user")
	defer span.End()

	request := new(users.UserId)
	request.Id, _ = strconv.Atoi(c.Param("id"))

	out, err := u.app.Queries.GetUserHandler.Handle(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", out))
}

func (u *userServer) UpdateUser(c *gin.Context) {
	ctx, span := u.tracer.Tracer.Start(c.Request.Context(), "http.update_user")
	defer span.End()

	request := new(users.User)
	if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, utils.ParseMessage(err))
		return
	}
	request.Id, _ = strconv.Atoi(c.Param("id"))

	err := u.app.Commands.UpdateUserHandler.Handle(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", nil))
}

func (u *userServer) DeleteUser(c *gin.Context) {
	ctx, span := u.tracer.Tracer.Start(c.Request.Context(), "http.delete_user")
	defer span.End()

	request := new(users.UserId)
	request.Id, _ = strconv.Atoi(c.Param("id"))

	err := u.app.Commands.DeleteUserHandler.Handle(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", nil))
}

func (u *userServer) Login(c *gin.Context) {
	ctx, span := u.tracer.Tracer.Start(c.Request.Context(), "http.login")
	defer span.End()

	request := new(users.Login)
	if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, utils.ParseMessage(errs.Invalidated.New(err.Error())))
		return
	}

	out, err := u.app.Queries.LoginHandler.Handle(ctx, request)
	if err != nil {
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", out))
}

func (u *userServer) RefreshToken(c *gin.Context) {
	ctx, span := u.tracer.Tracer.Start(c.Request.Context(), "http.refresh_token")
	defer span.End()

	request := new(users.RefreshToken)
	if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, utils.ParseMessage(errs.Invalidated.New(err.Error())))
		return
	}

	out, err := u.app.Queries.RefreshTokenHandler.Handle(ctx, request)
	if err != nil {
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", out))
}
