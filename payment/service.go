package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/bernardoazevedo/rinha-de-backend-2025/key"
)

func postPayment(payment Payment) (string, error) {
	paymentJson, err := json.Marshal(payment)
	if err != nil {
		return "", errors.New("error parsing payment")
	}

	paymentDefaultUrl := os.Getenv("PAYMENT_DEFAULT_URL")
	paymentFallbackUrl := os.Getenv("PAYMENT_FALLBACK_URL")
	url, _ := key.Get("url")

	postBody := bytes.NewBuffer(paymentJson)
	response, err := http.Post(url+"/payments", "application/json", postBody)
	if err != nil {
		
		if url == paymentDefaultUrl {
			url = paymentFallbackUrl
			key.Set("url", url)
		}
		return "", errors.New("error during request")
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("error parsing body")
	}

	return string(responseBody), nil
}
