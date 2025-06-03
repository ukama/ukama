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
  useGetNodesQuery,
  useRestartNodeMutation,
  useToggleRfStatusMutation,
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
import { NodeActionUI } from '@/components/NodeActionUI';
import NodeNetworkTab from '@/components/NodeNetworkTab';
import NodeOverviewTab from '@/components/NodeOverviewTab';
import NodeRadioTab from '@/components/NodeRadioTab';
import NodeResourcesTab from '@/components/NodeResourcesTab';
import NodeStatus from '@/components/NodeStatus';
import TabPanel from '@/components/TabPanel';
import {
  METRIC_RANGE_10800,
  NODE_ACTIONS_BUTTONS,
  NODE_ACTIONS_ENUM,
  NodePageTabs,
  STAT_STEP_29,
} from '@/constants';
import { useAppContext } from '@/context';
import MetricStatSubscription from '@/lib/MetricStatSubscription';
import { colors } from '@/theme';
import { TMetricResDto } from '@/types';
import {
  getNodeActionDescriptionByProgress,
  getNodeTabTypeByIndex,
  getUnixTime,
} from '@/utils';
import { Stack, Tab, Tabs } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useCallback, useEffect, useRef, useState } from 'react';

const NODE_UPTIME_KEY = 'unit_uptime';
interface INodePage {
  params: {
    id: string;
  };
}

