package main

import (
	"log"
	"net"

	"github.com/junaid9001/zvynt/auth/config"
	"github.com/junaid9001/zvynt/auth/db"
	grpcservices "github.com/junaid9001/zvynt/auth/grpc-services"
	"github.com/junaid9001/zvynt/proto/auth"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	postgresDB := db.InitDB(cfg)

	store := db.NewStore(postgresDB)

	lis, err := net.Listen("tcp", ":"+cfg.APP_PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()

	srv := grpcservices.NewAuthServiceServer(store, cfg)

	auth.RegisterAuthServiceServer(grpcServer, srv)

	grpcServer.Serve(lis)

}
