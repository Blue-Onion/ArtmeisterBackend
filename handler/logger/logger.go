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
	log    *Logger
	once   sync.Once
	logErr error
)

func GetLogger() (*Logger, error) {
	once.Do(func() {
		log, logErr = loadLogger(appName)
	})
	return log, logErr
}
func isPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false

}
func createFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
func initLogsFile() error {
	files := []string{
		"err.log",
		"info.log",
	}
	for _, f := range files {
		filePath := fmt.Sprintf("%s/%s", path, f)
		_, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				err = createFile(filePath)
				if err != nil {
					return err
				}
				return nil
			} else {
				return err
			}
		}
	}
	return nil
}
func loadLogger(str string) (*Logger, error) {
	ok := isPath(path)
	if !ok {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	err := initLogsFile()
	if err != nil {
		return nil, err
	}
	return &Logger{
		service: str,
	}, nil
}
func (*Logger) Info(str string) {

}
