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
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

const (
	testOrg       = "ukama"
	testNamespace = "ukama-messaging"
	testNode      = "uk-sa2643-tnode-v0-5db9"
	testMeshName  = "ukama-mesh-node-" + testNode
	testPodName   = testMeshName + "-6d4f9c7b8-x2k4p"
)

func templateDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ukama-mesh-svc",
			Namespace: testNamespace,
			Labels:    map[string]string{NameLabel: MeshName, "app.kubernetes.io/instance": "messaging"},
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyAlways,
					Containers:    []corev1.Container{{Name: "mesh", Image: "mesh:test"}},
				},
			},
		},
	}
}

func nodeDeployment(nodeId, image string) *appsv1.Deployment {
	m := &Manager{config: Config{Namespace: testNamespace, OrgName: testOrg}}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.objectName(nodeId),
			Namespace: testNamespace,
			Labels:    m.labels(nodeId),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: selector(nodeId)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: m.labels(nodeId)},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "mesh", Image: image}}},
			},
		},
	}
}

func meshPod(name, nodeId string, phase corev1.PodPhase, ready bool) *corev1.Pod {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
			Labels:    map[string]string{NameLabel: MeshName, ComponentLabel: Component, InstanceIdLabel: nodeId},
		},
		Status: corev1.PodStatus{Phase: phase},
	}
	if ready {
		pod.Status.PodIP = "10.0.110.252"
		pod.Status.Conditions = []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}}
	}
	return pod
}

func meshService(sel map[string]string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testMeshName,
			Namespace: testNamespace,
			Labels:    map[string]string{NameLabel: MeshName, ComponentLabel: Component, InstanceIdLabel: testNode},
		},
		Spec: corev1.ServiceSpec{Selector: sel, Ports: servicePorts(), Type: corev1.ServiceTypeClusterIP},
	}
}

func newManager(cs *fake.Clientset) *Manager {
	return NewManager(cs, Config{Namespace: testNamespace, OrgName: testOrg})
}

func nodeDeployments(t *testing.T, cs *fake.Clientset, nodeId string) []appsv1.Deployment {
	t.Helper()
	deployments, err := cs.AppsV1().Deployments(testNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: InstanceIdLabel + "=" + nodeId,
	})
	require.NoError(t, err)
	return deployments.Items
}

// The fake clientset runs no controllers, so tests create the pod a ReplicaSet would.
func startPod(t *testing.T, cs *fake.Clientset, ready bool) {
	t.Helper()
	phase := corev1.PodPending
	if ready {
		phase = corev1.PodRunning
	}
	_, err := cs.CoreV1().Pods(testNamespace).Create(context.Background(), meshPod(testPodName, testNode, phase, ready), metav1.CreateOptions{})
	require.NoError(t, err)
}

func markPodReady(t *testing.T, cs *fake.Clientset, name string) {
	t.Helper()
	pod, err := cs.CoreV1().Pods(testNamespace).Get(context.Background(), name, metav1.GetOptions{})
	require.NoError(t, err)
	pod.Status = meshPod(name, testNode, corev1.PodRunning, true).Status
	_, err = cs.CoreV1().Pods(testNamespace).UpdateStatus(context.Background(), pod, metav1.UpdateOptions{})
	require.NoError(t, err)
}

func countActions(cs *fake.Clientset, verb, resource string) int {
	n := 0
	for _, a := range cs.Actions() {
		if a.GetVerb() == verb && a.GetResource().Resource == resource {
			n++
		}
	}
	return n
}

func assertCode(t *testing.T, err error, want codes.Code) *status.Status {
	t.Helper()
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status, got %v", err)
	assert.Equal(t, want, st.Code(), st.Message())
	return st
}

func TestEnsure_NewNodeIsUnavailableUntilPodIsReady(t *testing.T) {
	cs := fake.NewSimpleClientset(templateDeployment())
	m := newManager(cs)

	st := assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assert.Equal(t, "mesh pod for node "+testNode+" is starting", st.Message())
	assert.Equal(t, 0, countActions(cs, "create", "pods"), "the Deployment creates the pod, not bootstrap")

	deployments := nodeDeployments(t, cs, testNode)
	require.Len(t, deployments, 1)
	d := deployments[0]
	assert.Equal(t, testMeshName, d.Name)
	assert.Equal(t, int32(1), *d.Spec.Replicas)
	assert.Equal(t, selector(testNode), d.Spec.Selector.MatchLabels)
	assert.Equal(t, appsv1.RecreateDeploymentStrategyType, d.Spec.Strategy.Type)
	assert.Equal(t, testNode, d.Spec.Template.Labels[InstanceIdLabel])
	assert.Equal(t, Component, d.Spec.Template.Labels[ComponentLabel])
	assert.Equal(t, corev1.RestartPolicyAlways, d.Spec.Template.Spec.RestartPolicy)
	assert.Equal(t, "mesh:test", d.Spec.Template.Spec.Containers[0].Image)

	svc, err := cs.CoreV1().Services(testNamespace).Get(context.Background(), testMeshName, metav1.GetOptions{})
	require.NoError(t, err, "the Service is created on the first request with the Deployment")
	assert.True(t, serviceMatches(svc, testNode))

	startPod(t, cs, false)
	st = assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assert.Equal(t, "mesh pod "+testPodName+" is not ready (phase Pending)", st.Message())

	markPodReady(t, cs, testPodName)
	assert.NoError(t, m.Ensure(context.Background(), testNode))
	assert.Len(t, nodeDeployments(t, cs, testNode), 1)
}

