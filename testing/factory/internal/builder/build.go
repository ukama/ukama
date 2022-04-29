package builder

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/ukama/ukama/testing/factory/internal"

	"github.com/ukama/ukama/testing/factory/internal/db"

	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

type BuildOps interface {
	Init(db db.NodeRepo, cb func(string, string) error)
	LaunchBuildJob(jobName *string, image *string, cmd *string, nodetype *string) error
	GetJobStatus(jobName string) int
	DeleteJob(jobName string) error
	ListBuildJobs()
	ListPods()
	WatcherForBuildJobs(db db.NodeRepo, cb func(string, string) error)
	LaunchAndMonitorBuild(jobName string, nodetype string) error
}

type Build struct {
	clientset        *kubernetes.Clientset
	currentNamespace string
}

func NewBuild() *Build {

	cset, err := connectToK8s()
	if err != nil {
		log.Fatalf("Build:: Can't connect to Kuberneets cluster. Err: %s", err.Error())
		return nil
	}

	ns, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		log.Fatalf("Build:: Can't read current namespace. Err: %s", err.Error())
		return nil
	}

	return &Build{
		clientset:        cset,
		currentNamespace: string(ns),
	}
}

/* Connect to Kubernetes cluster */
func connectToK8s() (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", internal.ServiceConf.Kubeconfig)
	if err != nil {
		log.Errorf("Build:: Failed to create K8s config. Error %s", err.Error())
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("Build:: Failed to create K8s clientset. Error %s", err.Error())
		return nil, err
	}

	return clientset, nil
}

/* Starting build job watcher routine */
func (b *Build) Init(db db.NodeRepo, cb func(string, string) error) {
	go b.WatcherForBuildJobs(db, cb)
}

/* Launch Build Job in K8 cluister */
func (b *Build) LaunchBuildJob(jobName *string, image *string, cmd *string, nodetype *string) error {

	jobs := b.clientset.BatchV1().Jobs(b.currentNamespace)

	/* Tries 4 time before matking it as fail.*/
	var backOffLimit int32 = 4

	/* Priviliged mode : mercy!! (because of dind for linuxkit build)
	Would be removed with our microCE
	*/
	var priviligemode bool = true
	var timetolive int32 = 60

	/* Add a time period for job to complete if not completetd within that time frame remove it.*/
	var activeDeadlineSeconds int64 = 90 * 60

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
							Command: strings.Split(*cmd, " "),
							SecurityContext: &v1.SecurityContext{
								Privileged: &priviligemode,
							},
							Env: []v1.EnvVar{
								{
									Name:  "UUID",
									Value: *jobName,
								},
								{
									Name:  "NODETYPE",
									Value: *nodetype,
								},
								{
									Name:  "GITUSR",
									Value: internal.ServiceConf.GitUser,
								},
								{
									Name:  "GITKEY",
									Value: internal.ServiceConf.GitPass,
								},
								{
									Name:  "DOCKER_USER",
									Value: internal.ServiceConf.Docker.User,
								},
								{
									Name:  "DOCKER_PASS",
									Value: internal.ServiceConf.Docker.Pass,
								},
								{
									Name:  "REPO_SERVER_URL",
									Value: internal.ServiceConf.RepoServerUrl,
								},
							},
						},
					},
					ImagePullSecrets: []v1.LocalObjectReference{
						{
							Name: "dregcred",
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
		log.Errorf("Build:: Failed to create Build job %s. Error %s", *jobName, err.Error())
		return err
	}

	log.Debugf("Build:: Created Build job %s successfully", *jobName)

	return nil

}

/* Debug List of Jobs */
func (b *Build) ListBuildJobs() {

	jobset := b.clientset.BatchV1().Jobs(b.currentNamespace)

	jobs, _ := jobset.List(context.TODO(), metav1.ListOptions{})
	for _, job := range jobs.Items {
		log.Debugf("Build:: Job Name %s Job Status %v.", job.Name, job.Status)
	}

}

/* Debug List of pods */
func (b *Build) ListPods() {

	pods, _ := b.clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	for _, pod := range pods.Items {
		log.Debugf("Build:: Pod Name %s Pod Status %v.", pod.Name, pod.Status)
	}
}

