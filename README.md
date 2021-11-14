# DISYSMandatory2

### protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative mutex/mutex.proto

# How to run the program
- To run the program start up more then 1 node by using "go run ."
- each node takes an input string as "%id %ownport %portofnextnode %isLastnode" 
- the last node starts with the token  

# Below is input for 3 nodes

- 0 8080 8090 false
- 1 8090 8100 false
- 2 8100 8080 true

# Implementation
- This implementation releases the access to the critical section after 1 second 
- Checkout feature/CLIExit branch for a little neater CLI experience

