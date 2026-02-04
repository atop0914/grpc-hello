package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "grpc-hello/proto/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:8080", "the address to connect to")
	name = flag.String("name", "World", "name to greet")
)

func main() {
	flag.Parse()

	// Set up a connection to the server
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Test basic greeting
	r, err := client.SayHello(ctx, &pb.HelloRequest{NameTest: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.TestMessage)

	// Test multiple greetings
	multiReq := &pb.HelloMultipleRequest{
		Names:         []string{"Alice", "Bob", "Charlie"},
		CommonMessage: "Welcome to our service!",
	}
	multiResp, err := client.SayHelloMultiple(ctx, multiReq)
	if err != nil {
		log.Printf("could not get multiple greetings: %v", err)
	} else {
		log.Printf("Got %d greetings", len(multiResp.Greetings))
		for _, greeting := range multiResp.Greetings {
			log.Printf("  - %s", greeting.TestMessage)
		}
	}

	// Test statistics
	statsReq := &pb.GreetingStatsRequest{}
	statsResp, err := client.GetGreetingStats(ctx, statsReq)
	if err != nil {
		log.Printf("could not get stats: %v", err)
	} else {
		log.Printf("Total requests: %d, Unique names: %d", statsResp.TotalRequests, statsResp.UniqueNames)
	}
}