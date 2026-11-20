/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"
)

const (
	testMeshPodIp    = "10.1.0.1"
	testOfflineRoute = "event.cloud.local.testorg.messaging.mesh.node.offline"
)

func meshPod(ip string, restarts int32, running bool) *corev1.Pod {
	state := corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{ExitCode: 137}}
	if running {
		state = corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mesh-pod",
			Labels: map[string]string{
				meshComponentLabel:  meshComponent,
				meshInstanceIdLabel: testValidNodeID,
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			PodIP: ip,
			ContainerStatuses: []corev1.ContainerStatus{
				{Name: "mesh", RestartCount: restarts, State: state},
			},
		},
	}
}

func TestMeshDied(t *testing.T) {
	failed := meshPod(testMeshPodIp, 0, true)
	failed.Status.Phase = corev1.PodFailed

	notReady := meshPod(testMeshPodIp, 0, true)
	notReady.Status.Conditions = []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionFalse}}

	tests := []struct {
		name   string
		oldPod *corev1.Pod
		newPod *corev1.Pod
		want   bool
	}{
		{"containerStopped", meshPod(testMeshPodIp, 0, true), meshPod(testMeshPodIp, 0, false), true},
		{"containerRestarted", meshPod(testMeshPodIp, 0, true), meshPod(testMeshPodIp, 1, true), true},
		{"podFailed", meshPod(testMeshPodIp, 0, true), failed, true},
		{"stillRunning", meshPod(testMeshPodIp, 0, true), meshPod(testMeshPodIp, 0, true), false},
		{"readinessLost", meshPod(testMeshPodIp, 0, true), notReady, false},
		{"containerStarted", meshPod(testMeshPodIp, 1, false), meshPod(testMeshPodIp, 1, true), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, meshDied(tc.oldPod, tc.newPod))
		})
	}
}

func TestMeshWatcherHandleMeshDown(t *testing.T) {
	ctx := context.Background()

	newWatcher := func(entry *pkg.NodeMeshMap, getErr error) (*MeshWatcher, *cmocks.MsgBusServiceClient, *[]string) {
		cleared := []string{}
		store := &fakeNnsStore{
			get: func(_ context.Context, _ string) (*pkg.NodeMeshMap, error) {
				if getErr != nil {
					return nil, getErr
				}
				e := *entry
				return &e, nil
			},
			updateMesh: func(_ context.Context, nodeId string, ip string, _ int32) error {
				cleared = append(cleared, nodeId+"="+ip)
				entry.MeshIp = ip
				return nil
			},
		}
		msgBus := &cmocks.MsgBusServiceClient{}
		return NewMeshWatcher(NewNnsServer(store, nil, nil), msgBus, testEventOrg), msgBus, &cleared
	}

	t.Run("publishesOfflineAndKeepsEntry", func(t *testing.T) {
		w, msgBus, cleared := newWatcher(storedNode(), nil)
		msgBus.On("PublishRequest", testOfflineRoute, mock.MatchedBy(func(e *epb.NodeOfflineEvent) bool {
			return e.GetNodeId() == testValidNodeID
		})).Return(nil).Once()

		w.handleMeshDown(ctx, meshPod(testMeshPodIp, 0, false))

		msgBus.AssertExpectations(t)
		assert.Equal(t, []string{testValidNodeID + "="}, *cleared)
	})

	t.Run("secondEventForSamePodPublishesNothing", func(t *testing.T) {
		w, msgBus, _ := newWatcher(storedNode(), nil)
		msgBus.On("PublishRequest", testOfflineRoute, mock.Anything).Return(nil).Once()

		w.handleMeshDown(ctx, meshPod(testMeshPodIp, 0, false))
		w.handleMeshDown(ctx, meshPod(testMeshPodIp, 1, true))

		msgBus.AssertNumberOfCalls(t, "PublishRequest", 1)
	})

	t.Run("nodeAlreadyDisconnected", func(t *testing.T) {
		entry := storedNode()
		entry.MeshIp = ""
		w, msgBus, cleared := newWatcher(entry, nil)

		w.handleMeshDown(ctx, meshPod(testMeshPodIp, 0, false))

		msgBus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
		assert.Empty(t, *cleared)
	})

	t.Run("nodeReconnectedThroughNewerPod", func(t *testing.T) {
		w, msgBus, cleared := newWatcher(storedNode(), nil)

		w.handleMeshDown(ctx, meshPod("10.1.0.99", 0, false))

		msgBus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
		assert.Empty(t, *cleared)
	})

	t.Run("unknownNode", func(t *testing.T) {
		w, msgBus, _ := newWatcher(nil, fmt.Errorf("%w: mesh|x", pkg.ErrNodeNotFound))

		w.handleMeshDown(ctx, meshPod(testMeshPodIp, 0, false))

		msgBus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
	})

	t.Run("podWithoutNodeLabel", func(t *testing.T) {
		w, msgBus, _ := newWatcher(storedNode(), nil)
		pod := meshPod(testMeshPodIp, 0, false)
		delete(pod.Labels, meshInstanceIdLabel)

		w.handleMeshDown(ctx, pod)

		msgBus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
	})

	t.Run("failedPublishLeavesEntryForRetry", func(t *testing.T) {
		w, msgBus, cleared := newWatcher(storedNode(), nil)
		msgBus.On("PublishRequest", testOfflineRoute, mock.Anything).Return(errors.New("msgclient down")).Once()

		w.handleMeshDown(ctx, meshPod(testMeshPodIp, 0, false))

		assert.Empty(t, *cleared)
	})
}
