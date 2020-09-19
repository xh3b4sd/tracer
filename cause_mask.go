package tracer

import (
	"fmt"
	"runtime"
)

func Cause(err error) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*Error)
	if ok {
		if e.Wrpd != nil {
			return e.Wrpd
		}
	}

	return err
}

func Mask(err error) error {
	if err == nil {
		return nil
	}

	return mask(err)
}

func Maskf(e *Error, f string, v ...interface{}) error {
	// The masking has to happen before annotating the error. Annotating the
	// error before masking it would manipulate the error which we track as
	// cause in Error.Wrpd.
	err := mask(e)

	// Only annotate the error once we masked it so that we do not manipulate
	// the cause after we tracked it.
	e = err.(*Error)
	e.Anno = fmt.Sprintf(f, v...)

	return e
}

func mask(err error) error {
	// In case we get some arbitrary error, we create our own Error type so that
	// we can properly work with it. The error we create ourselves gets simply
	// annotated with the error message provided by the arbitrary error type.
	e, ok := err.(*Error)
	if !ok {
		e = &Error{
			Anno: err.Error(),
		}
	}

	// In case we get our own Error type, we create a copy of it so that we do
	// not manipulate the originally wrapped pointer during consecutive masking.
	if ok {
		e = e.Copy()
	}

	// If we got some arbitrary error or our known Error type was not wrapped
	// yet, we want to wrap and fill it accordingly.
	if !ok || (ok && e.Wrpd == nil) {
		e.Type = fmt.Sprintf("%T", err)
		e.Wrpd = err
	}

	// In all cases we want to fill the stack so that we can inspect the stack
	// trace if we ever have to during debugging.
	{
		if e.Stck != "" {
			e.Stck += ","
		}

		_, file, line, _ := runtime.Caller(2)
		e.Stck += fmt.Sprintf("%s:%d", file, line)
	}

	return e
}
