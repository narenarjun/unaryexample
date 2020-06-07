package main

import (
	"context"
	"fmt"
	"log"
	"net"

	calculatorpbgen "github.com/narenarjun/unaryexample/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {}

func (*server) Sum(ctx context.Context, req *calculatorpbgen.SumRequest)(*calculatorpbgen.SumResponse, error){
	fmt.Println("Sum function was invoked with %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	sum :=  firstNumber + secondNumber
	res := &calculatorpbgen.SumResponse{
		SumResult: sum,
	}


	return res, nil
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