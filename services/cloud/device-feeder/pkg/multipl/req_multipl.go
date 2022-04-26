package multipl

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/cloud/device-feeder/pkg"
	"strings"
)

type requestMultiplier struct {
	registryClient RegistryClient
	queue          QueuePublisher
}

func NewRequestMultiplier(registryClient RegistryClient, queue QueuePublisher) pkg.RequestMultiplier {
	return &requestMultiplier{
		registryClient: registryClient,
		queue:          queue,
	}
}

func (r *requestMultiplier) Process(req *pkg.DevicesUpdateRequest) error {
	// "org.nodeId"
	segments := strings.Split(req.Target, ".")
	if len(segments) != 2 {
		return fmt.Errorf("Invalid format of target: %s", req.Target)
	}

	orgName := segments[0]
	nodeId := segments[1]

	if nodeId != "*" {
		return fmt.Errorf("device id in target is not supported")
	}

	nodes, err := r.registryClient.GetNodesList(orgName)
	if err != nil {
		return err
	}

	logrus.Infof("Creating requests for %d nodes", len(nodes))
	counter := 0
	for _, n := range nodes {
		err = r.queue.Publish(pkg.DevicesUpdateRequest{
			Target:     orgName + "." + n.NodeId,
			HttpMethod: req.HttpMethod,
			Path:       req.Path,
			Body:       req.Body,
		})

		if err != nil {
			logrus.Errorf("Failed to publish message to queue: %s", err)
			return fmt.Errorf("failed to publish message to queue")
		} else {
			counter++
		}
	}
	logrus.Infof("Created %d requests", counter)
	return nil
}
