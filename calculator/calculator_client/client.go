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

	// doServerStreaming(c)
	 doClientStreaming(c)
}

func doClientStreaming( c calculatorpbgen.CalculatorServiceClient){
	fmt.Println("starting to do a compute average Clinet Streaming RPC ....")

stream, err :=	c.ComputeAverage(context.Background())
	if err != nil{
		log.Fatalf("error while opening stream: %v", err)
	}

	numbers := []int32{2,6,12,26,42}
	for _, number := range numbers {
		fmt.Printf("sending number %v\n", number)
		stream.Send(&calculatorpbgen.ComputeAverageRequest{
			Number: number,
		})
	}
{
	res, err := stream.CloseAndRecv()

	if err != nil{
		log.Fatalf("Error while receiving Response: %v", err)
	}

	fmt.Printf("the average is %v", res.GetAverage())
}

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