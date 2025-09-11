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

var PostUrl 			string
var PaymentDefaultUrl 	string
var PaymentFallbackUrl 	string

func CheckHealth() (string, error) {
	var url string

	defaultHealth, err := check(PaymentDefaultUrl)
	if err != nil {
		defaultHealth.Failing = true
	}

	fallbackHealth, err := check(PaymentFallbackUrl)
	if err != nil {
		fallbackHealth.Failing = true
	}

	if defaultHealth.Failing && !fallbackHealth.Failing {
		url = PaymentFallbackUrl
	} else if !defaultHealth.Failing && fallbackHealth.Failing {
		url = PaymentDefaultUrl
	} else if !defaultHealth.Failing && !fallbackHealth.Failing {
		url = PaymentDefaultUrl
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
	PostUrl = checkUntilReturn()
	return PostUrl, nil
}

func checkUntilReturn() string {
	newUrl := ""
	var err error
	for i := 0; i < 3; i++ {
		newUrl, err = CheckHealth()
		if err == nil && newUrl != "" {
			break
		} else {
			// error checking, will try again
			logger.Add("error finding a service online, trying again..." + fmt.Sprintf("%d", i))
			for j := 0; j < i; j++ {
				time.Sleep(time.Second)
			}
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
