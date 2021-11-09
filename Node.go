package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ThomasITU/DISYSMandatory2/mutex"
	gRPC "google.golang.org/grpc"
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
	nextNodePort     int
	port         int
}


func main() {
	//get input id, ownport, next port
	node := node{id: 0, state: false, nextNodePort: 8090, port: 8080}

	go listen(node.port)

	ctx := context.Background()
	conn, err := gRPC.Dial("localhost:"+strconv.Itoa(node.nextNodePort), gRPC.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to: %s", strconv.Itoa(node.port))
	}
	c := mutex.NewMutexServiceClient(conn)


	if(node.id == 0 && node.state == true) {
		writeToLog(node.id,logFileName)
		c.Enter(ctx,node,conn)
	}


	
	// broadcast til de andre noder grpc.dial
	// grpc.send conn.send

}

func clientEnter(){
	conn, err := gRPC.Dial("localhost:"+strconv.Itoa(node.nextNodePort), gRPC.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to: %s", strconv.Itoa(node.port))
	}
	c := mutex.NewMutexServiceClient(conn)

	response := c.Enter(&mutex.{})
}

func (s *server) Enter(ctx context.Context, node *node) (*mutex.Response, error) {
	if (node.state){
		writeToLog(node.id, logFileName)
	}
	time.Sleep(1 * time.Second)
	clientEnter()

}

func (s *server) Exit(ctx context.Context)

func writeToLog(nodeID int, logName string) {
	f, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("Node nr. %s has entered the critical section at time: %v", nodeID, timestamp)
}

func listen(port int) {
	lis, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Could not listen to %s", port)
	}

	grpcServer := gRPC.NewServer()
	mutex.RegisterMutexServiceServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve on ")
	}
}
