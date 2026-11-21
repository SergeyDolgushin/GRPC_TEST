package main

import (
	"context"
	"log"
	"time"

	desc "github.com/SergeyDolgushin/GRPC_TEST/pkg/note_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	addr   = "raspberrypi5.local:50051"
	noteID = 12
)

func main() {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fialed to connect: %v", err)
	}

	defer conn.Close()

	c := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: noteID})
	if err != nil {
		log.Fatalf("failed get note by id: %v", err)
	}

	log.Printf("Note: \n%+v", r.Note)
}
