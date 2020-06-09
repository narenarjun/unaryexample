package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// doServerStreaming(c)

	doClientStreaming(c)

}

func doClientStreaming(c greetpbgen.GreetServiceClient){
	fmt.Println("starting to do a client streaming RPC ....")

	requests := []*greetpbgen.LongGreetRequest{
		&greetpbgen.LongGreetRequest{
			Greeting: &greetpbgen.Greeting{
				FirstName: "naren",
			},
		},
		&greetpbgen.LongGreetRequest{
			Greeting: &greetpbgen.Greeting{
				FirstName: "arjun",
			},
		},
		&greetpbgen.LongGreetRequest{
			Greeting: &greetpbgen.Greeting{
				FirstName: "karna",
			},
		},
		&greetpbgen.LongGreetRequest{
			Greeting: &greetpbgen.Greeting{
				FirstName: "nakul",
			},
		},
		&greetpbgen.LongGreetRequest{
			Greeting: &greetpbgen.Greeting{
				FirstName: "kamsha",
			},
		},
	}
	stream , err := c.LongGreet(context.Background())
	if err != nil{
		log.Fatalf("Error while executing longGreet: %v\n", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err :=	stream.CloseAndRecv()
	
	if err != nil{
		log.Fatalf("Error while receiving response from longGreet: %v", err)

	}
	fmt.Printf("LongGreet Response: %v\n", res)

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
