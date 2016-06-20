package logging

import (
	"io"
	"log"
	"os"
)

var (
	warningHandle, errorHandle, infoHandle io.Writer
)

// Declares the logging levels
var (
	Trace   *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Info    *log.Logger
)

// Log takes care of logging to file.
func Log(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	logfile string) {

	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Error.Println("Failed to open log file", ":", err)
	}

	Trace = log.New(file,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(file,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(file,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(file,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
