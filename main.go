package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"os"
	"strconv"
)

const (
	BCAST_INGRESS_PORT   = "GL_BCAST_INGRESS_PORT"
	BCAST_EGRESS_PORT    = "GL_BCAST_EGRESS_PORT"
	REQREP1_INGRESS_PORT = "GL_REQREP1_INGRESS_PORT"
	REQREP1_EGRESS_PORT  = "GL_REQREP1_EGRESS_PORT"
	REQREP2_INGRESS_PORT = "GL_REQREP2_INGRESS_PORT"
	REQREP2_EGRESS_PORT  = "GL_REQREP2_EGRESS_PORT"
)

func main() {
	env := FatalGetenv(BCAST_INGRESS_PORT)
	bcast_ingress_port, err := strconv.ParseUint(env, 0, 32)
	env = FatalGetenv(BCAST_EGRESS_PORT)
	bcast_egress_port, err := strconv.ParseUint(env, 0, 32)
	env = FatalGetenv(REQREP1_INGRESS_PORT)
	reqrep1_ingress_port, err := strconv.ParseUint(env, 0, 32)
	env = FatalGetenv(REQREP1_EGRESS_PORT)
	reqrep1_egress_port, err := strconv.ParseUint(env, 0, 32)
	env = FatalGetenv(REQREP2_INGRESS_PORT)
	reqrep2_ingress_port, err := strconv.ParseUint(env, 0, 32)
	env = FatalGetenv(REQREP2_EGRESS_PORT)
	reqrep2_egress_port, err := strconv.ParseUint(env, 0, 32)

	fmt.Println(bcast_ingress_port)
	fmt.Println(bcast_egress_port)
	fmt.Println(reqrep1_ingress_port)
	fmt.Println(reqrep1_egress_port)
	fmt.Println(reqrep2_ingress_port)
	fmt.Println(reqrep2_egress_port)

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

func FatalGetenv(env string) string {
	_env := os.Getenv(env)
	if len(_env) == 0 {
		log.Fatal("no " + env + " is set")
	}
	return _env
}
