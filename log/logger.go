package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

var _Logger *Logger

func OpenLog(l *LoggerConfig, options ...LoggerOption) {
	_Logger = &Logger{
		l: l,
	}
	for _, option := range options {
		option(l)
	}
	DefaultLogger(l)
	_Logger.newOutput()
}

func Default() *Logger {
	return _Logger
}

type Logger struct {
	errWriter io.WriteCloser
	writer    io.WriteCloser
	errLogger *log.Logger
	logger    *log.Logger
	l         *LoggerConfig
}

func (l *Logger) Fatal(data interface{}) {
	if l.errLogger == nil {
		log.Println("Err:default Logger doesn't work, use stdout")
		log.Println(data)
		return
	}
	l.errLogger.Fatalln(data)
}

func (l *Logger) Log(data interface{}) {
	if l.logger == nil {
		log.Println("Err:default Logger doesn't work, use stdout")
		log.Println(data)
		return
	}
	l.logger.Println(data)
}

func (l *Logger) outputPath(logName string) string {
	for _, postfix := range l.l.PostFix {
		logName = postfix(logName, l)
	}
	return fmt.Sprintf("%v/%v.log", l.l.LogDir, logName)
}

func (l *Logger) closeAll() {
	l.logger = nil
	l.errLogger = nil
	if l.writer != nil || l.errWriter != nil {
		// 尝试转化为文件
		ew, ok1 := l.errWriter.(*os.File)
		w, ok2 := l.writer.(*os.File)
		if ok1 && !isStdout(ew) {
			_ = l.errWriter.Close()
		}
		if ok2 && !isStdout(w) {
			_ = l.writer.Close()
		}

	}
	l.writer = nil
	l.errWriter = nil
}

func (l *Logger) Close() {
	l.closeAll()
	_Logger = nil
}

// 切换输出文件
func (l *Logger) newOutput() {
	// 关闭已经打开的日志
	l.closeAll()
	// 尝试打开文件
	writer, err := os.Create(l.outputPath(l.l.LogName))
	if err != nil {
		writer = os.Stdout
		fmt.Println(err)
	}
	l.writer = writer
	errWriter, err := os.Create(l.outputPath(l.l.ErrLogName))
	if err != nil {
		errWriter = os.Stdout
		fmt.Println(err)
	}
	l.writer = writer
	// 替换log对象
	l.logger = log.New(writer, "", log.Ldate|log.Ltime|log.Lshortfile)
	l.errLogger = log.New(errWriter, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func isStdout(file *os.File) bool {
	return file.Name() == "/dev/stdout"
}
