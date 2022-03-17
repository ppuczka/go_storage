package go_storage

import (
	"context"
	"log"
	"go_storage/core"
)

type Server struct {
  UnimplementedKeyVauleServer
  TransactionLogger core.TransactionLogger
}

func (s *Server) Get(ctx context.Context, r *GetRequest) (*GetResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	value, err := core.Get(r.Key)
	s.TransactionLogger.ReadEvents()
	return &GetResponse{Value: value}, err
}

func (s *Server) Put(ctx context.Context, r *PutRequest) (*PutResponse, error) {
	log.Printf("Received PUT request key=%v, value=%v", r.Key, r.Value)

	err := core.Put(r.Key, r.Value)
	s.TransactionLogger.WritePut(r.Key, r.Value)
	return &PutResponse{}, err
}

func (s *Server) Delete(ctx context.Context, r *DeleteRequest) (*DeleteResponse, error) {
	log.Printf("Received DELETE request key=%v", r.Key)

	err := core.Delete(r.Key)
	s.TransactionLogger.WriteDelete(r.Key)
	return &DeleteResponse{},err
}