const Page: React.FC<INodePage> = ({ params }) => {
  const { id } = params;
  const router = useRouter();
  const [metricFrom, setMetricFrom] = useState<number>(0);
  const [isEditNode, setIsEditNode] = useState<boolean>(false);
  const [nodeAction, setNodeAction] = useState({
    progress: 0,
    currentAction: '',
    actionInitiated: '',
  });
  const [graphType, setGraphType] = useState<Graphs_Type>(
    Graphs_Type.NodeHealth,
  );
  const [nodeUptime, setNodeUptime] = useState<number>(0);
  const [selectedTab, setSelectedTab] = useState<number>(0);
  const [metrics, setMetrics] = useState<MetricsRes>({ metrics: [] });
  const { user, setSnackbarMessage, env, subscriptionClient } = useAppContext();
  const subscriptionKeyRef = useRef<string | null>(null);
  const subscriptionControllerRef = useRef<AbortController | null>(null);

  const cleanupSubscription = useCallback(() => {
    if (subscriptionKeyRef.current) {
      PubSub.unsubscribe(subscriptionKeyRef.current);
      subscriptionKeyRef.current = null;
    }
    if (subscriptionControllerRef.current) {
      subscriptionControllerRef.current.abort();
      subscriptionControllerRef.current = null;
    }
  }, []);

  useEffect(() => {
    return () => {
      cleanupSubscription();
    };
  }, [cleanupSubscription]);

  const { data: nodesData, loading: nodesLoading } = useGetNodesQuery({
    skip: !id,
    fetchPolicy: 'network-only',
    variables: {
      data: {
        state: NodeStateEnum.Configured,
      },
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

  const currentNode = nodesData?.getNodes.nodes.find((n) => n.id === id);

  const [updateNode, { loading: updateNodeLoading }] = useUpdateNodeMutation({
    onCompleted: () => {
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
      setSnackbarMessage({
        id: 'restart-node-success-msg',
        message: 'Node restart initiated.',
        type: 'success',
        show: true,
      });
    },
    onError: () => {
      setNodeAction({
        progress: 0,
        currentAction: '',
        actionInitiated: '',
      });
      setSnackbarMessage({
        id: 'restart-node-err-msg',
        message: "Couldn't restart node.",
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

  const [
    getMetricStat,
    { data: statData, loading: statLoading, variables: statVar },
  ] = useGetMetricsStatLazyQuery({
    client: subscriptionClient,
    fetchPolicy: 'network-only',
    onCompleted: async (data) => {
      if (data.getMetricsStat.metrics.length > 0) {
        data.getMetricsStat.metrics.forEach((m) => {
          if (m.type === NODE_UPTIME_KEY) {
            setNodeUptime(m.value);
          }
        });

        const sKey = `stat-${user.orgName}-${user.id}-${Stats_Type.AllNode}-${statVar?.data.from ?? 0}`;
        cleanupSubscription();
        subscriptionKeyRef.current = sKey;

        const controller = await MetricStatSubscription({
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

  const [toggleRFStatus] = useToggleRfStatusMutation({
    fetchPolicy: 'network-only',
    onCompleted: (_, context) => {
      setSnackbarMessage({
        id: 'toggle-rf-status-success-msg',
        message: `RF status turned ${
          context?.variables?.data?.status ? 'On' : 'Off'
        } successfully.`,
        type: 'success',
        show: true,
      });
    },
    onError: (_, context) => {
      setSnackbarMessage({
        id: 'toggle-rf-status-error-msg',
        message: `Failed to turn RF status ${
          context?.variables?.data?.status ? 'On' : 'Off'
        }.`,
        type: 'error',
        show: true,
      });
    },
  });

  useEffect(() => {
    const to = getUnixTime();
    const from = to - STAT_STEP_29;
    if (!id) {
      setSnackbarMessage({
        id: 'node-not-found-msg',
        message: 'Node not found.',
        type: 'error',
        show: true,
      });
      router.back();
    } else if (id) {
      cleanupSubscription();

      getMetricStat({
        variables: {
          data: {
            to: to,
            nodeId: id,
            from: from,
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

    if (
      nodesData &&
      nodesData?.getNodes.nodes.length > 0 &&
      nodeAction.actionInitiated !== ''
    ) {
      const s = nodesData?.getNodes.nodes.find((n) => n.id === id);
      if (s && s?.status.connectivity !== nodeAction.currentAction) {
        if (nodeAction.actionInitiated === NODE_ACTIONS_ENUM.NODE_RESTART) {
          switch (s.status.connectivity) {
            case NodeConnectivityEnum.Offline:
              setNodeAction((prev) => ({
                ...prev,
                progress: prev.progress + 25,
                currentAction: s?.status.connectivity.toString() || '',
              }));
              return;
            case NodeConnectivityEnum.Online:
              setNodeAction((prev) => ({
                ...prev,
                progress: prev.progress + 25,
                currentAction: s?.status.connectivity.toString() || '',
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

  const handleStatSubscription = (_: any, data: string) => {
    const parsedData: TMetricResDto = JSON.parse(data);
    const { value, type, success } = parsedData.data.getMetricStatSub;
    if (success) {
      if (type === NODE_UPTIME_KEY) {
        setNodeUptime(Math.floor(value[1]));
      }
      PubSub.publish(`stat-${type}`, value);
    }
  };

  const handleNodeSelected = (node: Node) => {
    router.push(`/console/nodes/${node.id}`);
  };

  const handleEditNode = (str: string) => {
    setIsEditNode(false);
    updateNode({
      variables: {
        data: {
          id: currentNode?.id ?? '',
          name: str,
        },
      },
    });
  };

  const handleOverviewSectionChange = (type: Graphs_Type) => {
    setGraphType(type);
    setMetricFrom(() => getUnixTime() - METRIC_RANGE_10800);
  };

  const handleNetworkSectionChange = (type: Graphs_Type) => {
    setGraphType(type);
    setMetricFrom(() => getUnixTime() - METRIC_RANGE_10800);
  };

  const onTabSelected = (_: any, value: number) => {
    setSelectedTab(value);
    setGraphType(getNodeTabTypeByIndex(value));
    setMetricFrom(() => getUnixTime() - METRIC_RANGE_10800);
  };

  const handleNodeActionClick = (action: string) => {
    switch (action) {
      case NODE_ACTIONS_ENUM.NODE_RESTART:
        setNodeAction({
          progress: 0,
          currentAction: NODE_ACTIONS_ENUM.NODE_RESTART,
          actionInitiated: NODE_ACTIONS_ENUM.NODE_RESTART,
        });
        restartNode({
          variables: {
            data: {
              nodeId: currentNode?.id ?? '',
            },
          },
        });
        break;
      case NODE_ACTIONS_ENUM.NODE_RF_ON:
      case NODE_ACTIONS_ENUM.NODE_RF_OFF:
        if (currentNode?.id) {
          toggleRFStatus({
            variables: {
              data: {
                nodeId: currentNode.id,
                status: action === NODE_ACTIONS_ENUM.NODE_RF_ON,
              },
            },
          });
        }
        break;
      default:
        return;
    }
  };

  return (
    <Stack width={'100%'} height={'100%'} mt={1} spacing={1}>
      <NodeStatus
        uptime={nodeUptime}
        onAddNode={() => {}}
        selectedNode={currentNode}
        handleEditNodeClick={() => {
          setIsEditNode(true);
        }}
        handleNodeSelected={handleNodeSelected}
        nodes={nodesData?.getNodes.nodes ?? []}
        nodeActionOptions={NODE_ACTIONS_BUTTONS}
        handleNodeActionClick={handleNodeActionClick}
        isShowNodeAction={
          currentNode?.status.connectivity === NodeConnectivityEnum.Online &&
          !nodeAction.actionInitiated
        }
        loading={nodesLoading || updateNodeLoading || statLoading}
      />
      {currentNode?.status.connectivity === NodeConnectivityEnum.Online &&
      !nodeAction.actionInitiated ? (
        <div>
          <Tabs value={selectedTab} onChange={onTabSelected} sx={{ pb: 2 }}>
            {NodePageTabs.map(({ id, label, value }) => (
              <Tab
                key={id}
                label={label}
                id={`node-tab-${value}`}
                sx={{
                  display:
                    ((currentNode?.type === NodeTypeEnum.Hnode &&
                      label === 'Radio') ??
                    (currentNode?.type === NodeTypeEnum.Anode &&
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
                metrics={metrics}
                connectedUsers={'0'}
                metricFrom={metricFrom}
                statLoading={statLoading}
                isUpdateAvailable={false}
                onNodeSelected={() => {}}
                handleUpdateNode={() => {}}
                selectedNode={currentNode}
                metricsLoading={nodeMetricsLoading}
                getNodeSoftwareUpdateInfos={() => {}}
                handleOverviewSectionChange={handleOverviewSectionChange}
                nodeMetricsStatData={
                  statData?.getMetricsStat ?? { metrics: [] }
                }
              />
            </TabPanel>
            <TabPanel id={'node-network-tab'} value={selectedTab} index={1}>
              <NodeNetworkTab
                metrics={metrics}
                metricFrom={metricFrom}
                statLoading={statLoading}
                selectedNode={currentNode}
                loading={nodeMetricsLoading || statLoading}
                handleSectionChange={handleNetworkSectionChange}
                nodeMetricsStatData={
                  statData?.getMetricsStat ?? { metrics: [] }
                }
              />
            </TabPanel>
            <TabPanel id={'node-resources-tab'} value={selectedTab} index={2}>
              <NodeResourcesTab
                metrics={metrics}
                metricFrom={metricFrom}
                statLoading={statLoading}
                selectedNode={currentNode}
                loading={nodeMetricsLoading || statLoading}
                nodeMetricsStatData={
                  statData?.getMetricsStat ?? { metrics: [] }
                }
              />
            </TabPanel>
            <TabPanel id={'node-radio-tab'} value={selectedTab} index={3}>
              <NodeRadioTab
                metrics={metrics}
                metricFrom={metricFrom}
                statLoading={statLoading}
                selectedNode={currentNode}
                loading={nodeMetricsLoading || statLoading}
                nodeMetricsStatData={
                  statData?.getMetricsStat ?? { metrics: [] }
                }
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
        </div>
      ) : (
        <NodeActionUI
          value={nodeAction.progress}
          nodeType={currentNode?.type}
          action={nodeAction.actionInitiated || NODE_ACTIONS_ENUM.NODE_OFF}
          connectivity={
            (currentNode?.status?.connectivity as NodeConnectivityEnum) ||
            undefined
          }
          description={getNodeActionDescriptionByProgress(
            nodeAction.progress,
            nodeAction.actionInitiated ||
              currentNode?.status?.connectivity?.toString() ||
              '',
          )}
        />
      )}
      {isEditNode && (
        <EditNode
          title="Edit Node"
          isOpen={isEditNode}
          labelSuccessBtn="Save"
          labelNegativeBtn="Cancel"
          nodeName={currentNode?.name ?? ''}
          handleSuccessAction={handleEditNode}
          handleCloseAction={() => setIsEditNode(false)}
        />
      )}
    </Stack>
  );
};

export default Page;
