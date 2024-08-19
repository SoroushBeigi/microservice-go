package main

import (
	"context"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"

	pb "github.com/SoroushBeigi/microservice-go/service-consignment/proto/consignment"
	vesselProto "github.com/SoroushBeigi/microservice-go/service-vessel/proto/vessel"
)

const (
	defaultHost = "datastore:27017"
)

func main() {
	// Set-up micro instance
	service := micro.NewService(
		micro.Name("service.consignment"),
	)

	service.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.Background())

	consignmentCollection := client.Database("mymongo").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}
	vesselClient := vesselProto.NewVesselService("service.client", service.Client())
	h := &handler{repository, vesselClient}

	registerErr := pb.RegisterShippingServiceHandler(service.Server(), h)
	if registerErr != nil {
		logger.Fatal(registerErr)
		return
	}

	if err := service.Run(); err != nil {
		logger.Error(err)
	}
}
