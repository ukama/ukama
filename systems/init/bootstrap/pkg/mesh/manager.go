/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package mesh

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

const (
	NameLabel       = "app.kubernetes.io/name"
	ComponentLabel  = "app.kubernetes.io/component"
	InstanceIdLabel = "app.kubernetes.io/instance-id"
	MeshName        = "mesh"
	Component       = "mesh-node"
	ServicePort     = 8082
	WebsocketPort   = 18200

	defaultTemplateRefresh = 5 * time.Minute
)

type Config struct {
	Namespace       string
	OrgName         string
	TemplateRefresh time.Duration
}

type Manager struct {
	clientSet kubernetes.Interface
	config    Config

	mu         sync.Mutex
	template   *corev1.PodSpec
	templateAt time.Time
}

func NewManager(clientSet kubernetes.Interface, config Config) *Manager {
	if config.TemplateRefresh <= 0 {
		config.TemplateRefresh = defaultTemplateRefresh
	}
	return &Manager{clientSet: clientSet, config: config}
}

// Ensure never waits; a pod that is still starting returns Unavailable and the node's retry finds it Ready.
func (m *Manager) Ensure(ctx context.Context, nodeId string) error {
	pod, err := m.findPod(ctx, nodeId)
	if err != nil {
		return err
	}

	if pod == nil {
		pod, err = m.createPod(ctx, nodeId)
		if err != nil {
			return err
		}
	}

	if err := m.ensureService(ctx, nodeId); err != nil {
		return err
	}

	if !isPodReady(pod) {
		return status.Errorf(codes.Unavailable, "mesh pod %s is not ready (phase %s)", pod.Name, podPhase(pod))
	}
	return nil
}

// Pods from older bootstrap builds have the same labels, so a node's pod is found by label, not name.
func (m *Manager) findPod(ctx context.Context, nodeId string) (*corev1.Pod, error) {
	pods, err := m.clientSet.CoreV1().Pods(m.config.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s,%s=%s", ComponentLabel, Component, InstanceIdLabel, nodeId),
	})
	if err != nil {
		return nil, apiError("list mesh pods", err)
	}

	var best *corev1.Pod
	for i := range pods.Items {
		pod := &pods.Items[i]
		if !isPodActive(pod) {
			continue
		}
		if best == nil || podRank(pod) > podRank(best) {
			best = pod
		}
	}
	return best, nil
}

func (m *Manager) createPod(ctx context.Context, nodeId string) (*corev1.Pod, error) {
	spec, err := m.templateSpec(ctx)
	if err != nil {
		return nil, err
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: m.objectName(nodeId) + "-",
			Namespace:    m.config.Namespace,
			Labels:       m.labels(nodeId),
		},
		Spec: *spec,
	}

	created, err := m.clientSet.CoreV1().Pods(m.config.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, apiError("create mesh pod", err)
	}

	log.Infof("Created mesh pod %s for node %s", created.Name, nodeId)
	return created, nil
}

func (m *Manager) ensureService(ctx context.Context, nodeId string) error {
	name := m.objectName(nodeId)

	svc, err := m.clientSet.CoreV1().Services(m.config.Namespace).Get(ctx, name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = m.clientSet.CoreV1().Services(m.config.Namespace).Create(ctx, m.newService(nodeId), metav1.CreateOptions{})
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		if err != nil {
			return apiError("create mesh service", err)
		}
		log.Infof("Created mesh service %s for node %s", name, nodeId)
		return nil
	}
	if err != nil {
		return apiError("get mesh service", err)
	}

	if serviceMatches(svc, nodeId) {
		return nil
	}

	updated := svc.DeepCopy()
	updated.Spec.Ports = servicePorts()
	updated.Spec.Selector = selector(nodeId)
	if _, err := m.clientSet.CoreV1().Services(m.config.Namespace).Update(ctx, updated, metav1.UpdateOptions{}); err != nil {
		return apiError("update mesh service", err)
	}

	log.Infof("Updated mesh service %s for node %s", name, nodeId)
	return nil
}

