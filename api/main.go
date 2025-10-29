package main

import (
	"log"

	"github.com/bernardoazevedo/rinha-backend-2025/api/key"
	"github.com/bernardoazevedo/rinha-backend-2025/api/payment"
	"github.com/bernardoazevedo/rinha-backend-2025/api/summary"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	log.SetPrefix("main: ")
	log.SetFlags(0)

	key.GetNewClient()

	err := key.Set("url", "http://payment-processor-default:8080")
	if err != nil {
		panic(err)
	}

	r := router.New()
	r.POST("/payments", callPayments)
	r.GET("/payments-summary", callPaymentsSummary)

	log.Fatal(fasthttp.ListenAndServe(":1234", r.Handler))
}

func callPayments(ctx *fasthttp.RequestCtx) {
	payment.Payments(ctx)
}

func callPaymentsSummary(ctx *fasthttp.RequestCtx) {
	summary.PaymentsSummary(ctx)
}
