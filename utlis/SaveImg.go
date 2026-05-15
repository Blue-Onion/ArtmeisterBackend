package utlis

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type Save interface {
	SaveLocal()
}

func isPath(path string) bool {

	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false

}
func SaveLocal(file multipart.File, fileHeader *multipart.FileHeader, path string) (string, error) {
	is := isPath(path)
	if !is {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	filePath := fmt.Sprintf("%s/%s", path, fileHeader.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
