package interfaces

type Log interface {
	Errorln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Fatalf(format string, args ...interface{})
}