package payment

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func Payments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payment Payment

	responseBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(responseBody, &payment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	payment.RequestedAt = time.Now().UTC().Format(time.RFC3339Nano)

	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	paymentJson := string(paymentBytes)


	w.Header().Set("Content-Type", "application/json")

	err = queuePayment(paymentJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}