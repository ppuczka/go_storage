package go_storage

import (
	"context"
	"log"
	"go_storage/helpers"
)

type Server struct {
  UnimplementedKeyVauleServer
  TransactionLogger helpers.TransactionLogger
}

func (s *Server) Get(ctx context.Context, r *GetRequest) (*GetResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	value, err := helpers.Get(r.Key)
	s.TransactionLogger.ReadEvents()
	return &GetResponse{Value: value}, err
}

func (s *Server) Put(ctx context.Context, r *PutRequest) (*PutResponse, error) {
	log.Printf("Received PUT request key=%v, value=%v", r.Key, r.Value)

	err := helpers.Put(r.Key, r.Value)
	s.TransactionLogger.WritePut(r.Key, r.Value)
	return &PutResponse{}, err
}

func (s *Server) Delete(ctx context.Context, r *DeleteRequest) (*DeleteResponse, error) {
	log.Printf("Received DELETE request key=%v", r.Key)

	err := helpers.Delete(r.Key)
	s.TransactionLogger.WriteDelete(r.Key)
	return &DeleteResponse{},err
}