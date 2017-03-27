package supermq

type logger interface {
	Output(calldepth int, s string) error
}

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)
