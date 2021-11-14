# DISYSMandatory2

# How to run the program

- To run the program start up more then 1 node by using "go run ."
- each node takes an input string as "%id %ownport %portofnextnode %isLastnode" 
- the last node starts with the token  

# Below is input for 3 nodes

- 0 8080 8090 false
- 1 8090 8100 false
- 2 8100 8080 true

# Implementation

- To release access to the critical section input a string that is not the empty string 
