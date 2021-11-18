package pkg

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/common/msgbus"
)

const MissingRoutingKeyMessage = "missing routing key segment"

type HssQueue interface {
	SendImsiAddedEvent(imsi string)
	SendImsiUpdateEvent(imsi string)
	SendImsiDeleteEvent(imsi string)
}

type hssQueue struct {
	publisherClient   msgbus.Publisher
	connectionString  string
	routingKeyBuilder msgbus.RoutingKeyBuilder
}

func NewHssQueue(queueConfig config.Queue) *hssQueue {
	return &hssQueue{
		connectionString: queueConfig.Uri,
		routingKeyBuilder: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer("hss").
			SetObject("imsi"),
	}
}

func (h *hssQueue) getPub() (client msgbus.Publisher, err error) {
	if h.publisherClient == nil {
		logrus.Info("Connecting to the queue")
		h.publisherClient, err = msgbus.NewPublisherClient(h.connectionString)
		if err != nil {
			return nil, err
		}
	}

	return h.publisherClient, nil
}

func (h *hssQueue) SendImsiAddedEvent(imsi string) {
	logrus.Infoln("Sending imsi add event for imsi: ", imsi)
	rk, err := h.routingKeyBuilder.SetActionCreate().Build()
	if err != nil {

		logrus.Fatalf(MissingRoutingKeyMessage)
	}

	h.sendEventAndLogError(imsi, rk)
}

func (h *hssQueue) SendImsiUpdateEvent(imsi string) {
	logrus.Infoln("Sending imsi update event for imsi: ", imsi)
	rk, err := h.routingKeyBuilder.SetActionUpdate().Build()
	if err != nil {
		logrus.Fatalf(MissingRoutingKeyMessage)
	}
	h.sendEventAndLogError(imsi, rk)
}

func (h *hssQueue) SendImsiDeleteEvent(imsi string) {
	logrus.Infoln("Sending imsi delete event for imsi: ", imsi)
	rk, err := h.routingKeyBuilder.SetActionDelete().Build()
	if err != nil {
		logrus.Fatalf(MissingRoutingKeyMessage)
	}
	h.sendEventAndLogError(imsi, rk)
}

func (h *hssQueue) sendEventAndLogError(imsi string, rk string) {
	err := h.sendEvent(imsi, rk)
	if err != nil {
		logrus.Errorln("error sending event. Error: ", err)
	}
}

func (h *hssQueue) sendEvent(imsi string, rk string) error {

	conn, err := h.getPub()
	if err != nil {
		return errors.Wrap(err, "error sending event")
	}

	s := struct {
		Imsi string
	}{
		Imsi: imsi,
	}

	err = conn.PublishOnExchange(msgbus.DefaultExchange, rk, s)
	if err != nil {
		return errors.Wrap(err, "error sending message to queue")
	}
	return nil
}
