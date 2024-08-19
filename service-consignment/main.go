package main

import (
	pb "github.com/SoroushBeigi/microservice-go/service-consignment/proto/consignment"
	vesselProto "github.com/SoroushBeigi/microservice-go/service-vessel/proto/vessel"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"golang.org/x/net/context"
	"log"
)

type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo         Repository
	vesselClient vesselProto.VesselService
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	logger.Infof("Received CreateConsignment request: %v", req)
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	req.VesselId = vesselResponse.Vessel.Id

	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment
	logger.Infof("Sent CreateConsignment response: %v", res)
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	logger.Infof("Received GetConsignment request: %v", req)
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	logger.Infof("Sending GetConsignment response: %v", res)
	return nil
}

func main() {
	logger.DefaultLogger = logger.NewLogger(logger.WithLevel(logger.DebugLevel))
	repo := &ConsignmentRepository{}

	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselService("go.micro.srv.vessel", srv.Client())

	srv.Init()

	err := pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})
	if err != nil {
		log.Println(err)
		return
	}

	if err := srv.Run(); err != nil {
		log.Println(err)
	}
	log.Println("Server Started!")
}
