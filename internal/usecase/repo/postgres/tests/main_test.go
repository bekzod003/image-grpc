package tests

import (
	"context"
	"testing"

	"github.com/bekzod003/image-grpc/config"
	"github.com/bekzod003/image-grpc/pkg/database/client/postgresql"
)

var client postgresql.Client

func TestMain(m *testing.M) {
	cfg := config.GetConfig()
	var err error
	client, err = postgresql.NewClient(context.Background(), postgresql.ClientConfig{
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
		panic(err)
	}
	m.Run()
}
