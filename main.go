package main

import (
	"database/sql"
	"log"

	"github.com/sRRRs-7/GachaPon/api"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
	"github.com/sRRRs-7/GachaPon/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBdriver, config.DBsource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(conn)

	runGinServer(config, store)

	// go runGatewayServer(config, store)
	// runGrpcServer(config, store)

}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}

// func runGrpcServer(config utils.Config, store db.Store) {
// 	server, err := gapi.NewServer(config, store)
// 	if err != nil {
// 		log.Fatal("cannot create server", err)
// 	}

// 	grpcServer := grpc.NewServer()

// 	protobuf.RegisterBankAppServer(grpcServer, server)
// 	reflection.Register(grpcServer)

// 	listener, err := net.Listen("tcp", config.GrpcServerAddress)
// 	if err != nil {
// 		log.Fatal("cannot create listener:", err)
// 	}

// 	log.Printf("start gRPC server at %s", listener.Addr().String())
// 	err = grpcServer.Serve(listener)
// 	if err != nil {
// 		log.Fatal("cannot start gRPC server:", err)
// 	}
// }

// func runGatewayServer(config util.Config, store db.Store) {
// 	server, err := gapi.NewServer(config, store)
// 	if err != nil {
// 		log.Fatal("cannot create server", err)
// 	}

// 	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
// 		MarshalOptions: protojson.MarshalOptions{
// 			UseProtoNames: true,
// 		},
// 		UnmarshalOptions: protojson.UnmarshalOptions{
// 			DiscardUnknown: true,
// 		},
// 	})

// 	grpcMux := runtime.NewServeMux(jsonOption)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	err = protobuf.RegisterBankAppHandlerServer(ctx, grpcMux, server)
// 	if err != nil {
// 		log.Fatal("cannot register handler server:", err)
// 	}

// 	mux := http.NewServeMux()
// 	mux.Handle("/", grpcMux)

// 	listener, err := net.Listen("tcp", config.HttpServerAddress)
// 	if err != nil {
// 		log.Fatal("cannot create listener:", err)
// 	}

// 	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
// 	err = http.Serve(listener, mux)
// 	if err != nil {
// 		log.Fatal("cannot start HTTP gateway server;", err)
// 	}
// }


