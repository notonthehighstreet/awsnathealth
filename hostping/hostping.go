package hostping

import (
	"github.com/notonthehighstreet/awsnathealth/errhandling"
	"net"
	"time"

	"github.com/tatsushid/go-fastping"
)

// Ping returns a channel with true,false elements depending on the host availibility.
func Ping(hostname string, resp chan bool) {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "Ping")

	type response struct {
		addr *net.IPAddr
		rtt  time.Duration
	}

	p := fastping.NewPinger()
	p.Network("udp")
	netProto := "ip4:icmp"

	ra, err := net.ResolveIPAddr(netProto, hostname)
	if err != nil {
		panic(err)
	}

	results := make(map[string]*response)
	results[ra.String()] = nil
	p.AddIPAddr(ra)

	onRecv, onIdle := make(chan *response), make(chan bool)

	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		onRecv <- &response{addr: addr, rtt: t}
	}

	p.OnIdle = func() {
		onIdle <- true
	}

	p.MaxRTT = time.Second
	p.RunLoop()

	for {
		select {
		case res := <-onRecv:
			if _, ok := results[res.addr.String()]; ok {
				results[res.addr.String()] = res
			}
		case <-onIdle:
			for _, r := range results {
				if r == nil {
					resp <- false
				} else {
					resp <- true
				}
				results[ra.String()] = nil
			}
		}
	}
}
