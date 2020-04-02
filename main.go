package main

import (
	"context"
	"net"
	"sync"

	pb "github.com/fahruroze/kumbah/proto/pengiriman"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Pengiriman) (*pb.Pengiriman, error)
}

type Repository struct {
	mu          sync.RWMutex
	pengiriman2 []*pb.Pengiriman
}

//Create pengiriman baru

func (repo *Repository) Create(pengiriman *pb.Pengiriman) (*pb.Pengiriman, error) {
	repo.mu.Lock()
	updated := append(repo.pengiriman2, pengiriman)
	repo.pengiriman2 = updated
	repo.mu.Unlock()
	return pengiriman, nil
}

//Service harus mengimplementasikan semua
//yg digenerate protobuf

type service struct {
	repo repository
}

//CreatePengiriman

func (s *service) CreatePengiriman(ctx context.Context, req *pb.Pengiriman) (*pb.Response, error) {
	//simpan pengiriman kita
	pengiriman, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	//sesuaikan response dengan protobuf
	return &pb.Response{Created: true, Pengiriman: pengiriman}, nil
}

func main() {
	repo := &Repository{}

	//setup GRPC-Server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to Listen: %v", err)
	}
	s := grpc.NewServer()

	//daftarkan service kita ke gRPC server, itu akan menyingkronkan dengan interface yg digenerate protobuf
	pb.RegisterPengirimanServiceServer(s, &service{repo})

	//daftarkan reflection service di gRPC juga
	reflection.Register(s)

	//Jalankan port
	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
