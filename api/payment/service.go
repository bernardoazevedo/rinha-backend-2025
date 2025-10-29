package payment

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bernardoazevedo/rinha-backend-2025/api/key"
	paymentqueue "github.com/bernardoazevedo/rinha-backend-2025/api/paymentQueue"
	"github.com/valyala/fasthttp"
)

func postPayment(payment []byte) (bool, error) {
	var err error
	alreadyExistsPayment := false

	url, err := key.Get("url")
	if err != nil {
		// não tem url ou não consegui buscar, então uso a padrão
		url = "http://payment-processor-default:8080"
	}

	_, alreadyExistsPayment, err = post(url, payment)
	if alreadyExistsPayment {
		// já existe, não faço nada

	} else if err != nil {
		fmt.Printf("deu ruim, fila de novo: %s\n", string(payment))
		paymentqueue.Add(payment)
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

func queuePayment(payment []byte) error {
	paymentqueue.Add(payment)
	return nil
}
