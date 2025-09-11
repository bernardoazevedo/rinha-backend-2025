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

	"github.com/bernardoazevedo/rinha-de-backend-2025/logger"
)

var PostUrl string

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

func CheckHealth() (string, error) {
	paymentDefaultUrl := os.Getenv("PAYMENT_DEFAULT_URL")
	paymentFallbackUrl := os.Getenv("PAYMENT_FALLBACK_URL")

	url := paymentDefaultUrl

	defaultHealth, err := check(paymentDefaultUrl)
	if err != nil {
		defaultHealth.Failing = true
	}

	fallbackHealth, err := check(paymentFallbackUrl)
	if err != nil {
		fallbackHealth.Failing = true
	}

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

func CheckSetReturnUrl() (string, error) {
	onlineUrl, err := CheckHealth()

	if err == nil && onlineUrl != "" {
		PostUrl = onlineUrl
		return onlineUrl, nil

	} else {
		return "", errors.New("error finding a service online")
	}
}

func HealthWorker() {
	paymentDefaultUrl := os.Getenv("PAYMENT_DEFAULT_URL")
	paymentFallbackUrl := os.Getenv("PAYMENT_FALLBACK_URL")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		for {
			defaultHealth, err := check(paymentDefaultUrl)
			if err != nil {
				defaultHealth.Failing = true
			}
			log := fmt.Sprintf(" default: online: %t - minResponseTime: %d", !defaultHealth.Failing, defaultHealth.MinResponseTime)
			logger.Add(log)

			fallbackHealth, err := check(paymentFallbackUrl)
			if err != nil {
				fallbackHealth.Failing = true
			}
			log = fmt.Sprintf("fallback: online: %t - minResponseTime: %d", !fallbackHealth.Failing, fallbackHealth.MinResponseTime)
			logger.Add(log)

			time.Sleep(time.Second * 5)
		}
	}()

	logger.Add("Monitoring services health...")
	<-sigchan

	logger.Add("Killed, shutting down")
}
