package main

import (
	"context"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"

	pb "github.com/SoroushBeigi/microservice-go/service-vessel/proto/vessel"
)

func main() {
	service := micro.NewService(
		micro.Name("shippy.service.vessel"),
	)

	service.Init()

	uri := os.Getenv("DB_HOST")

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			logger.Error(err)
		}
	}(client, context.Background())

	vesselCollection := client.Database("shippy").Collection("vessels")
	repository := &MongoRepository{vesselCollection}

	h := &handler{repository}

	if err := pb.RegisterVesselServiceHandler(service.Server(), h); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
