package logger

type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	DebugFields(map[string]interface{}, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	InfoFields(map[string]interface{}, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	WarnFields(map[string]interface{}, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	ErrorFields(map[string]interface{}, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	FatalFields(map[string]interface{}, ...interface{})
}
