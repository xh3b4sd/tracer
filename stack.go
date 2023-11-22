package tracer

import "encoding/json"

func Stack(err error) string {
	if err == nil {
		return "null"
	}

	t, o := err.(*Error)
	if !o {
		return "[]"
	}

	b, e := json.Marshal(t.Stck)
	if e != nil {
		panic(e)
	}

	return string(b)
}
