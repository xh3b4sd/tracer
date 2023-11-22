package tracer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Panic is meant to be used in user facing applications like command line
// tools. Such applications usually propagate back runtime errors. In order to
// make error handling for these specific cases most convenient Panic might
// simply be called. The program entry point might be as simple as the following
// snippet.
//
//	func main() {
//	    err := mainE(context.Background())
//	    if err != nil {
//	        tracer.Panic(tracer.Mask(err))
//	    }
//	}
//
// The code snippet of the program entry point above might produce an output
// like below.
//
//	program panic at 2022-06-10 18:37:41.90837 +0000 UTC
//
//	    {
//	        "anno": "rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing dial tcp :7777: connect: connection refused\"",
//	        "stck": [
//	            "--REPLACED--/main.go:59",
//	            "--REPLACED--/main.go:23"
//	        ]
//	    }
//
//	exit status 1
func Panic(err error) {
	_, ok := err.(*Error)
	if !ok {
		panic(err)
	}

	fmt.Printf("program panic at %s\n", time.Now().UTC().String())
	fmt.Println()

	b := &bytes.Buffer{}
	err = json.Indent(b, []byte(mustJson(err)), "    ", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println("    " + b.String())
	fmt.Println()

	os.Exit(1)
}
