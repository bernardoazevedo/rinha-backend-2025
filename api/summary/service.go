package summary

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/valyala/fasthttp"
)

func getPaymentsSummary(paymentUrl string, from string, to string) (Summary, error) {
	var summary Summary

	path := "/admin/payments-summary"
	params := url.Values{
		"from": {from},
		"to":   {to},
	}
	url := paymentUrl + path + "?" + params.Encode()

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	req.Header.Add("X-Rinha-Token", "123")
	req.Header.SetRequestURI(url)
	response := fasthttp.AcquireResponse()

	err := fasthttp.Do(req, response)
	if err != nil {
		return summary, errors.New("error during request")
	}
	defer fasthttp.ReleaseRequest(req)

	responseBody := response.Body()
	err = json.Unmarshal(responseBody, &summary)
	if err != nil {
		return summary, errors.New("error parsing summary")
	}

	return summary, nil
}
