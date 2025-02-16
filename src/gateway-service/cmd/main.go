package main

import (
	"bantu-backend/src/gateway-service/container"
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	fmt.Println("Gateway Services started.")

	webContainer := container.NewWebContainer()
	// grpc server
	go func() {
		grpcAddress := fmt.Sprintf(
			"%s:%s",
			"0.0.0.0",
			webContainer.Env.App.GatewayGrpcPort,
		)
		netListen, err := net.Listen("tcp", grpcAddress)
		if err != nil {
			log.Fatalf("failed to listen %v", err)
		}
		if err := webContainer.Grpc.Serve(netListen); err != nil {
			log.Fatalf("failed to serve %v", err.Error())
		}
	}()
	// http server
	address := fmt.Sprintf(
		"%s:%s",
		"0.0.0.0",
		webContainer.Env.App.GatewayHttpPort,
	)
	listenAndServeErr := http.ListenAndServe(address, webContainer.Route.Router)
	if listenAndServeErr != nil {
		log.Fatalf("failed to serve HTTP: %v", listenAndServeErr)
	}
	fmt.Println("Gateway Services finished.")
}
