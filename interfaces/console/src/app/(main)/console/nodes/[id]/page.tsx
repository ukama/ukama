/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { use } from 'react';
import { NodeConnectivityEnum, NodeTypeEnum } from '@/client/graphql/generated';
import EditNode from '@/app/(main)/console/nodes/[id]/_components/EditNode';
import { NodeActionUI } from '@/app/(main)/console/nodes/[id]/_components/NodeActionUI';
import NodeNetworkTab from '@/app/(main)/console/nodes/[id]/_components/NodeNetworkTab';
import NodeOverviewTab from '@/app/(main)/console/nodes/[id]/_components/NodeOverviewTab';
import NodeRadioTab from '@/app/(main)/console/nodes/[id]/_components/NodeRadioTab';
import NodeResourcesTab from '@/app/(main)/console/nodes/[id]/_components/NodeResourcesTab';
import NodeSoftwareTab from '@/app/(main)/console/nodes/[id]/_components/NodeSoftwareTab';
import { useNodeDetailPage } from '@/app/(main)/console/nodes/[id]/_hooks/useNodeDetailPage';
import StatusBar from '@/app/(main)/console/_components/StatusBar';
import LoadingWrapper from '@/components/ui/LoadingWrapper';
import TabPanel from '@/components/ui/TabPanel';
import { NodePageTabs } from '@/constants';
import { colors } from '@/theme';
import { Stack, Tab, Tabs } from '@mui/material';

interface INodePage {
  params: Promise<{ id: string }>;
}

const Page: React.FC<INodePage> = ({ params }) => {
  const { id } = use(params);
  const vm = useNodeDetailPage(id);

  return (
    <Stack width={'100%'} height={'100%'} mt={1} spacing={1}>
      <StatusBar
        type="split"
        uptime={vm.nodeUptime}
        selected={vm.currentNode}
        handleEditClick={vm.onOpenEditNode}
        handleSelected={vm.handleNodeSelected}
        objs={vm.nodesData?.getNodes.nodes ?? []}
        actionOptions={vm.actionOptions}
        handleActionClick={vm.handleNodeActionClick}
        loading={vm.nodesLoading || vm.updateNodeLoading || vm.statLoading}
      />

      {vm.currentNode?.status.connectivity === NodeConnectivityEnum.Online &&
      !vm.nodeAction.actionInitiated ? (
        <div>
          <Tabs value={vm.selectedTab} onChange={vm.onTabSelected} sx={{ pb: 2 }}>
            {NodePageTabs.map(({ id, label, value }) => (
              <Tab
                key={id}
                label={label}
                id={`node-tab-${value}`}
                sx={{
                  display:
                    ((vm.currentNode?.type === NodeTypeEnum.Cnode ||
                      vm.currentNode?.type === NodeTypeEnum.Hnode) &&
                      label === 'Radio') ||
                    ((vm.currentNode?.type === NodeTypeEnum.Anode ||
                      vm.currentNode?.type === NodeTypeEnum.Cnode ||
                      vm.currentNode?.type === NodeTypeEnum.Hnode) &&
                      label === 'Network')
                      ? 'none'
                      : 'block',
                }}
              />
            ))}
          </Tabs>

          <LoadingWrapper
            radius="small"
            width={'100%'}
            isLoading={vm.nodesLoading || vm.updateNodeLoading}
            cstyle={{ backgroundColor: false ? colors.white : 'transparent' }}
          >
            <TabPanel id={'node-overview-tab'} value={vm.selectedTab} index={0}>
              <NodeOverviewTab
                nodeId={vm.id}
                metrics={vm.metrics}
                connectedUsers={'0'}
                metricFrom={vm.metricFrom}
                statLoading={vm.statLoading}
                isUpdateAvailable={false}
                onNodeSelected={() => {}}
                handleUpdateNode={() => {}}
                selectedNode={vm.currentNode}
                metricsLoading={vm.nodeMetricsLoading}
                getNodeSoftwareUpdateInfos={() => {}}
                handleOverviewSectionChange={vm.handleSectionChange}
                nodeMetricsStatData={vm.statData?.getMetricsStat ?? { metrics: [] }}
              />
            </TabPanel>

            <TabPanel id={'node-network-tab'} value={vm.selectedTab} index={1}>
              <NodeNetworkTab
                metrics={vm.metrics}
                metricFrom={vm.metricFrom}
                statLoading={vm.statLoading}
                selectedNode={vm.currentNode}
                loading={vm.nodeMetricsLoading || vm.statLoading}
                handleSectionChange={vm.handleSectionChange}
                nodeMetricsStatData={vm.statData?.getMetricsStat ?? { metrics: [] }}
              />
            </TabPanel>

            <TabPanel id={'node-resources-tab'} value={vm.selectedTab} index={2}>
              <NodeResourcesTab
                metrics={vm.metrics}
                metricFrom={vm.metricFrom}
                statLoading={vm.statLoading}
                selectedNode={vm.currentNode}
                loading={vm.nodeMetricsLoading || vm.statLoading}
                nodeMetricsStatData={vm.statData?.getMetricsStat ?? { metrics: [] }}
              />
            </TabPanel>

            <TabPanel id={'node-radio-tab'} value={vm.selectedTab} index={3}>
              <NodeRadioTab
                metrics={vm.metrics}
                metricFrom={vm.metricFrom}
                statLoading={vm.statLoading}
                selectedNode={vm.currentNode}
                loading={vm.nodeMetricsLoading || vm.statLoading}
                nodeMetricsStatData={vm.statData?.getMetricsStat ?? { metrics: [] }}
              />
            </TabPanel>

            <TabPanel id={'node-software-tab'} value={vm.selectedTab} index={4}>
              <NodeSoftwareTab
                loading={vm.softwaresLoading || vm.updateSoftwareLoading || vm.appsLoading}
                nodeApps={vm.softwaresData?.getSoftwares.software ?? []}
                handleUpdateAvailable={vm.handleUpdateAvailable}
              />
            </TabPanel>
          </LoadingWrapper>
        </div>
      ) : (
        <NodeActionUI
          value={vm.nodeAction.progress}
          nodeType={vm.currentNode?.type}
          action={vm.nodeAction.actionInitiated || ''}
          connectivity={
            (vm.currentNode?.status?.connectivity as NodeConnectivityEnum) || undefined
          }
          description={vm.nodeActionDescription}
        />
      )}

      {vm.isEditNode && (
        <EditNode
          title="Edit Node"
          isOpen={vm.isEditNode}
          labelSuccessBtn="Save"
          labelNegativeBtn="Cancel"
          nodeName={vm.currentNode?.name ?? ''}
          handleSuccessAction={vm.handleEditNode}
          handleCloseAction={vm.onCloseEditNode}
        />
      )}
    </Stack>
  );
};

export default Page;