/* Get Job status */
func (b *Build) GetJobStatus(jobName string) int {
	done := 0
	jobset := b.clientset.BatchV1().Jobs(b.currentNamespace)

	job, _ := jobset.Get(context.TODO(), jobName, metav1.GetOptions{})

	if job.Status.Active > 0 {
		log.Debugf("Build:: Job %s is still running", job.Name)

	} else {
		if job.Status.Succeeded > 0 {
			log.Infof("Build:: Job %s is completed successfully", job.Name)
			done = 1
		} else {
			log.Errorf("Build:: Job %s is failed", job.Name)
			done = -1
		}
	}

	return done
}

/* Go routine to start build process */
func (b *Build) LaunchAndMonitorBuild(jobName string, nodetype string) error {

	containerImage := internal.ServiceConf.BuilderImage

	entryCommand := "startup.sh"

	log.Infof("Starting build process for %s with id  %s", nodetype, jobName)

	err := b.LaunchBuildJob(&jobName, &containerImage, &entryCommand, &nodetype)
	if err != nil {
		log.Errorf("Build:: BuildJob fauiled for %s. Error: %s", jobName, err.Error())
		return err
	}
	return err
}

/* Delete job */
func (b *Build) DeleteJob(jobName string) error {

	jobset := b.clientset.BatchV1().Jobs(b.currentNamespace)

	err := jobset.Delete(context.TODO(), jobName, metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("Build:: Failed to delete job %s . Error: %s.", jobName, err.Error())
		return err
	}

	log.Debugf("Build:: Job %s delete requested.", jobName)
	return err
}

/* Watching for changes in job status */
func (b *Build) WatcherForBuildJobs(db db.NodeRepo, cb func(string, string) error) {

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
			log.Fatal(err)
		}
		ch := watcher.ResultChan()

		log.Debugf("Build:: Starting Build watcher routine.")

		for event := range ch {
			job, ok := event.Object.(*batchv1.Job)
			if !ok {
				log.Errorf("Build:: unexpected type")
				return
			}

			switch event.Type {
			case watch.Added:
				state := "assembly-begin"
				log.Infof("Build:: New Build job addded kind %s name %s created in namespace %s Active %d Completed %d.\n", job.Kind, job.Name, job.Namespace, job.Status.Active, job.Status.Succeeded)
				/* Updated database */
				err := db.UpdateNodeStatus(job.Name, state)
				if err != nil {
					log.Errorf("Build:: Error updating state of the node %s to %s.", job.Name, state)
				}

				/* Send Event on MessageBus */
				_ = cb(job.Name, state)

			case watch.Modified:
				state := "assembly-inprogress"
				log.Infof("Build:: Build job modified kind %s name %s created in namespace %s Active %d Completed %d Failed %d .\n",
					job.Kind, job.Name, job.Namespace, job.Status.Active, job.Status.Succeeded, job.Status.Failed)

				/* Since it single pod */
				if job.Status.Active > 0 {
					state = "assembly-inprogress"
				} else if job.Status.Succeeded > 0 {
					state = "assembly-completed"
					_ = b.DeleteJob(job.Name)
				} else if job.Status.Failed > 0 {
					state = "assembly-failure"
				}

				/* Updated database */
				err := db.UpdateNodeStatus(job.Name, state)
				if err != nil {
					log.Errorf("Build:: Error updating state of the node %s to %s.", job.Name, state)
				}

				/* Send Event on MessageBus */
				_ = cb(job.Name, state)

			case watch.Deleted:
				log.Infof("Build:: Build job deleted kind %s name %s created in namespace %s Active %d Completed %d Failed %d.",
					job.Kind, job.Name, job.Namespace, job.Status.Active, job.Status.Succeeded, job.Status.Failed)
				/* TODO:: Send Event on MessageBus */
				state := "shipped"
				/* Updated database */
				err := db.UpdateNodeStatus(job.Name, state)
				if err != nil {
					log.Errorf("Build:: Error updating state of the node %s to %s.", job.Name, state)
				}

				/* Send Event on MessageBus */
				_ = cb(job.Name, state)

			case watch.Error:
				log.Errorf("Build:: Build job error kind %s name %s created in namespace %s Active %d Completed %d  Failed %d.",
					job.Kind, job.Name, job.Namespace, job.Status.Active, job.Status.Succeeded, job.Status.Failed)
			}
		}

		log.Debugf("Build:: Exiting Build job watcher routine.")
	}

}
