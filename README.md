# tracer

The `tracer` package provides a traceable `error` instance that can be annotated
with arbitrary contextual information along the error handling path.

Error handling in Go does not provide any tracing functionality out of the box.
This makes error masking necessary so that we can comprehend what went wrong and
why some problem occured. If we do not have any way to understand where along
the code execution path inside of our business logic an error occurred, then
debugging and fixing any root cause of any given problem may take an unnecessary
large amount of time. Therefore, using `tracer`, errors can be masked and stack
traces can be printed.

### Match

A typical `error.go` in any package might look like the following example. Note
a couple of best practices to align with, for simplicity and consistency
reasons.

- Keep error types private so that nobody outside your package can mess with it.
- Keep error matchers public so that anyone can match against your package errors.
- Keep error matcher implementations simple by using `errors.Is(a, b)`.
- Keep the order of errors and matchers alphabetical for easier navigation.
- Provide useful error descriptions explaining why errors were likely to occur.
- Provide contextual information along the error handling path to help the humans.

```golang
package foo

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var notFoundError = &tracer.Error{
	Description: "this thing was not found because of xyz",
}

func IsNotFound(err error) bool {
	return errors.Is(err, notFoundError)
}
```

### Mask

Below is a **bad** example to illustrate how to not do error handling.

```golang
return err
```

Below are **good** examples to illustrate how to do error handling.

```golang
return tracer.Mask(err)
```

```golang
return tracer.Mask(err, tracer.Context{Key: "resource", Val: "identifier"})
```

### Print

Use `tracer.Json(err)` in order to print the JSON repesentation of the error
instance at hand, similar to the indented example shown below. The first item of
the trace array identifies the first instance of error masking. That location
should typically represent the root cause.

```
{
  "context": [
    {
      "key": "resource",
      "value": "identifier"
    }
  ],
  "description": "some useful error description",
  "trace": [
    "--REPLACED--/error_test.go:119",
    "--REPLACED--/error_test.go:124"
  ]
}
```

### Panic

Use `tracer.Panic(tracer.Mask(err))` in program entry points of command line
tools in order to conveniently produce consistent error messages upon unexpected
program failure.

```golang
package main

import (
	"github.com/xh3b4sd/tracer"
)

var (
	alreadyExistsError = &tracer.Error{
		Description: "that thing does already exist because of xyz",
		Context: []tracer.Context{
			{Key: "code", Val: "alreadyExistsError"},
		},
	}
)

func main() {
	err := mainE()
	if err != nil {
		tracer.Panic(tracer.Mask(err)) // line 19
	}
}

func mainE() error {
	err := tracer.Mask(alreadyExistsError) // line 24
	return tracer.Mask(err) // line 25
}
```

```
program panic at 2025-07-17 19:22:58.39201 +0000 UTC

    {
        "context": [
            {
                "key": "code",
                "value": "alreadyExistsError"
            }
        ],
        "description": "that thing does already exist because of xyz",
        "trace": [
            "--REPLACED--/main.go:24",
            "--REPLACED--/main.go:25",
            "--REPLACED--/main.go:19"
        ]
    }

exit status 1
```
