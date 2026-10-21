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
	"fmt"
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
)

func templateDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ukama-mesh-svc",
			Namespace: testNamespace,
			Labels:    map[string]string{NameLabel: MeshName},
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

// The fake clientset ignores GenerateName, which the API server would expand.
func newClientset(objs ...runtime.Object) *fake.Clientset {
	cs := fake.NewSimpleClientset(objs...)
	var n atomic.Int32
	cs.PrependReactor("create", "pods", func(a k8stesting.Action) (bool, runtime.Object, error) {
		pod := a.(k8stesting.CreateAction).GetObject().(*corev1.Pod)
		if pod.Name == "" && pod.GenerateName != "" {
			pod.Name = fmt.Sprintf("%s%05d", pod.GenerateName, n.Add(1))
		}
		return false, nil, nil
	})
	return cs
}

func newManager(cs *fake.Clientset) *Manager {
	return NewManager(cs, Config{Namespace: testNamespace, OrgName: testOrg})
}

func nodePods(t *testing.T, cs *fake.Clientset, nodeId string) []corev1.Pod {
	t.Helper()
	pods, err := cs.CoreV1().Pods(testNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: InstanceIdLabel + "=" + nodeId,
	})
	require.NoError(t, err)
	return pods.Items
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
	cs := newClientset(templateDeployment())
	m := newManager(cs)

	st := assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assert.Contains(t, st.Message(), "is not ready (phase Pending)")

	pods := nodePods(t, cs, testNode)
	require.Len(t, pods, 1)
	assert.Equal(t, corev1.RestartPolicyOnFailure, pods[0].Spec.RestartPolicy)

	svc, err := cs.CoreV1().Services(testNamespace).Get(context.Background(), testMeshName, metav1.GetOptions{})
	require.NoError(t, err, "the Service is created on the first request with the pod")
	assert.True(t, serviceMatches(svc, testNode))

	assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assert.Len(t, nodePods(t, cs, testNode), 1, "a retry before the pod is Ready must not create a second pod")

	markPodReady(t, cs, pods[0].Name)
	assert.NoError(t, m.Ensure(context.Background(), testNode))
}

func TestEnsure_ReadyMeshOnlyReads(t *testing.T) {
	cs := newClientset(
		templateDeployment(),
		meshPod(testMeshName+"-w76pt", testNode, corev1.PodRunning, true),
		meshService(selector(testNode)),
	)
	m := newManager(cs)
	cs.ClearActions()

	require.NoError(t, m.Ensure(context.Background(), testNode))

	for _, a := range cs.Actions() {
		assert.Contains(t, []string{"list", "get"}, a.GetVerb(), "unexpected %s %s", a.GetVerb(), a.GetResource().Resource)
	}
}

func TestEnsure_ReadyPodWithoutServiceGetsService(t *testing.T) {
	cs := newClientset(templateDeployment(), meshPod(testMeshName+"-w76pt", testNode, corev1.PodRunning, true))
	m := newManager(cs)

	require.NoError(t, m.Ensure(context.Background(), testNode))

	_, err := cs.CoreV1().Services(testNamespace).Get(context.Background(), testMeshName, metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, 0, countActions(cs, "create", "pods"))
}

func TestEnsure_FixesServiceSelector(t *testing.T) {
	cs := newClientset(
		templateDeployment(),
		meshPod(testMeshName+"-w76pt", testNode, corev1.PodRunning, true),
		meshService(map[string]string{ComponentLabel: Component, InstanceIdLabel: "uk-sa2643-tnode-v0-0000"}),
	)
	m := newManager(cs)

	require.NoError(t, m.Ensure(context.Background(), testNode))

	svc, err := cs.CoreV1().Services(testNamespace).Get(context.Background(), testMeshName, metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, selector(testNode), svc.Spec.Selector)
}

func TestEnsure_FinishedPodIsReplaced(t *testing.T) {
	cs := newClientset(
		templateDeployment(),
		meshPod(testMeshName+"-w76pt", testNode, corev1.PodFailed, false),
		meshService(selector(testNode)),
	)
	m := newManager(cs)

	assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)

	assert.Len(t, nodePods(t, cs, testNode), 2)
	assert.Equal(t, 1, countActions(cs, "create", "pods"))
}

func TestEnsure_ThrottledCreateIsUnavailable(t *testing.T) {
	cs := newClientset(templateDeployment())
	var calls atomic.Int32
	cs.PrependReactor("create", "pods", func(k8stesting.Action) (bool, runtime.Object, error) {
		if calls.Add(1) == 1 {
			return true, nil, errors.New("client rate limiter Wait returned an error: rate: Wait(n=1) would exceed context deadline")
		}
		return false, nil, nil
	})
	m := newManager(cs)

	st := assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assert.Contains(t, st.Message(), "failed to create mesh pod")
	assert.Contains(t, st.Message(), "rate limiter")
	assert.Empty(t, nodePods(t, cs, testNode))

	st = assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assert.Contains(t, st.Message(), "is not ready")
	assert.Len(t, nodePods(t, cs, testNode), 1, "the next request creates the pod")
}

func TestEnsure_MissingTemplateIsInternal(t *testing.T) {
	cs := newClientset()
	m := newManager(cs)

	st := assertCode(t, m.Ensure(context.Background(), testNode), codes.Internal)
	assert.Equal(t, "no mesh deployment found in namespace "+testNamespace, st.Message())
}

func TestEnsure_ForbiddenIsInternal(t *testing.T) {
	cs := newClientset(templateDeployment())
	cs.PrependReactor("create", "pods", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewForbidden(schema.GroupResource{Resource: "pods"}, "", errors.New("RBAC: access denied"))
	})
	m := newManager(cs)

	st := assertCode(t, m.Ensure(context.Background(), testNode), codes.Internal)
	assert.Contains(t, st.Message(), "failed to create mesh pod")
}

func TestEnsure_TemplateIsCachedAcrossNodes(t *testing.T) {
	cs := newClientset(templateDeployment())
	m := newManager(cs)

	assertCode(t, m.Ensure(context.Background(), testNode), codes.Unavailable)
	assertCode(t, m.Ensure(context.Background(), "uk-sa2643-anode-v0-5db9"), codes.Unavailable)

	assert.Equal(t, 1, countActions(cs, "list", "deployments"))
	assert.Equal(t, 2, countActions(cs, "create", "pods"))
}
