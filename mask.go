package tracer

import (
	"fmt"
	"runtime"
)

func Mask(e error, c ...Context) error {
	if e == nil {
		return nil
	}

	t, o := e.(*Error)
	if !o {
		return mask(&Error{
			Context:     c,
			Description: e.Error(),
			cause:       e,
		}, c...)
	}

	return mask(t, c...)
}

func mask(e *Error, c ...Context) *Error {
	var n *Error

	if e.cause == nil {
		n = e.Copy()
		n.cause = e
	} else {
		n = e
	}

	if len(c) != 0 {
		n.Context = append(n.Context, c...)
	}

	{
		_, f, l, _ := runtime.Caller(2)
		n.trace = append(n.trace, fmt.Sprintf("%s:%d", f, l))
	}

	return n
}
