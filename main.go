package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"os"
	"strconv"
)

func main() {
	bcast_ingress_port, err := strconv.ParseUint(os.Getenv("GREENLINE_BCAST_INGRESS_PORT"), 0, 32)
	fmt.Println(bcast_ingress_port)
	fmt.Println("starting")
	ingress, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		fmt.Println("creating ingress socket: " + err.Error())
		os.Exit(1)
	}
	ingress.SetSubscribe("")
	defer ingress.Close()

	egress, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		fmt.Println("creating egress socket: " + err.Error())
		os.Exit(1)
	}
	egress.SetLinger(1)
	defer egress.Close()

	ingress.Bind("tcp://127.0.0.1:9002")
	egress.Bind("tcp://127.0.0.1:9003")

	poller := zmq.NewPoller()
	poller.Add(ingress, zmq.POLLIN)

	fmt.Println("greenline alive")
	for {
		sockets, _ := poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case ingress:
				for {
					msg, _ := s.Recv(0)
					if more, _ := s.GetRcvmore(); more {
						egress.Send(msg, zmq.SNDMORE)
					} else {
						egress.Send(msg, 0)
						break
					}
				}
			}
		}
	}
}
