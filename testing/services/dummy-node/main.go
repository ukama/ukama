package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/ukama/ukamaX/common/msgbus"
	"google.golang.org/protobuf/proto"

	"github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
	commonpb "github.com/ukama/ukamaX/common/pb/gen/ukamaos/mesh"
	"github.com/wagslane/go-rabbitmq"
)

func main() {

	go func() {
		StartEchoServer()
	}()

	go func() {
		StartNetworkUpdates()
	}()

	StartSslEchoServer()

}

// sends "mesh.link" events to queue
func StartNetworkUpdates() {
	queueHost := os.Getenv("QUEUEURI")
	if queueHost == "" {
		queueHost = "amqp://guest:guest@rabbitmq"
	}
	logrus.Infof("Connecting to queue: %s", queueHost)

	nodeId := os.Getenv("NODE_ID")
	if nodeId == "" {
		panic("NODE_ID env var not set")
	}

	publisher, err := rabbitmq.NewPublisher(queueHost, amqp.Config{},
		rabbitmq.WithPublisherOptionsLogging)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	for {
		myIp := os.Getenv("MY_POD_IP")
		if myIp == "" {
			panic("MY_POD_IP env var not set")
		}
		link := commonpb.Link{
			NodeId: &nodeId,
			Ip:     &myIp,
		}

		b, err := proto.Marshal(&link)
		if err != nil {
			logrus.Errorf("Failed to marshal message. Error: %+v", err)
		}

		err = publisher.Publish(b, []string{string(msgbus.DeviceConnectedRoutingKey)},
			rabbitmq.WithPublishOptionsExchange(msgbus.DefaultExchange))
		if err != nil {
			logrus.Errorf("Failed to publish message. Error: %+v", err)
		}

		logrus.WithFields(logrus.Fields{
			"message": link.String(),
		}).Info("Published message")
		time.Sleep(time.Minute)
	}

}

func StartEchoServer() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

		logrus.Infoln("Request:")
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			logrus.Info(err)
		}

		fmt.Print(string(b))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infoln("Config Request:")
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			logrus.Info(err)
		}

		fmt.Print(string(b))
	})

	logrus.Infoln("Starting http server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// metrics
// var reg = prometheus.NewRegistry()

// var (
// 	activeUe = promauto.NewGauge((prometheus.GaugeOpts{
// 		Name: "epc_active_ue",
// 		Help: "The total of active users",
// 	}))
// )
