package logger

import (
	"fmt"
	"os"
	"sync"
)

var path string = "logs"
var appName string = "Mozaik"

type Logger struct {
	service string
}

var (
	log  *Logger
	once sync.Once
)

func GetLogger() *Logger {
	once.Do(func() {
		log = loadLogger(appName)
	})
	return log
}
func isPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false

}
func loadLogger(str string) *Logger {
	ok := isPath(path)
	if !ok {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
	}
	return &Logger{
		service: str,
	}
}
func (*Logger) Info(str string) {

}
