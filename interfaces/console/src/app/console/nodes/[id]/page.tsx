/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import {
  Node,
  NodeConnectivityEnum,
  NodeStateEnum,
  NodeTypeEnum,
  useGetNodesByStateQuery,
  useRestartNodeMutation,
  useUpdateNodeMutation,
} from '@/client/graphql/generated';
import {
  Graphs_Type,
  MetricsRes,
  Stats_Type,
  useGetMetricByTabLazyQuery,
  useGetMetricsStatLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import EditNode from '@/components/EditNode';
import LoadingWrapper from '@/components/LoadingWrapper';
import NodeNetworkTab from '@/components/NodeNetworkTab';
import NodeOverviewTab from '@/components/NodeOverviewTab';
import NodeRadioTab from '@/components/NodeRadioTab';
import NodeResourcesTab from '@/components/NodeResourcesTab';
import NodeStatus from '@/components/NodeStatus';
import TabPanel from '@/components/TabPanel';
import {
  METRIC_RANGE_3600,
  NODE_ACTIONS_BUTTONS,
  NodePageTabs,
} from '@/constants';
import { useAppContext } from '@/context';
import MetricSubscription from '@/lib/MetricSubscription';
import { colors } from '@/theme';
import { TMetricResDto } from '@/types';
import { getNodeTabTypeByIndex, getUnixTime } from '@/utils';
import { Stack, Tab, Tabs } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

const NODE_UPTIME_KEY = 'unit_uptime';
interface INodePage {
  params: {
    id: string;
  };
}

const Page: React.FC<INodePage> = ({ params }) => {
  const { id } = params;
  const router = useRouter();
  const [isEditNode, setIsEditNode] = useState<boolean>(false);
  const [metricFrom, setMetricFrom] = useState<number>(0);
  const [graphType, setGraphType] = useState<Graphs_Type>(
    Graphs_Type.NodeHealth,
  );
  const [nodeUptime, setNodeUptime] = useState<number>(0);
  const [metrics, setMetrics] = useState<MetricsRes>({ metrics: [] });
  const [selectedTab, setSelectedTab] = useState<number>(0);
  const { user, setSnackbarMessage, env, subscriptionClient } = useAppContext();
  const [selectedNode, setSelectedNode] = useState<Node | undefined>(undefined);

  const { data: nodesData, loading: nodesLoading } = useGetNodesByStateQuery({
    skip: !id,
    fetchPolicy: 'cache-and-network',
    variables: {
      data: {
        connectivity: NodeConnectivityEnum.Online,
        state: NodeStateEnum.Configured,
      },
    },
    onCompleted: (data) => {
      if (data.getNodesByState.nodes.length > 0) {
        const node =
          data.getNodesByState.nodes.find((n) => n.id === id) ?? undefined;
        setSelectedNode(node);
      }
    },
    onError: (err) => {
      setSnackbarMessage({
        id: 'node-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [updateNode, { loading: updateNodeLoading }] = useUpdateNodeMutation({
    onCompleted: (data) => {
      setSelectedNode(data.updateNode);
      setSnackbarMessage({
        id: 'update-node-success-msg',
        message: 'Node updated successfully.',
        type: 'success',
        show: true,
      });
    },
    onError: (err) => {
      setSnackbarMessage({
        id: 'update-node-err-msg',
        message: err.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [restartNode, { loading: restartNodeLoading }] = useRestartNodeMutation(
    {
      fetchPolicy: 'network-only',
      onCompleted: () => {
        setSnackbarMessage({
          id: 'restart-node-success-msg',
          message: 'Node restart initiated.',
          type: 'success',
          show: true,
        });
      },
      onError: (err) => {
        setSnackbarMessage({
          id: 'restart-node-err-msg',
          message: "Couldn't restart node.",
          type: 'error',
          show: true,
        });
      },
    },
  );

  const [
    getNodeMetricByTab,
    { loading: nodeMetricsLoading, variables: nodeMetricsVariables },
  ] = useGetMetricByTabLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      setMetrics(data.getMetricByTab);
    },
  });

  const [
    getMetricsStat,
    { loading: nodeMetricsStatLoading, data: nodeMetricsStatData },
  ] = useGetMetricsStatLazyQuery({
    client: subscriptionClient,
    onCompleted: (data) => {
      if (data.getMetricsStat.metrics?.length > 0) {
        const uptimeStat = data?.getMetricsStat.metrics.find(
          (m) => m.type === NODE_UPTIME_KEY,
        );
        setNodeUptime(uptimeStat?.value ?? 0);
      }
    },
  });

  useEffect(() => {
    if (!id) {
      setSnackbarMessage({
        id: 'node-not-found-msg',
        message: 'Node not found.',
        type: 'error',
        show: true,
      });
      router.back();
    } else {
      const now = getUnixTime();
      getMetricsStat({
        variables: {
          data: {
            to: now,
            nodeId: id,
            userId: user.id,
            orgName: user.orgName,
            type: Stats_Type.Overview,
            from: now - METRIC_RANGE_3600,
          },
        },
      });
    }
  }, [id]);

  useEffect(() => {
    if (metricFrom > 0 && nodeMetricsVariables?.data?.from !== metricFrom) {
      getMetricsStat({
        variables: {
          data: {
            nodeId: id,
            userId: user.id,
            from: metricFrom,
            orgName: user.orgName,
            type:
              graphType === Graphs_Type.NodeHealth ||
              graphType === Graphs_Type.Subscribers
                ? Stats_Type.Overview
                : graphType === Graphs_Type.NetworkBackhaul ||
                    graphType === Graphs_Type.NetworkCellular
                  ? Stats_Type.Network
                  : graphType === Graphs_Type.Resources
                    ? Stats_Type.Resources
                    : Stats_Type.Radio,
            to: metricFrom + METRIC_RANGE_3600,
          },
        },
      });
      const psKey = `metric-${user.orgName}-${user.id}-${graphType}-${metricFrom}`;
      getNodeMetricByTab({
        variables: {
          data: {
            nodeId: id,
            userId: user.id,
            type: graphType,
            from: metricFrom,
            to: metricFrom + METRIC_RANGE_3600,
            orgName: user.orgName,
            withSubscription: true,
          },
        },
      }).then(() => {
        MetricSubscription({
          nodeId: id,
          key: psKey,
          type: graphType,
          userId: user.id,
          from: metricFrom,
          url: env.METRIC_URL,
          orgName: user.orgName,
        });
      });

      PubSub.subscribe(psKey, handleNotification);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    metricFrom,
    nodeMetricsVariables?.data?.from,
    getNodeMetricByTab,
    graphType,
  ]); // Added all missing dependencies

  const handleNotification = (_: any, data: string) => {
    const parsedData: TMetricResDto = JSON.parse(data);
    const { msg, type, value, nodeId, success } =
      parsedData.data.getMetricByTabSub;
    if (success) {
      if (type === NODE_UPTIME_KEY) {
        setNodeUptime(value[value.length - 1][1]);
      }
      PubSub.publish(type, value.slice(0, 30));
    }
  };

  const handleNodeSelected = (node: Node) => {
    setSelectedNode(node);
    router.push(`/console/nodes/${node.id}`);
  };

  const handleEditNode = (str: string) => {
    setIsEditNode(false);
    updateNode({
      variables: {
        data: {
          id: selectedNode?.id ?? '',
          name: str,
        },
      },
    });
  };

  const handleOverviewSectionChange = (type: Graphs_Type) => {
    setGraphType(type);
    setMetricFrom(() => getUnixTime() - METRIC_RANGE_3600);
  };

  const handleNetworkSectionChange = (type: Graphs_Type) => {
    setGraphType(type);
    setMetricFrom(() => getUnixTime() - METRIC_RANGE_3600);
  };

  const onTabSelected = (_: any, value: number) => {
    setSelectedTab(value);
    setGraphType(getNodeTabTypeByIndex(value));
    setMetricFrom(() => getUnixTime() - METRIC_RANGE_3600);
  };

  const handleNodeActionClick = (action: string) => {
    switch (action) {
      case 'node-restart':
        restartNode({
          variables: {
            data: {
              nodeId: selectedNode?.id ?? '',
            },
          },
        });
        break;
      default:
        return;
    }
  };

  return (
    <Stack width={'100%'} mt={1} spacing={1}>
      <NodeStatus
        uptime={nodeUptime}
        onAddNode={() => {}}
        loading={nodesLoading || updateNodeLoading}
        selectedNode={selectedNode}
        handleEditNodeClick={() => {
          setIsEditNode(true);
        }}
        handleNodeSelected={handleNodeSelected}
        nodeActionOptions={NODE_ACTIONS_BUTTONS}
        handleNodeActionClick={handleNodeActionClick}
        nodes={nodesData?.getNodesByState.nodes ?? []}
      />

      <Tabs value={selectedTab} onChange={onTabSelected} sx={{ pb: 2 }}>
        {NodePageTabs.map(({ id, label, value }) => (
          <Tab
            key={id}
            label={label}
            id={`node-tab-${value}`}
            sx={{
              display:
                ((selectedNode?.type === NodeTypeEnum.Hnode &&
                  label === 'Radio') ??
                (selectedNode?.type === NodeTypeEnum.Anode &&
                  label === 'Network'))
                  ? 'none'
                  : 'block',
            }}
          />
        ))}
      </Tabs>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={nodesLoading || updateNodeLoading}
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <TabPanel id={'node-overview-tab'} value={selectedTab} index={0}>
          <NodeOverviewTab
            nodeId={id}
            loading={false}
            metrics={metrics}
            connectedUsers={'0'}
            metricFrom={metricFrom}
            isUpdateAvailable={false}
            onNodeSelected={() => {}}
            handleUpdateNode={() => {}}
            selectedNode={selectedNode}
            getNodeSoftwareUpdateInfos={() => {}}
            handleOverviewSectionChange={handleOverviewSectionChange}
            metricsLoading={nodeMetricsLoading || nodeMetricsStatLoading}
            nodeMetricsStatData={
              nodeMetricsStatData?.getMetricsStat ?? { metrics: [] }
            }
          />
        </TabPanel>
        <TabPanel id={'node-network-tab'} value={selectedTab} index={1}>
          <NodeNetworkTab
            metrics={metrics}
            metricFrom={metricFrom}
            selectedNode={selectedNode}
            nodeMetricsStatData={
              nodeMetricsStatData?.getMetricsStat ?? { metrics: [] }
            }
            handleSectionChange={handleNetworkSectionChange}
            loading={nodeMetricsLoading || nodeMetricsStatLoading}
          />
        </TabPanel>
        <TabPanel id={'node-resources-tab'} value={selectedTab} index={2}>
          <NodeResourcesTab
            metrics={metrics}
            metricFrom={metricFrom}
            selectedNode={selectedNode}
            nodeMetricsStatData={
              nodeMetricsStatData?.getMetricsStat ?? { metrics: [] }
            }
            loading={nodeMetricsLoading || nodeMetricsStatLoading}
          />
        </TabPanel>
        <TabPanel id={'node-radio-tab'} value={selectedTab} index={3}>
          <NodeRadioTab
            metrics={metrics}
            metricFrom={metricFrom}
            selectedNode={selectedNode}
            nodeMetricsStatData={
              nodeMetricsStatData?.getMetricsStat ?? { metrics: [] }
            }
            loading={nodeMetricsLoading || nodeMetricsStatLoading}
          />
        </TabPanel>
        {/* <TabPanel id={'node-software-tab'} value={selectedTab} index={4}>
          <NodeSoftwareTab
            loading={nodeAppsLoading}
            nodeApps={nodeAppsRes?.getNodeApps.apps ?? []}
          />
        </TabPanel> */}
        {/* <TabPanel id={'node-schematic-tab'} value={selectedTab} index={5}>
          <NodeSchematicTab
            getSearchValue={() => {}}
            schematicsSpecsData={SPEC_DATA}
            nodeTitle={selectedNode?.name ?? 'Node'}
            loading={false}
          />
        </TabPanel> */}
      </LoadingWrapper>
      {isEditNode && (
        <EditNode
          title="Edit Node"
          isOpen={isEditNode}
          labelSuccessBtn="Save"
          labelNegativeBtn="Cancel"
          nodeName={selectedNode?.name ?? ''}
          handleSuccessAction={handleEditNode}
          handleCloseAction={() => setIsEditNode(false)}
        />
      )}
    </Stack>
  );
};

export default Page;
