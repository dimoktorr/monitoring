package app

import (
	"github.com/caarlos0/env/v6"
	"github.com/dimoktorr/monitoring/internal/pkg/persistent/repository"
	"github.com/dimoktorr/monitoring/internal/pkg/persistent/storage"
	"time"
)

type Config struct {
	Prometheus Prometheus        `envPrefix:"PROMETHEUS_"`
	Service    Service           `envPrefix:"SERVICE_"`
	Tracing    Tracing           `envPrefix:"TRACING_"`
	Database   repository.Config `envPrefix:"DATABASE_"`
	Redis      storage.Config    `envPrefix:"REDIS_"`
}

type Service struct {
	Host                   string        `env:"HOST" envDefault:"127.0.0.1"`
	GRPCPort               string        `env:"GRPC_PORT" envDefault:"8090"`
	HTTPPort               string        `env:"HTTP_PORT" envDefault:"8080"`
	ShutdownContextTimeout time.Duration `env:"SHUTDOWN_CONTEXT_TIMEOUT_DURATION" envDefault:"5s"`
}

type Prometheus struct {
	Host string `env:"HOST" validate:"required"`
	Port string `env:"PORT" validate:"required"`
}

type Tracing struct {
	ImsSystemID int64  `env:"IMS_SYSTEM_ID" envDefault:"5737" validate:"required"`
	SolObjectID int64  `env:"SOL_OBJECT_ID" envDefault:"1926407" validate:"required"`
	URL         string `env:"URL" validate:"required"`
	Environment string `env:"ENVIRONMENT" validate:"required"`
	ServiceName string `env:"SERVICE_NAME" validate:"required"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
