package main

import (
	"net"
	"time"

	"github.com/martingallagher/imgapi/service"
	"google.golang.org/grpc"
)

func getClient(config *service.Config) (*grpc.ClientConn, service.ImgAPIClient, error) {
	conn, err := grpc.Dial(
		config.Address,
		grpc.WithInsecure(),
		grpc.WithDialer(func(string, time.Duration) (net.Conn, error) {
			return net.Dial(config.Network, config.Address)
		}),
	)

	if err != nil {
		return nil, nil, err
	}

	return conn, service.NewImgAPIClient(conn), nil
}
