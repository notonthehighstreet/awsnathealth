package httptools

import (
	"github.com/notonthehighstreet/awsnathealth/errhandling"
	"net"
	"net/http"
	"time"
)

var timeout = time.Duration(5 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// RespCode return the reponse code
func RespCode(url string) int {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "RespCode")

	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	resp, err := client.Get(url)
	resp.Body.Close()
	if err != nil {
		panic(err)
	} else {
		return resp.StatusCode
	}
}
