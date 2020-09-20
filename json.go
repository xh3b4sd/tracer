package tracer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func JSON(err error) string {
	// If the given error is nil, we simply return an empty JSON object.
	if err == nil {
		return "{}"
	}

	// If the given error is our Error type we can simply serialize a JSON
	// string based on it.
	{
		e, ok := err.(*Error)
		if ok {
			return mustJSONMarshal(e)
		}
	}

	// If the given error is some arbitary error type we simply return an
	// annotated JSON object.
	{
		e := &Error{
			Anno: err.Error(),
			Type: fmt.Sprintf("%T", err),
		}

		return mustJSONMarshal(e)
	}
}

func escape(s string) string {
	return strconv.Quote(s)
}

func jsonError(e *Error) string {
	var s string

	s += "{"

	if e.Anno != "" {
		s += "\""
		s += "anno"
		s += "\""
		s += ":"
		s += escape(e.Anno)
	}

	if e.Desc != "" {
		if !strings.HasSuffix(s, "{") {
			s += ","
		}

		s += "\""
		s += "desc"
		s += "\""
		s += ":"
		s += escape(e.Desc)
	}

	if e.Docs != "" {
		if !strings.HasSuffix(s, "{") {
			s += ","
		}

		s += "\""
		s += "docs"
		s += "\""
		s += ":"
		s += escape(e.Docs)
	}

	if e.Kind != "" {
		if !strings.HasSuffix(s, "{") {
			s += ","
		}

		s += "\""
		s += "kind"
		s += "\""
		s += ":"
		s += escape(e.Kind)
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
		s += escape(e.Type)
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

		s += "\""
		s += c[0]
		s += ":"
		s += c[1]
		s += "\""

		if i+1 < len(frames) {
			s += ", "
		}
	}
	s += "]"

	return s
}

func mustJSONMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}

	return string(b)
}
