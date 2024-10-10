package configmap

import (
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

type Config struct {
	Host string `env:"CONFIGMAP_GRPC_HOST,required"`
	Port int    `env:"CONFIGMAP_GRPC_PORT,required"`
}

type ClientBuilder struct {
	cfg    *Config
	client proto.ConfigMapClient
	mx     sync.Mutex
}

func NewClientBuilder(cfg *Config) *ClientBuilder {
	return &ClientBuilder{
		cfg: cfg,
		mx:  sync.Mutex{},
	}
}
func (e *ClientBuilder) Client() (proto.ConfigMapClient, error) {
	e.mx.Lock()
	defer e.mx.Unlock()
	if e.client == nil {
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials())}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port), dialOptions...)
		if err != nil {
			return nil, err
		}
		e.client = proto.NewConfigMapClient(conn)
	}
	return e.client, nil
}
