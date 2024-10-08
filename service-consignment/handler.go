package main

import (
	"context"

	pb "github.com/SoroushBeigi/microservice-go/service-consignment/proto/consignment"
	vesselProto "github.com/SoroushBeigi/microservice-go/service-vessel/proto/vessel"
	"github.com/pkg/errors"
)

type handler struct {
	repository
	vesselClient vesselProto.VesselService
}

func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	vesselResponse, err := s.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if vesselResponse == nil {
		return errors.New("error fetching vessel, returned nil")
	}

	if err != nil {
		return err
	}

	req.VesselId = vesselResponse.Vessel.Id

	if err = s.repository.Create(ctx, MarshalConsignment(req)); err != nil {
		return err
	}

	res.Created = true
	res.Consignment = req
	return nil
}

func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments, err := s.repository.GetAll(ctx)
	if err != nil {
		return err
	}
	res.Consignments = UnmarshalConsignmentCollection(consignments)
	return nil
}
