package httptools

import (
	"aws_nat/errhandling"
	"io"
	"net/http"
	"os"
)

func hostname() string {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "hostname")

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
