package main

import (
	"context"
	"fmt"
	"io"
	"log"

	calculatorpbgen "github.com/narenarjun/unaryexample/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main(){
	fmt.Println("starting calculator client ....")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	defer conn.Close()

	c := calculatorpbgen.NewCalculatorServiceClient(conn)

	// doSum(c)

	doServerStreaming(c)
}


func doServerStreaming( c calculatorpbgen.CalculatorServiceClient){
	fmt.Println("starting to do a Prime Decomposition Server Streaming RPC ....")
	req := &calculatorpbgen.PrimeNumberDecompositionRequest{
		Number:24 ,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(),req)

	if err != nil {
		log.Fatalf("error while requesting calculator grpc: %v\n", err)

	}

	for {
		res, err := stream.Recv()
		if err == io.EOF{
			break
		}
		if err != nil {
			log.Fatalf("error occured while streaming: %v\n",err )
		}

		fmt.Println(res.GetPrimeFactor())

	}

}


func doSum( c calculatorpbgen.CalculatorServiceClient){
	fmt.Println("starting to do a Sum unary RPC ....")
	req := &calculatorpbgen.SumRequest{
		FirstNumber: 40,
		SecondNumber: 2,
	}
	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("error while requesting calculator grpc: %v\n", err)

	}

	log.Printf("Result form the Calculator: %v\n",res.SumResult)
}