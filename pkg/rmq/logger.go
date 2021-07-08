package rmq

type Logger interface {
	Debugf(format string, data ...interface{})
	Infof(format string, data ...interface{})
	Warnf(format string, data ...interface{})
	Errorf(format string, data ...interface{})
	Fatalf(format string, data ...interface{})

	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)

	WithField(name string, value interface{}) Logger
	WithFields(pairs ...interface{}) Logger

	WithErr(err error) Logger
	WithStr(key string, val string) Logger
}
