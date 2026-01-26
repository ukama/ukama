package utils

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	PodNamePrefix = "mesh-node"
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
	podNamePrefix := GetPodNamePrefix(config.OrgName)

	// Check for existing healthy pod first
	existingPod, err := findExistingMeshPod(ctx, namespace, node.MeshPodName, podNamePrefix, clientSet)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to check existing pods: %v", err)
	}

	// If a healthy pod exists, sync database if needed and return
	if existingPod != nil {
		return syncNodePodName(node, existingPod.Name, nodeRepo)
	}

	// No healthy pod exists, create a new one
	return createMeshPod(ctx, namespace, GetPodName(config.OrgName, node.NodeId), node, clientSet, nodeRepo)
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

// syncNodePodName updates the database if the pod name differs from what's stored.
func syncNodePodName(node *db.Node, podName string, nodeRepo db.NodeRepo) error {
	if node.MeshPodName == podName {
		log.Debugf("Mesh pod already exists and synced for node %s: %s", node.NodeId, podName)
		return nil
	}

	log.Infof("Syncing mesh pod name in database for node %s: %s -> %s", node.NodeId, node.MeshPodName, podName)
	
	if node.MeshPodName == "" {
		// Node doesn't exist in DB yet
		if err := nodeRepo.CreateNode(&db.Node{
			NodeId:      node.NodeId,
			MeshPodName: podName,
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to create node record: %v", err)
		}
	} else {
		// Update existing node
		if err := nodeRepo.UpdateNode(&db.Node{
			NodeId:      node.NodeId,
			MeshPodName: podName,
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to update node record: %v", err)
		}
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

	// Create the pod
	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: podName,
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

	log.Infof("Created mesh pod %s for node %s", createdPod.Name, node.NodeId)

	// Update database
	return syncNodePodName(node, createdPod.Name, nodeRepo)
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