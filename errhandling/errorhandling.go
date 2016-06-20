package errhandling

import (
	"fmt"
	"runtime"
)

//Error type struct
type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// ErrorChannel is an exporterd common error Channel, all error and panic event should be reouted into it.
var ErrorChannel = make(chan error)

// CatchPanic captures the panic event and routes it on into the ErrorChannel chan.
func CatchPanic(err *error, functionName string) {
	if r := recover(); r != nil {
		ErrorChannel <- Error{fmt.Sprintf("%s : PANIC Defered : %v\n", functionName, r)}

		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		ErrorChannel <- Error{fmt.Sprintf("%s : Stack Trace : %s", functionName, string(buf))}

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	} else if err != nil && *err != nil {
		ErrorChannel <- Error{fmt.Sprintf("%s : ERROR : %v\n", functionName, *err)}
	}
}
