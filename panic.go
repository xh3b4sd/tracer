package tracer

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Panic is meant to be used in user facing applications like command line
// tools. Such applications usually propagate back runtime errors. In order to
// make error handling for these specific cases most convenient, Panic might
// simply be called like shown below. The program entry point might be as simple
// as the following snippet.
//
//	func main() {
//	    err := mainE()
//	    if err != nil {
//	        tracer.Panic(tracer.Mask(err))
//	    }
//	}
//
// The code snippet of the program entry point above might produce an output
// similar to the example below.
//
//	program panic at 2025-07-17 19:22:58.39201 +0000 UTC
//
//	    {
//	        "context": [
//	            {
//	                "key": "code",
//	                "value": "alreadyExistsError"
//	            }
//	        ],
//	        "description": "that thing does already exist because of xyz",
//	        "trace": [
//	            "--REPLACED--/main.go:24",
//	            "--REPLACED--/main.go:25",
//	            "--REPLACED--/main.go:19"
//	        ]
//	    }
//
//	exit status 1
func Panic(err error) {
	t, ok := err.(*Error)
	if !ok {
		panic(err)
	}

	fmt.Printf("program panic at %s\n", time.Now().UTC().String())
	fmt.Println()
	fmt.Println("    " + string(musJsn(t)))
	fmt.Println()

	os.Exit(1)
}

func musJsn(v any) []byte {
	jsn, err := json.MarshalIndent(v, "    ", "    ")
	if err != nil {
		panic(err)
	}

	return jsn
}
