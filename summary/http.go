package summary

import (
	"encoding/json"
	"net/http"
)

func PaymentsSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	summary := make(map[string]Summary)
	paymentDefaultUrl := "http://payment-processor-default:8080"
	paymentFallbackUrl := "http://payment-processor-fallback:8080"

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	defaultSummary, err := getPaymentsSummary(paymentDefaultUrl, from, to)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	summary["default"] = defaultSummary
	
	fallbackSummary, err := getPaymentsSummary(paymentFallbackUrl, from, to)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	summary["fallback"] = fallbackSummary

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	summaryJson, err := json.Marshal(summary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	w.Write(summaryJson)
}
