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

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	evt "github.com/ukama/ukama/systems/common/events"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"
)

const (
	meshComponentLabel  = "app.kubernetes.io/component"
	meshInstanceIdLabel = "app.kubernetes.io/instance-id"
	meshComponent       = "mesh-node"
)

type MeshWatcher struct {
	nns          *NnsServer
	msgbus       mb.MsgBusServiceClient
	offlineRoute string
}

func NewMeshWatcher(nns *NnsServer, msgBus mb.MsgBusServiceClient, orgName string) *MeshWatcher {
	return &MeshWatcher{
		nns:          nns,
		msgbus:       msgBus,
		offlineRoute: msgbus.PrepareRoute(orgName, evt.EventRoutingKey[evt.EventNodeOffline]),
	}
}

func (w *MeshWatcher) Start(ctx context.Context, clientSet kubernetes.Interface, namespace string) error {
	factory := informers.NewSharedInformerFactoryWithOptions(clientSet, 0,
		informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(o *metav1.ListOptions) {
			o.LabelSelector = meshComponentLabel + "=" + meshComponent
		}))

	informer := factory.Core().V1().Pods().Informer()

	err := informer.SetWatchErrorHandler(func(_ *cache.Reflector, err error) {
		log.Warnf("Mesh pod watch in namespace %s failed: %v", namespace, err)
	})
	if err != nil {
		return err
	}

	_, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod, okOld := oldObj.(*corev1.Pod)
			newPod, okNew := newObj.(*corev1.Pod)
			if okOld && okNew && meshDied(oldPod, newPod) {
				w.handleMeshDown(ctx, newPod)
			}
		},
		DeleteFunc: func(obj interface{}) {
			if tombstone, ok := obj.(cache.DeletedFinalStateUnknown); ok {
				obj = tombstone.Obj
			}
			if pod, ok := obj.(*corev1.Pod); ok {
				w.handleMeshDown(ctx, pod)
			}
		},
	})
	if err != nil {
		return err
	}

	factory.Start(ctx.Done())
	log.Infof("Watching mesh pods in namespace %s", namespace)

	return nil
}

func meshDied(oldPod, newPod *corev1.Pod) bool {
	if isPodFinished(newPod) && !isPodFinished(oldPod) {
		return true
	}

	before := make(map[string]corev1.ContainerStatus, len(oldPod.Status.ContainerStatuses))
	for _, c := range oldPod.Status.ContainerStatuses {
		before[c.Name] = c
	}

	for _, c := range newPod.Status.ContainerStatuses {
		prev, ok := before[c.Name]
		if !ok {
			continue
		}
		if c.RestartCount > prev.RestartCount {
			return true
		}
		if prev.State.Running != nil && c.State.Running == nil {
			return true
		}
	}

	return false
}

func isPodFinished(pod *corev1.Pod) bool {
	return pod.Status.Phase == corev1.PodFailed || pod.Status.Phase == corev1.PodSucceeded
}

func (w *MeshWatcher) handleMeshDown(ctx context.Context, pod *corev1.Pod) {
	nodeId := pod.Labels[meshInstanceIdLabel]
	if nodeId == "" {
		return
	}

	node, err := w.nns.nns.Get(ctx, nodeId)
	if err != nil {
		if !errors.Is(err, pkg.ErrNodeNotFound) {
			log.Errorf("Mesh pod %s died but node %s lookup failed: %v", pod.Name, nodeId, err)
		}
		return
	}

	if node.MeshIp == "" || (pod.Status.PodIP != "" && node.MeshIp != pod.Status.PodIP) {
		return
	}

	log.Warnf("Mesh pod %s died while node %s was connected. Publishing offline.", pod.Name, nodeId)

	err = w.msgbus.PublishRequest(w.offlineRoute, &epb.NodeOfflineEvent{NodeId: nodeId})
	if err != nil {
		log.Errorf("Failed to publish offline for node %s: %v", nodeId, err)
		return
	}

	if _, err = w.nns.markMeshDown(ctx, nodeId, node.MeshIp); err != nil {
		log.Errorf("Failed to mark node %s as having no mesh. Error: %v", nodeId, err)
	}
}
