# DISYSMandatory2

## protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative mutex/mutex.proto

--
To run the program start up 3 node with - go run . 
each node takes an input string as "%id %ownport %portofnextnode %isLastnode"  

below is input for 3 clients

0 8080 8090 false
1 8090 8100 false
2 8100 8080 true
