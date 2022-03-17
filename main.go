package main

import (
	"log"
	"go_storage/config"
	"go_storage/core"
	"go_storage/transact"
	"go_storage/frontend"	
)


func main() {
	properties, err := config.NewConfigFilePropertiesLoader("config.yml")
	if err != nil {
		log.Fatal("error while loading application properties: ", err)
	}
	tl, err := transact.NewTransactionLogger(properties)
	if err != nil {
		log.Fatal("err")
	}	
	store := core.NewKeyValueStore(tl)
	store.Restore()
	
	fe, err := frontend.NewFrontEnd(properties.AppConnfig().FrondEndType)
	if err != nil {
		log.Fatal("err")
	}	
	log.Fatal(fe.Start(store))

	// PROTOBUF STUFF
	
	// server := grpc.NewServer()
	// pb.RegisterKeyVauleServer(server, &pb.Server{TransactionLogger: logger})
	// lis, err := net.Listen("tcp", ":50051")
	// log.Printf("GRPC server listening on 50051 port")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	
	// if err := server.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }	
}