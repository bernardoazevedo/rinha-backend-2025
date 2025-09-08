package logger

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func getFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func Add(message string) (int, error) {
	date := actualDateYMD()
	fileName := "tmp/" + date + ".txt"

	file, err := getFile(fileName)
	if err != nil {
		return 0, errors.New("error opening " + fileName + ": " + err.Error())
	}
	defer file.Close()

	byteMessage := []byte(message + "\n")
	bytes, err := file.Write(byteMessage)
	if err != nil {
		return bytes, err
	}

	return bytes, nil
}

func formatDateYMD(data time.Time) string {
	year := data.Year()
	month := data.Month()
	day := data.Day()

	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

func actualDateYMD() string {
	return formatDateYMD(time.Now())
}

func actualDateHMS() string {
	hour, minute, second := time.Now().Clock()
	return fmt.Sprintf("%d:%d:%d", hour, minute, second)
}
