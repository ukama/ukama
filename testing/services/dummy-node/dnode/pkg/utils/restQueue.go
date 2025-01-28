package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitMQURL  = "amqp://guest:guest@10.1.0.3:15672/" // Replace with your RabbitMQ URL
	exchangeName = "amq.topic"                          // Name of the exchange
	exchangeType = "topic"                              // Exchange type: direct, fanout, topic, etc.
)

// RabbitMQ connection struct
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Create a new RabbitMQ connection
func NewRabbitMQ() (*RabbitMQ, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare the exchange
	err = channel.ExchangeDeclare(
		exchangeName, // exchange name
		exchangeType, // exchange type
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

// Publish a message to the RabbitMQ exchange
func (r *RabbitMQ) Publish(message []byte, routingKey string) error {
	return r.channel.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
}

// Close the RabbitMQ connection and channel
func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

func RestPublish(msg []byte, routingKey string) {
	// Initialize RabbitMQ connection
	rabbitMQ, err := NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Initialize Gin router
	router := gin.Default()

	// REST endpoint to publish a message
	router.POST("/publish", func(c *gin.Context) {

		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := rabbitMQ.Publish(msg, routingKey); err != nil {
			log.Printf("Failed to publish message: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message"})
			return
		}

		log.Printf("Message published to exchange '%s' with routing key '%s': %+v", exchangeName, routingKey, msg)
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Message published"})
	})

	// Start the HTTP server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
