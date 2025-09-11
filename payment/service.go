package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bernardoazevedo/rinha-de-backend-2025/health"
)

func postPayment(payment Payment) (string, error) {
	paymentJson, err := json.Marshal(payment)
	if err != nil {
		return "", errors.New("error parsing payment")
	}

	url := health.PostUrl

	var response *http.Response
	var statusCode int
	var errorMessage string
	
	postBody := bytes.NewBuffer(paymentJson)

	for i := 0; i < 3; i++ {
		response, err = http.Post(url+"/payments", "application/json", postBody)
		if response != nil {
			statusCode = response.StatusCode
		} else {
			statusCode = 400
		}

		if err != nil {
			errorMessage = fmt.Sprintf("[%d] "+err.Error(), statusCode)
			// return response, false, errors.New(errorMessage)

		} else if statusCode == http.StatusUnprocessableEntity {
			return "This payment already exists", nil

		} else if statusCode != 200 {
			errorMessage = fmt.Sprintf("[%d] status != 200", statusCode)
			// return response, false, errors.New(errorMessage)

		} else { //success
			return "Success", nil
		}

		url, _ = health.CheckSetReturnUrl()
	}

	return "", errors.New(errorMessage)
}
