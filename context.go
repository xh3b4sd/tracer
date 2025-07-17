package tracer

type Context struct {
	Key string `json:"key"`
	Val any    `json:"val"`
}
