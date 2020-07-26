package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/vmlellis/grpc-go-learning/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Printf("Sum function was invoked with %v", req)
	result := req.GetFirstNumber() + req.GetSecondNumber()
	res := &calculatorpb.SumResponse{SumResult: result}
	return res, nil
}

func main() {
	fmt.Println("Calculator Server...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
