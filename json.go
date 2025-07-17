package tracer

import (
	"encoding/json"
	"fmt"
)

// Json returns the marshaled JSON string of type *Error, or "{}", or the result
// of err.Error().
//
//	{
//	  "context": [
//	    {
//	      "key": "re-source",
//	      "val": "id"
//	    }
//	  ],
//	  "description": "test error two description",
//	  "trace": [
//	    "--REPLACED--/error_test.go:189",
//	    "--REPLACED--/error_test.go:190"
//	  ]
//	}
func Json(err error) string {
	if err == nil {
		return "{}"
	}

	t, o := err.(*Error)
	if !o {
		t = &Error{
			Context: []Context{
				{Key: "type", Val: fmt.Sprintf("%T", err)},
			},
			Description: err.Error(),
		}
	}

	b, e := json.Marshal(t)
	if e != nil {
		panic(e)
	}

	return string(b)
}
