package main

import (
	"context"
	"fmt"
	micro "github.com/micro/go-micro/v2"
	"log"
	"os"
	pb "shippy-service-consignment/example.com/project/protos/consignment"
	vesselProto "shippy-service-consignment/example.com/project/protos/vessel"
)

const (
	defaultHost = "datastore:27017"
)

func main() {
	// Set-up micro instance
	service := micro.NewService(
		micro.Name("shippy.service.consignment"),
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
	defer client.Disconnect(context.Background())

	consignmentCollection := client.Database("shippy").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}
	vesselClient := vesselProto.NewVesselService("shippy.service.client", service.Client())
	h := &handler{repository, vesselClient}

	// Register handlers
	pb.RegisterShippingServiceHandler(service.Server(), h)

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
