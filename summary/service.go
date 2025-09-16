package summary

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

func getPaymentsSummary(paymentUrl string, from string, to string) (Summary, error) {
	var summary Summary
	client := &http.Client{}

	path := "/admin/payments-summary"
	params := url.Values{
		"from": {from},
		"to":   {to},
	}
	url := paymentUrl + path + "?" + params.Encode()

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return summary, errors.New("error starting request: " + err.Error())
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Rinha-Token", "123")

	response, err := client.Do(request)
	if err != nil {
		return summary, errors.New("error during request: " + err.Error())
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return summary, errors.New("error parsing response: " + err.Error())
	}

	err = json.Unmarshal(responseBody, &summary)
	if err != nil {
		return summary, errors.New("error parsing summary: " + err.Error())
	}

	return summary, nil
}
