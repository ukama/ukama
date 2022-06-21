package pkg

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/ukama/ukama/services/cloud/mailer/pkg/metrics"
	"github.com/ukama/ukama/services/common/errors"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/wagslane/go-rabbitmq"
)

const deadLetterExchangeHeaderName = "x-dead-letter-exchange"

type Mailer struct {
	queueConf *QueueConfig
	Mail      *Sender
	serviceId string
}

func NewMailer(queueConf *QueueConfig, mail *Sender) (*Mailer, error) {
	return &Mailer{
		queueConf: queueConf,
		Mail:      mail,
		serviceId: os.Getenv("POD_NAME"),
	}, nil
}

// Start queue listening
func (m *Mailer) Start() error {
	err := m.declareQueueTopology(m.queueConf.Uri)
	if err != nil {
		logrus.Fatalf("Failed to declare queue topology: %s", err)
	}

	client, err := rabbitmq.NewConsumer(m.queueConf.Uri, rabbitmq.Config{},
		rabbitmq.WithConsumerOptionsLogger(logrus.WithField("service", ServiceName)),
		rabbitmq.WithConsumerOptionsReconnectInterval(5*time.Second))

	if err != nil {
		logrus.Fatalf("error creating queue consumer. Error: %s", err.Error())
	}
	queueArgs := map[string]interface{}{
		deadLetterExchangeHeaderName: m.getDeadLetterName(),
	}
	err = client.StartConsuming(m.incomingMessageHandler, m.queueConf.QueueName, []string{},
		rabbitmq.WithConsumeOptionsConsumerName(m.serviceId),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQueueArgs(queueArgs),
	)

	if err != nil {
		logrus.Fatalf("Error subscribing for a queue messages. Error: %+v", err)
	}

	logrus.Info("Listening for messages...")
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	logrus.Info("Shutting down...")
	err = client.Close()
	if err != nil {
		logrus.Errorf("Error closing consumer: %s", err.Error())
	}

	return nil
}

func (q *Mailer) getDeadLetterName() string {
	return q.queueConf.QueueName + "-dead-letter"
}

func (q *Mailer) declareQueueTopology(queueUri string) error {
	conn, err := amqp.Dial(queueUri)
	if err != nil {
		return errors.Wrap(err, "error connecting to queue")
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "error creating amqp channel")
	}
	defer ch.Close()

	// declare dead letter exchange
	err = ch.ExchangeDeclare(q.getDeadLetterName(), amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "error declaring dead letter exchange")
	}

	// declare and bind dead letter queue
	_, err = ch.QueueDeclare(q.getDeadLetterName(), true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "error declaring dead letter queue")
	}

	err = ch.QueueBind(q.getDeadLetterName(), "", q.getDeadLetterName(), false, nil)
	if err != nil {
		return errors.Wrap(err, "error binding dead letter queue")
	}

	return nil
}

func (m *Mailer) incomingMessageHandler(delivery rabbitmq.Delivery) rabbitmq.Action {
	mail := msgbus.MailMessage{}
	err := json.Unmarshal(delivery.Body, &mail)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %s. Error: %s", string(delivery.Body), err)
		metrics.EmailSentFailureRequestMetric()
		return rabbitmq.NackDiscard
	}

	err = m.Mail.SendEmail(&mail)
	if err != nil {
		logrus.Errorf("Failed to send email: %s. Error: %s", string(delivery.Body), err)
		metrics.EmailSentFailureRequestMetric()
		// beware that delivery tag is Channel scoped so it won't work for multiple consumers
		if delivery.DeliveryTag >= m.queueConf.RetryAttempts {
			logrus.Errorf("Failed to send email: %s. Error: %s. Discarding message", string(delivery.Body), err)
			return rabbitmq.NackDiscard
		}

		return rabbitmq.NackRequeue
	}

	metrics.EmailSentSuccessfulRequestMetric()
	return rabbitmq.Ack
}

func (m *Mailer) Close() {

}
