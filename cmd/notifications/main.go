package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/marwan562/fintech-ecosystem/pkg/messaging"
	"github.com/marwan562/fintech-ecosystem/pkg/monitoring"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	NotificationsSent = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "notifications_sent_total",
		Help: "Total number of notifications sent.",
	}, []string{"type", "status"})
)

type NotificationTask struct {
	Type    string `json:"type"` // "email" or "sms"
	To      string `json:"to"`
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body"`
}

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://user:password@localhost:5672/"
	}

	client, err := messaging.NewRabbitMQClient(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer client.Close()

	queue, err := client.DeclareQueue("notifications")
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Start Metrics Server
	monitoring.StartMetricsServer(":8084")

	log.Println("Notifications Service started. Waiting for tasks...")

	client.Consume(queue.Name, func(body []byte) error {
		var task NotificationTask
		if err := json.Unmarshal(body, &task); err != nil {
			return err
		}

		log.Printf("Sending %s to %s...", task.Type, task.To)
		log.Printf("Body: %s", task.Body)

		// In a real system, call SendGrid, Twilio, etc.
		log.Printf("Successfully sent %s to %s", task.Type, task.To)
		NotificationsSent.WithLabelValues(task.Type, "success").Inc()

		return nil
	})

	// Keep main running
	select {}
}
