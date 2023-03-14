package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bekzod003/image-grpc/config"
	"github.com/bekzod003/image-grpc/genproto/file_processing"
	"github.com/bekzod003/image-grpc/internal/controller/grpc/server"
	"github.com/bekzod003/image-grpc/internal/controller/grpc/service"
	"github.com/bekzod003/image-grpc/internal/usecase"
	"github.com/bekzod003/image-grpc/internal/usecase/repo/postgres"
	"github.com/bekzod003/image-grpc/pkg/database/client/postgresql"
	"github.com/bekzod003/image-grpc/pkg/logger"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	log := logger.NewLogger(cfg.ServiceName, cfg.LoggerLevel)

	psqlClient, err := postgresql.NewClient(ctx, postgresql.ClientConfig{
		Login:    cfg.PostgreSQL.User,
		Password: cfg.PostgreSQL.Password,
		Host:     cfg.PostgreSQL.Host,
		Port:     cfg.PostgreSQL.Port,
		DBName:   cfg.PostgreSQL.DBName,
		PoolConfig: postgresql.PoolConfig{
			MaxConns:                 cfg.PostgreSQL.PoolConfig.MaxConns,
			MaxConnIdleMinutes:       cfg.PostgreSQL.PoolConfig.MaxConnIdleMinute,
			MaxConnLifetimeMinutes:   cfg.PostgreSQL.PoolConfig.MaxConnLifetimeMinute,
			HealthCheckPeriodMinutes: cfg.PostgreSQL.PoolConfig.HealthCheckPeriodMinute,
		},
	})
	if err != nil {
		log.Fatal("error while creating postgresql client", logger.Error(err))
		return
	}

	fileRepo := postgres.NewFileRepo(psqlClient)
	fileUseCase := usecase.NewFileProcessingUseCase(fileRepo, log)
	fileProcessingGRPCService := service.NewFileProcessingService(log, fileUseCase)

	s := server.NewGRPCServer(cfg, &log)

	listener, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("error while creating listener", logger.Error(err))
	}

	file_processing.RegisterFileProcessingServer(s, fileProcessingGRPCService)

	log.Info("service is running...", zap.String("port", cfg.GRPC.Port))
	if err = s.Serve(listener); err != nil {
		log.Fatal("error while serving", logger.Error(err))
		return
	}

	c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Info("Gracefully shutting down...")

	psqlClient.Close()
	s.GracefulStop()
}
