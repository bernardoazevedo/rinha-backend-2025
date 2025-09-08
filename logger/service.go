package logger

import (
	"errors"
	"os"

	"github.com/bernardoazevedo/rinha-de-backend-2025/dates"
)

func getFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func Add(message string) (int, error) {
	date := dates.YearMonthDay()
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

