package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/ukama/ukama/testing/network/internal"
	"github.com/ukama/ukama/testing/network/internal/db"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ControllerOps interface {
	ControllerInit()
	ListNodes()
	GetNodeStatus(nodeId string) error
	WatcherForNodes() error
	PowerOnNode(nodeId string) error
	PowerOffNode(nodeId string) error
	CreateNode(name, command string, args []string, ntype string) error
}

type Controller struct {
	repo db.VNodeRepo
	cs   *kubernetes.Clientset
	ns   string
	m    msgbus.Publisher
}

func NewController(d db.VNodeRepo) *Controller {
	cset, err := connectToK8s()
	if err != nil {
		logrus.Fatalf("Build:: Can't connect to Kuberneets cluster. Err: %s", err.Error())
		return nil
	}

	/* For test */
	pods, err := cset.CoreV1().Pods("kube-system").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("error getting pods: %v\n", err)
		return nil
	}
	for _, pod := range pods.Items {
		logrus.Tracef("Pod name: %s\n", pod.Name)
	}

	// ns, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	// if err != nil {
	// 	logrus.Fatalf("Build:: Can't read current namespace. Err: %s", err.Error())
	// 	return nil
	// }
	msgC, err := msgbus.NewPublisherClient(internal.ServiceConfig.RabbitUri)
	if err != nil {
		logrus.Errorf("error getting message publisher: %s\n", err.Error())
		return nil
	}

	return &Controller{
		cs:   cset,
		ns:   internal.ServiceConfig.Namespace,
		m:    msgC,
		repo: d,
	}
}

/* Connect to Kubernetes cluster */
func connectToK8s() (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", internal.ServiceConfig.Kubeconfig)
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

/* Return data volume name forn Id strng */
func getDataVolumeName(id string) string {
	return "dv-" + id
}

/* Return node name from ID */
func getVirtNodeName(id string) string {
	return "vn-" + id
}

/* Return VM name from ID */
func getVirtNodeId(name string) string {
	return strings.Trim(name, "vn-")
}

/* Starting build job watcher routine */
func (c *Controller) ControllerInit() {
	go c.WatcherForNodes()
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

func (c *Controller) CreateNode(name, command string, ntype string) error {

	/* Image URL */
	imageName := fmt.Sprintf("%s:%s", internal.ServiceConfig.NodeImage, name)

	/* Virtual Node Name */
	vnName := getVirtNodeName(name)

	/* Virt Nodes Instance Labels */
	labels := map[string]string{
		"node": vnName,
		"type": ntype,
		"app":  "virtual-node",
		"org":  "ukama",
	}

	// /* Data Volume */
	// dataVolumeName := getDataVolumeName(name)

	// /* Data Volume Labels */
	// dvlabels := map[string]string{
	// 	"node": name,
	// 	"app":  "virtnode-datavolume",
	// 	"type": ntype,
	// 	"org":  "ukama",
	// }

	/* Pod spec */
	podSpec := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vnName,
			Namespace: c.ns,
			Labels:    labels,
		},
		Spec: v1.PodSpec{
			RestartPolicy: v1.RestartPolicyOnFailure,
			Containers: []v1.Container{
				v1.Container{
					Name:    name,
					Image:   imageName,
					Command: []string{command},
					Args:    []string{},
					Env: []v1.EnvVar{
						{
							Name:  "UUID",
							Value: name,
						},
						{
							Name:  "NODETYPE",
							Value: ntype,
						},
					},
				},
			},
		},
	}

	_, err := c.cs.CoreV1().
		Pods(c.ns).
		Create(context.TODO(), podSpec, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("PowerOn failure for node %s. Error: %s", name, err.Error())
		return err
	}

	logrus.Debugf("PowerOn success for node %s", name)

	return err

}

/* Debug List of pods */
func (c *Controller) ListNodes() {

	pods, _ := c.cs.CoreV1().Pods(c.ns).List(context.TODO(), metav1.ListOptions{})

	for _, pod := range pods.Items {
		logrus.Debugf("Node Name %s Node Status %v.", pod.Name, pod.Status)
	}
}

/* Get Job status */
func (c *Controller) GetNodeStatus(nodeId string) int {
	done := 0

	return done
}

/* Go routine to start build process */
func (c *Controller) PowerOnNode(nodeId string, nodetype string) error {

	//containerImage := internal.ServiceConfig.BuilderImage

	entryCommand := "startup.sh"

	err := c.CreateNode(nodeId, entryCommand, nodetype)
	if err != nil {
		logrus.Errorf("Create Node innstance failed for %s. Error: %s", nodeId, err.Error())
		return err
	}
	return err
}

/* Delete job */
func (c *Controller) PowerOffNode(nodeId string) error {

	logrus.Debugf("Node%s powerOff requested.", nodeId)
	return nil
}

/* Watching for changes in job status */
func (c *Controller) WatcherForNodes() {

	// nodeSelctor, _ := labels.NewRequirement("app", selection.In, []string{"virtual-node"})
	// selector := labels.NewSelector()
	// selector = selector.Add(*nodeSelctor)
	timeout := int64(60 * 60 * 240) // 24*10 hours= 10 days Test purpose

	for {

		watcher, err := c.cs.CoreV1().Events(c.ns).Watch(context.TODO(), metav1.ListOptions{
			Watch:          true,
			LabelSelector:  "app: virtual-node",
			TimeoutSeconds: &timeout,
		})
		if err != nil {
			return
		}

		ch := watcher.ResultChan()

		for event := range ch {
			node, ok := event.Object.(*v1.Pod)
			if !ok {
				logrus.Errorf("Controller:: unexpected event type on watcher")
				return
			}

			switch event.Type {
			case watch.Added:
				logrus.Debugf("NodeId %s PoweredOn", node.Name)
			case watch.Modified:
				logrus.Debugf("NodeId %s Updated", node.Name)
			case watch.Deleted:
				logrus.Debugf("NodeId %s PoweredOff", node.Name)
			case watch.Error:
				logrus.Debugf("NodeId %s Failure", node.Name)
			}
		}
	}
	logrus.Debugf("Controller:: Exiting VirtualNode watcher routine.")
}

/*
func (c *Controller) PublishEvent(uuid string, state string) error {

	evtMsg := &spec.EvtUpdateVirtnode{
		Uuid:   uuid,
		Status: state,
	}

	// Marshal
	data, err := proto.Marshal(evtMsg)
	if err != nil {
		logrus.Errorf("Router:: fail marshal: %s", err.Error())
		return err
	}
	logrus.Debugf("Router:: Proto data for message is %+v and MsgClient %+v", data, c.m)

	// Publish a message
	err = c.m.Publish(data, msgbus.DeviceQ.Queue, msgbus.DeviceQ.Exchange, msgbus.EventVirtNodeUpdateStatus, msgbus.DeviceQ.ExchangeType)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	return nil
}
*/
