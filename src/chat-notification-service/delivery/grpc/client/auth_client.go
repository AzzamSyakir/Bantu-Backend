package client

import (
	"bantu-backend/grpc/pb"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GatewayServiceClient struct {
	Client pb.GatewayServiceClient
}

func InitGatewayServiceClient(url string) GatewayServiceClient {
	cc, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := GatewayServiceClient{
		Client: pb.NewGatewayServiceClient(cc),
	}
	fmt.Println("init auth grpc service", url)
	return c
}
func (c *GatewayServiceClient) LogoutWithUserId(req *pb.ByUserId) (*pb.Empty, error) {
	resp, err := c.Client.LogoutWithUserId(context.Background(), req)
	if err != nil {
		log.Fatal("failed to LogoutWithUserId: %w", err)
	}
	return resp, nil
}
