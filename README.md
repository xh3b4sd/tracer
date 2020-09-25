# tracer

Simple error handling and stack traces for fast debugging. Error handling in go
does not provide any tracing functionality out of the box. This makes error
masking necessary so that you can comprehend what went wrong and why some
problem occured. When you do not have any way to understand where along the code
execution of your business logic an error occurred then debugging and fixing the
root cause of a problem takes an unnecessary big amount of time. Therefore,
using `tracer` errors can be masked and stack traces can be printed.



### Errors And Matchers

A typical `error.go` in any package might look like the following example. Note
a couple of best practices to align with for simplicity and consistency reasons.

* Keep error types private so that nobody outside your package can mess with it.
* Keep error matchers public so that anyone can match against your package errors.
* Keep the variable name and `Kind` consistent for easy tracking during debugging.
* Keep error matcher implementations simple by using `errors.Is(a, b)`.
* Keep the order errors and matchers alphabetical for easier navigation.

```golang
package foo

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var invalidConfigError = &tracer.Error{
	Kind: "invalidConfigError",
}

func IsInvalidConfig(err error) bool {
	return errors.Is(err, invalidConfigError)
}

var notFoundError = &tracer.Error{
	Kind: "notFoundError",
}

func IsNotFound(err error) bool {
	return errors.Is(err, notFoundError)
}
```



### Matching In Code

Below is a **bad** example to illustrate how not to do error handling.

```golang
return err
```

Below is a **good** example to illustrate how to do error handling.

```golang
return tracer.Mask(err)
```



### Stack Trace Printing

Use `tracer.JSON(err)` to get the JSON repesentation of an error of type
`*tracer.Error` like the example below shows.

```json
{
	"anno": "some useful annotation",
	"kind": "testError",
	"stck": [
		"--REPLACED--/json_test.go:111",
		"--REPLACED--/json_test.go:112"
	],
	"type": "*tracer.Error"
}
```

Use `tracer.Panic(err)` in program entry points of command line tools in order
to conveniently produce consistent error messages upon unexpected program
failure.

```golang
func main() {
    err := mainE(context.Background())
    if err != nil {
        tracer.Panic(err)
    }
}
```