func TestEnsure_OverlappingRequestsCreateOneDeployment(t *testing.T) {
	cs := fake.NewSimpleClientset(templateDeployment())
	m := newManager(cs)

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
		}()
	}
	wg.Wait()

	assert.Equal(t, 2, countActions(cs, "create", "deployments"), "both requests tried to create")
	assert.Len(t, nodeDeployments(t, cs, testNode), 1, "the fixed name lets only one through")
}

func TestEnsure_ReadyMeshOnlyReads(t *testing.T) {
	cs := fake.NewSimpleClientset(
		templateDeployment(),
		nodeDeployment(testNode, "mesh:test"),
		meshPod(testPodName, testNode, corev1.PodRunning, true),
		meshService(selector(testNode)),
	)
	m := newManager(cs)
	cs.ClearActions()

	require.NoError(t, m.Ensure(context.Background(), testNode))

	for _, a := range cs.Actions() {
		assert.Contains(t, []string{"list", "get"}, a.GetVerb(), "unexpected %s %s", a.GetVerb(), a.GetResource().Resource)
	}
}

func TestEnsure_ReadyBarePodFromOlderBuildIsKept(t *testing.T) {
	cs := fake.NewSimpleClientset(templateDeployment(), meshPod(testMeshName+"-w76pt", testNode, corev1.PodRunning, true))
	m := newManager(cs)

	require.NoError(t, m.Ensure(context.Background(), testNode))

	_, err := cs.CoreV1().Services(testNamespace).Get(context.Background(), testMeshName, metav1.GetOptions{})
	assert.NoError(t, err, "a missing Service is created for the existing pod")
	assert.Equal(t, 0, countActions(cs, "create", "deployments"), "a Deployment next to the bare pod would put two pods behind the Service")
}

func TestEnsure_FixesServiceSelector(t *testing.T) {
	cs := fake.NewSimpleClientset(
		templateDeployment(),
		nodeDeployment(testNode, "mesh:test"),
		meshPod(testPodName, testNode, corev1.PodRunning, true),
		meshService(map[string]string{ComponentLabel: Component, InstanceIdLabel: "uk-sa2643-tnode-v0-0000"}),
	)
	m := newManager(cs)

	require.NoError(t, m.Ensure(context.Background(), testNode))

	svc, err := cs.CoreV1().Services(testNamespace).Get(context.Background(), testMeshName, metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, selector(testNode), svc.Spec.Selector)
}

func TestEnsure_FailedPodIsLeftToTheDeployment(t *testing.T) {
	cs := fake.NewSimpleClientset(
		templateDeployment(),
		nodeDeployment(testNode, "mesh:test"),
		meshPod(testPodName, testNode, corev1.PodFailed, false),
		meshService(selector(testNode)),
	)
	m := newManager(cs)

	st := assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assert.Equal(t, "mesh pod for node "+testNode+" is starting", st.Message())
	assert.Equal(t, 0, countActions(cs, "create", "pods"))
	assert.Equal(t, 0, countActions(cs, "delete", "pods"))
	assert.Len(t, nodeDeployments(t, cs, testNode), 1)
}

func TestEnsure_ThrottledCreateIsUnavailable(t *testing.T) {
	cs := fake.NewSimpleClientset(templateDeployment())
	var calls atomic.Int32
	cs.PrependReactor("create", "deployments", func(k8stesting.Action) (bool, runtime.Object, error) {
		if calls.Add(1) == 1 {
			return true, nil, errors.New("client rate limiter Wait returned an error: rate: Wait(n=1) would exceed context deadline")
		}
		return false, nil, nil
	})
	m := newManager(cs)

	st := assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assert.Contains(t, st.Message(), "failed to create mesh deployment")
	assert.Contains(t, st.Message(), "rate limiter")
	assert.Empty(t, nodeDeployments(t, cs, testNode))

	assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assert.Len(t, nodeDeployments(t, cs, testNode), 1, "the next request creates the Deployment")
}

func TestEnsure_MissingTemplateIsInternal(t *testing.T) {
	cs := fake.NewSimpleClientset()
	m := newManager(cs)

	st := assertCode(t, m.Ensure(context.Background(), testNode), codes.Internal)
	assert.Equal(t, "no mesh deployment found in namespace "+testNamespace, st.Message())
}

func TestEnsure_ForbiddenIsInternal(t *testing.T) {
	cs := fake.NewSimpleClientset(templateDeployment())
	cs.PrependReactor("create", "deployments", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewForbidden(schema.GroupResource{Group: "apps", Resource: "deployments"}, "", errors.New("RBAC: access denied"))
	})
	m := newManager(cs)

	st := assertCode(t, m.Ensure(context.Background(), testNode), codes.Internal)
	assert.Contains(t, st.Message(), "failed to create mesh deployment")
}

func TestEnsure_TemplateLookupSkipsNodeDeployments(t *testing.T) {
	cs := fake.NewSimpleClientset(templateDeployment(), nodeDeployment("uk-sa2643-anode-v0-0001", "mesh:stale"))
	m := newManager(cs)

	assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)

	deployments := nodeDeployments(t, cs, testNode)
	require.Len(t, deployments, 1)
	assert.Equal(t, "mesh:test", deployments[0].Spec.Template.Spec.Containers[0].Image)
}

func TestEnsure_TemplateIsCachedAcrossNodes(t *testing.T) {
	cs := fake.NewSimpleClientset(templateDeployment())
	m := newManager(cs)

	assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assertCode(t, m.Ensure(context.Background(), "uk-sa2643-anode-v0-5db9"), codes.Unavailable)

	assert.Equal(t, 1, countActions(cs, "list", "deployments"))
	assert.Equal(t, 2, countActions(cs, "create", "deployments"))
}
