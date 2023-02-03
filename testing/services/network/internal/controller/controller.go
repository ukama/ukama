package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/errors"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/testing/services/network/internal"
	"github.com/ukama/ukama/testing/services/network/internal/db"
	spec "github.com/ukama/ukama/testing/services/network/specs/controller/spec"
	"google.golang.org/protobuf/proto"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type NodeState string

const (
	VNodeBooting string = "BootingUp"
	VNodeActive  string = "Running"
	VNodeHalted  string = "Halted"
	VNodeFaulty  string = "Faulty"
	VNodeUnkown  string = "Unkown"
)

type ControllerOps interface {
	ControllerInit() error
	ListNodes()
	GetNodeRuntimeStatus(nodeId string) (*string, error)
	PowerOnNode(nodeId string, org string) error
	PowerOffNode(nodeId string) error
	CreateNode(name string, image string, command []string, ntype string, org string) error
	WatcherForNodes(ctx context.Context, cb func(string, string) error) error
}

type Controller struct {
	repo db.VNodeRepo
	cs   kubernetes.Interface
	ns   string
	m    msgbus.Publisher
}

func NewController(d db.VNodeRepo) *Controller {
	cset, err := connectToK8s()
	if err != nil {
		logrus.Fatalf("Build:: Can't connect to Kuberneets cluster. Err: %s", err.Error())
		return nil
	}

	msgC, err := msgbus.NewPublisherClient(internal.ServiceConfig.Queue.Uri)
	if err != nil {
		logrus.Errorf("error getting message publisher: %s\n", err.Error())
		return nil
	}

	ns := "default"

	if internal.ServiceConfig.Namespace != "" {
		ns = internal.ServiceConfig.Namespace
	}

	return &Controller{
		cs:   cset,
		ns:   ns,
		m:    msgC,
		repo: d,
	}
}

/* Connect to Kubernetes cluster */
func connectToK8s() (kubernetes.Interface, error) {

	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Errorf("Build:: Failed to create K8s config. Error %s", err.Error())
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Errorf("Build:: Failed to create K8s clientset. Error %s", err.Error())
		return nil, err
	}

	return clientset, nil
}

/* Return node name from ID */
func getVirtNodeName(id string) string {
	return "vn-" + id
}

/* Return VM name from ID */
func getVirtNodeId(name string) string {
	return strings.Trim(name, "vn-")
}

/* Starting build virtual node watcher routine */
func (c *Controller) ControllerInit() error {

	/* For listing already running virtual nodes in network  */
	pods, err := c.cs.CoreV1().Pods(c.ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("Network:: Error getting pods: %v\n", err)
		return err
	}
	for _, pod := range pods.Items {
		logrus.Tracef("Network:: Pod Name: %s\n", pod.Name)
	}

	return c.WatcherForNodes(context.TODO(), c.PublishEvent)
}

/* Get Node status from Pod phase */
func getNodeRuntimeStatus(phase v1.PodPhase) string {
	var state string
	switch phase {
	case v1.PodPending:
		state = VNodeBooting
	case v1.PodRunning:
		state = VNodeActive
	case v1.PodSucceeded:
		state = VNodeHalted
	case v1.PodFailed:
		state = VNodeFaulty
	case v1.PodUnknown:
		state = VNodeUnkown
	default:
		state = VNodeUnkown
	}
	return state
}

/* Pod spec

name: <node-id>
metadata:
	labels:
		node: <nodeid>
		type: <hnode|anode|tnode>
		org: <orgname>
		category: <vm|container>
spec:
	container:

	volume:
*/

func (c *Controller) CreateNode(nodeId string, image string, command []string, ntype string, org string) error {

	/* Virtual Node Name */
	vnName := getVirtNodeName(nodeId)

	/* Virt Nodes Instance Labels */
	labels := map[string]string{
		"node": vnName,
		"type": ntype,
		"app":  "virtual-node",
		"org":  org,
	}

	/* Pod spec */
	podSpec := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vnName,
			Namespace: c.ns,
			Labels:    labels,
		},
		Spec: v1.PodSpec{
			RestartPolicy: v1.RestartPolicyAlways,
			Containers: []v1.Container{
				v1.Container{
					Name:  nodeId,
					Image: image,
					//Command: command,
					Args: []string{},
					Env: []v1.EnvVar{
						{
							Name:  "UUID",
							Value: nodeId,
						},
						{
							Name:  "NODETYPE",
							Value: ntype,
						},
					},
				},
			},
			Hostname:                      nodeId,
			TerminationGracePeriodSeconds: &internal.ServiceConfig.TerminationGracePeriodSeconds,
			ActiveDeadlineSeconds:         &internal.ServiceConfig.ActiveDeadlineSeconds,
		},
	}

	/* Fixing context time 5 s */
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses

	_, err := c.cs.CoreV1().
		Pods(c.ns).
		Create(ctx, podSpec, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("PowerOn failure for node %s. Error: %s", nodeId, err.Error())
		return err
	}

	logrus.Infof("PowerOn poweron initiated for node %s", nodeId)

	return err

}

/* Debug List of pods */
func (c *Controller) ListNodes() {

	pods, err := c.cs.CoreV1().Pods(c.ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("Failed to get virtulal node runtime list. Error %s", err.Error())
	}

	for _, pod := range pods.Items {
		logrus.Debugf("Node Name %s Node Status %v.", pod.Name, pod.Status)
	}
}

