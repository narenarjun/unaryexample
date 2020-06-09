package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	calculatorpbgen "github.com/narenarjun/unaryexample/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {}

func (*server) Sum(ctx context.Context, req *calculatorpbgen.SumRequest)(*calculatorpbgen.SumResponse, error){
	fmt.Printf("Sum function was invoked with %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	sum :=  firstNumber + secondNumber
	res := &calculatorpbgen.SumResponse{
		SumResult: sum,
	}


	return res, nil
}

func (*server) 	PrimeNumberDecomposition(req *calculatorpbgen.PrimeNumberDecompositionRequest,stream calculatorpbgen.CalculatorService_PrimeNumberDecompositionServer) error{
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v", req)

	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpbgen.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased %v\n", divisor)
		}
	}
	return nil 

}
 
func (*server) ComputeAverage(stream calculatorpbgen.CalculatorService_ComputeAverageServer) error{
	fmt.Println("Received computeAverage RPC ")
	sum := int32(0)
	count := 0

	for {
	 req, err := stream.Recv()
	 if err == io.EOF {
		 average := float64(sum) / float64(count)
		 return stream.SendAndClose(&calculatorpbgen.ComputeAverageResponse{
			 Average: average,
			 })
	 }

	 if err != nil {
		 log.Fatalf("Error while reading client stream: %v", err)
		 return err
	 }
	 sum += req.GetNumber()
	 count++


	}
}

func (*server) FindMaximum(stream calculatorpbgen.CalculatorService_FindMaximumServer) error{
	fmt.Printf("Received FindMaximum RPC")

	maximum := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		number := req.GetNumber()
		if number > maximum{
			maximum = number
		sendErr :=	stream.Send(&calculatorpbgen.FindMaximumResponse{
				Maximum: maximum ,
			})
		if sendErr != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}	
		}
	}
}

func main(){
	fmt.Println("starting calculator rpc...")

	listn, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil{
		log.Fatalf("Failed to connect %v\n", err)
	}

	s := grpc.NewServer()
	calculatorpbgen.RegisterCalculatorServiceServer(s,&server{})

	if err := s.Serve(listn); err != nil{
		log.Fatalf("failed to serve : %v\n", err)
	}
}