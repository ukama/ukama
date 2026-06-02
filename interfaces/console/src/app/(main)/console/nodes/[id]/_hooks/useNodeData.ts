/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import {
  NodeStateEnum,
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
  useGetMetricByTabLazyQuery,
  useGetMetricsStatLazyQuery,
  GetMetricsStatQuery,
} from '@/client/graphql/generated/subscriptions';
import { METRIC_RANGE_10800 } from '@/constants';
import { useUIContext, useUserContext } from '@/context';
import { TNodeActionState } from '@/types';
import { ApolloClient, NormalizedCacheObject } from '@apollo/client';
import { useMemo } from 'react';

interface UseNodeDataParams {
  id: string;
  graphType: Graphs_Type;
  metricFrom: number;
  subscriptionClient: ApolloClient<NormalizedCacheObject>;
  onMetricsFetched: (data: MetricsRes) => void;
  onStatFetched: (data: ReturnType<typeof useGetMetricsStatLazyQuery>[1]['data']) => void;
  setNodeUptime: (uptime: number) => void;
  setNodeAction: React.Dispatch<React.SetStateAction<TNodeActionState & { currentAction: string; actionInitiated: string }>>;
}

export function useNodeData({
  id,
  graphType,
  metricFrom,
  subscriptionClient,
  onMetricsFetched,
  onStatFetched,
  setNodeUptime: _setNodeUptime,
  setNodeAction,
}: UseNodeDataParams) {
  const { setSnackbarMessage } = useUIContext();
  const { user } = useUserContext();

  const notify = (msgId: string, message: string, type: 'success' | 'error' | 'warning' | 'info') =>
    setSnackbarMessage({ id: msgId, message, type, show: true });

  const { data: nodesData, loading: nodesLoading } = useGetNodesQuery({
    skip: !id,
    fetchPolicy: 'cache-and-network',
    variables: { data: { state: NodeStateEnum.Configured } },
    onError: (err) => notify('node-msg', err.message, 'error'),
  });

  const currentNode = useMemo(
    () => nodesData?.getNodes.nodes.find((n) => n.id === id),
    [nodesData, id],
  );

  const [updateNode, { loading: updateNodeLoading }] = useUpdateNodeMutation({
    onCompleted: () =>
      notify('update-node-success-msg', 'Node updated successfully.', 'success'),
    onError: (err) => notify('update-node-err-msg', err.message, 'error'),
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
    fetchPolicy: 'cache-and-network',
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
    onCompleted: (data) => onMetricsFetched(data.getMetricByTab),
  });

  const [
    getMetricStat,
    { data: statData, loading: statLoading, variables: statVar },
  ] = useGetMetricsStatLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      if (onStatFetched) {
        onStatFetched(data as GetMetricsStatQuery);
      }
    },
  });

  const fetchMetricByTab = () => {
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

  const handleEditNode = (name: string) => {
    updateNode({ variables: { data: { id: currentNode?.id ?? '', name } } });
  };

  const handleRestartNode = () => {
    restartNode({ variables: { data: { nodeId: currentNode?.id ?? '' } } });
  };

  return {
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
    handleEditNode,
    handleRestartNode,
  };
}
