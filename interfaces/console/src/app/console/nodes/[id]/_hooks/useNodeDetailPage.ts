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
  NodeStateEnum,
  NodeTypeEnum,
  SoftwareStatusEnum,
  useGetAppsQuery,
  useGetNodesQuery,
  useRestartNodeMutation,
  useSoftwareQuery,
  useUpdateNodeMutation,
  useUpdateSoftwareMutation,
} from '@/client/graphql/generated';
import {
  Graphs_Type,
  MetricsRes,
  Stats_Type,
  useGetMetricByTabLazyQuery,
  useGetMetricsStatLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import {
  METRIC_RANGE_10800,
  NODE_ACTIONS_BUTTONS,
  NODE_ACTIONS_ENUM,
  NODE_KPIS,
  STAT_STEP_29,
} from '@/constants';
import { useEnvContext, useUserContext, useUIContext } from '@/context';
import MetricStatSubscription from '@/features/subscriptions/MetricStatSubscription';
import { TMetricResDto, TNodeActionState, TStatusBarObj } from '@/types';
import {
  getNodeActionDescriptionByProgress,
  getNodeTabTypeByIndex,
  getNodeTypeFromId,
  getUnixTime,
  nodeTypeEnumToString,
} from '@/utils';
import { useRouter } from 'next/navigation';
import { useCallback, useEffect, useRef, useState } from 'react';

export function useNodeDetailPage(id: string) {
  const nodeType = nodeTypeEnumToString(getNodeTypeFromId(id) as NodeTypeEnum);
  const router = useRouter();

  const [metricFrom, setMetricFrom] = useState<number>(0);
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
  const [graphType, setGraphType] = useState<Graphs_Type>(
    Graphs_Type.NodeHealth,
  );
  const [nodeUptime, setNodeUptime] = useState<number>(0);
  const [selectedTab, setSelectedTab] = useState<number>(0);
  const [metrics, setMetrics] = useState<MetricsRes>({ metrics: [] });

  const { env, subscriptionClient } = useEnvContext();
  const { user } = useUserContext();
  const { setSnackbarMessage } = useUIContext();
  const subscriptionKeyRef = useRef<string | null>(null);
  const subscriptionControllerRef = useRef<{ cancel: () => void } | null>(null);

  const notify = (msgId: string, message: string, type: string) =>
    setSnackbarMessage({ id: msgId, message, type, show: true });

  const cleanupSubscription = useCallback(() => {
    if (subscriptionKeyRef.current) {
      PubSub.unsubscribe(subscriptionKeyRef.current);
      subscriptionKeyRef.current = null;
    }
    if (subscriptionControllerRef.current) {
      subscriptionControllerRef.current.cancel();
      subscriptionControllerRef.current = null;
    }
  }, []);

  useEffect(() => () => cleanupSubscription(), [cleanupSubscription]);

  const { data: nodesData, loading: nodesLoading } = useGetNodesQuery({
    skip: !id,
    fetchPolicy: 'network-only',
    variables: { data: { state: NodeStateEnum.Configured } },
    onError: (err) => notify('node-msg', err.message, 'error'),
  });

  const currentNode = nodesData?.getNodes.nodes.find((n) => n.id === id);

  const [updateNode, { loading: updateNodeLoading }] = useUpdateNodeMutation({
    onCompleted: () =>
      notify(
        'update-node-success-msg',
        'Node updated successfully.',
        'success',
      ),
    onError: (err) => notify('update-node-err-msg', err.message, 'error'),
    refetchQueries: ['GetNodes'],
  });

  const [restartNode] = useRestartNodeMutation({
    fetchPolicy: 'network-only',
    onCompleted: () => {
      setNodeAction((prev) => ({
        ...prev,
        progress: prev.progress + 25,
        currentAction: 'loading',
      }));
      notify('restart-node-success-msg', 'Node restart initiated.', 'success');
    },
    onError: () => {
      setNodeAction({
        progress: 0,
        currentAction: '',
        actionInitiated: '',
        action: '',
        isActive: false,
      });
      notify('restart-node-err-msg', "Couldn't restart node.", 'error');
    },
  });

  const { loading: appsLoading } = useGetAppsQuery({
    fetchPolicy: 'cache-and-network',
  });

  const {
    loading: softwaresLoading,
    data: softwaresData,
    refetch: refetchSoftwares,
  } = useSoftwareQuery({
    fetchPolicy: 'network-only',
    variables: {
      data: { name: '', nodeId: id, status: SoftwareStatusEnum.Unknown },
    },
  });

  const [updateSoftware, { loading: updateSoftwareLoading }] =
    useUpdateSoftwareMutation({
      fetchPolicy: 'network-only',
      onCompleted: () => {
        refetchSoftwares();
        notify(
          'update-software-success-msg',
          'Software updated successfully.',
          'success',
        );
      },
      onError: (err) =>
        notify('update-software-error-msg', err.message, 'error'),
    });

  const [
    getNodeMetricByTab,
    { loading: nodeMetricsLoading, variables: nodeMetricsVariables },
  ] = useGetMetricByTabLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => setMetrics(data.getMetricByTab),
  });

  const [
    getMetricStat,
    { data: statData, loading: statLoading, variables: statVar },
  ] = useGetMetricsStatLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      if (data.getMetricsStat.metrics.length > 0) {
        data.getMetricsStat.metrics.forEach((m) => {
          if (m.type === NODE_KPIS.NODE_UPTIME[nodeType][0].id) {
            setNodeUptime(m.value);
          }
        });

        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.AllNode}-${statVar?.data.from ?? 0}`;
        cleanupSubscription();
        subscriptionKeyRef.current = sKey;

        const controller = MetricStatSubscription({
          key: sKey,
          nodeId: id,
          userId: user.id,
          url: env.METRIC_URL,
          orgName: user.orgName,
          type: Stats_Type.AllNode,
          from: statVar?.data.from ?? 0,
        });

        subscriptionControllerRef.current = controller;
        PubSub.subscribe(sKey, handleStatSubscription);
      }
    },
  });

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

  useEffect(() => {
    const to = getUnixTime();
    const from = to - STAT_STEP_29;
    if (!id) {
      notify('node-not-found-msg', 'Node not found.', 'error');
      router.back();
    } else {
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
  }, [id, user.id, user.orgName, cleanupSubscription]);

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
  }, [nodesData]);

  useEffect(() => {
    if (metricFrom > 0 && nodeMetricsVariables?.data?.from !== metricFrom) {
      getNodeMetricByTab({
        variables: {
          data: {
            step: 30,
            nodeId: id,
            userId: user.id,
            type: graphType,
            from: metricFrom,
            orgName: user.orgName,
            withSubscription: false,
            to: metricFrom + METRIC_RANGE_10800,
          },
        },
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    metricFrom,
    nodeMetricsVariables?.data?.from,
    getNodeMetricByTab,
    graphType,
  ]);

  const handleStatSubscription = (_: unknown, data: string) => {
    const parsedData: TMetricResDto = JSON.parse(data);
    const { value, type, success } = parsedData.data.getMetricStatSub;
    if (success) {
      if (type === NODE_KPIS.NODE_UPTIME[nodeType][0].id) {
        setNodeUptime(Math.floor(value[1]));
      }
      PubSub.publish(`stat-${type}`, value);
    }
  };

  const handleNodeSelected = (obj: TStatusBarObj) => {
    router.push(`/console/nodes/${obj.id}`);
  };

  const handleEditNode = (name: string) => {
    setIsEditNode(false);
    updateNode({ variables: { data: { id: currentNode?.id ?? '', name } } });
  };

  const handleSectionChange = (type: Graphs_Type) => {
    setGraphType(type);
    setMetricFrom(() => getUnixTime() - METRIC_RANGE_10800);
  };

  const onTabSelected = (_: unknown, value: number) => {
    setSelectedTab(value);
    setGraphType(getNodeTabTypeByIndex(value));
    setMetricFrom(() => getUnixTime() - METRIC_RANGE_10800);
  };

  const handleNodeActionClick = (action: string, _: boolean) => {
    if (action === NODE_ACTIONS_ENUM.NODE_RESTART) {
      setNodeAction({
        progress: 0,
        currentAction: NODE_ACTIONS_ENUM.NODE_RESTART,
        actionInitiated: NODE_ACTIONS_ENUM.NODE_RESTART,
        action: NODE_ACTIONS_ENUM.NODE_RESTART,
        isActive: true,
      });
      restartNode({ variables: { data: { nodeId: currentNode?.id ?? '' } } });
    }
  };

  const handleUpdateAvailable = (
    name: string,
    desiredVersion: string,
    nodeId: string,
  ) => {
    updateSoftware({
      variables: { data: { name, nodeId, tag: desiredVersion } },
    });
  };

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
