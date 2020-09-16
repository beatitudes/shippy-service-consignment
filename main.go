package main

import (
	"context"
	"log"
	"sync"

	pb "github.com/beatitudes/shippy-service-consignment/proto/consignment"
	vesselProto "github.com/beatitudes/shippy-service-vessel/proto/vessel"
	"github.com/micro/go-micro/v2"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	mu          sync.RWMutex
	consigments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consigments, consignment)
	repo.consigments = updated
	repo.mu.Unlock()
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consigments
}

type vesselService struct {
	repo         repository
	vasselClient vesselProto.VesselService
}

func (s *vesselService) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	vesselResponse, err := s.vasselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return err
	}
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)

	req.VesselId = vesselResponse.Vessel.Id

	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consigment = consignment
	return nil
}

func (s *vesselService) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	service := micro.NewService(
		micro.Name("shippy.service.consignment"),
	)

	service.Init()

	vesselClient := vesselProto.NewVesselService("shippy.service.vessel", service.Client())

	if err := pb.RegisterShippingServiceHandler(service.Server(), &vesselService{repo, vesselClient}); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
