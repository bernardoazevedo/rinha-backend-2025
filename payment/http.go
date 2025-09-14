package payment

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Payments(c *gin.Context) {
	var payment Payment

	err := c.BindJSON(&payment)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error binding payment"})
		return
	}
	payment.RequestedAt = time.Now().UTC().Format(time.RFC3339Nano)

	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error binding payment back"})
		return
	}
	paymentJson := string(paymentBytes)


	err = queuePayment(paymentJson)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}
	c.IndentedJSON(http.StatusOK, "")
}