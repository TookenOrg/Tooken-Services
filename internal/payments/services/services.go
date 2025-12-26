// internal/blockchain/services/service.go
package services

import "github.com/gin-gonic/gin"

type Service struct {
	// Vos fields blockchain (db, eth client, etc.)
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AddPayment(c *gin.Context) {
	return
}
