package payment

import (
	"encoding/json"
	"time"

	"github.com/valyala/fasthttp"
)

func Payments(ctx *fasthttp.RequestCtx) {
	var payment Payment

	post := ctx.PostBody()

	err := json.Unmarshal(post, &payment)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	payment.RequestedAt = time.Now().UTC().Format(time.RFC3339Nano)

	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	paymentJson := string(paymentBytes)

	ctx.SetContentType("application/json")

	// err = queuePayment(paymentJson)
	_, err = postPayment(paymentJson)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}
