package main

import (
	"context"
	"fmt"
	"io"
	"log"

	greetpbgen "github.com/narenarjun/unaryexample/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello, !! i'm the client ")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpbgen.NewGreetServiceClient(conn)

	// doUnary(c)

	doServerStreaming(c)

}

func doServerStreaming(c greetpbgen.GreetServiceClient){
	fmt.Println("starting to do a server streaming RPC ....")

req := &greetpbgen.GreetManyTimesRequest{
	Greeting: &greetpbgen.Greeting{
		FirstName: "Naren",
		LastName: "Dev",
	},
}

	resStream , err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error occured in streaming %v", err)
	}
for  {
	
	msg, err := resStream.Recv()
	if err == io.EOF{
		// we've reached the end of the stream 
		break
	}
	if err != nil{
		log.Fatalf("error while reading the stream: %v", err)
	}

	log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
}

}



func doUnary(c greetpbgen.GreetServiceClient) {
	fmt.Println("starting to do a Unary RPC ....")
	req := &greetpbgen.GreetRequest{
		Greeting: &greetpbgen.Greeting{
			FirstName: "naren",
			LastName:  "Arjun",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error while logging Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}
