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
  useGetNodesByStateQuery,
  useUpdateNodeMutation,
} from '@/client/graphql/generated';
import {
  Graphs_Type,
  MetricsRes,
  useGetMetricByTabLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import EditNode from '@/components/EditNode';
import LoadingWrapper from '@/components/LoadingWrapper';
import NodeNetworkTab from '@/components/NodeNetworkTab';
import NodeOverviewTab from '@/components/NodeOverviewTab';
import NodeRadioTab from '@/components/NodeRadioTab';
import NodeResourcesTab from '@/components/NodeResourcesTab';
import NodeStatus from '@/components/NodeStatus';
import TabPanel from '@/components/TabPanel';
import { NODE_ACTIONS_BUTTONS, NodePageTabs } from '@/constants';
import { useAppContext } from '@/context';
import MetricSubscription from '@/lib/MetricSubscription';
import { colors } from '@/theme';
import { TMetricResDto } from '@/types';
import { getNodeTabTypeByIndex, getUnixTime } from '@/utils';
import { Stack, Tab, Tabs } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

const METRIC_RANGE_3600 = 3600;

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
  const [metrics, setMetrics] = useState<MetricsRes>({ metrics: [] });
  const [selectedTab, setSelectedTab] = useState<number>(0);
  const { user, setSnackbarMessage, env, subscriptionClient } = useAppContext();
  const [selectedNode, setSelectedNode] = useState<Node | undefined>(undefined);

  useEffect(() => {
    if (!id) {
      setSnackbarMessage({
        id: 'node-not-found-msg',
        message: 'Node not found.',
        type: 'error',
        show: true,
      });
      router.back();
    }
  }, []);

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

  useEffect(() => {
    if (metricFrom > 0 && nodeMetricsVariables?.data?.from !== metricFrom) {
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

  const onTabSelected = (_: any, value: number) => {
    setSelectedTab(value);
    setGraphType(getNodeTabTypeByIndex(value));
    setMetricFrom(() => getUnixTime() - METRIC_RANGE_3600);
  };

  const handleNodeActionClick = (action: string) => {};

  return (
    <Stack width={'100%'} mt={1} spacing={1}>
      <NodeStatus
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
            sx={
              {
                // display:
                //   ((selectedNode?.type === NodeTypeEnum.Hnode &&
                //     label === 'Radio') ??
                //   (selectedNode?.type === NodeTypeEnum.Anode &&
                //     label === 'Network'))
                //     ? 'none'
                //     : 'block',
              }
            }
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
            metricsLoading={nodeMetricsLoading}
            getNodeSoftwareUpdateInfos={() => {}}
            handleOverviewSectionChange={handleOverviewSectionChange}
          />
        </TabPanel>
        <TabPanel id={'node-network-tab'} value={selectedTab} index={1}>
          <NodeNetworkTab
            nodeId={id}
            metrics={metrics}
            metricFrom={metricFrom}
            loading={nodeMetricsLoading}
          />
        </TabPanel>
        <TabPanel id={'node-resources-tab'} value={selectedTab} index={2}>
          <NodeResourcesTab
            nodeId={id}
            metrics={metrics}
            metricFrom={metricFrom}
            selectedNode={selectedNode}
            loading={nodeMetricsLoading}
          />
        </TabPanel>
        <TabPanel id={'node-radio-tab'} value={selectedTab} index={3}>
          <NodeRadioTab
            nodeId={id}
            metrics={metrics}
            metricFrom={metricFrom}
            loading={nodeMetricsLoading}
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
