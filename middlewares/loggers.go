package middlewares

import "github.com/gyozatech/noodlog"

var (
	infoLog  infoLogger
	printLog printLogger
	errLog   errorLogger
)

type printLogger interface {
	Println(message ...interface{})
}

type infoLogger interface {
	Info(message ...interface{})
}

type errorLogger interface {
	Error(message ...interface{})
}

// SetLogger sets the logger for the middlewares
func SetLogger(l interface{}) {
	iLogger, isInfoLogger := l.(infoLogger)
	if isInfoLogger {
		infoLog = iLogger
	}

	pLogger, isPrintLogger := l.(printLogger)
	if isPrintLogger {
		printLog = pLogger
	}

	eLogger, isErrorLogger := l.(errorLogger)
	if isErrorLogger {
		errLog = eLogger
	}
}

// InitDefaultLogger sets the default logger to "github.com/gyozatech/noodlog" is no other logger has been set
func InitDefaultLogger() {
	SetLogger(noodlog.NewLogger().EnableTraceCaller())
}

func logInfo(message ...interface{}) {
	if infoLog != nil {
		infoLog.Info(message)
		return
	} else if printLog != nil {
		printLog.Println(message)
		return
	}
	InitDefaultLogger()
	infoLog.Info(message)
}

func logError(message ...interface{}) {
	if errLog != nil {
		errLog.Error(message)
		return
	}
	InitDefaultLogger()
	errLog.Error(message)
}
