package main

import (
	"fmt"
	"log"
	"time"

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

func uniqid(prefix string) string {
	now := time.Now()
	sec := now.Unix()
	usec := now.UnixNano() % 0x100000
	return fmt.Sprintf("%s%08x%05x", prefix, sec, usec)
}
