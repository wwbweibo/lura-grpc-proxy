package main

import (
	"context"
	"log"
	"net"

	"github.com/wwbweibo/lura-grpc-proxy/cmd/test-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	proto.UnimplementedEchoServiceServer
}

func (srv Server) Get(ctx context.Context, req *proto.EchoGetRequest) (*proto.Response, error) {
	return &proto.Response{Resp: req.PathName + req.QueryName}, nil
}
func (srv Server) Post(ctx context.Context, req *proto.EchoPostRequest) (*proto.Response, error) {
	return &proto.Response{Resp: req.PathName + req.QueryName + req.BodyName}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	proto.RegisterEchoServiceServer(s, &Server{})
	reflection.Register(s)
	s.Serve(lis)
}
