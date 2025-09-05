package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func postPayment(payment Payment, url string) (string, error) {
	paymentJson, err := json.Marshal(payment)
	if err != nil {
		return "", errors.New("error parsing payment")
	}

	postBody := bytes.NewBuffer(paymentJson)
	response, err := http.Post(url, "application/json", postBody)
	if err != nil {
		return "", errors.New("error during request")
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("error parsing body")
	}

	return string(responseBody), nil
}