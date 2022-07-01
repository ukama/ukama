package builder

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ukama/ukama/testing/services/factory/internal"
	"google.golang.org/protobuf/proto"

	"github.com/ukama/ukama/testing/services/factory/internal/nmr"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/msgbus"
	spec "github.com/ukama/ukama/testing/services/factory/specs/factory/spec"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type BuildOps interface {
	BuildInit()

	GetJobStatus(jobName string) int
	DeleteJob(jobName string) error
	ListBuildJobs()
	ListPods()
	WatcherForBuildJobs() error
	LaunchAndMonitorBuild(jobName string, node internal.Node) error
	LaunchBuildJob(jobName *string, image *string, cmd []string, nodetype *string, jNodeMetaData []byte) error
}

type Build struct {
	clientset        *kubernetes.Clientset
	currentNamespace string
	fd               *nmr.NMR
	m                msgbus.Publisher
}

type NodeMetaData struct {
	NodeInfo   internal.Node     `json:"nodeInfo"`
	NodeConfig []internal.Module `json:"nodeConfig"`
}

func NewMsgBus() *msgbus.MsgClient {
	msgClient := &msgbus.MsgClient{}
	msgClient.ConnectToBroker(internal.ServiceConfig.RabbitUri)
	return msgClient
}

func NewBuild(d *nmr.NMR) *Build {

	cset, err := connectToK8s()
	if err != nil {
		logrus.Fatalf("Build:: Can't connect to Kuberneets cluster. Err: %s", err.Error())
		return nil
	}

	/* For listing already running virtual nodes's  */
	pods, err := cset.CoreV1().Pods(internal.ServiceConfig.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("Build:: error getting pods: %v\n", err)
		return nil
	}
	for _, pod := range pods.Items {
		logrus.Tracef("Build:: Virtual Node : %s\n", pod.Name)
	}

	msgC, err := msgbus.NewPublisherClient(internal.ServiceConfig.RabbitUri)
	if err != nil {
		logrus.Errorf("Build:: error getting message publisher: %s\n", err.Error())
		return nil
	}

	ns := "default"

	if internal.ServiceConfig.Namespace != "" {
		ns = internal.ServiceConfig.Namespace
	}

	return &Build{
		clientset:        cset,
		currentNamespace: ns,
		fd:               d,
		m:                msgC,
	}
}

/* Connect to Kubernetes cluster */
func connectToK8s() (*kubernetes.Clientset, error) {

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

/* Starting build job watcher routine */
func (b *Build) BuildInit() {
	go b.WatcherForBuildJobs()
}

/* Launch Build Job in K8 cluster */
func (b *Build) LaunchBuildJob(jobName *string, image *string, cmd []string, nodetype *string, jNodeMetaData []byte) error {

	jobs := b.clientset.BatchV1().Jobs(b.currentNamespace)

	/* Tries 4 time before matking it as fail.*/
	var backOffLimit int32 = internal.ServiceConfig.BackOffLimit

	/* Priviliged mode : mercy!! (because of dind for linuxkit build)
	Would be removed with our microCE
	*/
	var priviligemode bool = true
	var timetolive int32 = internal.ServiceConfig.TimeToLive

	/* Add a time period for job to complete if not completetd within that time frame remove it.*/
	var activeDeadlineSeconds int64 = internal.ServiceConfig.ActiveDeadLineSeconds

	/* Job spec */
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *jobName,
			Namespace: b.currentNamespace,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &timetolive,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    *jobName,
							Image:   *image,
							Command: cmd,
							SecurityContext: &v1.SecurityContext{
								Privileged: &priviligemode,
							},
							Env: []v1.EnvVar{
								{
									Name:  "VNODE_ID",
									Value: *jobName,
								},
								{
									Name:  "NODETYPE",
									Value: *nodetype,
								},
								{
									Name:  "VNODE_METADATA",
									Value: string(jNodeMetaData),
								},
								{
									Name:  "GITUSR",
									Value: internal.ServiceConfig.GitUser,
								},
								{
									Name:  "GITKEY",
									Value: internal.ServiceConfig.GitPass,
								},
								{
									Name:  "DOCKER_USER",
									Value: internal.ServiceConfig.Docker.User,
								},
								{
									Name:  "DOCKER_PASS",
									Value: internal.ServiceConfig.Docker.Pass,
								},
								{
									Name:  "REPO_SERVER_URL",
									Value: internal.ServiceConfig.VNodeRepoServerUrl,
								},
								{
									Name:  "REPO_NAME",
									Value: internal.ServiceConfig.VNodeRepoName,
								},
								{
									Name:  "AWS_ACCESS_KEY_ID",
									Value: internal.ServiceConfig.AwsKey,
								},
								{
									Name:  "AWS_SECRET_ACCESS_KEY",
									Value: internal.ServiceConfig.AwsSecret,
								},
							},

							EnvFrom: []v1.EnvFromSource{
								{
									ConfigMapRef: &v1.ConfigMapEnvSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: internal.ServiceConfig.CmRef,
										},
									},
								},
								{
									SecretRef: &v1.SecretEnvSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: internal.ServiceConfig.SecRef,
										},
									},
								},
							},
						},
					},
					ImagePullSecrets: []v1.LocalObjectReference{
						{
							Name: internal.ServiceConfig.BuilderRegCred,
						},
					},
					RestartPolicy: v1.RestartPolicyOnFailure,
				},
			},
			BackoffLimit:          &backOffLimit,
			ActiveDeadlineSeconds: &activeDeadlineSeconds,
		},
	}

	/* Create Job */
	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("Build:: Failed to create Build job %s. Error %s", *jobName, err.Error())
		return err
	}

	logrus.Debugf("Build:: Created Build job %s successfully", *jobName)

	return nil

}

