package service

import (
	"github.com/kumparan/go-lib/redcachekeeper"
	"github.com/kumparan/kumnats"
	"gitlab.kumparan.com/yowez/skeleton-service/repository"
	"gitlab.kumparan.com/yowez/skeleton-service/worker"
)

// Service :nodoc:
type Service struct {
	nats        kumnats.NATS
	cacheKeeper redcachekeeper.Keeper
	helloRepo   repository.HelloRepository
	worker      worker.Worker
}

// RegisterNATS :nodoc:
func (s *Service) RegisterNATS(n kumnats.NATS) {
	s.nats = n
}

// GetNATS :nodoc:
func (s *Service) GetNATS() kumnats.NATS {
	return s.nats
}

// NewHelloService :nodoc:
func NewHelloService() *Service {
	return new(Service)
}

// RegisterHelloRepository :nodoc:
func (s *Service) RegisterHelloRepository(r repository.HelloRepository) {
	s.helloRepo = r
}

// RegisterCacheKeeper :nodoc:
func (s *Service) RegisterCacheKeeper(k redcachekeeper.Keeper) {
	s.cacheKeeper = k
}

// RegisterWorker :nodoc:
func (s *Service) RegisterWorker(w worker.Worker) {
	s.worker = w
}
