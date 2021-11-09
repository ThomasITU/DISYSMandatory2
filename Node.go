package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/ThomasITU/DISYSMandatory2/mutex"
	gRPC "google.golang.org/grpc"
)

// lav en node struct med info og de andre
const (
	logFileName = "serverLog"
)

const (
	RELEASED int = 0
	WANTED       = 1
	HELD         = 2
)

type server struct {
	mutex.UnimplementedMutexServiceServer
}

type node struct {
	id           int
	state        int
	ports        []int
	port         int
	timeStamp    int
	requestQueue chan request
}

type request struct {
	id        int
	timestamp int
}

func main() {
	node := node{id: 0, ports: getAllPorts(), port: getAllPorts()[0], timeStamp: 0}

	go listen(node.port)

	request := mutex.RequestCriticalSection{Id: int32(node.id)}
	ctx := context.Background()

	// broadcast til de andre noder grpc.dial
	// grpc.send conn.send

}

func getAllPorts() []int {
	s := make([]int, 5)

	file, err := os.Open("ports.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {
		s[i], err = strconv.Atoi(scanner.Text())
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return s
}

func requestAccess(ctx context.Context, request *mutex.RequestCriticalSection, node *node) {
	for otherNode := range node.ports {

		conn, err := gRPC.Dial("localhost:"+strconv.Itoa(otherNode), gRPC.WithInsecure())
		if err != nil {

		}
		c := mutex.NewMutexServiceClient(conn)

		c.Enter(ctx, request)
	}
}

func (s *server) Enter(ctx context.Context, otherNodeRequest *mutex.RequestCriticalSection, node *node) (*mutex.Response, error) {
	id := otherNodeRequest.GetId()
	timeStamp := otherNodeRequest.GetTimestamp()

	if node.state == 2 || node.state == 1 && node.timeStamp < int(timeStamp) {
		req := request{id: int(id), timestamp: int(timeStamp)}
		node.requestQueue <- req
	}else{
		response := mutex.Response{node.id, }
		return 
	}

}

func writeToLog(nodeID int, timestamp int, logName string) {
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
