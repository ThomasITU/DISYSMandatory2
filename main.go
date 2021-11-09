package main

import(
	Node "github.com/ThomasITU/DISYSMandatory2/node"
)

var (
	ports []int
)

func main() {
	ports = make([]int, 5)
	ports = append(ports, 8080, 8090, 8100, 8110, 8120)
	for i := 0; i < 5; i++ {
		if i%5 == 0 {
			go startNode(i, ports[i], ports[0], true)
		} else {
			go startNode(i, ports[i], ports[i+1], false)
		}
	}

}

func startNode(id int, port int, nextPort int, hasToken bool) {
	go Node.start(id, port, nextPort, hasToken)	
	
}
