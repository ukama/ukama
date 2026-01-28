package utils

import (
	"context"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	PodNamePrefix       = "mesh-node"
	PodIPWaitTimeout    = 60 * time.Second
	PodIPPollInterval   = 2 * time.Second
)

type NodeMeshInfo struct {
	NodeId      string
	MeshPodIp   string
	MeshPodPort int32
}

func getPodNamePrefix(orgName string) string {
	return orgName + "-" + PodNamePrefix
}

func getPodName(orgName, nodeId string) string {
	return getPodNamePrefix(orgName) + "-" + nodeId
}

func SpawnReplica(ctx context.Context, node NodeMeshInfo, config *pkg.Config, clientSet *kubernetes.Clientset) error {
	namespace := config.OrgName + "-" + config.MeshNamespace
	podNamePrefix := getPodName(config.OrgName, node.NodeId)

	// Check for existing healthy pod first
	existingPod, err := findExistingMeshPod(ctx, namespace, podNamePrefix, clientSet)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to check existing pods: %v", err)
	}

	// If a healthy pod exists, validate IP and return
	if existingPod != nil {
		// If NNS return empty mesh IP, update it asynchronously
		if node.MeshPodIp == "" {
			log.Debugf("Pod %s exists but has no IP yet, updating asynchronously", existingPod.Name)
			go updatePodIPAsync(context.Background(), namespace, existingPod.Name, node.NodeId, clientSet)
		}

		// Check if NNS returned mesh IP matches the pod IP
		if  node.MeshPodIp == existingPod.Status.PodIP {
			log.Debugf("Mesh pod already exists and IP matched for node %s: %s (IP: %s)", node.NodeId, existingPod.Name, existingPod.Status.PodIP)
			return nil
		}
		
		return nil
	}

	// No healthy pod exists, create a new one
	return createMeshPod(ctx, namespace, podNamePrefix, node, clientSet)
}

// findExistingMeshPod looks for an existing healthy mesh pod for the node.
// Returns nil if no healthy pod is found.
func findExistingMeshPod(ctx context.Context, namespace, podNamePrefix string, clientSet *kubernetes.Clientset) (*corev1.Pod, error) {
	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/component=mesh-node",
	})
	if err != nil {
		return nil, err
	}

	for i := range pods.Items {
		pod := &pods.Items[i]

		// Skip pods that don't match the pod name prefix
		if !isPodForNode(pod, podNamePrefix) {
			continue
		}

		// Check if pod is healthy (Running or Pending)
		if isPodHealthy(pod) {
			log.Debugf("Found healthy mesh pod for node: %s (phase: %s)", pod.Name, pod.Status.Phase)
			return pod, nil
		}

		// Pod exists but is unhealthy - log and continue to create new one
		log.Warnf("Found unhealthy mesh pod %s (phase: %s), will create new one", pod.Name, pod.Status.Phase)
	}

	return nil, nil
}

// isPodForNode checks if the pod belongs to our node.
func isPodForNode(pod *corev1.Pod, podNamePrefix string) bool {
	return strings.HasPrefix(pod.Name, podNamePrefix)
}

// isPodHealthy checks if a pod is in a healthy state.
func isPodHealthy(pod *corev1.Pod) bool {
	switch pod.Status.Phase {
	case corev1.PodRunning, corev1.PodPending:
		// Also check if pod is being deleted
		return pod.DeletionTimestamp == nil
	default:
		return false
	}
}

// updatePodIPOnly updates only the IP address for a node in the database.
// This is used for asynchronous IP updates.
func updatePodIPOnly(nodeId, podIP string) error {
	log.Infof("Generate mesh pod IP update event for node %s: %s", nodeId, podIP)
	return nil
}

// updatePodIPAsync waits for a pod to get an IP address and updates it in the database asynchronously.
// This function runs in a goroutine and handles IP updates without blocking the main flow.
func updatePodIPAsync(ctx context.Context, namespace, podName, nodeId string, clientSet *kubernetes.Clientset) {
	// Use a separate context with timeout for async operation
	asyncCtx, cancel := context.WithTimeout(context.Background(), PodIPWaitTimeout)
	defer cancel()

	log.Debugf("Starting async IP update for pod %s (node %s)", podName, nodeId)

	podIP, err := waitForPodIP(asyncCtx, namespace, podName, clientSet)
	if err != nil {
		log.Warnf("Failed to get pod IP for %s (node %s) asynchronously: %v", podName, nodeId, err)
		// Don't update with empty IP - it will be retried on next sync or pod check
		return
	}

	if podIP == "" {
		log.Warnf("Pod %s (node %s) has empty IP, skipping update", podName, nodeId)
		return
	}

	// Update IP in database asynchronously
	if err := updatePodIPOnly(nodeId, podIP); err != nil {
		log.Errorf("Failed to update pod IP asynchronously for node %s: %v", nodeId, err)
	} else {
		log.Infof("Successfully updated pod IP asynchronously for node %s: %s", nodeId, podIP)
	}
}

// createMeshPod creates a new mesh pod for the node.
func createMeshPod(ctx context.Context, namespace, podName string, node NodeMeshInfo, clientSet *kubernetes.Clientset) error {
	// Get template deployment
	podSpec, err := getTemplatePodSpec(ctx, namespace, clientSet)
	if err != nil {
		return err
	}

	// Create the pod (add trailing - for cleaner random suffix separation)
	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: podName + "-",
			Namespace:    namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":      "mesh",
				"app.kubernetes.io/component": "mesh-node",
				"app.kubernetes.io/node-id":  node.NodeId,
			},
		},
		Spec: *podSpec,
	}

	createdPod, err := clientSet.CoreV1().Pods(namespace).Create(ctx, newPod, metav1.CreateOptions{})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create mesh pod: %v", err)
	}

	log.Infof("Created mesh pod %s for node %s", createdPod.Name, node.NodeId)

	// Update IP asynchronously in background
	go updatePodIPAsync(context.Background(), namespace, createdPod.Name, node.NodeId, clientSet)

	return nil
}

// waitForPodIP waits for a pod to be assigned an IP address.
func waitForPodIP(ctx context.Context, namespace, podName string, clientSet *kubernetes.Clientset) (string, error) {
	timeout := time.After(PodIPWaitTimeout)
	ticker := time.NewTicker(PodIPPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-timeout:
			return "", status.Errorf(codes.DeadlineExceeded, "timeout waiting for pod %s to get IP", podName)
		case <-ticker.C:
			pod, err := clientSet.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
			if err != nil {
				log.Warnf("Error getting pod %s: %v", podName, err)
				continue
			}

			if pod.Status.PodIP != "" {
				return pod.Status.PodIP, nil
			}

			// Check if pod failed
			if pod.Status.Phase == corev1.PodFailed {
				return "", status.Errorf(codes.Internal, "pod %s failed: %s", podName, pod.Status.Reason)
			}

			log.Debugf("Pod %s phase: %s, waiting for IP...", podName, pod.Status.Phase)
		}
	}
}

// getTemplatePodSpec retrieves the pod spec from the mesh deployment template.
func getTemplatePodSpec(ctx context.Context, namespace string, clientSet *kubernetes.Clientset) (*corev1.PodSpec, error) {
	deployments, err := clientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=mesh",
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list deployments: %v", err)
	}

	if len(deployments.Items) == 0 {
		return nil, status.Errorf(codes.NotFound, "no mesh deployment found in namespace %s", namespace)
	}

	podSpec := deployments.Items[0].Spec.Template.Spec.DeepCopy()
	podSpec.RestartPolicy = corev1.RestartPolicyOnFailure

	return podSpec, nil
}