package main

import (
	"log"

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

	//init workers

	router := gin.Default()
	router.POST("/payments", payment.Payments)
	router.GET("/payments-summary", summary.PaymentsSummary)

	router.Run(":1234")
}
