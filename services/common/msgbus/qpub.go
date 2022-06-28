package msgbus

import (
	"encoding/json"

	"github.com/wagslane/go-rabbitmq"
)

type QPub interface {
	Publish(payload any, routingKey string) error
	PublishToQueue(queueName string, payload any) error
	Close() error
}

// QPub is a simplified AMQP client that publishes messages to a default exchange
// Client reconnect in case of connection loss.
type qPub struct {
	publisher   *rabbitmq.Publisher
	serviceName string
	instanceId  string
}

func NewQPub(queueUri string, serviceName string, instanceId string) (*qPub, error) {
	publisher, err := rabbitmq.NewPublisher(queueUri, rabbitmq.Config{})
	if err != nil {
		return nil, err
	}

	return &qPub{
		publisher:   publisher,
		serviceName: serviceName,
		instanceId:  instanceId,
	}, nil
}

// Publish publishes a message in json format to the default topic exchange with a routing key specified
func (q *qPub) Publish(payload any, routingKey string) error {

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = q.publisher.Publish(b, []string{routingKey},
		rabbitmq.WithPublishOptionsHeaders(map[string]interface{}{
			"source-service": q.serviceName,
			"instance-id":    q.instanceId,
		}),
		rabbitmq.WithPublishOptionsExchange(DefaultExchange))

	if err != nil {
		return err
	}

	return nil
}

func (q *qPub) Close() error {
	return q.publisher.Close()
}

func (q *qPub) PublishToQueue(queueName string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = q.publisher.Publish(b, []string{queueName},
		rabbitmq.WithPublishOptionsHeaders(map[string]interface{}{
			"source-service": q.serviceName,
			"instance-id":    q.instanceId,
		}))

	if err != nil {
		return err
	}

	return nil
}