/* Get Job status */
func (c *Controller) GetNodeRuntimeStatus(nodeId string) (*string, error) {

	/* Virtual Node Name */
	vnName := getVirtNodeName(nodeId)

	pod, err := c.cs.CoreV1().Pods(c.ns).Get(context.TODO(), vnName, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("Failed to get status for Node %s. Error: %s", nodeId, err.Error())
		return nil, err
	}

	state := getNodeRuntimeStatus(pod.Status.Phase)

	return &state, nil
}

/* Go routine to start build process */
func (c *Controller) PowerOnNode(nodeId string, org string) error {

	containerImage := internal.ServiceConfig.NodeImage + ":" + nodeId

	entryCommand := internal.ServiceConfig.NodeCmd

	nodeType := ukama.GetNodeType(nodeId)
	if nodeType == nil {
		logrus.Errorf("NodeId %s is not a expected uakma node Id.", nodeId)
		return fmt.Errorf("%s not expected nodeid format", nodeId)
	}

	logrus.Debugf("Starting node %s with Image %s and start up %s", nodeId, containerImage, entryCommand)
	err := c.CreateNode(nodeId, containerImage, entryCommand, *nodeType, org)
	if err != nil {
		logrus.Errorf("Create Node instance failed for %s. Error: %s", nodeId, err.Error())
		return err
	}
	return err
}

/* Delete job */
func (c *Controller) PowerOffNode(nodeId string) error {

	/* Virtual Node Name */
	vnName := getVirtNodeName(nodeId)

	logrus.Infof("Node %s powerOff requested.", nodeId)
	err := c.cs.CoreV1().Pods(c.ns).Delete(context.Background(), vnName, metav1.DeleteOptions{})
	if err != nil {
		logrus.Errorf("Delete Node failed for %s. Error: %s", nodeId, err.Error())
		return err
	}

	return nil
}

/* Watching for changes in virtual nodes */
func (c *Controller) WatcherForNodes(ctx context.Context, cb func(string, string) error) error {

	watcher, err := c.cs.CoreV1().Pods(c.ns).Watch(ctx, metav1.ListOptions{
		LabelSelector: "app=virtual-node",
	})
	if err != nil {
		return errors.Wrap(err, "cannot create Pod event watcher")
	}

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				if e.Object == nil {
					return
				}

				pod, ok := e.Object.(*v1.Pod)
				if !ok {
					continue
				}

				switch e.Type {

				case watch.Added:

					switch pod.Status.Phase {
					case v1.PodPending:
						state := db.VNodeOn.String()
						logrus.Infof("BootingUp: Node %s ", pod.Name)

						/* Update database */

						err := c.repo.Update(getVirtNodeId(pod.Name), state)
						if err != nil {
							logrus.Errorf("Error updating state of the node %s to %s.", pod.Name, state)
						}

						/* Send Event on MessageBus */
						err = cb(getVirtNodeId(pod.Name), state)
						if err != nil {
							logrus.Warningf("Failed to publish Virtual node event %s for %s", state, pod.Name)
						}

					default:
						logrus.Infof("Unkown Node state for %s during PowerOn.", pod.Name)

					}

				case watch.Deleted:
					for _, cst := range pod.Status.ContainerStatuses {
						if cst.State.Terminated != nil {
							state := db.VNodeOff.String()
							logrus.Infof("Poweroff : Node %s PoweredOff at %v Details: ExitCode: %d Reason: %s Message %s ", pod.Name, cst.State.Terminated.FinishedAt, cst.State.Terminated.ExitCode, cst.State.Terminated.Reason, cst.State.Terminated.Message)

							/* Update database */

							err := c.repo.Update(getVirtNodeId(pod.Name), state)
							if err != nil {
								logrus.Errorf("Error updating state of the node %s to %s.", pod.Name, state)
							}

							/* Send Event on MessageBus */
							err = cb(getVirtNodeId(pod.Name), state)
							if err != nil {
								logrus.Warningf("Failed to publish Virtual node event %s for %s", state, pod.Name)
							}
						}
					}

				case watch.Modified:

					for _, cst := range pod.Status.ContainerStatuses {
						if cst.State.Running != nil {
							logrus.Tracef("NodeState: Node %s running since %v", pod.Name, cst.State.Running.StartedAt)
						}
					}

				default:
					continue
				}

			case <-ctx.Done():
				watcher.Stop()
				return
			}
		}
	}()

	return nil
}

func (c *Controller) PublishEvent(uuid string, state string) error {

	evtMsg := &spec.EvtUpdateVirtnode{
		Uuid:   uuid,
		Status: state,
	}

	// Routing key
	key := msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(internal.ServiceName).SetEventType().SetObject("virtnode").SetAction("update").MustBuild()
	routingKey := msgbus.RoutingKey(key)

	// Marshal
	data, err := proto.Marshal(evtMsg)
	if err != nil {
		logrus.Errorf("Router:: fail marshal: %s", err.Error())
		return err
	}
	logrus.Debugf("Router:: Proto data for message is %+v \n MsgClient %+v", data, c.m)

	// Publish a message
	err = c.m.Publish(data, msgbus.DeviceQ.Queue, msgbus.DeviceQ.Exchange, routingKey, msgbus.DeviceQ.ExchangeType)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	return nil
}
