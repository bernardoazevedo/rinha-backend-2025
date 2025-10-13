package main

import (
	"log"
	"net/http"

	"github.com/bernardoazevedo/rinha-backend-2025/health"
	"github.com/bernardoazevedo/rinha-backend-2025/payment"
	"github.com/bernardoazevedo/rinha-backend-2025/summary"
)

func main() {
	log.SetPrefix("main: ")
	log.SetFlags(0)

	health.PostUrl = "http://payment-processor-default:8080"

	go health.HealthWorker()
	go payment.PaymentWorker()
	go payment.PaymentWorker()

	http.HandleFunc("/payments", payment.Payments)
	http.HandleFunc("/payments-summary", summary.PaymentsSummary)

	log.Fatal(http.ListenAndServe(":1234", nil))
}
