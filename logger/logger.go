package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	f        *os.File
	filename string
	loggers  map[string]*log.Logger
}

var Log *Logger

var once sync.Once

func OpenLog(filename string) {
	once.Do(func() {
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		Log.filename = filename
		Log.f = f
		// loggers may already have been created (at global scope level);
		// make sure their output is set correctly
		for _, v := range Log.loggers {
			v.SetOutput(f)
		}
	})
}

func CloseLog() {
	if Log.f != nil {
		Log.f.Sync()
		Log.f.Close()
	}
}

func (l *Logger) GetLogger(prefix string) *log.Logger {
	prefix = "[" + prefix + "] "
	if _, ok := l.loggers[prefix]; !ok {
		l.loggers[prefix] = log.New(l.f, prefix, log.Ldate|log.Ltime|log.Llongfile)
	}
	return l.loggers[prefix]
}

func GetLogger(prefix string) *log.Logger {
	return Log.GetLogger(prefix)
}

func init() {
	Log = &Logger{loggers: make(map[string]*log.Logger)}
}
