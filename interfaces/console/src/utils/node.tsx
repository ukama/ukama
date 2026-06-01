/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import {
  Node,
  NodeConnectivityEnum,
  NodeTypeEnum,
  Nodes,
  SitesResDto,
} from '@/client/graphql/generated';
import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import { NODE_ACTIONS_ENUM } from '@/constants';
import colors from '@/theme/colors';
import { StatusType, StyleOutput, TNodeSiteTree } from '@/types';
import Battery50Icon from '@mui/icons-material/Battery50';
import BatteryAlertIcon from '@mui/icons-material/BatteryAlert';
import BatteryChargingFullIcon from '@mui/icons-material/BatteryChargingFull';
import RouterIcon from '@mui/icons-material/Router';
import SignalCellular1BarIcon from '@mui/icons-material/SignalCellular1Bar';
import SignalCellular2BarIcon from '@mui/icons-material/SignalCellular2Bar';
import SignalCellularAltIcon from '@mui/icons-material/SignalCellularAlt';
import SignalCellularConnectedNoInternet4BarIcon from '@mui/icons-material/SignalCellularConnectedNoInternet4Bar';
import SignalCellularOffIcon from '@mui/icons-material/SignalCellularOff';

export const getNodeTabTypeByIndex = (index: number) => {
  switch (index) {
    case 0:
      return Graphs_Type.NodeHealth;
    case 1:
      return Graphs_Type.NetworkCellular;
    case 2:
      return Graphs_Type.Resources;
    case 3:
      return Graphs_Type.Radio;
    case 4:
      return Graphs_Type.Subscribers;
    default:
      return Graphs_Type.NodeHealth;
  }
};

export const structureNodeSiteDate = (nodes: Nodes, sites: SitesResDto) => {
  const t: TNodeSiteTree[] = [];

  sites.sites.forEach((site) => {
    nodes.nodes.forEach((node: Node) => {
      if (node.site.siteId === site.id && node.type === NodeTypeEnum.Tnode) {
        t.push({
          id: site.id,
          name: `${site.name} (Site)`,
          nodeId: node.id,
          nodeType: node.type,
          nodeName: `${node.name} (Node)`,
        });
      }
    });
  });

  return t;
};

export const NodeEnumToString = (type: NodeTypeEnum): string => {
  switch (type) {
    case NodeTypeEnum.Tnode:
      return 'Tower Node';
    case NodeTypeEnum.Anode:
      return 'Amplifier Node';
    case NodeTypeEnum.Hnode:
      return 'Home Node';
    case NodeTypeEnum.Cnode:
      return 'Controller Node';
    default:
      return 'Unknown';
  }
};

export const getNodeTypeFromId = (id: string) => {
  if (id.includes('tnode')) return NodeTypeEnum.Tnode;
  if (id.includes('anode')) return NodeTypeEnum.Anode;
  if (id.includes('hnode')) return NodeTypeEnum.Hnode;
  if (id.includes('cnode')) return NodeTypeEnum.Cnode;
  return null;
};

export const nodeTypeEnumToString = (nodeType: NodeTypeEnum) => {
  switch (nodeType) {
    case NodeTypeEnum.Tnode:
      return 'tnode';
    case NodeTypeEnum.Anode:
      return 'anode';
    case NodeTypeEnum.Hnode:
      return 'hnode';
    case NodeTypeEnum.Cnode:
      return 'cnode';
    default:
      return 'tnode';
  }
};

export const getNodeActionDescriptionByProgress = (
  progress: number,
  action: string,
) => {
  if (action === NODE_ACTIONS_ENUM.NODE_RESTART) {
    switch (progress) {
      case 25:
        return 'Node restart initiated...';
      case 50:
        return 'Node is offline...';
      case 75:
        return 'Node is back online...';
      case 100:
        return 'Node is ready to use.';
      default:
        return '';
    }
  }
  if (action === NodeConnectivityEnum.Online) {
    return 'Node is online...';
  }
  if (action === NodeConnectivityEnum.Offline) {
    return 'Node is offline...';
  }
  return '';
};

export const getConnectionStyles = (connectionStatus: string) => {
  switch (connectionStatus) {
    case 'Online':
      return {
        color: colors.green,
        icon: <RouterIcon sx={{ color: colors.green }} />,
      };
    case 'Offline':
      return {
        color: colors.red,
        icon: <RouterIcon sx={{ color: colors.red }} />,
      };
    case 'Warning':
      return {
        color: colors.orange,
        icon: <RouterIcon sx={{ color: colors.orange }} />,
      };
    default:
      return {
        color: colors.green,
        icon: <RouterIcon sx={{ color: colors.green }} />,
      };
  }
};

