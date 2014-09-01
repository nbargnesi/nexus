package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"os"
	"strconv"
)

const (
	BCAST_INGRESS_PORT = "GL_BCAST_INGRESS_PORT"
	BCAST_EGRESS_PORT  = "GL_BCAST_EGRESS_PORT"
	RR1_INGRESS_PORT   = "GL_RR1_INGRESS_PORT"
	RR1_EGRESS_PORT    = "GL_RR1_EGRESS_PORT"
	RR2_INGRESS_PORT   = "GL_RR2_INGRESS_PORT"
	RR2_EGRESS_PORT    = "GL_RR2_EGRESS_PORT"
)

func main() {
	env := Getenv(BCAST_INGRESS_PORT)
	bcast_ingress_port := AsPort(env)
	env = Getenv(BCAST_EGRESS_PORT)
	bcast_egress_port := AsPort(env)
	env = Getenv(RR1_INGRESS_PORT)
	rr1_ingress_port := AsPort(env)
	env = Getenv(RR1_EGRESS_PORT)
	rr1_egress_port := AsPort(env)
	env = Getenv(RR2_INGRESS_PORT)
	rr2_ingress_port := AsPort(env)
	env = Getenv(RR2_EGRESS_PORT)
	rr2_egress_port := AsPort(env)

	fmt.Println("starting")

	// CREATE EACH SOCKET...
	sub_ingress := NewSocket(zmq.SUB)
	sub_ingress.SetSubscribe("")
	defer sub_ingress.Close()
	pub_egress := NewSocket(zmq.PUB)
	pub_egress.SetLinger(1)
	defer pub_egress.Close()

	rr1_ingress := NewSocket(zmq.ROUTER)
	defer rr1_ingress.Close()
	rr1_egress := NewSocket(zmq.DEALER)
	defer rr1_egress.Close()

	rr2_ingress := NewSocket(zmq.ROUTER)
	defer rr2_ingress.Close()
	rr2_egress := NewSocket(zmq.DEALER)
	defer rr2_egress.Close()

	// ... AND BIND
	Bind(sub_ingress, "tcp", "127.0.0.1", bcast_ingress_port)
	Bind(pub_egress, "tcp", "127.0.0.1", bcast_egress_port)
	Bind(rr1_ingress, "tcp", "127.0.0.1", rr1_ingress_port)
	Bind(rr1_egress, "tcp", "127.0.0.1", rr1_egress_port)
	Bind(rr2_ingress, "tcp", "127.0.0.1", rr2_ingress_port)
	Bind(rr2_egress, "tcp", "127.0.0.1", rr2_egress_port)

	poller := zmq.NewPoller()
	poller.Add(sub_ingress, zmq.POLLIN)
	poller.Add(rr1_ingress, zmq.POLLIN)
	poller.Add(rr2_ingress, zmq.POLLIN)
	poller.Add(rr1_egress, zmq.POLLIN)
	poller.Add(rr2_egress, zmq.POLLIN)

	fmt.Println("greenline alive")
	for {
		sockets, _ := poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case sub_ingress:
				log.Println("processing broadcast message")
				for {
					msg, err := s.Recv(0)
					if err != nil {
						log.Fatalf("broadcast more: %s", err.Error())
					}
					more, err := s.GetRcvmore()
					if err != nil {
						log.Fatalf("broadcast recv more: %s", err.Error())
					}
					if more {
						pub_egress.Send(msg, zmq.SNDMORE)
					} else {
						pub_egress.Send(msg, 0)
						break
					}
				}
			case rr1_ingress:
				log.Println("processing rr1 request")
				for {
					msg, err := s.Recv(0)
					if err != nil {
						log.Fatalf("rr1 ingress: %s", err.Error())
					}
					more, err := s.GetRcvmore()
					if err != nil {
						log.Fatalf("rr1 ingress recv more: %s", err.Error())
					}
					if more {
						rr1_egress.Send(msg, zmq.SNDMORE)
					} else {
						rr1_egress.Send(msg, 0)
						break
					}
				}
			case rr2_ingress:
				for {
					msg, err := s.Recv(0)
					if err != nil {
						log.Fatalf("rr2 ingress: %s", err.Error())
					}
					more, err := s.GetRcvmore()
					if err != nil {
						log.Fatalf("rr2 ingress recv more: %s", err.Error())
					}
					if more {
						rr2_egress.Send(msg, zmq.SNDMORE)
					} else {
						rr2_egress.Send(msg, 0)
						break
					}
				}
			case rr1_egress:
				log.Println("processing rr1 response")
				for {
					msg, err := s.Recv(0)
					if err != nil {
						log.Fatalf("rr1 egress: %s", err.Error())
					}
					more, err := s.GetRcvmore()
					if err != nil {
						log.Fatalf("rr1 egress recv more: %s", err.Error())
					}
					if more {
						rr1_ingress.Send(msg, zmq.SNDMORE)
					} else {
						rr1_ingress.Send(msg, 0)
						break
					}
				}
			case rr2_egress:
				for {
					msg, err := s.Recv(0)
					if err != nil {
						log.Fatalf("rr2 egress: %s", err.Error())
					}
					more, err := s.GetRcvmore()
					if err != nil {
						log.Fatalf("rr2 egress recv more: %s", err.Error())
					}
					if more {
						rr2_ingress.Send(msg, zmq.SNDMORE)
					} else {
						rr2_ingress.Send(msg, 0)
						break
					}
				}
			}
		}
	}
}

func Getenv(env string) string {
	_env := os.Getenv(env)
	if len(_env) == 0 {
		log.Fatal("no " + env + " is set")
	}
	return _env
}

func AsPort(env string) (port int) {
	port, err := strconv.Atoi(env)
	if err != nil {
		log.Fatalf("invalid port: %s", env)
	}
	return
}

func NewSocket(ztype zmq.Type) (socket *zmq.Socket) {
	socket, err := zmq.NewSocket(ztype)
	if err != nil {
		log.Fatalf("failed creating socket type %d: %s", ztype, err.Error())
	}
	return
}

func Bind(socket *zmq.Socket, transport string, address string, port int) {
	endpoint := fmt.Sprintf("%s://%s:%d", transport, address, port)
	err := socket.Bind(endpoint)
	if err != nil {
		log.Fatalf("failed binding %s: %s", endpoint, err.Error())
	}
}
