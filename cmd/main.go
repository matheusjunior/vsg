package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"matheus.com/vgs/configs"
	"matheus.com/vgs/internal/controller"
	"matheus.com/vgs/internal/engine"
	"matheus.com/vgs/internal/logger"
	"matheus.com/vgs/internal/messaging"
	"matheus.com/vgs/internal/repo"
	"matheus.com/vgs/internal/sender"
	"matheus.com/vgs/internal/service"
)

func bootstrapServer() *echo.Echo {
	server := echo.New()
	server.Use(middleware.RemoveTrailingSlash())
	server.Use(middleware.Recover())
	server.Use(middleware.Gzip())
	return server
}

func bootstrapMetrics(server *echo.Echo) *echo.Echo {
	metrics := echo.New()
	metrics.HideBanner = true
	prom := prometheus.NewPrometheus("echo", nil)
	server.Use(prom.HandlerFunc)
	prom.SetMetricsPath(metrics)
	return metrics
}

func main() {
	userQueueUrl := configs.GetUserQueueUrl()
	voucherQueueUrl := configs.GetVoucherQueueUrl()

	userPublisher := messaging.NewSQSPublisher(userQueueUrl)
	voucherListener := messaging.NewSQSListener(userQueueUrl)
	voucherPublisher := messaging.NewSQSPublisher(voucherQueueUrl)
	notificationListener := messaging.NewSQSListener(voucherQueueUrl)
	engine := engine.NewEvaluatorEngine()

	mongoClient := repo.NewMongoConnection()
	postgresConnection := repo.NewPostgresConnection()
	userRepo := repo.NewUserRepo(postgresConnection)
	voucherRepo := repo.NewVoucherRepo(mongoClient)
	userSvc := service.NewUserService(userRepo)
	userMatcherSvc := service.NewUserMatcherService(userRepo, userPublisher, engine)
	voucherCreatorSvc := service.NewVoucherCreatorService(voucherListener, voucherPublisher, voucherRepo)
	emailSender := sender.NewEmailSender()
	notifierSvc := service.NewNotifierService(notificationListener, emailSender)
	userController := controller.NewUserController(userSvc)
	userMatcherController := controller.NewUserMatcherController(userMatcherSvc)

	voucherCreatorSvc.Start()
	notifierSvc.Start()
	server := bootstrapServer()
	server.GET("/users/:id", userController.GetById)
	server.POST("/users", userController.Create)
	server.POST("/vouchers/generate", userMatcherController.Match)
	metrics := bootstrapMetrics(server)

	go func() {
		logger.Logger().Info("Starting HTTP server")
		err := server.Start(configs.GetHost())
		logger.Logger().Fatal(errors.Wrap(err, "failed to start server"))
	}()

	go func() {
		logger.Logger().Info("Starting HTTP metrics server")
		err := metrics.Start(configs.GetMetricsHost())
		logger.Logger().Fatal(errors.Wrap(err, "failed to start metrics server"))
	}()

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	switch <-signalChannel {
	case os.Interrupt:
		logger.Logger().Warn("SIGINT signal received, stopping...")
	case syscall.SIGTERM:
		logger.Logger().Warn("SIGTERM signal received, stopping...")
	}

	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		logger.Logger().Fatal(err)
	}
}
