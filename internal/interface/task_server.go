package _interface

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"strings"
	"test_ina_bank/config"
	"test_ina_bank/internal/common/errs"
	"test_ina_bank/internal/domain/tasks"
	"test_ina_bank/internal/usecase/app"
	"test_ina_bank/pkg/middleware"
	"test_ina_bank/pkg/telemetry"
	"test_ina_bank/pkg/utils"
)

type taskServer struct {
	tracer *telemetry.OtelSdk
	app    app.TaskApplication
}

func NewTaskServer(tracer *telemetry.OtelSdk, app app.TaskApplication) *taskServer {
	return &taskServer{tracer: tracer, app: app}
}

// InsertTask godoc
// @Summary InsertTask API
// @Description Insert New Task
// @Tags tasks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param task body tasks.Task true "Task Data"
// @Success 201 {object} utils.ResponseMessage{status=int,message=string,data=nil}
// @Failure 400 {object} utils.ErrorMessage{status=int,message=string,data=nil}
// @Failure 500 {object} utils.ErrorMessage{status=int,message=string,data=nil}
// @Router /task [post]
func (t *taskServer) InsertTask(c *gin.Context) {
	ctx, span := t.tracer.Tracer.Start(c.Request.Context(), "http.insert_task")
	defer span.End()

	task := new(tasks.TaskModel)
	if err := c.ShouldBindWith(&task, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, utils.ParseMessage(errs.Invalidated.New(err.Error())))
		return
	}

	authUser, err := middleware.GetAuthUser(c)
	if err != nil {
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
		return
	}

	if task.UserId == 0 {
		task.UserId = authUser.UserId
	}

	request, err := task.UnmarshalTask()
	if err != nil {
		if strings.Contains(err.Error(), config.Config.ApplicationMessage[config.Config.General.CurrentLanguage].InvalidStatus.Key) {
			c.JSON(errs.GetHttpCode(err), utils.ParseMessage(errs.Invalidated.New(config.Config.ApplicationMessage[config.Config.General.CurrentLanguage].InvalidStatus.Message)))
			return
		}
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
		return
	}

	err = t.app.Commands.InsertTaskHandler.Handle(ctx, request)
	if err != nil {
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusCreated, utils.ParseResponse(http.StatusCreated, "success", nil))
}

// ListTask godoc
// @Summary ListTask API
// @Description List All Task
// @Tags tasks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} utils.ResponseMessage{status=int,message=string,data=tasks.ListTask}
// @Failure 400 {object} utils.ErrorMessage{status=int,message=string,data=nil}
// @Failure 500 {object} utils.ErrorMessage{status=int,message=string,data=nil}
// @Router /task [get]
func (t *taskServer) ListTask(c *gin.Context) {
	ctx, span := t.tracer.Tracer.Start(c.Request.Context(), "http.list_task")
	defer span.End()

	out, err := t.app.Queries.ListTaskHandler.Handle(ctx, &tasks.Task{})
	if err != nil {
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", out))
}

// GetTask godoc
// @Summary GetTask API
// @Description Get Task by Id
// @Tags tasks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int64 true "Task ID"
// @Success 200 {object} utils.ResponseMessage{status=int,message=string,data=tasks.Task}
// @Failure 400 {object} utils.ErrorMessage{status=int,message=string,data=nil}
// @Failure 500 {object} utils.ErrorMessage{status=int,message=string,data=nil}
// @Router /task/{id} [get]
func (t *taskServer) GetTask(c *gin.Context) {
	ctx, span := t.tracer.Tracer.Start(c.Request.Context(), "http.get_task")
	defer span.End()

	request := new(tasks.TaskId)
	request.Id, _ = strconv.Atoi(c.Param("id"))

	out, err := t.app.Queries.GetTaskHandler.Handle(ctx, request)
	if err != nil {
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", out))
}

// UpdateTask godoc
// @Summary UpdateTask API
// @Description Update Task by Id
// @Tags tasks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int64 true "Task ID"
// @Param task body tasks.Task true "Task Data"
// @Success 200 {object} utils.ResponseMessage{status=int,message=string,data=nil}
// @Failure 400 {object} utils.ErrorMessage{status=int,message=string,data=nil}
// @Failure 500 {object} utils.ErrorMessage{status=int,message=string,data=nil}
// @Router /task/{id} [put]
func (t *taskServer) UpdateTask(c *gin.Context) {
	ctx, span := t.tracer.Tracer.Start(c.Request.Context(), "http.update_task")
	defer span.End()

	task := new(tasks.TaskModel)
	if err := c.ShouldBindWith(&task, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, utils.ParseMessage(errs.Invalidated.New(err.Error())))
		return
	}
	authUser, err := middleware.GetAuthUser(c)
	if err != nil {
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
	}

	if task.UserId == 0 {
		task.UserId = authUser.UserId
	}

	request, err := task.UnmarshalTask()
	if err != nil {
		if strings.Contains(err.Error(), config.Config.ApplicationMessage[config.Config.General.CurrentLanguage].InvalidStatus.Key) {
			c.JSON(errs.GetHttpCode(err), utils.ParseMessage(errs.Invalidated.New(config.Config.ApplicationMessage[config.Config.General.CurrentLanguage].InvalidStatus.Message)))
			return
		}
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
		return
	}
	request.Id, _ = strconv.Atoi(c.Param("id"))

	err = t.app.Commands.UpdateTaskHandler.Handle(ctx, request)
	if err != nil {
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", nil))
}

// DeleteTask godoc
// @Summary DeleteTask API
// @Description Delete Task by Id
// @Tags tasks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int64 true "Task ID"
// @Success 200 {object} utils.ResponseMessage{status=int,message=string,data=nil}
// @Failure 400 {object} utils.ErrorMessage{status=int,message=string,data=nil}
// @Failure 500 {object} utils.ErrorMessage{status=int,message=string,data=nil}
// @Router /task/{id} [delete]
func (t *taskServer) DeleteTask(c *gin.Context) {
	ctx, span := t.tracer.Tracer.Start(c.Request.Context(), "http.delete_task")
	defer span.End()

	request := new(tasks.TaskId)
	request.Id, _ = strconv.Atoi(c.Param("id"))

	err := t.app.Commands.DeleteTaskHandler.Handle(ctx, request)
	if err != nil {
		c.JSON(errs.GetHttpCode(err), utils.ParseMessage(err))
		return
	}

	c.JSON(http.StatusOK, utils.ParseResponse(http.StatusOK, "success", nil))
}
