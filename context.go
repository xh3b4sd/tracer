package tracer

type Context struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
