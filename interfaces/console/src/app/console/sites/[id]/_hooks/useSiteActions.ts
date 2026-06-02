/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { Node, NodeTypeEnum } from '@/client/graphql/generated';
import { ToggleInternetSwitchMutation } from '@/client/graphql/generated';
import { ToggleRfStatusMutation } from '@/client/graphql/generated';
import { ToggleServiceMutation } from '@/client/graphql/generated';
import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import { NODE_ACTIONS_ENUM } from '@/constants';
import { useUIContext } from '@/context';
import { ActiveView, KPIType, TSiteActionToggle, TStatusBarObj } from '@/types';
import { kpiToGraphType } from '@/utils';
import { MutationFunction } from '@apollo/client';
import { useRouter } from 'next/navigation';
import { Dispatch, SetStateAction, useCallback } from 'react';

interface UseSiteActionsProps {
  id: string;
  nodes: Node[];
  setSiteActionData: Dispatch<SetStateAction<TSiteActionToggle[]>>;
  setActiveView: Dispatch<SetStateAction<ActiveView>>;
  resetMetrics: () => void;
  toggleRFStatus: MutationFunction<
    ToggleRfStatusMutation,
    { data: { nodeId: string; status: boolean } }
  >;
  toggleService: MutationFunction<
    ToggleServiceMutation,
    { data: { nodeId: string; status: boolean } }
  >;
  updateSwitchPort: MutationFunction<
    ToggleInternetSwitchMutation,
    { data: { port: number; siteId: string; status: boolean } }
  >;
}

export function useSiteActions({
  id,
  nodes,
  setSiteActionData,
  setActiveView,
  resetMetrics,
  toggleRFStatus,
  toggleService,
  updateSwitchPort,
}: UseSiteActionsProps) {
  const router = useRouter();
  const { setSnackbarMessage } = useUIContext();

  const notify = useCallback(
    (msgId: string, message: string, type: string) =>
      setSnackbarMessage({ id: msgId, message, type, show: true }),
    [setSnackbarMessage],
  );

  const handleViewChange = useCallback(
    (kpiType: string): void => {
      setActiveView({
        graphType: kpiToGraphType[kpiType] || Graphs_Type.Solar,
        kpi: kpiType as KPIType,
      });
      resetMetrics();
    },
    [setActiveView, resetMetrics],
  );

  const handleSwitchChange = useCallback(
    async (portNumber: number, currentStatus: boolean) => {
      const newStatus = !currentStatus;
      try {
        const result = await updateSwitchPort({
          variables: {
            data: { port: portNumber, siteId: id, status: newStatus },
          },
        });
        if (result?.data?.toggleInternetSwitch?.success) {
          notify(
            'update-switch-success',
            `Port ${portNumber} status updated to ${newStatus ? 'On' : 'Off'}`,
            'success',
          );
        }
      } catch (err) {
        notify(
          'update-site-error',
          err instanceof Error ? err.message : 'Unknown error',
          'error',
        );
      }
    },
    [id, updateSwitchPort, notify],
  );

  const handleSiteChange = useCallback(
    (newSiteId: string) => {
      router.push('/console/sites/' + newSiteId);
    },
    [router],
  );

  const handleActionClick = useCallback(
    (actionId: string, value: boolean) => {
      const tnodeId =
        nodes.find((node) => node.id.includes(NodeTypeEnum.Tnode))?.id ?? '';
      switch (actionId) {
        case NODE_ACTIONS_ENUM.TOGGLE_RADIO:
          toggleRFStatus({
            variables: { data: { nodeId: tnodeId, status: value } },
          });
          break;
        case NODE_ACTIONS_ENUM.TOGGLE_SERVICE:
          toggleService({
            variables: { data: { nodeId: tnodeId, status: value } },
          });
          break;
      }
      setSiteActionData((prev) =>
        prev.map((item) => (item.id === actionId ? { ...item, value } : item)),
      );
    },
    [nodes, toggleRFStatus, toggleService, setSiteActionData],
  );

  const handleSelected = useCallback(
    (obj: TStatusBarObj) => handleSiteChange(obj.id),
    [handleSiteChange],
  );

  return {
    handleViewChange,
    handleSwitchChange,
    handleSiteChange,
    handleActionClick,
    handleSelected,
  };
}
