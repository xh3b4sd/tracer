package tracer

import (
	"encoding/json"
	"strings"
)

type Error struct {
	Anno string   `json:"anno,omitempty"`
	Code string   `json:"code,omitempty"`
	Desc string   `json:"desc,omitempty"`
	Docs string   `json:"docs,omitempty"`
	Kind string   `json:"kind,omitempty"`
	Stck []string `json:"stck,omitempty"`
	Wrpd error    `json:"-"`
}

func (e *Error) Copy() *Error {
	c := &Error{
		Anno: e.Anno,
		Code: e.Code,
		Desc: e.Desc,
		Docs: e.Docs,
		Kind: e.Kind,
		Stck: e.Stck,
		Wrpd: e.Wrpd,
	}

	return c
}

func (e *Error) Error() string {
	if e.Kind == "" && e.Anno != "" {
		return e.Anno
	}

	kind := toStringCase(e.Kind)

	if e.Kind != "" && e.Anno == "" {
		return kind
	}

	return strings.TrimSuffix(kind, " error") + ": " + e.Anno
}

func (e *Error) GoString() string {
	return mustJson(e)
}

func (e *Error) Is(x error) bool {
	return cause(e) == cause(x)
}

func (e *Error) Unwrap() error {
	return e.Wrpd
}

func cause(err error) error {
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

func mustJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}

	return string(b)
}
