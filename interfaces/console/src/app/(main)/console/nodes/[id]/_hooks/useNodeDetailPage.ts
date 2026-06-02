/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import {
  NodeConnectivityEnum,
  NodeTypeEnum,
} from '@/client/graphql/generated';
import { GetMetricsStatQuery, MetricsRes, Stats_Type } from '@/client/graphql/generated/subscriptions';
import { nodeMetricsVar } from '@/client/vars';
import {
  NODE_ACTIONS_BUTTONS,
  NODE_ACTIONS_ENUM,
  NODE_KPIS,
  STAT_STEP_29,
} from '@/constants';
import { useEnvContext, useUserContext, useUIContext } from '@/context';
import { TNodeActionState } from '@/types';
import {
  getNodeActionDescriptionByProgress,
  getNodeTypeFromId,
  getUnixTime,
  nodeTypeEnumToString,
} from '@/utils';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useNodeActions } from './useNodeActions';
import { useNodeData } from './useNodeData';
import { useNodeMetrics } from './useNodeMetrics';
import { useSubscriptionManager } from '@/features/subscriptions/useSubscriptionManager';
import { useRouter } from 'next/navigation';

export function useNodeDetailPage(id: string) {
  const nodeType = nodeTypeEnumToString(getNodeTypeFromId(id) as NodeTypeEnum);
  const router = useRouter();

  const { subscriptionClient } = useEnvContext();
  const { user } = useUserContext();
  const { setSnackbarMessage } = useUIContext();

  // Shared state that is used across hooks
  const [isEditNode, setIsEditNode] = useState<boolean>(false);
  const [nodeAction, setNodeAction] = useState<
    TNodeActionState & { currentAction: string; actionInitiated: string }
  >({
    progress: 0,
    currentAction: '',
    actionInitiated: NODE_ACTIONS_ENUM.NODE_LOADING,
    action: NODE_ACTIONS_ENUM.NODE_LOADING,
    isActive: false,
  });

  const { subscribe, unsubscribe } = useSubscriptionManager();

  const subscriptionKeyRef = useRef<string | null>(null);
  const subscriptionControllerRef = useRef<{ cancel: () => void } | null>(null);

  const notify = (msgId: string, message: string, type: string) =>
    setSnackbarMessage({ id: msgId, message, type, show: true });

  const cleanupSubscription = useCallback(() => {
    if (subscriptionKeyRef.current) {
      unsubscribe(subscriptionKeyRef.current);
      subscriptionKeyRef.current = null;
    }
    if (subscriptionControllerRef.current) {
      subscriptionControllerRef.current.cancel();
      subscriptionControllerRef.current = null;
    }
  }, [unsubscribe]);

  useEffect(() => () => cleanupSubscription(), [cleanupSubscription]);

  // Metrics hook — provides metric state and subscription helpers
  const {
    metricFrom,
    setMetricFrom,
    graphType,
    nodeUptime,
    setNodeUptime,
    selectedTab,
    metrics,
    handleStatSubscription,
    startStatSubscription,
  } = useNodeMetrics({ id, nodeType, getMetricStat: () => {}, cleanupSubscription, subscribe });

  // Data hook — GraphQL queries/mutations
  const {
    nodesData,
    nodesLoading,
    currentNode,
    updateNodeLoading,
    appsLoading,
    softwaresLoading,
    softwaresData,
    updateSoftwareLoading,
    statData,
    statLoading,
    statVar,
    nodeMetricsLoading,
    getMetricStat,
    fetchMetricByTab,
    handleUpdateAvailable,
    handleEditNode: doEditNode,
    handleRestartNode,
  } = useNodeData({
    id,
    graphType,
    metricFrom,
    subscriptionClient,
    onMetricsFetched: (data: MetricsRes) => nodeMetricsVar(data),
    onStatFetched: (data) => {
      if (!data) return;
      const metrics = (data as GetMetricsStatQuery).getMetricsStat?.metrics ?? [];
      if (metrics.length > 0) {
        metrics.forEach((m: { type: string; value: number }) => {
          if (m.type === NODE_KPIS.NODE_UPTIME[nodeType][0].id) {
            setNodeUptime(m.value);
          }
        });
        const from = statVar?.data.from ?? 0;
        startStatSubscription(from, from, subscriptionKeyRef, subscriptionControllerRef);
        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.AllNode}-${from}`;
        PubSub.subscribe(sKey, handleStatSubscription);
        subscribe(sKey, () => PubSub.unsubscribe(sKey));
      }
    },
    setNodeUptime,
    setNodeAction,
  });

  // Actions hook — event handlers for UI interactions
  const {
    handleNodeSelected,
    handleEditNode,
    handleSectionChange,
    onTabSelected,
    handleNodeActionClick,
  } = useNodeActions({
    currentNodeId: currentNode?.id,
    setIsEditNode,
    setNodeAction,
    setMetricFrom,
    handleEditNode: doEditNode,
    handleRestartNode,
  });

  // Reset nodeAction when node comes back online
  useEffect(() => {
    if (currentNode?.status.connectivity === NodeConnectivityEnum.Online) {
      setNodeAction({
        progress: 0,
        currentAction: '',
        actionInitiated: '',
        action: '',
        isActive: false,
      });
    }
  }, [currentNode]);

  // Initial stat fetch on mount / id change
  useEffect(() => {
    const to = getUnixTime();
    const from = to - STAT_STEP_29;
    if (!id) {
      notify('node-not-found-msg', 'Node not found.', 'error');
      router.back();
    } else {
      // Reset stale metrics when node changes
      nodeMetricsVar({ metrics: [] });
      cleanupSubscription();
      getMetricStat({
        variables: {
          data: {
            to,
            nodeId: id,
            from,
            userId: user.id,
            step: STAT_STEP_29,
            orgName: user.orgName,
            withSubscription: true,
            type: Stats_Type.AllNode,
          },
        },
      });
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id, user.id, user.orgName, cleanupSubscription]);

  // Track node restart progress via connectivity transitions
  useEffect(() => {
    let intervalId: NodeJS.Timeout | null = null;

    if (nodesData?.getNodes.nodes.length && nodeAction.actionInitiated !== '') {
      const s = nodesData.getNodes.nodes.find((n) => n.id === id);
      if (s && s.status.connectivity !== nodeAction.currentAction) {
        if (nodeAction.actionInitiated === NODE_ACTIONS_ENUM.NODE_RESTART) {
          switch (s.status.connectivity) {
            case NodeConnectivityEnum.Offline:
              setNodeAction((prev) => ({
                ...prev,
                progress: prev.progress + 25,
                currentAction: s.status.connectivity.toString(),
              }));
              return;
            case NodeConnectivityEnum.Online:
              setNodeAction((prev) => ({
                ...prev,
                progress: prev.progress + 25,
                currentAction: s.status.connectivity.toString(),
              }));
              break;
          }

          intervalId = setInterval(() => {
            setNodeAction((prev) => {
              const newProgress =
                prev.progress === 100 ? 0 : prev.progress + 25;
              return {
                ...prev,
                progress: newProgress,
                currentAction:
                  newProgress === 0 ? '' : s.status.connectivity.toString(),
                actionInitiated: newProgress === 0 ? '' : prev.actionInitiated,
                action: newProgress === 0 ? '' : prev.actionInitiated,
                isActive: newProgress !== 0,
              };
            });
          }, 5000);
        }
      }
    }

    return () => {
      if (intervalId) {
        setNodeAction({
          progress: 0,
          currentAction: '',
          actionInitiated: '',
          action: '',
          isActive: false,
        });
        clearInterval(intervalId);
      }
    };
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [nodesData]);

  // Refetch metrics when metricFrom or graphType changes
  useEffect(() => {
    fetchMetricByTab();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [metricFrom, graphType]);

  const nodeActionDescription = getNodeActionDescriptionByProgress(
    nodeAction.progress,
    nodeAction.actionInitiated ||
      currentNode?.status?.connectivity?.toString() ||
      '',
  );

  return {
    id,
    currentNode,
    nodesData,
    nodesLoading,
    updateNodeLoading,
    statLoading,
    statData,
    metrics,
    nodeUptime,
    metricFrom,
    graphType,
    nodeMetricsLoading,
    selectedTab,
    isEditNode,
    nodeAction,
    nodeActionDescription,
    appsLoading,
    softwaresLoading,
    softwaresData,
    updateSoftwareLoading,
    actionOptions: NODE_ACTIONS_BUTTONS,
    onOpenEditNode: () => setIsEditNode(true),
    onCloseEditNode: () => setIsEditNode(false),
    handleNodeSelected,
    handleEditNode,
    handleSectionChange,
    onTabSelected,
    handleNodeActionClick,
    handleUpdateAvailable,
  };
}
