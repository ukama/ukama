package controller

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/ukama/ukama/testing/network/internal"
	"github.com/ukama/ukama/testing/network/internal/db"
	spec "github.com/ukama/ukama/testing/network/specs/controller/spec"
	"google.golang.org/protobuf/proto"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ControllerOps interface {
	ControllerInit() error
	ListNodes()
	GetNodeRuntimeStatus(nodeId string) (*string, error)
	PowerOnNode(nodeId string) error
	PowerOffNode(nodeId string) error
	CreateNode(name string, image string, command []string, ntype string) error
	WatcherForNodes(ctx context.Context, cb func(string, string) error) error
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
	pods, err := cset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
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

// /* Return data volume name forn Id strng */
// func getDataVolumeName(id string) string {
// 	return "dv-" + id
// }

/* Return node name from ID */
func getVirtNodeName(id string) string {
	return "vn-" + id
}

/* Return VM name from ID */
func getVirtNodeId(name string) string {
	return strings.Trim(name, "vn-")
}

/* Starting build job watcher routine */
func (c *Controller) ControllerInit() error {
	return c.WatcherForNodes(context.TODO(), c.repo, c.PublishEvent)
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

func (c *Controller) CreateNode(nodeId string, image string, command []string, ntype string) error {

	/* Image URL */
	//imageName := fmt.Sprintf("%s:%s", internal.ServiceConfig.NodeImage, name)

	/* Virtual Node Name */
	vnName := getVirtNodeName(nodeId)

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
					Name:    nodeId,
					Image:   image,
					Command: command,
					Args:    []string{},
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
		},
	}

	_, err := c.cs.CoreV1().
		Pods(c.ns).
		Create(context.TODO(), podSpec, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("PowerOn failure for node %s. Error: %s", nodeId, err.Error())
		return err
	}

	logrus.Debugf("PowerOn success for node %s", nodeId)

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
func (c *Controller) GetNodeRuntimeStatus(nodeId string) (*string, error) {

	/* Virtual Node Name */
	var state string
	vnName := getVirtNodeName(nodeId)

	pod, err := c.cs.CoreV1().Pods(c.ns).Get(context.TODO(), vnName, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("Failed to get status for Node %s. Error: %s", nodeId, err.Error())
		return nil, err
	}

	switch pod.Status.Phase {
	case v1.PodPending:
		state = "BootingUp"
	case v1.PodRunning:
		state = "Running"
	case v1.PodSucceeded:
		state = "Halted"
	case v1.PodFailed:
		state = "Failed"
	case v1.PodUnknown:
		state = "Failed"
	}
	return &state, nil
}

/* Go routine to start build process */
func (c *Controller) PowerOnNode(nodeId string) error {

	//containerImage := internal.ServiceConfig.BuilderImage

	entryCommand := []string{"sh", "-c", "echo \"Hello, Kubernetes!\" && sleep 3600"}
	nodeType := "hnode"
	image := "busybox:1.28"
	err := c.CreateNode(nodeId, image, entryCommand, nodeType)
	if err != nil {
		logrus.Errorf("Create Node innstance failed for %s. Error: %s", nodeId, err.Error())
		return err
	}
	return err
}

/* Delete job */
func (c *Controller) PowerOffNode(nodeId string) error {

	/* Virtual Node Name */
	vnName := getVirtNodeName(nodeId)

	logrus.Debugf("Node %s powerOff requested.", nodeId)
	err := c.cs.CoreV1().Pods(c.ns).Delete(context.Background(), vnName, metav1.DeleteOptions{})
	if err != nil {
		logrus.Errorf("Delete Node failed for %s. Error: %s", nodeId, err.Error())
		return err
	}

	return nil
}

/* Watching for changes in job status */
func (c *Controller) WatcherForNodes(ctx context.Context, d db.VNodeRepo, cb func(string, string) error) error {

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

				// logrus.WithFields(logrus.Fields{
				// 	"action":     e.Type,
				// 	"namespace":  pod.Namespace,
				// 	"name":       pod.Name,
				// 	"phase":      pod.Status.Phase,
				// 	"reason":     pod.Status.Reason,
				// 	"message":    pod.Status.Message,
				// 	"container#": len(pod.Status.ContainerStatuses),
				// }).Debug("Event notified")

				switch e.Type {

				case watch.Added:

					switch pod.Status.Phase {
					case v1.PodPending:
						state := db.VNodeOn.String()
						logrus.Infof("BootingUp: Node %s ", pod.Name)

						/* Updated= database */
						err := d.Update(getVirtNodeId(pod.Name), state)
						if err != nil {
							logrus.Errorf("Error updating state of the node %s to %s.", pod.Name, state)
						}

						/* Send Event on MessageBus */
						_ = cb(getVirtNodeId(pod.Name), state)

					default:
						logrus.Infof("Unkown Node state for %s during PowerOn.", pod.Name)

					}

				case watch.Deleted:
					for _, cst := range pod.Status.ContainerStatuses {
						if cst.State.Terminated != nil {
							state := db.VNodeOff.String()
							logrus.Infof("Poweroff : Node %s PoweredOff at %v Details: ExitCode: %d Reason: %s Message %s ", pod.Name, cst.State.Terminated.FinishedAt, cst.State.Terminated.ExitCode, cst.State.Terminated.Reason, cst.State.Terminated.Message)

							/* Updated= database */
							err := d.Update(getVirtNodeId(pod.Name), state)
							if err != nil {
								logrus.Errorf("Error updating state of the node %s to %s.", pod.Name, state)
							}

							/* Send Event on MessageBus */
							_ = cb(getVirtNodeId(pod.Name), state)
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
