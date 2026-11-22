package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/SergeyDolgushin/GRPC_TEST/internal/repository"
	"github.com/SergeyDolgushin/GRPC_TEST/internal/repository/note"
	desc "github.com/SergeyDolgushin/GRPC_TEST/pkg/note_v1"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedNoteV1Server
	noteRepository repository.NoteRepository
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.noteRepository.Create(ctx, req.GetInfo())
	if err != nil {
		return nil, err
	}

	log.Printf("inserted note with id: %d", id)
	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	noteObj, err := s.noteRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, title: %s", noteObj.Id, noteObj.Info.Title)

	return &desc.GetResponse{
		Note: noteObj,
	}, nil
}

func main() {

	// ctx := context.Background()
	dbDSN := readEnv()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("filade to listen: %v", err)
	}

	db, erdb := sql.Open("mysql", dbDSN)
	if erdb != nil {
		log.Fatalf("DB Open Error, %v", erdb)
	}
	defer db.Close()

	noteRepo := note.NewRepository(db)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{noteRepository: noteRepo})

	log.Printf("Server listening at adr: %v", lis.Addr().String())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func readEnv() string {
	envFile := flag.String("env", ".env", "path to .env file")
	log.Printf("File %v", envFile)
	flag.Parse()
	flag.CommandLine.Parse(os.Args[1:])

	err := godotenv.Load(*envFile)
	if err != nil {
		log.Printf("Warning: .env file not found at %s: %v", *envFile, err)
		// Продолжаем работу — используем системные переменные
	}

	dbDSN := os.Getenv("MYSQL_DSN")
	fmt.Println("Using DSN from:", *envFile)
	fmt.Println("DB DSN:", dbDSN)

	return dbDSN
}