export const getBatteryStyles = (batteryStatus: string) => {
  switch (batteryStatus) {
    case 'Charged':
      return {
        color: colors.green,
        icon: <BatteryChargingFullIcon sx={{ color: colors.green }} />,
      };
    case 'Medium':
      return {
        color: colors.orange,
        icon: <Battery50Icon sx={{ color: colors.orange }} />,
      };
    case 'Low':
      return {
        color: colors.red,
        icon: <BatteryAlertIcon sx={{ color: colors.red }} />,
      };
    default:
      return {
        color: colors.green,
        icon: <BatteryChargingFullIcon sx={{ color: colors.green }} />,
      };
  }
};

export const getSignalStyles = (signalStrength: string) => {
  switch (signalStrength) {
    case 'Strong':
      return {
        color: colors.green,
        icon: <SignalCellularAltIcon sx={{ color: colors.green }} />,
      };
    case 'Medium':
      return {
        color: colors.orange,
        icon: <SignalCellular2BarIcon sx={{ color: colors.orange }} />,
      };
    case 'Weak':
      return {
        color: colors.red,
        icon: <SignalCellular1BarIcon sx={{ color: colors.red }} />,
      };
    default:
      return {
        color: colors.green,
        icon: <SignalCellularAltIcon sx={{ color: colors.green }} />,
      };
  }
};

export const getStatusStyles = (type: StatusType, value: number): StyleOutput => {
  if (type === 'uptime') {
    return value <= 0
      ? { color: colors.red, icon: <RouterIcon sx={{ color: colors.red }} /> }
      : {
          color: colors.green,
          icon: <RouterIcon sx={{ color: colors.green }} />,
        };
  }

  if (type === 'battery') {
    if (value < 20) {
      return {
        color: colors.red,
        icon: <BatteryAlertIcon sx={{ color: colors.red }} />,
      };
    } else if (value < 40) {
      return {
        color: colors.orange,
        icon: <Battery50Icon sx={{ color: colors.orange }} />,
      };
    } else if (value > 60) {
      return {
        color: colors.green,
        icon: <BatteryChargingFullIcon sx={{ color: colors.green }} />,
      };
    }
  }

  if (type === 'signal') {
    if (value < 10) {
      return {
        color: colors.red,
        icon: <SignalCellularOffIcon sx={{ color: colors.red }} />,
      };
    } else if (value < 70) {
      return {
        color: colors.orange,
        icon: (
          <SignalCellularConnectedNoInternet4BarIcon
            sx={{ color: colors.orange }}
          />
        ),
      };
    } else {
      return {
        color: colors.green,
        icon: <SignalCellularAltIcon sx={{ color: colors.green }} />,
      };
    }
  }

  return {
    color: colors.green,
    icon: <RouterIcon sx={{ color: colors.green }} />,
  };
};

export const kpiToGraphType: Record<string, Graphs_Type> = {
  solar: Graphs_Type.Solar,
  battery: Graphs_Type.Battery,
  controller: Graphs_Type.Controller,
  main_backhaul: Graphs_Type.MainBackhaul,
  backhaul: Graphs_Type.MainBackhaul,
  switch: Graphs_Type.Switch,
  node: Graphs_Type.Site,
};

export const graphTypeToSection: Record<Graphs_Type | string, string> = {
  [Graphs_Type.Solar]: 'SOLAR',
  [Graphs_Type.Battery]: 'BATTERY',
  [Graphs_Type.Controller]: 'CONTROLLER',
  [Graphs_Type.MainBackhaul]: 'MAIN_BACKHAUL',
  [Graphs_Type.Switch]: 'SWITCH',
  [Graphs_Type.Site]: 'SITE',
};

export const getSectionFromKPI = (kpi: string) => {
  switch (kpi) {
    case 'solar':
      return 'SOLAR';
    case 'battery':
      return 'BATTERY';
    case 'controller':
      return 'CONTROLLER';
    case 'backhaul':
      return 'MAIN_BACKHAUL';
    case 'switch':
      return 'SWITCH';
    case 'node':
      return 'NODE';
    default:
      return 'SOLAR';
  }
};
