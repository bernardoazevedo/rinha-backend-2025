package summary

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func PaymentsSummary(c *gin.Context) {
	paymentDefaultUrl := os.Getenv("PAYMENT_DEFAULT_URL")
	paymentFallbackUrl := os.Getenv("PAYMENT_FALLBACK_URL")

	from := c.Query("from")
	to := c.Query("to")

	defaultSummary, err := getPaymentsSummary(paymentDefaultUrl, from, to)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fallbackSummary, err := getPaymentsSummary(paymentFallbackUrl, from, to)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"default": defaultSummary, "fallback": fallbackSummary})
}

