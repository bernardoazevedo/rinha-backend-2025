package health

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"
)

var PostUrl string

func CheckHealth() (string, error) {
	paymentDefaultUrl := "http://payment-processor-default:8080"
	paymentFallbackUrl := "http://payment-processor-fallback:8080"

	url := paymentDefaultUrl

	defaultCh, fallbackCh := make(chan Health), make(chan Health)

	go func() {
		defaultHealth, err := check(paymentDefaultUrl)
		if err != nil {
			defaultHealth.Failing = true
		}
		defaultCh <- defaultHealth
	}()

	go func() {
		fallbackHealth, err := check(paymentFallbackUrl)
		if err != nil {
			fallbackHealth.Failing = true
		}
		fallbackCh <- fallbackHealth
	}()

	defaultHealth, fallbackHealth := <-defaultCh, <-fallbackCh

	if (defaultHealth.MinResponseTime <= fallbackHealth.MinResponseTime) && !defaultHealth.Failing {
		url = paymentDefaultUrl
	} else {
		url = paymentFallbackUrl
	}

	return url, nil
}

func check(url string) (Health, error) {
	var health Health

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	req.Header.SetRequestURI(url + "/payments/service-health")
	response := fasthttp.AcquireResponse()

	err := fasthttp.Do(req, response)
	if err != nil {
		return health, errors.New("error during request")
	}
	defer fasthttp.ReleaseRequest(req)

	responseBody := response.Body()
	err = json.Unmarshal(responseBody, &health)
	if err != nil {
		return health, errors.New("error parsing health")
	}

	return health, nil
}

func CheckSetReturnUrl() (string, error) {
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

	if newUrl != "" {
		PostUrl = newUrl
	}

	return newUrl, nil
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
