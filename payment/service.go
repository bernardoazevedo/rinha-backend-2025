package payment

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bernardoazevedo/rinha-backend-2025/health"
	paymentqueue "github.com/bernardoazevedo/rinha-backend-2025/paymentQueue"
	"github.com/valyala/fasthttp"
)

func postPayment(payment []byte) (bool, error) {
	var err error
	alreadyExistsPayment := false

	url := health.PostUrl

	for i := 0; i < 3; i++ {
		_, alreadyExistsPayment, err = post(url, payment)

		if alreadyExistsPayment {
			// saio fora
			break

		} else if err != nil {
			// tento até a 3 vez
			url, err = health.CheckSetReturnUrl()
			if err != nil {
				fmt.Printf("erro [%d] ao checar url: "+err.Error()+"\n", i)

				for j := 0; j < i; j++ { // espera 1s * numRequisicao => 1s, 2s, 3s
					time.Sleep(time.Second)
				}
			}

		} else {
			// defer response.CloseBodyStream()
			// deu bom, saio fora
			break
		}
	}

	return alreadyExistsPayment, err
}

func post(url string, body []byte) (*fasthttp.Response, bool, error) {
	var statusCode int

	req := fasthttp.AcquireRequest()
	req.SetBody(body)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.Header.SetRequestURI(url + "/payments")
	response := fasthttp.AcquireResponse()

	err := fasthttp.Do(req, response)
	statusCode = response.StatusCode()
	if err != nil {
		statusCode = 400
	}
	defer fasthttp.ReleaseRequest(req)

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

func queuePayment(paymentJson string) error {
	paymentqueue.Push(paymentJson)
	return nil
}
