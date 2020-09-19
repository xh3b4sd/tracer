package tracer

type Error struct {
	Anno string `json:"anno,omitempty"`
	Desc string `json:"desc,omitempty"`
	Docs string `json:"docs,omitempty"`
	Kind string `json:"kind,omitempty"`
	Stck string `json:"stck,omitempty"`
	Type string `json:"type,omitempty"`
	Wrpd error  `json:"-"`
}

func (e *Error) Copy() *Error {
	c := &Error{
		Anno: e.Anno,
		Desc: e.Desc,
		Docs: e.Docs,
		Kind: e.Kind,
		Stck: e.Stck,
		Type: e.Type,
		Wrpd: e.Wrpd,
	}

	return c
}

func (e *Error) Error() string {
	if e.Anno == "" {
		return toStringCase(e.Kind)
	}

	return e.Anno
}

func (e *Error) GoString() string {
	return JSON(e)
}

func (e *Error) Is(x error) bool {
	return Cause(e) == Cause(x)
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return []byte(jsonError(e)), nil
}
