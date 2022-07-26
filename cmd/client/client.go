package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/tuliofergulha/fullcycle-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Túlio Fergulha",
		Email: "tfergulha@gmail.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make to gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Túlio Fergulha",
		Email: "tfergulha@gmail.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make to gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}

		fmt.Println("Status: ", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "t1",
			Name:  "Túlio Fergulha",
			Email: "test@test.com.br",
		},
		&pb.User{
			Id:    "t2",
			Name:  "Túlio Fergulha 2",
			Email: "test2@test.com.br",
		},
		&pb.User{
			Id:    "t3",
			Name:  "Túlio Fergulha 3",
			Email: "test3@test.com.br",
		},
		&pb.User{
			Id:    "t4",
			Name:  "Túlio Fergulha 4",
			Email: "test4@test.com.br",
		},
		&pb.User{
			Id:    "t5",
			Name:  "Túlio Fergulha 5",
			Email: "test5@test.com.br",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "t1",
			Name:  "Túlio Fergulha",
			Email: "test@test.com.br",
		},
		&pb.User{
			Id:    "t2",
			Name:  "Túlio Fergulha 2",
			Email: "test2@test.com.br",
		},
		&pb.User{
			Id:    "t3",
			Name:  "Túlio Fergulha 3",
			Email: "test3@test.com.br",
		},
		&pb.User{
			Id:    "t4",
			Name:  "Túlio Fergulha 4",
			Email: "test4@test.com.br",
		},
		&pb.User{
			Id:    "t5",
			Name:  "Túlio Fergulha 5",
			Email: "test5@test.com.br",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}

			fmt.Printf("Recebendo user %v com status %v \n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
