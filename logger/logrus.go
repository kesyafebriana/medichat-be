package logger

import "github.com/sirupsen/logrus"

type LoggerLogrus struct {
	l *logrus.Logger
}

func FromLogrus(l *logrus.Logger) Logger {
	return &LoggerLogrus{
		l: l,
	}
}

func (l *LoggerLogrus) Debug(args ...interface{}) {
	l.l.Debug(args...)
}

func (l *LoggerLogrus) Debugf(format string, args ...interface{}) {
	l.l.Debugf(format, args...)
}

func (l *LoggerLogrus) DebugFields(fieds map[string]interface{}, args ...interface{}) {
	l.l.WithFields(fieds).Debug(args...)
}

func (l *LoggerLogrus) Info(args ...interface{}) {
	l.l.Info(args...)
}

func (l *LoggerLogrus) Infof(format string, args ...interface{}) {
	l.l.Infof(format, args...)
}

func (l *LoggerLogrus) InfoFields(fieds map[string]interface{}, args ...interface{}) {
	l.l.WithFields(fieds).Info(args...)
}

func (l *LoggerLogrus) Warn(args ...interface{}) {
	l.l.Warn(args...)
}

func (l *LoggerLogrus) Warnf(format string, args ...interface{}) {
	l.l.Warnf(format, args...)
}

func (l *LoggerLogrus) WarnFields(fieds map[string]interface{}, args ...interface{}) {
	l.l.WithFields(fieds).Warn(args...)
}

func (l *LoggerLogrus) Error(args ...interface{}) {
	l.l.Error(args...)
}

func (l *LoggerLogrus) Errorf(format string, args ...interface{}) {
	l.l.Errorf(format, args...)
}

func (l *LoggerLogrus) ErrorFields(fieds map[string]interface{}, args ...interface{}) {
	l.l.WithFields(fieds).Error(args...)
}

func (l *LoggerLogrus) Fatal(args ...interface{}) {
	l.l.Fatal(args...)
}

func (l *LoggerLogrus) Fatalf(format string, args ...interface{}) {
	l.l.Fatalf(format, args...)
}

func (l *LoggerLogrus) FatalFields(fieds map[string]interface{}, args ...interface{}) {
	l.l.WithFields(fieds).Fatal(args...)
}
