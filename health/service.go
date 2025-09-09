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

	"github.com/bernardoazevedo/rinha-de-backend-2025/key"
	"github.com/bernardoazevedo/rinha-de-backend-2025/logger"
)

func CheckHealth() (string, error) {
	paymentDefaultUrl := os.Getenv("PAYMENT_DEFAULT_URL")
	paymentFallbackUrl := os.Getenv("PAYMENT_FALLBACK_URL")

	url := paymentDefaultUrl

	defaultHealth, err := check(paymentDefaultUrl)
	if err != nil {
		defaultHealth.Failing = true
	}
	// logger.Add("\tfailing? " + fmt.Sprint(defaultHealth.Failing) + " - minResponseTime: " + fmt.Sprint(defaultHealth.MinResponseTime))

	fallbackHealth, err := check(paymentFallbackUrl)
	if err != nil {
		fallbackHealth.Failing = true
	}
	// logger.Add("\tfailing? " + fmt.Sprint(fallbackHealth.Failing) + " - minResponseTime: " + fmt.Sprint(fallbackHealth.MinResponseTime))

	if defaultHealth.Failing && !fallbackHealth.Failing {
		url = paymentFallbackUrl
	} else if !defaultHealth.Failing && fallbackHealth.Failing {
		url = paymentDefaultUrl
	} else if !defaultHealth.Failing && !fallbackHealth.Failing {
		url = paymentDefaultUrl
	} else { // both offline
		return "", errors.New("no payment service online, try again in a few moments")
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
	url := ""
	newUrl := returnOnlineUrl()

	if newUrl != url {
		url = newUrl

		err := key.Set("url", url)
		if err != nil {
			return url, errors.New("error saving url: " + err.Error())
		} else {
			return url, nil
		}
	}

	return url, nil
}

func returnOnlineUrl() string {
	newUrl := ""
	var err error
	// for {
		newUrl, err = CheckHealth()
		if err == nil && newUrl != "" {
			//success
			// break
		} else {
			// error checking, will try again
			logger.Add("error finding a service online, trying again...")
			time.Sleep(time.Second / 2)
		}
	// }
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
