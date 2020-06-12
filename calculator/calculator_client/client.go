package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	calculatorpbgen "github.com/narenarjun/unaryexample/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	//  doClientStreaming(c)

	// doBiDiStreaming(c)

	doErrorUnary(c)
}

func doBiDiStreaming( c calculatorpbgen.CalculatorServiceClient){
	fmt.Println("starting to do a FindMaximum BiDi Streaming RPC ....")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream and calling findMaximm: %v", err)
	}

	waitc := make(chan struct{})

	// goroutine for sending

	go func ()  {
		numbers := []int32{1,6,4,9,12,2,69}

		for _, number := range numbers {
			fmt.Printf("sending number: %v\n",number)
			stream.Send(&calculatorpbgen.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()
	}()

	// goroutine for receive

	go func() {
		for {
		res,err := stream.Recv()
		if err == io.EOF{
			break
		}
		
		if err != nil{
			log.Fatalf("Problem while reading client: %v", err)
			break 
		}
		maximum := res.GetMaximum()
		log.Printf("received a new maximum of... : %v\n", maximum)
	}
		close(waitc)
	}()
	<-waitc
}

func doClientStreaming( c calculatorpbgen.CalculatorServiceClient){
	fmt.Println("starting to do a compute average Client Streaming RPC ....")

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



func doErrorUnary( c calculatorpbgen.CalculatorServiceClient){
	fmt.Println("starting to do a SquareRoot unary RPC ....")
	
	//  correct call
	doSquareRootCall(c,144)

	// error call
	doSquareRootCall(c,-16)



}

func doSquareRootCall (c calculatorpbgen.CalculatorServiceClient, n int32){
	res, err :=	c.SquareRoot(context.Background(),&calculatorpbgen.SquareRootRequest{
		Number: n,
	})

	if err != nil {
	 respErr, ok :=	status.FromError(err)
	 if ok {
		//  error thrown by grpc server
		fmt.Printf("Error message from Server: %v\n",respErr.Message())
		fmt.Println(respErr.Code())

		if respErr.Code() == codes.InvalidArgument{
			fmt.Println("A negative number was sent")
			return
		}
	 } else {
		 log.Fatalf("Big Strange Error when calling SquareRoot: %v", err)
		 return
	 }

	}

	fmt.Printf("Result of square root is %v : %v \n", n, res.GetNumberRoot())

}