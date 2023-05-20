package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-gateway/myservice/proto"
	"log"
	"net/http"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "")
)

func main() {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := proto.RegisterMyServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalln(err)
	}

	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatalln(err)
	}
}

//func pingHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
//
//}
