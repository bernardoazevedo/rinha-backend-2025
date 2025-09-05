package payment

import (
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

	response, err := postPayment(payment)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": response})
}
