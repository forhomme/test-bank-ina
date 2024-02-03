package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"test_ina_bank/config"
	_interface "test_ina_bank/internal/interface"
	"test_ina_bank/internal/service"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/middleware"
	"test_ina_bank/pkg/sqlhandler"
	"test_ina_bank/pkg/telemetry"
)

const (
	DBInitFile = "./migration/init.sql"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := baselogger.NewLogger()
	config.InitConfig(ctx, logger)

	tracer := telemetry.NewOtelSdk(config.Config)

	dBHandler := sqlhandler.NewSqlHandler(logger, config.Config)
	err := InitDB(dBHandler, DBInitFile)
	if err != nil {
		logger.Fatalf("Database init: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(middleware.DefaultStructuredLogger())
	route.Use(gin.Recovery())
	route.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:  []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))
	route.Use(secure.New(secure.DefaultConfig()))

	userApps := service.NewUserService(logger, dBHandler, tracer)
	userServer := _interface.NewUserServer(tracer, userApps)

	taskApps := service.NewTaskService(logger, dBHandler, tracer)
	taskServer := _interface.NewTaskServer(tracer, taskApps)

	api := route.Group("/api")
	api.POST("/login", userServer.Login)
	api.POST("/refresh-token", userServer.RefreshToken)
	api.GET("/user", userServer.ListUser)
	api.GET("/user/:id", userServer.GetUser)
	api.POST("/user", userServer.InsertUser)
	api.PUT("/user/:id", userServer.UpdateUser)
	api.DELETE("/user/:id", userServer.DeleteUser)

	taskRoute := api.Group("/task")
	taskRoute.Use(middleware.JwtAuthMiddleware())
	taskRoute.GET("", taskServer.ListTask)
	taskRoute.GET("/:id", taskServer.GetTask)
	taskRoute.POST("", taskServer.InsertTask)
	taskRoute.PUT("/:id", taskServer.UpdateTask)
	taskRoute.DELETE("/:id", taskServer.DeleteTask)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Config.Server.Port),
		Handler: route,
	}
	ctx = handleShutdown(ctx, logger, dBHandler, srv, tracer)
	go func() {
		// service connections
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()
}

// InitDB Method for initializing Database.
func InitDB(sqlHandler sqlhandler.SqlHandler, dbfilepath string) error {

	query, err := os.ReadFile(filepath.Clean(dbfilepath))
	if err != nil {
		return err
	}
	if err = sqlHandler.MultiExec(string(query)); err != nil {
		return err
	}
	return nil
}

func handleShutdown(ctx context.Context, logger *baselogger.Logger, sql sqlhandler.SqlHandler, server *http.Server, tracer *telemetry.OtelSdk) context.Context {
	ctx, done := context.WithCancel(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer done()

		<-quit
		signal.Stop(quit)
		close(quit)

		if err := sql.Close(); err != nil {
			logger.Errorf("could not gracefully shutdown the sql connection")
		} else {
			logger.Info("sql connection is shutting down")
		}

		if err := server.Shutdown(ctx); err != nil {
			logger.Errorf("could not gracefully shutdown the api server")
		} else {
			logger.Info("api server is shutting down")
		}

		if err := tracer.Shutdown(ctx); err != nil {
			logger.Errorf("could not gracefully shutdown the tracer")
		} else {
			logger.Info("tracer is shutting down")
		}
	}()
	return ctx
}
