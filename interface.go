package tracer

type Interface interface {
	Error() string
	Is(error) bool
	Unwrap() error
}
