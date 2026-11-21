package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/SergeyDolgushin/GRPC_TEST/pkg/note_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedNoteV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Note id: #{req.GetId()}")

	return &desc.GetResponse{
		Note: &desc.Note{
			Id: req.GetId(),
			Info: &desc.NoteInfo{
				Title:    "FFFF",
				Content:  "12312312",
				Author:   "sdasdad",
				IsPublic: true,
			},
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("filade to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{})

	log.Printf("server listening at %v}", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: #{err}")
	}
}
