package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"grpc-gateway/myservice/proto"
	"io"
	"log"
	"net"
	"os"
)

type server struct {
	proto.UnimplementedMyServiceServer
}

func (*server) Echo(ctx context.Context, in *proto.SimpleMessage) (*proto.SimpleMessage, error) {
	fmt.Println(in)
	return in, nil
}
func (*server) EchoUpload(stream proto.MyService_EchoUploadServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		err := errors.New("获取元数据失败")
		log.Println(err)
		return err
	}

	fileName := md["file_name"][0]
	fmt.Println("server receives: " + fileName)

	filePath := "myservice/server/upload" + fileName
	dst, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	defer dst.Close()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return err
		}

		dst.Write(req.Content[:req.Size])
	}

	stream.SendAndClose(&proto.UploadResponse{
		Path: filePath,
	})

	return err
}

var (
	port = flag.Int("port", 50051, "")
)

func main() {
	flag.Parse()
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	proto.RegisterMyServiceServer(s, &server{})
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
