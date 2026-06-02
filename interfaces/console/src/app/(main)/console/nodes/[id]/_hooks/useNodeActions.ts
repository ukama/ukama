/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { activeGraphTypeVar, activeNodeTabVar } from '@/client/vars';
import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import { METRIC_RANGE_10800, NODE_ACTIONS_ENUM } from '@/constants';
import { TNodeActionState, TStatusBarObj } from '@/types';
import { getNodeTabTypeByIndex, getUnixTime } from '@/utils';
import { useRouter } from 'next/navigation';
import React from 'react';

interface UseNodeActionsParams {
  currentNodeId: string | undefined;
  setIsEditNode: React.Dispatch<React.SetStateAction<boolean>>;
  setNodeAction: React.Dispatch<
    React.SetStateAction<
      TNodeActionState & { currentAction: string; actionInitiated: string }
    >
  >;
  setMetricFrom: React.Dispatch<React.SetStateAction<number>>;
  handleEditNode: (name: string) => void;
  handleRestartNode: () => void;
}

export function useNodeActions({
  setIsEditNode,
  setNodeAction,
  setMetricFrom,
  handleEditNode: doEditNode,
  handleRestartNode: doRestartNode,
}: UseNodeActionsParams) {
  const router = useRouter();

  const handleNodeSelected = (obj: TStatusBarObj) => {
    router.push(`/console/nodes/${obj.id}`);
  };

  const handleEditNode = (name: string) => {
    setIsEditNode(false);
    doEditNode(name);
  };

  const handleSectionChange = (type: Graphs_Type) => {
    activeGraphTypeVar(type);
    setMetricFrom(() => getUnixTime() - METRIC_RANGE_10800);
  };

  const onTabSelected = (_: unknown, value: number) => {
    activeNodeTabVar(value);
    activeGraphTypeVar(getNodeTabTypeByIndex(value));
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
      doRestartNode();
    }
  };

  return {
    handleNodeSelected,
    handleEditNode,
    handleSectionChange,
    onTabSelected,
    handleNodeActionClick,
  };
}
