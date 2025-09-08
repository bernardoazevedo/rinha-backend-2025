package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bernardoazevedo/rinha-de-backend-2025/health"
	"github.com/bernardoazevedo/rinha-de-backend-2025/key"
	"github.com/bernardoazevedo/rinha-de-backend-2025/logger"
	"github.com/bernardoazevedo/rinha-de-backend-2025/payment"
	"github.com/bernardoazevedo/rinha-de-backend-2025/summary"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.SetPrefix("main: ")
	log.SetFlags(0)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	paymentDefaultUrl := os.Getenv("PAYMENT_DEFAULT_URL")

	key.Set("url", paymentDefaultUrl)

	// go health.HealthWorker()
	go func() {
		url := ""
		newUrl := ""
		var err error

		for {
			newUrl, err = health.CheckSetReturnUrl()
			if err != nil {
				fmt.Println(err.Error())
			}

			if newUrl != url {
				url = newUrl
				fmt.Println("New url: " + newUrl)
				logger.Add("\t\tworker: " + newUrl)
			}
			time.Sleep(time.Second * 5)
		}
	}()

	router := gin.Default()
	router.POST("/payments", payment.Payments)
	router.GET("/payments-summary", summary.PaymentsSummary)

	router.Run(":1234")
}
