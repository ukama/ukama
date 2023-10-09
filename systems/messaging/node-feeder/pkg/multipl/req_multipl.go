package multipl

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"

	"github.com/ukama/ukama/systems/messaging/node-feeder/pkg"
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

	func (r *requestMultiplier) Process(req *cpb.NodeFeederMsg) error {
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

		nodes, err := r.registryClient.GetNodesList()
		if err != nil {
			return err
		}

		logrus.Infof("Creating requests for %d nodes", len(nodes))
		counter := 0
		for _, n := range nodes {
			err = r.queue.Publish(&cpb.NodeFeederMsg{
				Target:     orgName + "." + n.Id,
				HTTPMethod: req.HTTPMethod,
				Path:       req.Path,
				Msg:       req.Msg,
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
