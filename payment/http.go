package payment

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Payments(c *gin.Context) {
	var payment Payment
	paymentDefaultUrl := os.Getenv("PAYMENT_DEFAULT_URL")
	paymentFallbackUrl := os.Getenv("PAYMENT_FALLBACK_URL")

	err := c.BindJSON(&payment)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error binding payment"})
		return
	}
	payment.RequestedAt = time.Now().UTC().Format(time.RFC3339Nano)

	response, err := postPayment(payment, paymentDefaultUrl+"/payments")
	if err != nil {

		response, err = postPayment(payment, paymentFallbackUrl+"/payments")
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": response})
}