/* Debug List of Jobs */
func (b *Build) ListBuildJobs() {

	jobset := b.clientset.BatchV1().Jobs(b.currentNamespace)

	jobs, _ := jobset.List(context.TODO(), metav1.ListOptions{})
	for _, job := range jobs.Items {
		logrus.Debugf("Build:: Job Name %s Job Status %v.", job.Name, job.Status)
	}

}

/* Debug List of pods */
func (b *Build) ListPods() {

	pods, err := b.clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Debugf("Build:: Failed to get pod status. Error %s.", err.Error())
	}

	for _, pod := range pods.Items {
		logrus.Debugf("Build:: Pod Name %s Pod Status %v.", pod.Name, pod.Status)
	}
}

/* Get Job status */
func (b *Build) GetJobStatus(jobName string) int {
	done := 0
	jobset := b.clientset.BatchV1().Jobs(b.currentNamespace)

	job, _ := jobset.Get(context.TODO(), jobName, metav1.GetOptions{})

	if job.Status.Active > 0 {
		logrus.Debugf("Build:: Job %s is still running", job.Name)

	} else {
		if job.Status.Succeeded > 0 {
			logrus.Infof("Build:: Job %s is completed successfully", job.Name)
			done = 1
		} else {
			logrus.Errorf("Build:: Job %s is failed", job.Name)
			done = -1
		}
	}

	return done
}

/* Go routine to start build process */
func (b *Build) LaunchAndMonitorBuild(jobName string, node internal.Node) error {

	containerImage := internal.ServiceConfig.BuilderImage

	entryCommand := internal.ServiceConfig.BuilderCmd

	logrus.Debugf("Build:: Starting build process for node %s with details: %+v", jobName, node)

	/* Marshal the node info json */
	nodeMetaData := NodeMetaData{
		NodeInfo:   node,
		NodeConfig: node.Modules,
	}

	jNodeMetaData, err := json.Marshal(nodeMetaData)
	if err != nil {
		logrus.Errorf("Build:: Failed to add node with nodeID %s. Error %s", node.NodeID, err.Error())
		return fmt.Errorf("failed to pass nodeID %s info to worker. Error %s", node.NodeID, err.Error())
	}
	logrus.Debugf("Build:: Node meta data: %+v", string(jNodeMetaData))

	err = b.LaunchBuildJob(&jobName, &containerImage, entryCommand, &node.Type, jNodeMetaData)
	if err != nil {
		logrus.Errorf("Build:: BuildJob fauiled for %s. Error: %s", jobName, err.Error())
		return err
	}
	return err
}

/* Delete job */
func (b *Build) DeleteJob(jobName string) error {

	jobset := b.clientset.BatchV1().Jobs(b.currentNamespace)

	err := jobset.Delete(context.TODO(), jobName, metav1.DeleteOptions{})
	if err != nil {
		logrus.Errorf("Build:: Failed to delete job %s . Error: %s.", jobName, err.Error())
		return err
	}

	logrus.Debugf("Build:: Job %s delete requested.", jobName)
	return err
}

