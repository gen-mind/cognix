package main

import (
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	"github.com/caarlos0/env/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

const (
	EnvK8s           = "k8s"
	EnvDocker        = "docker"
	dockerConfigPath = "/configmap"
)

type Config struct {
	Port        int    `env:"CONFIGMAP_GRPC_PORT,required"`
	Environment string `env:"CONFIGMAP_ENVIRONMENT,required"`
}

var Module = fx.Options(
	fx.Provide(
		readConfig,
		newServer),
	fx.Invoke(RunServer))

func readConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	utils.InitLogger(false)
	return cfg, nil

}

func newServer(cfg *Config) (proto.ConfigMapServer, error) {
	switch cfg.Environment {
	case EnvDocker:
		return NewDockerServer(dockerConfigPath)
	case EnvK8s:
		return NewK8SServer()
	}
	return nil, fmt.Errorf("invalid envieronment")
}

func RunServer(lc fx.Lifecycle, cfg *Config,
	grpcServer proto.ConfigMapServer) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		zap.S().Errorf("Failed to listen: %s ", err.Error())
		return err
	}

	service := grpc.NewServer()
	proto.RegisterConfigMapServer(service, grpcServer)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err = service.Serve(listener); err != nil {
					zap.S().Errorf("Failed to serve: %s", err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			service.Stop()
			return nil
		},
	})
	return nil
}
