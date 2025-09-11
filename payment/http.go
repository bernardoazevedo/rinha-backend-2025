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


	// tentando enviar diretamente
	alreadyExistsPayment, err := postPayment(paymentJson)
	if alreadyExistsPayment {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "this payment already exists"})
		return

	} else if err != nil { 
		// deu erro, coloco na fila pra poder tentar de novo
		c.IndentedJSON(http.StatusOK, gin.H{"message": "de primeira n deu, vou botar na fila"})
		err = queuePayment(payment)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		return

	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "success!"})
	}
}
