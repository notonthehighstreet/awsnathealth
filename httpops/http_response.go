package httpops

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

var timeout = time.Duration(5 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func RespCode(url string) int {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Whoops: ", e)
		}
	}()

	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	} else {
		return resp.StatusCode
	}
}

// func main() {
// 	respcode := httpRespCode("http://www.notonthehighstreet.com/admin")
// 	fmt.Print(respcode)
// }
