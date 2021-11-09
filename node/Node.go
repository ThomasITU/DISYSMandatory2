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

type Server struct {
	mutex.UnimplementedMutexServiceServer
	this node
}

type node struct {
	id           int
	state        bool
	nextNodePort int
	port         int
}

func main() {
	//get input id, ownport, next port

	var id, port, nextPort int
	var hasToken bool
	fmt.Scanln(&id, &port, &nextPort, &hasToken)
	node := node{id: id, state: hasToken, nextNodePort: nextPort, port: port}
	fmt.Printf("node id: %v, node port: %v, nextNodePort: %v, state: %t", node.id, node.port, node.nextNodePort, node.state)

	go listen(node.port)

	ctx := context.Background()

	if node.state == true {
		writeToLog(node.id, logFileName)
		PassToken(ctx, &node)
	}

	fmt.Scanln()
}

func PassToken(ctx context.Context, node *node) {
	//address       = "localhost:8080"
	address := fmt.Sprintf("localhost:%v", node.nextNodePort)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to: %s", strconv.Itoa(node.port))
	}
	nextNode := mutex.NewMutexServiceClient(conn)

	if _, err := nextNode.Token(ctx, &mutex.EmptyRequest{}); err != nil {
		log.Println(err)
	} else {
		log.Println("No errors")
	}
}

func (s *Server) Token(ctx context.Context, empty *mutex.EmptyRequest) (*mutex.EmptyResponse, error) {
	if s.this.state {
		writeToLog(s.this.id, logFileName)
		s.this.state = false
	}
	log.Printf("I'm node id: %v", s.this.id)
	//fmt.Printf()
	time.Sleep(1 * time.Second)
	PassToken(ctx, &s.this)
	return &mutex.EmptyResponse{}, nil
}

func writeToLog(nodeID int, logName string) {
	f, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("Node nr. %v has entered the critical section", nodeID)
}

func listen(port int) {
	lis, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Could not listen to %v", port)
	}

	grpcServer := grpc.NewServer()
	mutex.RegisterMutexServiceServer(grpcServer, &Server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve on ")
	}
}
