package msgbus

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/protobuf/proto"
)

type QPub interface {
	Publish(payload any, routingKey string) error
	PublishProto(payload proto.Message, routingKey string) error
	PublishToQueue(queueName string, payload any) error
	Close() error
}

// QPub is a simplified AMQP client that publishes messages to a default exchange
// Client reconnect in case of connection loss.
type qPub struct {
	conn        *rabbitmq.Conn
	publisher   *rabbitmq.Publisher
	serviceName string
	instanceId  string
}

func NewQPub(queueUri string, serviceName string, instanceId string) (*qPub, error) {
	conn, err := rabbitmq.NewConn(
		queueUri,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Infof("Error creating publisher %s.", err.Error())
		return nil, err
	}

	publisher, err := rabbitmq.NewPublisher(conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}

	return &qPub{
		conn:        conn,
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

func (q *qPub) PublishProto(payload proto.Message, routingKey string) error {

	b, err := proto.Marshal(payload)
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
	q.conn.Close()
	return nil
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
