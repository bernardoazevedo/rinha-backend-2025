package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/bernardoazevedo/rinha-de-backend-2025/health"
	"github.com/bernardoazevedo/rinha-de-backend-2025/key"
	"github.com/bernardoazevedo/rinha-de-backend-2025/logger"
)

func postPayment(payment Payment) (string, error) {
	logger.Add("init - " + payment.CorrelationId)

	paymentJson, err := json.Marshal(payment)
	if err != nil {
		return "", errors.New("error parsing payment")
	}

	url, _ := key.Get("url")

	postBody := bytes.NewBuffer(paymentJson)

	var response *http.Response
	for {
		// logger.Add("\tpost: "+url)
		response, err = http.Post(url+"/payments", "application/json", postBody)
		
		if err != nil || response.StatusCode != 200 {
			url, _ = health.CheckSetReturnUrl()
		} else {
			break //success
		}
		defer response.Body.Close()
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("error parsing body")
	}

	return string(responseBody), nil
}
