package file

import (
	"fmt"
	"github.com/yaza-putu/crud-fiber/internal/pkg/unique"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

// ToPublic folder
func ToPublic(file *multipart.FileHeader, dest string, randomName bool) (string, error) {
	src, err := file.Open()
	defer src.Close()

	if err != nil {
		return "", err
	}

	// Destination
	fileName := file.Filename
	if randomName {
		split := strings.Split(file.Filename, ".")
		fileName = fmt.Sprintf("%s.%s", unique.Uid(13), split[len(split)-1])
	}

	destPath := fmt.Sprintf("public/%s/%s", dest, fileName)
	_, err = os.Stat(fmt.Sprintf("public/%s"))
	if err != nil {
		err = os.Mkdir(fmt.Sprintf("public/%s", dest), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	dst, err := os.Create(destPath)
	defer dst.Close()
	if err != nil {
		return "", err
	}

	// store to destination
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return dest, nil
}
