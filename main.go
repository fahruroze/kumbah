package main

import (
	"sync"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Pengiriman) (*pb.Pengiriman, error)
}

type Repository struct {
	mu          sync.RWMutex
	pengirimans []*pb.Pengiriman
}

func (repo *Repository) Create(pengiriman) {

}
