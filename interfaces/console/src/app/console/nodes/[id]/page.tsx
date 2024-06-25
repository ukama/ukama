/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { metricsClient } from '@/client/client';
import {
  Node,
  NodeTypeEnum,
  useGetNodeAppsLazyQuery,
  useGetNodeQuery,
  useGetNodesLazyQuery,
  useUpdateNodeMutation,
} from '@/client/graphql/generated';
import {
  Graphs_Type,
  MetricsRes,
  useGetMetricByTabLazyQuery,
} from '@/client/graphql/generated/metrics';
import EditNode from '@/components/EditNode';
import LoadingWrapper from '@/components/LoadingWrapper';
import NodeNetworkTab from '@/components/NodeNetworkTab';
import NodeOverviewTab from '@/components/NodeOverviewTab';
import NodeRadioTab from '@/components/NodeRadioTab';
import NodeResourcesTab from '@/components/NodeResourcesTab';
import NodeSchematicTab from '@/components/NodeSchematicTab';
import NodeSoftwareTab from '@/components/NodeSoftwareTab';
import NodeStatus from '@/components/NodeStatus';
import TabPanel from '@/components/TabPanel';
import { NODE_ACTIONS_BUTTONS, NodePageTabs } from '@/constants';
import { useAppContext } from '@/context';
import colors from '@/theme/colors';
import { getNodeTabTypeByIndex, getUnixTime } from '@/utils';
import { Stack, Tab, Tabs } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

const SPEC_DATA = [
  { id: 'pdf-1', title: 'PDF with Technical Specs', readingTime: '2mint' },
  { id: 'pdf-2', title: 'PDF with Technical Specs', readingTime: '2mint' },
  { id: 'pdf-3', title: 'PDF with Technical Specs', readingTime: '2mint' },
  { id: 'pdf-3', title: 'PDF with Technical Specs', readingTime: '2mint' },
];

export default function Page({ params }: { params: { id: string } }) {
  const router = useRouter();
  const [isEditNode, setIsEditNode] = useState<boolean>(false);
  const [metricFrom, setMetricFrom] = useState<number>(0);
  const [graphType, setGraphType] = useState<Graphs_Type>(
    Graphs_Type.NodeHealth,
  );
  const [metrics, setMetrics] = useState<MetricsRes>({ metrics: [] });
  const [selectedTab, setSelectedTab] = useState<number>(0);
  const [selectedNode, setSelectedNode] = useState<Node | undefined>(undefined);
  const { setSnackbarMessage } = useAppContext();

  const [
    getNodes,
    { data: getNodesData, loading: getNodesLoading, refetch: refetchNodes },
  ] = useGetNodesLazyQuery({
    fetchPolicy: 'cache-first',
  });

  const { data: getNodeData, loading: getNodeLoading } = useGetNodeQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      data: {
        id: params.id as string,
      },
    },
    onCompleted: (data) => {
      setSelectedNode(data.getNode);
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

  const [
    getNodeMetricByTab,
    {
      data: nodeMetricsData,
      loading: nodeMetricsLoading,
      variables: nodeMetricsVariables,
    },
  ] = useGetMetricByTabLazyQuery({
    client: metricsClient,
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      setMetrics(data.getMetricByTab);
    },
  });

  const [updateNode, { loading: updateNodeLoading }] = useUpdateNodeMutation({
    onCompleted: (data) => {
      setSelectedNode(data.updateNode);
      refetchNodes();
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

  const [getApps, { data: nodeAppsRes, loading: nodeAppsLoading }] =
    useGetNodeAppsLazyQuery({
      fetchPolicy: 'cache-and-network',
      onError: (err) => {
        setSnackbarMessage({
          id: 'node-apps-err-msg',
          message: err.message,
          type: 'error',
          show: true,
        });
      },
    });

  useEffect(() => {
    getNodes({
      variables: {
        data: {
          isFree: false,
        },
      },
    });
  }, []);

  useEffect(() => {
    if (selectedTab === 4) {
      getApps({
        variables: {
          data: {
            type: NodeTypeEnum.Hnode,
          },
        },
      });
    }
  }, [selectedTab]);

  useEffect(() => {
    if (metricFrom > 0 && nodeMetricsVariables?.data?.from !== metricFrom) {
      getNodeMetricByTab({
        variables: {
          data: {
            orgId: 'ukama',
            userId: 'salman',
            from: metricFrom,
            type: graphType,
            to: metricFrom + 120,
            withSubscription: true,
            nodeId: 'uk-test36-hnode-a1-00ff',
          },
        },
      });
    }
  }, [metricFrom]);

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
    setMetricFrom(() => getUnixTime() - 120);
  };

  const onTabSelected = (_: any, value: number) => {
    setSelectedTab(value);
    setGraphType(getNodeTabTypeByIndex(value));
    setMetricFrom(() => getUnixTime() - 120);
  };

  return (
    <Stack width={'100%'} mt={1} spacing={1}>
      <NodeStatus
        nodes={getNodesData?.getNodes.nodes ?? []}
        loading={getNodeLoading}
        onAddNode={() => {}}
        selectedNode={selectedNode}
        handleNodeActionClick={() => {}}
        handleNodeSelected={handleNodeSelected}
        handleNodeActionItemSelected={() => {}}
        nodeActionOptions={NODE_ACTIONS_BUTTONS}
        handleEditNodeClick={() => setIsEditNode(true)}
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
                //   (selectedNode?.type === 'HOME' && label === 'Radio') ??
                //   (selectedNode?.type === 'AMPLIFIER' && label === 'Network')
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
        isLoading={false}
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <TabPanel id={'node-overview-tab'} value={selectedTab} index={0}>
          <NodeOverviewTab
            metrics={metrics}
            metricFrom={metricFrom}
            isUpdateAvailable={true}
            selectedNode={selectedNode}
            metricsLoading={nodeMetricsLoading}
            handleOverviewSectionChange={handleOverviewSectionChange}
            handleUpdateNode={() => {}}
            connectedUsers={'0'}
            onNodeSelected={() => {}}
            uptime={0}
            getNodeSoftwareUpdateInfos={() => {}}
            loading={false}
          />
        </TabPanel>
        <TabPanel id={'node-network-tab'} value={selectedTab} index={1}>
          <NodeNetworkTab
            metrics={metrics}
            metricFrom={metricFrom}
            loading={nodeMetricsLoading}
          />
        </TabPanel>
        <TabPanel id={'node-resources-tab'} value={selectedTab} index={2}>
          <NodeResourcesTab
            metrics={metrics}
            metricFrom={metricFrom}
            selectedNode={selectedNode}
            loading={nodeMetricsLoading}
          />
        </TabPanel>
        <TabPanel id={'node-radio-tab'} value={selectedTab} index={3}>
          <NodeRadioTab
            metrics={metrics}
            metricFrom={metricFrom}
            loading={nodeMetricsLoading}
          />
        </TabPanel>
        <TabPanel id={'node-software-tab'} value={selectedTab} index={4}>
          <NodeSoftwareTab
            loading={nodeAppsLoading}
            nodeApps={nodeAppsRes?.getNodeApps.apps ?? []}
          />
        </TabPanel>
        <TabPanel id={'node-schematic-tab'} value={selectedTab} index={5}>
          <NodeSchematicTab
            getSearchValue={() => {}}
            schematicsSpecsData={SPEC_DATA}
            nodeTitle={selectedNode?.name ?? 'Node'}
            loading={false}
          />
        </TabPanel>
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
}
