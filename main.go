// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// See http://formwork-io.github.io/ for more.

package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
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
	env := getenv(BCAST_INGRESS_PORT)
	bcast_ingress_port := asPort(env)
	env = getenv(BCAST_EGRESS_PORT)
	bcast_egress_port := asPort(env)
	env = getenv(RR1_INGRESS_PORT)
	rr1_ingress_port := asPort(env)
	env = getenv(RR1_EGRESS_PORT)
	rr1_egress_port := asPort(env)
	env = getenv(RR2_INGRESS_PORT)
	rr2_ingress_port := asPort(env)
	env = getenv(RR2_EGRESS_PORT)
	rr2_egress_port := asPort(env)

	pprint("starting")

	// CREATE EACH SOCKET...
	sub_ingress := newSocket(zmq.SUB)
	sub_ingress.SetSubscribe("")
	defer sub_ingress.Close()
	pub_egress := newSocket(zmq.PUB)
	pub_egress.SetLinger(1)
	defer pub_egress.Close()

	rr1_ingress := newSocket(zmq.ROUTER)
	defer rr1_ingress.Close()
	rr1_egress := newSocket(zmq.DEALER)
	defer rr1_egress.Close()
	rr2_ingress := newSocket(zmq.ROUTER)
	defer rr2_ingress.Close()
	rr2_egress := newSocket(zmq.DEALER)
	defer rr2_egress.Close()

	// ... AND BIND
	bind(sub_ingress, "tcp", "0.0.0.0", bcast_ingress_port)
	bind(pub_egress, "tcp", "0.0.0.0", bcast_egress_port)
	bind(rr1_ingress, "tcp", "0.0.0.0", rr1_ingress_port)
	bind(rr1_egress, "tcp", "0.0.0.0", rr1_egress_port)
	bind(rr2_ingress, "tcp", "0.0.0.0", rr2_ingress_port)
	bind(rr2_egress, "tcp", "0.0.0.0", rr2_egress_port)

	poller := zmq.NewPoller()
	poller.Add(sub_ingress, zmq.POLLIN)
	poller.Add(rr1_ingress, zmq.POLLIN)
	poller.Add(rr2_ingress, zmq.POLLIN)
	poller.Add(rr1_egress, zmq.POLLIN)
	poller.Add(rr2_egress, zmq.POLLIN)

	pprint("greenline alive")
	exitchan := make(chan os.Signal, 0)
	signal.Notify(exitchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-exitchan
		out("received %s signal, exiting.\n", sig.String())
		os.Exit(0)
	}()

	pprint("greenline ready")
	for {
		sockets, _ := poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case sub_ingress:
				pprint("processing broadcast message")
				for {
					msg, err := s.Recv(0)
					if err != nil {
						die("broadcast more: %s", err.Error())
					}
					more, err := s.GetRcvmore()
					if err != nil {
						die("broadcast recv more: %s", err.Error())
					}
					if more {
						pub_egress.Send(msg, zmq.SNDMORE)
					} else {
						pub_egress.Send(msg, 0)
						break
					}
				}
			case rr1_ingress:
				pprint("processing rr1 request")
				for {
					msg, err := s.Recv(0)
					if err != nil {
						die("rr1 ingress: %s", err.Error())
					}
					more, err := s.GetRcvmore()
					if err != nil {
						die("rr1 ingress recv more: %s", err.Error())
					}
					if more {
						rr1_egress.Send(msg, zmq.SNDMORE)
					} else {
						rr1_egress.Send(msg, 0)
						break
					}
				}
			case rr2_ingress:
				pprint("processing rr2 request")
				for {
					msg, err := s.Recv(0)
					if err != nil {
						die("rr2 ingress: %s", err.Error())
					}
					more, err := s.GetRcvmore()
					if err != nil {
						die("rr2 ingress recv more: %s", err.Error())
					}
					if more {
						rr2_egress.Send(msg, zmq.SNDMORE)
					} else {
						rr2_egress.Send(msg, 0)
						break
					}
				}
			case rr1_egress:
				pprint("processing rr1 response")
				for {
					msg, err := s.Recv(0)
					if err != nil {
						die("rr1 egress: %s", err.Error())
					}
					more, err := s.GetRcvmore()
					if err != nil {
						die("rr1 egress recv more: %s", err.Error())
					}
					if more {
						rr1_ingress.Send(msg, zmq.SNDMORE)
					} else {
						rr1_ingress.Send(msg, 0)
						break
					}
				}
			case rr2_egress:
				pprint("processing rr2 response")
				for {
					msg, err := s.Recv(0)
					if err != nil {
						die("rr2 egress: %s", err.Error())
					}
					more, err := s.GetRcvmore()
					if err != nil {
						die("rr2 egress recv more: %s", err.Error())
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

func getenv(env string) string {
	_env := os.Getenv(env)
	if len(_env) == 0 {
		die("no " + env + " is set")
	}
	return _env
}

func asPort(env string) (port int) {
	port, err := strconv.Atoi(env)
	if err != nil {
		die("invalid port: %s", env)
	} else if port < 1 || port > 65535 {
		die("invalid port: %s", env)
	}
	return
}

func newSocket(ztype zmq.Type) (socket *zmq.Socket) {
	socket, err := zmq.NewSocket(ztype)
	if err != nil {
		die("failed creating socket type %d: %s", ztype, err.Error())
	}
	return
}

func bind(socket *zmq.Socket, transport string, address string, port int) {
	endpoint := fmt.Sprintf("%s://%s:%d", transport, address, port)
	out("Binding socket %d... ", port)
	err := socket.Bind(endpoint)
	if err != nil {
		die("failed binding %s: %s", endpoint, err.Error())
	}
	fmt.Println("done.")
}

func makeMsg(msg string, args ...interface{}) string {
	const layout = "%d%02d%02d-%02d-%02d-%02d greenline[%d]: %s"
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	seconds := now.Second()
	pid := os.Getpid()
	arg := fmt.Sprintf(msg, args...)
	ret := fmt.Sprintf(layout, year, month, day, hour, minute, seconds, pid, arg)
	return ret
}

func pprint(msg string, args ...interface{}) {
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stdout, msg+"\n")
}

func out(msg string, args ...interface{}) {
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stdout, msg)
	os.Stdout.Sync()
}

func die(msg string, args ...interface{}) {
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stderr, msg+"\n")
	os.Exit(1)
}
