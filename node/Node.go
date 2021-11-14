package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/ThomasITU/DISYSMandatory2/mutex"
	"google.golang.org/grpc"
)

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

	var id, port, nextPort int
	var hasToken bool
	fmt.Scanln(&id, &port, &nextPort, &hasToken)
	node := node{id: id, state: hasToken, nextNodePort: nextPort, port: port}

	server := Server{this: node}
	go listen(server.this.port, &server)

	if server.this.state {
		startTokenRing(&server)
	}

	for {
		input := waitForInput(&server.this)
		if len(input) > 0 {
			AccessWanted(&server)
			for server.this.state {
			}
		}
	}
}

func startTokenRing(s *Server) {
	s.this.state = false
	go s.Token(context.Background(), &mutex.EmptyRequest{})
}

func AccessWanted(s *Server) {
	s.this.state = true
	fmt.Println("Waiting to gain access")
}

func PassToken(node *node) {
	ctx := context.Background()
	address := fmt.Sprintf("localhost:%v", node.nextNodePort)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to: %s", strconv.Itoa(node.port))
	}
	nextNode := mutex.NewMutexServiceClient(conn)

	if _, err := nextNode.Token(ctx, &mutex.EmptyRequest{}); err != nil {
		log.Println(err)
	}
}

func (s *Server) Token(ctx context.Context, empty *mutex.EmptyRequest) (*mutex.EmptyResponse, error) {
	if s.this.state {
		enterMsg := fmt.Sprintf("Node: %v has entered the critical section", s.this.id)
		writeToLog(enterMsg, logFileName)

		leaveMsg := fmt.Sprintf("Node: %v has left the critical section", s.this.id)
		fmt.Println("To leave the section input any string that is not the empty string")
		for {
			input := waitForInput(&s.this)
			if len(input) > 0 {
				writeToLog(leaveMsg, logFileName)
				s.this.state = false
				break
			}
		}

	}

	PassToken(&s.this)
	return &mutex.EmptyResponse{}, nil
}

func waitForInput(n *node) string {
	var input string
	fmt.Printf("nodeID: %v - ", n.id)
	fmt.Scanln(&input)
	return input
}

func writeToLog(msg string, logName string) {
	file, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Printf(msg)
}

func listen(port int, s *Server) {
	lis, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Could not listen to %v", port)
	}

	grpcServer := grpc.NewServer()
	mutex.RegisterMutexServiceServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve on ")
	}
}
