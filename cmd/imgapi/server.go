package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/martingallagher/imgapi/service"
	"google.golang.org/grpc"
)

func startServer(config *service.Config) {
	imgapiServer, err := service.NewServer(config)

	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen(config.Network, config.Address)

	if err != nil {
		log.Fatal(err)
	}

	defer l.Close() // nolint: errcheck

	server := grpc.NewServer()

	defer server.Stop()

	service.RegisterImgAPIServer(server, imgapiServer)

	log.Println("Server starting")

	go func() {
		log.Fatal(server.Serve(l))
	}()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt)

	<-sigs

	log.Println("Server shutdown")
}
