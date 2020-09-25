package tracer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// Panic is meant to be used in user facing applications like command line
// tools. Such applications usually propagate back runtime errors. In order to
// make error handling for these specific cases most convenient Panic might
// simply be called. The program entry point might be as simple as the following
// snippet.
//
//     func main() {
//         err := mainE(context.Background())
//         if err != nil {
//             tracer.Panic(err)
//         }
//     }
//
// The code snippet of the program entry point above might produce an output
// like below.
//
//     program panic
//
//         {
//             "anno": "rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing dial tcp :7777: connect: connection refused\"",
//             "stck": [
//                 "--REPLACED--/main.go:59",
//                 "--REPLACED--/main.go:23"
//             ],
//             "type": "*status.Error"
//         }
//
//     exit status 1
//
func Panic(err error) {
	fmt.Println("program panic")
	fmt.Println()

	b := &bytes.Buffer{}
	err = json.Indent(b, []byte(JSON(err)), "    ", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println("    " + b.String())
	fmt.Println()

	os.Exit(1)
}
