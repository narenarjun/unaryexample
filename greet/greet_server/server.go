package main

import (
	"context"
	"fmt"
	"io"
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

func (*server) LongGreet(stream  greetpbgen.GreetService_LongGreetServer) error{
	fmt.Println("LongGreet function was invoked with a streaming request ")

	result := "hello"
	for {

	req, err := stream.Recv()

	if err == io.EOF{
		//  finished receiving all the client requests
	return	stream.SendAndClose(&greetpbgen.LongGreetResponse{
			Result: result,
		})
	}
	if err != nil{
		log.Fatalf("Error while reading client stream %v", err)
	}
	firstName :=	req.GetGreeting().GetFirstName()
	result += " " +firstName + "!"
	}
	
}


func (*server) 	GreetEveryone(stream greetpbgen.GreetService_GreetEveryoneServer) error{
	fmt.Println("GreetEveryone function was invoked with a streaming request ")

	for{
		req, err := stream.Recv()
		if err == io.EOF{
			return nil
		}

		if err != nil{
			log.Fatalf("Error while reading client stream request: %v ", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "hello " + firstName + "!"
{
		err := stream.Send(
			&greetpbgen.GreetEveryoneResponse{
				Result: result,},
		)

		if err != nil {
			log.Fatalf("Error while sending data to client: %v ", err)
			return err
		}
	}

	}
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
