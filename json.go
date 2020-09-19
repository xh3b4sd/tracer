package tracer

import (
	"encoding/json"
	"fmt"
	"strings"
)

func JSON(err error) string {
	// If the given error is nil, we simply return an empty JSON object.
	if err == nil {
		return "{}"
	}

	// If the given error is our Error type we can simply serialize a JSON
	// string based on it.
	e, ok := err.(*Error)
	if ok {
		b, err := json.Marshal(e)
		if err != nil {
			panic(err.Error())
		}

		return string(b)
	}

	// If the given error is some arbitary error type we simply return an
	// annotated JSON object.
	{
		e := &Error{
			Anno: err.Error(),
			Type: fmt.Sprintf("%T", err),
		}

		b, err := json.Marshal(e)
		if err != nil {
			panic(err.Error())
		}

		return string(b)
	}
}

func jsonError(e *Error) string {
	var s string

	s += "{"

	if e.Anno != "" {
		s += "\""
		s += "anno"
		s += "\""
		s += ":"
		s += "\""
		s += e.Anno
		s += "\""
	}

	if e.Desc != "" {
		if !strings.HasSuffix(s, "{") {
			s += ","
		}

		s += "\""
		s += "desc"
		s += "\""
		s += ":"
		s += "\""
		s += e.Desc
		s += "\""
	}

	if e.Docs != "" {
		if !strings.HasSuffix(s, "{") {
			s += ","
		}

		s += "\""
		s += "docs"
		s += "\""
		s += ":"
		s += "\""
		s += e.Docs
		s += "\""
	}

	if e.Kind != "" {
		if !strings.HasSuffix(s, "{") {
			s += ","
		}

		s += "\""
		s += "kind"
		s += "\""
		s += ":"
		s += "\""
		s += e.Kind
		s += "\""
	}

	// Note that all struct properties like anno, desc and type are strings,
	// while stck is a JSON array when being marshalled.
	if e.Stck != "" {
		if !strings.HasSuffix(s, "{") {
			s += ","
		}

		s += "\""
		s += "stck"
		s += "\""
		s += ":"
		s += jsonStck(e.Stck)
	}

	// Note that the struct property type is the last property of the Error type
	// object. Therefore the trailing comma is omitted.
	if e.Type != "" {
		if !strings.HasSuffix(s, "{") {
			s += ","
		}

		s += "\""
		s += "type"
		s += "\""
		s += ":"
		s += "\""
		s += e.Type
		s += "\""
	}

	s += "}"

	return s
}

func jsonStck(stck string) string {
	if stck == "" {
		return "\"\""
	}

	var s string

	frames := strings.Split(stck, ",")

	s += "["
	for i, f := range frames {
		c := strings.Split(f, ":")

		s += "{"
		s += "\""
		s += "file"
		s += "\""
		s += ":"
		s += "\""
		s += c[0]
		s += "\""
		s += ","
		s += "\""
		s += "line"
		s += "\""
		s += ":"
		s += c[1]
		s += "}"

		if i+1 < len(frames) {
			s += ", "
		}
	}
	s += "]"

	return s
}
