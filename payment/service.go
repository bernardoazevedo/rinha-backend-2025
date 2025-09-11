package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bernardoazevedo/rinha-de-backend-2025/health"
	"github.com/bernardoazevedo/rinha-de-backend-2025/logger"
	paymentqueue "github.com/bernardoazevedo/rinha-de-backend-2025/paymentQueue"
)

func postPayment(paymentJson string) (bool, error) {
	var response *http.Response
	var err error
	alreadyExistsPayment := false

	url := health.PostUrl

	for i := 0; i < 3; i++ {
		response, alreadyExistsPayment, err = post(url, paymentJson)

		if alreadyExistsPayment {
			// saio fora
			break

		} else if err != nil {
			// tento até a 3 vez
			url, err = health.CheckSetReturnUrl()
			if err != nil {
				logger.Add(fmt.Sprintf("erro [%d] ao checar url: ", i) + err.Error())

				for j := 0; j < i; j++ { // espera 1s * numRequisicao => 1s, 2s, 3s
					time.Sleep(time.Second)
				}
			}

		} else {
			defer response.Body.Close()
			// deu bom, saio fora
			break
		}
	}

	return alreadyExistsPayment, err
}

func post(url string, paymentJson string) (*http.Response, bool, error) {
	var statusCode int

	postBody := bytes.NewBufferString(paymentJson)

	response, err := http.Post(url+"/payments", "application/json", postBody)
	if response != nil {
		statusCode = response.StatusCode
	} else {
		statusCode = 400
	}

	if err != nil {
		message := fmt.Sprintf("[%d] "+err.Error(), statusCode)
		return response, false, errors.New(message)

	} else if statusCode == http.StatusUnprocessableEntity {
		message := fmt.Sprintf("[%d] payment already exists", statusCode)
		return response, true, errors.New(message)

	} else if statusCode != 200 {
		message := fmt.Sprintf("[%d] status != 200", statusCode)
		return response, false, errors.New(message)

	} else { //success
		return response, false, nil
	}
}

func queuePayment(payment Payment) error {
	paymentJson, err := json.Marshal(payment)
	if err != nil {
		return errors.New("error parsing payment")
	}

	err = paymentqueue.Add(string(paymentJson))
	if err != nil {
		return errors.New("error adding to queue")
	}

	return nil
}
