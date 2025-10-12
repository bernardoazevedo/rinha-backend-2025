package summary

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func PaymentsSummary(ctx *fasthttp.RequestCtx) {
	summary := make(map[string]Summary)
	paymentDefaultUrl := "http://payment-processor-default:8080"
	paymentFallbackUrl := "http://payment-processor-fallback:8080"

	queryArgs := ctx.QueryArgs()
	from := string(queryArgs.Peek("from"))
	to := string(queryArgs.Peek("to"))

	defaultSummary, err := getPaymentsSummary(paymentDefaultUrl, from, to)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	summary["default"] = defaultSummary

	fallbackSummary, err := getPaymentsSummary(paymentFallbackUrl, from, to)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	summary["fallback"] = fallbackSummary

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)

	summaryJson, err := json.Marshal(summary)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.Write(summaryJson)
}
