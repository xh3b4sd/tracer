package tracer

import "encoding/json"

// Error provides a traceable error instance that can be annotated with
// arbitrary contextual information along the error handling path.
type Error struct {
	Context     []Context
	Description string

	cause error
	trace []string
}

// Copy creates a runtime copy of the underlying *Error{} instance.
func (e *Error) Copy() *Error {
	return &Error{
		Context:     append([]Context{}, e.Context...),
		Description: e.Description,

		cause: e.cause,
		trace: append([]string{}, e.trace...),
	}
}

// Error returns the error's description or "ERROR".
func (e *Error) Error() string {
	if e.Description != "" {
		return e.Description
	}

	return "ERROR"
}

func (e *Error) Is(x error) bool {
	return cause(e) == cause(x)
}

// MarshalJSON returns the JSON representation of a non nil *Error type, or {}.
func (e *Error) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("{}"), nil
	}

	return json.Marshal(struct {
		Context     []Context `json:"context,omitempty"`
		Description string    `json:"description,omitempty"`
		Trace       []string  `json:"trace,omitempty"`
	}{
		Context:     e.Context,
		Description: e.Error(),
		Trace:       e.trace,
	})
}

// Unwrap returns the error's root cause. That is the first masked error.
func (e *Error) Unwrap() error {
	return e.cause
}

func cause(err error) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*Error)
	if ok && e.cause != nil {
		return e.cause
	}

	return err
}
