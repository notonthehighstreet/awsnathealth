package httpops

import (
	"io"
	"net/http"
	"os"
)

func hostname() string {
	name, err := os.Hostname()

	if err != nil {
		panic(err)
	}
	return name
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, hostname()+" "+"is alive")
}

// HttpdHandler brings up a httpd hander and lisens on a specific port
func HttpdHandler(port string) {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":"+port, nil)
}
