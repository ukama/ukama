package pkg

import (
	"encoding/json"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/cloud/device-feeder/pkg/global"
	"github.com/ukama/ukamaX/cloud/device-feeder/pkg/metrics"
	"github.com/ukama/ukamaX/common/msgbus"
	"github.com/wagslane/go-rabbitmq"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const deadLetterExchangeName = "device-feeder.dead-letter"
const deadLetterExchangeHeaderName = "x-dead-letter-exchange"
const errorCreatingWaitingQueueErr = "error declaring waiting queue"
const deadLetterRoutingKeyHeaderName = "x-dead-letter-routing-key"

type QueueListener struct {
	consumer       *rabbitmq.Consumer
	serviceId      string
	requestMult    RequestMultiplier
	requestExec    RequestExecutor
	maxRetryCount  int64
	retryPeriodSec int
	listenerConfig ListenerConfig
}

type RequestMultiplier interface {
	Process(body *DevicesUpdateRequest) error
}

type DevicesUpdateRequest struct {
	Target     string `json:"target"` // Target devices in form of "organization.network.device-id". Device id and network could be wildcarded
	HttpMethod string `json:"httpMethod"`
	Path       string `json:"path"`
	Body       string `json:"body"`
}

func NewQueueListener(queueUri string, serviceId string, requestMult RequestMultiplier, requestExec RequestExecutor, conf ListenerConfig) (*QueueListener, error) {
	consumer, err := rabbitmq.NewConsumer(queueUri, amqp.Config{},
		rabbitmq.WithConsumerOptionsLogger(log.WithField("service", "rabbitmq")))
	if err != nil {
		return nil, errors.Wrap(err, "error creating queue consumer")
	}

	q := &QueueListener{
		consumer:       &consumer,
		serviceId:      serviceId,
		requestMult:    requestMult,
		requestExec:    requestExec,
		listenerConfig: conf,
	}

	err = q.declareQueueTopology(queueUri)
	if err != nil {
		return nil, errors.Wrap(err, "error declaring queue topology")
	}

	return q, nil
}

func (q *QueueListener) declareQueueTopology(queueUri string) error {
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

	err = ch.ExchangeDeclare(deadLetterExchangeName, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "error declaring dead letter exchange")
	}

	waitingQueue, err := q.createWaitingQueue(ch)
	if err != nil {
		if err.(*amqp.Error).Code == amqp.PreconditionFailed {
			log.Warnf("Did you change waiting queue message TTL? If so then delete queue manually")
		}

		return errors.Wrap(err, errorCreatingWaitingQueueErr)
	}

	err = ch.QueueBind(waitingQueue.Name, "*", deadLetterExchangeName, false, nil)
	if err != nil {
		return errors.Wrap(err, "error binding waiting-queue")
	}

	// data feeder queue
	dataFeederQueue, err := ch.QueueDeclare(
		"device-feeder", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		map[string]interface{}{
			deadLetterExchangeHeaderName:   deadLetterExchangeName,
			deadLetterRoutingKeyHeaderName: string(msgbus.DeviceFeederRequestRoutingKey),
		}, // arguments
	)
	if err != nil {
		return errors.Wrap(err, errorCreatingWaitingQueueErr)
	}
	err = ch.QueueBind(dataFeederQueue.Name, string(msgbus.DeviceFeederRequestRoutingKey), msgbus.DefaultExchange, false, nil)

	if err != nil {
		return errors.Wrap(err, errorCreatingWaitingQueueErr)
	}

	return nil
}

func (q *QueueListener) createWaitingQueue(ch *amqp.Channel) (amqp.Queue, error) {

	// TODO: set TTL via policy
	waitingQueue, err := ch.QueueDeclare(
		"device-feeder.waiting-queue", // name
		true,                          // durable
		false,                         // delete when unused
		false,                         // exclusive
		false,                         // no-wait
		map[string]interface{}{
			"x-message-ttl":                q.retryPeriodSec * 1000,
			deadLetterExchangeHeaderName:   msgbus.DefaultExchange,
			deadLetterRoutingKeyHeaderName: string(msgbus.DeviceFeederRequestRoutingKey),
		},
	)
	return waitingQueue, err
}

func (q *QueueListener) StartQueueListening() (err error) {

	queueArgs := map[string]interface{}{
		deadLetterExchangeHeaderName:   deadLetterExchangeName,
		deadLetterRoutingKeyHeaderName: string(msgbus.DeviceFeederRequestRoutingKey),
	}

	err = q.consumer.StartConsuming(q.incomingMessageHandler, global.QueueName, []string{string(msgbus.DeviceFeederRequestRoutingKey)},
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsConsumerName(q.serviceId),
		rabbitmq.WithConsumeOptionsQueueArgs(queueArgs),
		rabbitmq.WithConsumeOptionsBindingExchangeName(msgbus.DefaultExchange),
		rabbitmq.WithConsumeOptionsBindingExchangeKind(amqp.ExchangeTopic),
		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
		rabbitmq.WithConsumeOptionsConcurrency(q.listenerConfig.Threads))

	if err != nil {
		log.Errorf("Error subscribing for a queue messages. Error: %+v", err)
		return err
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	return nil
}

func (q *QueueListener) incomingMessageHandler(delivery rabbitmq.Delivery) rabbitmq.Action {

	if q.isRetryLimitReached(delivery) {
		metrics.RecordFailedRequestMetric()
		return rabbitmq.Ack
	}

	err := q.processRequest(delivery)
	if err != nil {
		log.Errorf("Error processing request. Error: %+v", err)
		metrics.RecordFailedRequestMetric()
		return rabbitmq.NackDiscard
	}
	metrics.RecordSuccessfulRequestMetric()

	return rabbitmq.Ack
}

func (q *QueueListener) isRetryLimitReached(delivery rabbitmq.Delivery) bool {
	const deathHeader = "x-death"
	if delivery.Headers[deathHeader] == nil {
		return false
	}
	death := delivery.Headers[deathHeader].([]interface{})
	for _, d := range death {
		vals := d.(amqp.Table)
		if vals == nil {
			log.Errorf("Unexpected format of death header")
			return false
		}

		if vals["exchange"] == deadLetterExchangeName {
			count := vals["count"].(int64)
			if count > q.maxRetryCount {
				log.Infof("Retry limit reached for message: %v, target: %v", delivery.MessageId, delivery.Headers[global.OptionalTargetHeaderName])
				return true
			} else {
				log.Infof("Retry count: %v, target: %v", count, delivery.Headers[global.OptionalTargetHeaderName])
				return false
			}
		}
	}

	log.Warning("Cannot get retry count from message headers")
	return false
}

// return error only if it could be fixed by retry
// malformed request should not be considered as error
func (q *QueueListener) processRequest(delivery rabbitmq.Delivery) error {
	request := &DevicesUpdateRequest{}
	err := json.Unmarshal(delivery.Body, request)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		return nil
	}

	log.Infof("Received request: %+v", request)
	if strings.HasSuffix(request.Target, "*") {
		log.Infof("Wildcarded target: %s", request.Target)
		err = q.requestMult.Process(request)
		return err
	} else {
		log.Infof("Direct node target: %s", request.Target)
		err = q.requestExec.Execute(request)

		if err != nil {
			if dErr, ok := err.(Device4xxServerError); ok {
				log.Warningf("Request failed but won't be retried. Error %v", dErr)
				return nil
			}
		}

		return err
	}
}
