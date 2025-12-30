package config

const (
	DebugLevel   = "debug"
	InfoLevel    = "info"
	FatalLevel   = "fatal"
	WarnLevel    = "warn"
	ErrorLevel   = "error"
	PanicLevel   = "panic"
	DefaultLevel = "debug"
)

type ZapConfig struct {
	Level            string         `mapstruct:"level"`
	Development      bool           `mapstruct:"development"`
	Encoding         string         `mapstruct:"encoding"`
	OutputPaths      []string       `mapstruct:"outputPaths"`
	ErrorOutputPaths []string       `mapstruct:"errorOutputPaths"`
	InitialFields    map[string]any `mapstruct:"initialFields"`
}
