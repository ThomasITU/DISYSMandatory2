package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ThomasITU/DISYSMandatory2/mutex"
	"google.golang.org/grpc"
)

// lav en node struct med info og de andre
const (
	logFileName = "serverLog"
)

type server struct {
	mutex.UnimplementedMutexServiceServer
}

type node struct {
	id           int
	state        bool
	nextNodePort int
	port         int
}

func start(id int, port int, nextPort int,hasToken bool) {
	//get input id, ownport, next port
	node := node{id: id, state: hasToken, nextNodePort: nextPort, port: port}

	go listen(node.port)

	ctx := context.Background()

	if node.state == true {
		writeToLog(node.id, logFileName)
		PassToken(ctx, &node)
	}

	fmt.Scanln()

	for {

	}

	// broadcast til de andre noder grpc.dial
	// grpc.send conn.send

}

func PassToken(ctx context.Context,node *node) {
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(node.nextNodePort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to: %s", strconv.Itoa(node.port))
	}
	nextNode := mutex.NewMutexServiceClient(conn)

	if _,err  := nextNode.Token(ctx, &mutex.EmptyRequest{}); err != nil {}
}

func (s *server) Token(ctx context.Context, node *node) (*mutex.EmptyResponse, error) {
	if node.state {
		writeToLog(node.id, logFileName)
		node.state = false
	}
	time.Sleep(1 * time.Second)
	PassToken(ctx, node)
	return &mutex.EmptyResponse{}, nil
}

func writeToLog(nodeID int, logName string) {
	f, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("Node nr. %s has entered the critical section at time: %v", nodeID)
}

func listen(port int) {
	lis, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Could not listen to %s", port)
	}

	grpcServer := grpc.NewServer()
	mutex.RegisterMutexServiceServer(grpcServer, &erver{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve on ")
	}
}
