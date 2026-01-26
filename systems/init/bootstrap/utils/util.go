package utils

import (
	"context"
	"errors"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	PodNamePrefix       = "mesh-node"
	PodIPWaitTimeout    = 60 * time.Second
	PodIPPollInterval   = 2 * time.Second
)

func GetPodNamePrefix(orgName string) string {
	return orgName + "-" + PodNamePrefix
}

func GetPodName(orgName, nodeId string) string {
	return orgName + "-" + PodNamePrefix + "-" + nodeId
}

func SpawnReplica(ctx context.Context, node *db.Node, config *pkg.Config, clientSet *kubernetes.Clientset, nodeRepo db.NodeRepo) error {
	// Input validation
	if node == nil {
		return status.Error(codes.InvalidArgument, "node cannot be nil")
	}
	if node.NodeId == "" {
		return status.Error(codes.InvalidArgument, "node ID cannot be empty")
	}

	namespace := config.OrgName + "-" + config.MeshNamespace
	// Use full pod name prefix including nodeId to match only pods for THIS node
	podNamePrefix := GetPodName(config.OrgName, node.NodeId)

	// Check for existing healthy pod first
	existingPod, err := findExistingMeshPod(ctx, namespace, node.MeshPodName, podNamePrefix, clientSet)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to check existing pods: %v", err)
	}

	// If a healthy pod exists, sync database if needed and return
	if existingPod != nil {
		return syncNodePodInfo(node, existingPod.Name, existingPod.Status.PodIP, nodeRepo)
	}

	// No healthy pod exists, create a new one
	return createMeshPod(ctx, namespace, podNamePrefix, node, clientSet, nodeRepo)
}

// findExistingMeshPod looks for an existing healthy mesh pod for the node.
// Returns nil if no healthy pod is found.
func findExistingMeshPod(ctx context.Context, namespace, meshPodName, podNamePrefix string, clientSet *kubernetes.Clientset) (*corev1.Pod, error) {
	// Use label selector to filter pods more efficiently
	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/component=mesh-node",
	})
	if err != nil {
		return nil, err
	}

	for i := range pods.Items {
		pod := &pods.Items[i]

		// Skip pods that don't match our node
		if !isPodForNode(pod, meshPodName, podNamePrefix) {
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
func isPodForNode(pod *corev1.Pod, meshPodName, podNamePrefix string) bool {
	if meshPodName != "" && pod.Name == meshPodName {
		return true
	}
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

// syncNodePodInfo updates the database with pod name and IP if they differ from what's stored.
// Uses an upsert pattern: checks if node exists in DB first, then updates or creates accordingly.
func syncNodePodInfo(node *db.Node, podName, podIP string, nodeRepo db.NodeRepo) error {
	// Check if pod info is already synced
	if node.MeshPodName == podName && node.MeshPodIp == podIP {
		log.Debugf("Mesh pod already exists and synced for node %s: %s (IP: %s)", node.NodeId, podName, podIP)
		return nil
	}

	// Check if node exists in database
	existingNode, err := nodeRepo.GetNode(node.NodeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Node doesn't exist, will create it
			existingNode = nil
		} else {
			return status.Errorf(codes.Internal, "failed to check if node exists: %v", err)
		}
	}

	if existingNode == nil {
		// Node doesn't exist in DB, create it
		log.Infof("Creating new node record for node %s with pod %s (IP: %s)", node.NodeId, podName, podIP)
		if err := nodeRepo.CreateNode(&db.Node{
			Id:          uuid.NewV4(),
			NodeId:      node.NodeId,
			MeshPodName: podName,
			MeshPodIp:   podIP,
			MeshPodPort: 8082,
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to create node record: %v", err)
		}
		log.Infof("Successfully created node record for node %s", node.NodeId)
	} else {
		// Node exists in DB, update it
		log.Infof("Updating mesh pod info in database for node %s: %s (IP: %s) -> %s (IP: %s)", 
			node.NodeId, existingNode.MeshPodName, existingNode.MeshPodIp, podName, podIP)
		if err := nodeRepo.UpdateNode(&db.Node{
			Id:          existingNode.Id,
			NodeId:      node.NodeId,
			MeshPodName: podName,
			MeshPodIp:   podIP,
			MeshPodPort: node.MeshPodPort,
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to update node record: %v", err)
		}
		log.Infof("Successfully updated node record for node %s", node.NodeId)
	}

	return nil
}

// createMeshPod creates a new mesh pod for the node.
func createMeshPod(ctx context.Context, namespace, podName string, node *db.Node, clientSet *kubernetes.Clientset, nodeRepo db.NodeRepo) error {
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
				"app.kubernetes.io/instance":  node.NodeId,
			},
		},
		Spec: *podSpec,
	}

	createdPod, err := clientSet.CoreV1().Pods(namespace).Create(ctx, newPod, metav1.CreateOptions{})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create mesh pod: %v", err)
	}

	log.Infof("Created mesh pod %s for node %s, waiting for IP assignment...", createdPod.Name, node.NodeId)

	// Wait for the pod to get an IP address
	podIP, err := waitForPodIP(ctx, namespace, createdPod.Name, clientSet)
	if err != nil {
		log.Warnf("Failed to get pod IP for %s: %v (will be updated on next call)", createdPod.Name, err)
		// Still save the pod name even if IP is not available yet
		podIP = ""
	}

	log.Infof("Mesh pod %s for node %s has IP: %s", createdPod.Name, node.NodeId, podIP)

	// Update database with pod info
	return syncNodePodInfo(node, createdPod.Name, podIP, nodeRepo)
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