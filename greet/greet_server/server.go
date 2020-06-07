package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	greetpbgen "github.com/narenarjun/unaryexample/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpbgen.GreetRequest) (*greetpbgen.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "hello" + firstName
	res := &greetpbgen.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpbgen.GreetManyTimesRequest, stream greetpbgen.GreetService_GreetManyTimesServer)  error{
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)
	firstName := req.Greeting.GetFirstName()
	for i := 0; i < 10; i++ {
		result := "hello" + firstName + "number" + strconv.Itoa(i)
		res := &greetpbgen.GreetManyTimesResponse{
			Result: result ,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func main() {
	fmt.Println("hello world")

	// ! 50051 is the default port for grpc
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failes to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpbgen.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
