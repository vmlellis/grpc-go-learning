package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vmlellis/grpc-go-learning/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Calculator Client...")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")

	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 10,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.GetSumResult())
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{Number: 12000}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stram
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		fmt.Println(msg.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	numbers := []int32{1, 2, 3, 4, 7, 9, 12, 13}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v", err)
	}

	for _, number := range numbers {
		req := &calculatorpb.ComputeAverageRequest{Number: number}
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving reponse ComputeAverage: %v", err)
	}

	fmt.Printf("The Average is: %v\n", res.GetAverage())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream FindMaximum: %v", err)
		return
	}

	waitc := make(chan struct{})

	// send go routine
	go sendNumbersToFindMaximum(stream)

	// receive go routine
	go receiveNumbersFromFindMaximum(stream, waitc)

	<-waitc
}

func sendNumbersToFindMaximum(stream calculatorpb.CalculatorService_FindMaximumClient) {
	numbers := []int32{1, 5, 3, 6, 2, 20}
	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		req := &calculatorpb.FindMaximumRequest{
			Number: number,
		}
		stream.Send(req)
		time.Sleep(10 * time.Millisecond)
	}
	stream.CloseSend()
}

func receiveNumbersFromFindMaximum(stream calculatorpb.CalculatorService_FindMaximumClient, waitc chan<- struct{}) {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading server stream: %v", err)
			break
		}
		fmt.Printf("New maximum: %v\n", res.GetMaximum())
	}
	close(waitc)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Streaming RPC...")

	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println("Error message form server:", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
			return
		}

		log.Fatalf("Big Error calling SquareRoot: %v", err)
		return
	}
	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
}
