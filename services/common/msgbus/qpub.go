package msgbus

import (
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
)

type QPub interface {
	Publish(payload any, routingKey string) error
}

// QPub is a simplified AMQP client that publishes messages to a default exchange
// Client reconnect in case of connection loss.
type qPub struct {
	publisher *rabbitmq.Publisher
}

func NewQPub(queueUri string) (*qPub, error) {
	publisher, err := rabbitmq.NewPublisher(queueUri, rabbitmq.Config{})
	if err != nil {
		return nil, err
	}

	return &qPub{
		publisher: publisher,
	}, nil
}

func (q *qPub) Publish(payload any, routingKey string) error {

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = q.publisher.Publish(b, []string{routingKey},
		rabbitmq.WithPublishOptionsHeaders(map[string]interface{}{}),
		rabbitmq.WithPublishOptionsExchange(DefaultExchange))

	if err != nil {
		return err
	}

	return nil
}
