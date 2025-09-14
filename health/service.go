package health

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

)

var PostUrl string

func CheckHealth() (string, error) {
	paymentDefaultUrl := "http://payment-processor-default:8080"
	paymentFallbackUrl := "http://payment-processor-fallback:8080"

	url := paymentDefaultUrl

	defaultHealth, err := check(paymentDefaultUrl)
	if err != nil {
		defaultHealth.Failing = true
	}

	fallbackHealth, err := check(paymentFallbackUrl)
	if err != nil {
		fallbackHealth.Failing = true
	}

	if (defaultHealth.MinResponseTime <= fallbackHealth.MinResponseTime) && !defaultHealth.Failing {
		url = paymentDefaultUrl
	} else {
		url = paymentFallbackUrl
	}

	return url, nil
}

func check(url string) (Health, error) {
	var health Health
	response, err := http.Get(url + "/payments/service-health")
	if err != nil {
		return health, errors.New("error during request")
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return health, errors.New("error parsing body")
	}

	err = json.Unmarshal(responseBody, &health)
	if err != nil {
		return health, errors.New("error parsing health")
	}

	return health, nil
}

func CheckSetReturnUrl() (string, error) {
	newUrl := checkUntilReturn()

	if newUrl != "" {
		PostUrl = newUrl
	}

	return newUrl, nil
}

func checkUntilReturn() string {
	newUrl := ""
	var err error
	for {
		newUrl, err = CheckHealth()
		if err == nil && newUrl != "" {
			break
		} else {
			// error checking, will try again
			fmt.Println("error finding a service online, trying again...")
			time.Sleep(time.Second / 2)
		}
	}
	return newUrl
}

func HealthWorker() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		url := ""
		newUrl := ""
		var err error

		for {
			newUrl, err = CheckSetReturnUrl()
			if err != nil {
				fmt.Println(err.Error())
			}

			if newUrl != url {
				url = newUrl
				fmt.Println("New url: " + newUrl)
			}
			time.Sleep(time.Second * 5)
		}
	}()

	fmt.Println("Monitoring services health...")
	<-sigchan

	fmt.Println("Killed, shutting down")
}
