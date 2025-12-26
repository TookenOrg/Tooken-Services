package services

import (
	"context"
	"log"
)

func StartWorkers(ctx context.Context) {
	log.Println("🔄 Background workers started (Redis/RabbitMQ)")

	// Simuler workers
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Println("💼 Worker tick: blockchain, trading, payments...")
		}
	}
}
