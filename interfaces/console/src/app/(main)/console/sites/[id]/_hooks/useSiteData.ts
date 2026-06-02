/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import {
  GetSitesDocument,
  GetSitesQuery,
  Node,
  NodeStateEnum,
  NodeTypeEnum,
  SiteDto,
  Timeframe_Filter,
  useGetHealthReportQuery,
  useGetNodesForSiteLazyQuery,
  useToggleInternetSwitchMutation,
  useToggleRfStatusMutation,
  useToggleServiceMutation,
} from '@/client/graphql/generated';
import { NODE_ACTIONS_ENUM } from '@/constants';
import { useUIContext } from '@/context';
import { TSiteActionToggle } from '@/types';
import { stringToBoolean } from '@/utils';
import { useApolloClient } from '@apollo/client';
import { AlertColor } from '@mui/material';
import { Dispatch, SetStateAction, useCallback } from 'react';

export function useSiteData(
  id: string,
  activeSite: SiteDto,
  nodes: Node[],
  setNodes: Dispatch<SetStateAction<Node[]>>,
  setNodesFetched: Dispatch<SetStateAction<boolean>>,
  setSiteActionData: Dispatch<SetStateAction<TSiteActionToggle[]>>,
) {
  const { setSnackbarMessage } = useUIContext();
  const apolloClient = useApolloClient();

  const notify = useCallback(
    (msgId: string, message: string, type: 'success' | 'error' | 'warning' | 'info') =>
      setSnackbarMessage({ id: msgId, message, type, show: true }),
    [setSnackbarMessage],
  );

  // Read sites from Apollo cache — the parent page already fetched them
  const siteData = apolloClient.readQuery<GetSitesQuery>({
    query: GetSitesDocument,
    variables: { data: {} },
  });

  const [fetchNodesForSite] = useGetNodesForSiteLazyQuery({
    fetchPolicy: 'cache-first',
    onCompleted: (data) => {
      const filtered = data.getNodesForSite.nodes.filter(
        (node) =>
          node.site.siteId === activeSite.id &&
          node.status.state === NodeStateEnum.Configured,
      ) as Node[];
      setNodes(filtered);
      setNodesFetched(true);
    },
  });

  const [toggleRFStatus, { loading: toggleRFStatusLoading }] =
    useToggleRfStatusMutation({
      fetchPolicy: 'network-only',
      onCompleted: (_, ctx) => {
        notify(
          'toggle-rf-status-success-msg',
          `RF status turned ${ctx?.variables?.data?.status ? 'On' : 'Off'} successfully.`,
          'success',
        );
      },
      onError: (_, ctx) => {
        notify(
          'toggle-rf-status-error-msg',
          `Failed to turn RF status ${ctx?.variables?.data?.status ? 'On' : 'Off'}.`,
          'error',
        );
      },
    });

  const [toggleService, { loading: toggleServiceLoading }] =
    useToggleServiceMutation({
      fetchPolicy: 'network-only',
      onCompleted: (_, ctx) => {
        notify(
          'toggle-service-status-success-msg',
          `Service status turned ${ctx?.variables?.data?.status ? 'On' : 'Off'} successfully.`,
          'success',
        );
      },
      onError: (_, ctx) => {
        notify(
          'toggle-service-status-error-msg',
          `Failed to turn service status ${ctx?.variables?.data?.status ? 'On' : 'Off'}.`,
          'error',
        );
      },
    });

  const [updateSwitchPort] = useToggleInternetSwitchMutation({
    onError: (err) => notify('update-node-err-msg', err.message, 'error'),
  });

  const { loading: healthLoading } = useGetHealthReportQuery({
    variables: {
      data: {
        id: '',
        timestamp: '',
        timeframe: Timeframe_Filter.Latest,
        nodeId:
          nodes.find((node) => node.id.includes(NodeTypeEnum.Tnode))?.id || '',
      },
    },
    onCompleted: (data) => {
      if (data.getHealthReport.system.length > 0) {
        const actions: TSiteActionToggle[] = [];
        data.getHealthReport.system.forEach((system) => {
          if (system.name === 'radio') {
            actions.push({
              id: NODE_ACTIONS_ENUM.TOGGLE_RADIO,
              key: 'radio',
              value: stringToBoolean(system.value),
            });
          }
          if (system.name === 'service') {
            actions.push({
              id: NODE_ACTIONS_ENUM.TOGGLE_SERVICE,
              key: 'service',
              value: stringToBoolean(system.value),
            });
          }
        });
        setSiteActionData(actions);
      }
    },
    onError: (err) =>
      notify('fetching-health-report-msg', err.message, 'error'),
  });

  return {
    siteData,
    fetchNodesForSite,
    toggleRFStatus,
    toggleRFStatusLoading,
    toggleService,
    toggleServiceLoading,
    updateSwitchPort,
    healthLoading,
    notify,
  };
}
