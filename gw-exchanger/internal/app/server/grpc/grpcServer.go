package grpc

import (
	"exchanger/internal/app/service"
	"gorm.io/gorm"
	proto "proto-exchange/gen/exchange"
)

type Server struct {
	service.ExchangeService
	proto.ExchangeServiceServer
}

func NewServer(db *gorm.DB) *Server {
	return &Server{
		ExchangeService:       service.ExchangeService{Db: db},
		ExchangeServiceServer: nil,
	}
}