func (m *Manager) templateSpec(ctx context.Context) (*corev1.PodSpec, error) {
	m.mu.Lock()
	if m.template != nil && time.Since(m.templateAt) < m.config.TemplateRefresh {
		spec := m.template.DeepCopy()
		m.mu.Unlock()
		return spec, nil
	}
	m.mu.Unlock()

	deployments, err := m.clientSet.AppsV1().Deployments(m.config.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: NameLabel + "=" + MeshName,
	})

	m.mu.Lock()
	defer m.mu.Unlock()

	if err != nil {
		if m.template != nil {
			log.Warnf("Using cached mesh pod template: %v", err)
			return m.template.DeepCopy(), nil
		}
		return nil, apiError("list mesh deployments", err)
	}
	if len(deployments.Items) == 0 {
		return nil, status.Errorf(codes.Internal, "no mesh deployment found in namespace %s", m.config.Namespace)
	}

	spec := deployments.Items[0].Spec.Template.Spec.DeepCopy()
	spec.RestartPolicy = corev1.RestartPolicyOnFailure
	m.template = spec
	m.templateAt = time.Now()
	return spec.DeepCopy(), nil
}

func (m *Manager) objectName(nodeId string) string {
	return m.config.OrgName + "-" + Component + "-" + nodeId
}

func (m *Manager) labels(nodeId string) map[string]string {
	return map[string]string{
		NameLabel:       MeshName,
		ComponentLabel:  Component,
		InstanceIdLabel: nodeId,
	}
}

func (m *Manager) newService(nodeId string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.objectName(nodeId),
			Namespace: m.config.Namespace,
			Labels:    m.labels(nodeId),
		},
		Spec: corev1.ServiceSpec{
			Selector: selector(nodeId),
			Ports:    servicePorts(),
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
}

// Retrying cannot fix missing RBAC or a rejected object, so those are Internal rather than Unavailable.
func apiError(action string, err error) error {
	code := codes.Unavailable
	if apierrors.IsForbidden(err) || apierrors.IsUnauthorized(err) || apierrors.IsInvalid(err) || apierrors.IsBadRequest(err) {
		code = codes.Internal
	}
	return status.Errorf(code, "failed to %s: %v", action, err)
}

func selector(nodeId string) map[string]string {
	return map[string]string{
		ComponentLabel:  Component,
		InstanceIdLabel: nodeId,
	}
}

func servicePorts() []corev1.ServicePort {
	return []corev1.ServicePort{
		{
			Name:       "http",
			Protocol:   corev1.ProtocolTCP,
			Port:       ServicePort,
			TargetPort: intstr.FromInt32(ServicePort),
		},
		{
			Name:       "websocket",
			Protocol:   corev1.ProtocolTCP,
			Port:       WebsocketPort,
			TargetPort: intstr.FromInt32(WebsocketPort),
		},
	}
}

func serviceMatches(svc *corev1.Service, nodeId string) bool {
	want := selector(nodeId)
	if len(svc.Spec.Selector) != len(want) {
		return false
	}
	for k, v := range want {
		if svc.Spec.Selector[k] != v {
			return false
		}
	}

	ports := servicePorts()
	if len(svc.Spec.Ports) != len(ports) {
		return false
	}
	for _, p := range ports {
		found := false
		for _, have := range svc.Spec.Ports {
			if have.Name == p.Name && have.Port == p.Port && have.TargetPort == p.TargetPort &&
				(have.Protocol == p.Protocol || have.Protocol == "") {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func isPodActive(pod *corev1.Pod) bool {
	if pod.DeletionTimestamp != nil {
		return false
	}
	return pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodPending || pod.Status.Phase == ""
}

func podRank(pod *corev1.Pod) int {
	switch {
	case isPodReady(pod):
		return 3
	case pod.Status.Phase == corev1.PodRunning:
		return 2
	default:
		return 1
	}
}

func podPhase(pod *corev1.Pod) corev1.PodPhase {
	if pod.Status.Phase == "" {
		return corev1.PodPending
	}
	return pod.Status.Phase
}

func isPodReady(pod *corev1.Pod) bool {
	if pod.Status.Phase != corev1.PodRunning || pod.Status.PodIP == "" {
		return false
	}
	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.PodReady {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}
