package pkg

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/wagslane/go-rabbitmq"
	"os"
	"os/signal"
	"syscall"
)

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

	client, err := rabbitmq.NewConsumer(m.queueConf.Uri, rabbitmq.Config{},
		rabbitmq.WithConsumerOptionsLogger(logrus.WithField("service", ServiceName)))
	if err != nil {
		logrus.Fatalf("error creating queue consumer. Error: %s", err.Error())
	}

	err = client.StartConsuming(m.incomingMessageHandler, m.queueConf.QueueName, []string{},
		rabbitmq.WithConsumeOptionsConsumerName(m.serviceId),
		rabbitmq.WithConsumeOptionsQueueDurable,
	)

	if err != nil {
		logrus.Fatalf("Error subscribing for a queue messages. Error: %+v", err)
	}

	logrus.Info("Listening for messages...")
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	return nil
}

func (m *Mailer) incomingMessageHandler(delivery rabbitmq.Delivery) rabbitmq.Action {
	mail := msgbus.MailMessage{}
	err := json.Unmarshal(delivery.Body, &mail)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %s. Error: %s", string(delivery.Body), err)
		return rabbitmq.NackDiscard
	}

	err = m.Mail.SendEmail(&mail)
	if err != nil {
		logrus.Errorf("Failed to send email: %s. Error: %s", string(delivery.Body), err)
		return rabbitmq.NackRequeue
	}

	return rabbitmq.Ack
}