/* Watching for changes in job status */
func (b *Build) WatcherForBuildJobs() {

	/* List Jobs*/
	b.ListBuildJobs()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/* watcher time out happens after every 30 to 60 mints so need to start it again*/
	for {
		timeout := int64(60 * 60 * 240) // 24*10 hours= 10 days Test purpose

		/* Watch future changes to Build jobs */
		watcher, err := b.clientset.BatchV1().Jobs(b.currentNamespace).Watch(ctx, metav1.ListOptions{TimeoutSeconds: &timeout})
		if err != nil {
			logrus.Fatal(err)
		}
		ch := watcher.ResultChan()

		logrus.Debugf("Build:: Starting Build watcher routine.")

		for event := range ch {
			job, ok := event.Object.(*batchv1.Job)
			if !ok {
				logrus.Errorf("Build:: unexpected type")
				return
			}

			switch event.Type {
			case watch.Added:
				state := "StatusPendingAssembly"
				logrus.Infof("Build:: New Build job addded kind %s name %s created in namespace %s Active %d Completed %d.\n", job.Kind, job.Name, job.Namespace, job.Status.Active, job.Status.Succeeded)

				/* Updated database */
				err := b.fd.NmrUpdateNodeStatus(job.Name, state)
				if err != nil {
					logrus.Errorf("Build:: Error updating state of the node %s to %s.", job.Name, state)
				}

				/* Send Event on MessageBus */
				err = b.PublishEvent(job.Name, state)
				if err != nil {
					logrus.Errorf("Build:: Error publishing updated state of the node %s to %s.", job.Name, state)
				}

			case watch.Modified:
				state := "StatusUnderAssembly"
				logrus.Infof("Build:: Build job modified kind %s name %s created in namespace %s Active %d Completed %d Failed %d .\n",
					job.Kind, job.Name, job.Namespace, job.Status.Active, job.Status.Succeeded, job.Status.Failed)

				/* Since it single pod */
				if job.Status.Active > 0 {
					state = "StatusUnderAssembly"
				} else if job.Status.Succeeded > 0 {
					state = "StatusAssemblyCompleted"
					_ = b.DeleteJob(job.Name)
				} else if job.Status.Failed > 0 {
					state = "StatusProductionTestFail"
				}

				/* Updated database */
				err := b.fd.NmrUpdateNodeStatus(job.Name, state)
				if err != nil {
					logrus.Errorf("Build:: Error updating state of the node %s to %s.", job.Name, state)
				}

				/* Send Event on MessageBus */
				err = b.PublishEvent(job.Name, state)
				if err != nil {
					logrus.Errorf("Build:: Error publishing updated state of the node %s to %s.", job.Name, state)
				}

			case watch.Deleted:
				state := "StatusProductionTestFail"
				logrus.Infof("Build:: Build job deleted kind %s name %s created in namespace %s Active %d Completed %d Failed %d.",
					job.Kind, job.Name, job.Namespace, job.Status.Active, job.Status.Succeeded, job.Status.Failed)
				/* TODO:: Send Event on MessageBus */
				if job.Status.Succeeded > 0 {
					state = "StatusNodeIntransit"
					err = b.DeleteJob(job.Name)
					if err != nil {
						logrus.Errorf("Build:: Error deleteing job %s. Error %s.", job.Name, err.Error())
					}
				}
				/* Updated database */
				err := b.fd.NmrUpdateNodeStatus(job.Name, state)
				if err != nil {
					logrus.Errorf("Build:: Error updating state of the node %s to %s.", job.Name, state)
				}

				/* Send Event on MessageBus */
				err = b.PublishEvent(job.Name, state)
				if err != nil {
					logrus.Errorf("Build:: Error publishing updated state of the node %s to %s.", job.Name, state)
				}

			case watch.Error:
				logrus.Errorf("Build:: Build job error kind %s name %s created in namespace %s Active %d Completed %d  Failed %d.",
					job.Kind, job.Name, job.Namespace, job.Status.Active, job.Status.Succeeded, job.Status.Failed)
			}
		}

		logrus.Debugf("Build:: Exiting Build job watcher routine.")
	}

}

func (b *Build) PublishEvent(uuid string, state string) error {

	evtMsg := &spec.EvtUpdateVirtnode{
		Uuid:   uuid,
		Status: state,
	}

	// Marshal
	data, err := proto.Marshal(evtMsg)
	if err != nil {
		logrus.Errorf("Build:: fail marshal: %s", err.Error())
		return err
	}
	logrus.Debugf("Build:: Proto data for message is %+v and MsgClient %+v", data, b.m)

	// Publish a message
	err = b.m.Publish(data, msgbus.DeviceQ.Queue, msgbus.DeviceQ.Exchange, msgbus.EventVirtNodeUpdateStatus, msgbus.DeviceQ.ExchangeType)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	return nil
}
