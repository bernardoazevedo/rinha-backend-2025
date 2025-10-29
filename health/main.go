package main

import (
	"github.com/bernardoazevedo/rinha-backend-2025/api/health"
	"github.com/bernardoazevedo/rinha-backend-2025/api/key"
)

func main() {
	key.GetNewClient()
	health.HealthWorker()
}
