/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { ColumnsWithOptions, MenuItemType } from '@/types';

export const STAT_STEP_29 = 29;
export const METRIC_RANGE_3600 = 3600;
export const METRIC_RANGE_10800 = 10800;

export const NODE_ACTIONS_ENUM = {
  NODE_LOADING: 'node-loading',
  NODE_RESTART: 'node-restart',
  NODE_RADIO_ON: 'node-radio-on',
  NODE_RADIO_OFF: 'node-radio-off',
  NODE_SERVICE_ON: 'node-service-on',
  NODE_SERVICE_OFF: 'node-service-off',
  TOGGLE_RADIO: 'toggle-radio',
  TOGGLE_SERVICE: 'toggle-service',
};

export const NODE_ACTIONS_BUTTONS = [
  {
    id: NODE_ACTIONS_ENUM.NODE_RESTART,
    name: 'Restart',
    consent: 'Are you sure you want to restart node?',
    type: 'button',
  },
];

export const NODE_TABLE_COLUMNS: ColumnsWithOptions[] = [
  { id: 'id', label: 'Node #', minWidth: 160 },
  { id: 'type', label: 'Type', minWidth: 180 },
  { id: 'connectivity', label: 'Connectivity', minWidth: 140 },
  { id: 'site', label: 'Site', minWidth: 140 },
  { id: 'actions', label: 'Actions', align: 'right', minWidth: 80 },
];

export const NODE_TABLE_MENU: MenuItemType[] = [
  {
    id: 2,
    Icon: null,
    title: 'Turn node off',
    route: 'node-off',
    color: colors.redMatt,
  },
  {
    id: 3,
    Icon: null,
    title: 'Restart node',
    route: 'restart-node',
    color: colors.redMatt,
  },
  {
    id: 4,
    Icon: null,
    title: 'Restart RF',
    route: 'restart-rf',
    color: colors.redMatt,
  },
];

export const NodePageTabs = [
  { id: 'node-tab-0', label: 'Overview', value: 0 },
  { id: 'node-tab-1', label: 'Network', value: 1 },
  { id: 'node-tab-2', label: 'Resources', value: 2 },
  { id: 'node-tab-3', label: 'Radio', value: 3 },
  { id: 'node-tab-4', label: 'Software', value: 4 },
];

export const MANAGE_NODE_POOL_COLUMN: ColumnsWithOptions[] = [
  { id: 'id', label: 'Node #', minWidth: 160 },
  { id: 'type', label: 'Type', minWidth: 120 },
  { id: 'connectivity', label: 'Connectivity', minWidth: 120 },
  {
    id: 'state',
    label: 'State',
    minWidth: 120,
    options: {
      isSortable: true,
    },
  },
  { id: 'site', label: 'Site', minWidth: 120 },
  { id: 'createdAt', label: 'Date installed', minWidth: 140 },
];

export const MASK_BY_TYPE = {
  hnode: '{uk- }######{ -hnode- }##{ - }####',
  anode: '{uk- }######{ -\\anode- }##{ - }}####',
  tnode: '{uk- }######{ -tnode- }##{ - }}####',
};

export const MASK_PLACEHOLDERS = {
  hnode: 'uk- ______ -hnode- __ - ____',
  anode: 'uk- ______ -anode- __ - ____',
  tnode: 'uk- ______ -tnode- __ - ____',
};

export const NODE_IMAGES = {
  tnode:
    'https://ukama-site-assets.s3.amazonaws.com/images/ukama_tower_node.png',
  anode:
    'https://ukama-site-assets.s3.amazonaws.com/images/ukama_amplifier_node.png',
  hnode:
    'https://ukama-site-assets.s3.amazonaws.com/images/ukama_home_node.png',
};

export const NODE_KPIS = {
  HOME: {
    stats: [
      {
        unit: '',
        show: true,
        format: 'number',
        name: 'Network uptime',
        id: 'network_uptime',
        description: 'Network uptime',
        threshold: null,
      },
      {
        unit: '',
        show: true,
        name: 'Sales',
        format: 'decimal',
        id: 'package_sales',
        description: 'Network sales',
        threshold: null,
      },
      {
        unit: '',
        show: true,
        format: 'number',
        name: 'Data volume',
        id: 'data_usage',
        description: 'Network data volume',
        threshold: null,
      },
      {
        unit: '',
        show: true,
        format: 'number',
        name: 'Network subscribers',
        id: 'node_active_subscribers',
        description: 'Network active subscribers',
        threshold: null,
      },
    ],
  },
  NODE_UPTIME: {
    tnode: [
      {
        unit: 's',
        show: true,
        format: 'number',
        name: 'Node uptime',
        id: 'uptime',
        description: 'Node uptime in seconds',
        threshold: null,
        tickInterval: undefined,
        tickPositions: undefined,
      },
    ],
    anode: [
      {
        unit: 's',
        show: true,
        format: 'number',
        name: 'Node uptime',
        id: 'uptime',
        description: 'Node uptime in seconds',
        threshold: null,
        tickInterval: undefined,
        tickPositions: undefined,
      },
    ],
    hnode: [
      {
        unit: 's',
        show: true,
        format: 'number',
        name: 'Node uptime',
        id: 'uptime',
      },
    ],
    cnode: [
      {
        unit: 's',
        show: true,
        format: 'number',
        name: 'Node uptime',
        id: 'uptime',
        description: 'Node uptime in seconds',
        threshold: null,
        tickInterval: undefined,
        tickPositions: undefined,
      },
    ],
  },
  HEALTH: {
    tnode: [
      {
        unit: '%',
        show: true,
        name: 'Load',
        format: 'number',
        id: 'memory',
        description: 'Node Load index',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: '°',
        show: true,
        format: 'number',
        name: 'Hardware',
        id: 'cpu_temperature',
        description: 'Hardware health index (Temperature)',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
    ],
    anode: [
      {
        unit: '°',
        show: true,
        format: 'number',
        name: 'FEM 1 Temp.',
        id: 'fem1_temperature',
        description: 'FEM 1 Temperature',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: '°',
        show: true,
        format: 'number',
        name: 'FEM 2 Temp.',
        id: 'fem2_temperature',
        description: 'FEM 2 Temperature',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
    ],
    cnode: [
      {
        unit: '%',
        show: true,
        name: 'Load',
        format: 'number',
        id: 'memory',
        description: 'Node Load index',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
    ],
    hnode: [],
  },
  SUBSCRIBER: {
    tnode: [
      {
        unit: '',
        name: 'Active',
        format: 'number',
        id: 'subscribers_active',
        description: 'Current active subscriber in a network',
        threshold: null,
        tickInterval: 20,
        tickPositions: undefined,
      },
    ],
  },
  NETWORK_CELLULAR: {
    tnode: [
      {
        unit: 'mbps',
        name: 'Uplink',
        description: '',
        format: 'number',
        id: 'cellular_uplink',
        tickInterval: 5,
        tickPositions: undefined,
        threshold: { min: 0, normal: 5, max: 30 },
      },
      {
        unit: 'mbps',
        description: '',
        name: 'Downlink',
        format: 'number',
        id: 'cellular_downlink',
        tickInterval: 20,
        tickPositions: undefined,
        threshold: { min: 0, normal: 60, max: 160 },
      },
    ],
  },
  NETWORK_BACKHAUL: {
    tnode: [
      {
        unit: 'mbps',
        description: '',
        name: 'Uplink',
        format: 'number',
        id: 'backhaul_uplink',
        tickInterval: 5,
        tickPositions: undefined,
        threshold: { min: 0, normal: 10, max: 200 },
      },
      {
        unit: 'mbps',
        description: '',
        name: 'Downlink',
        format: 'number',
        id: 'backhaul_downlink',
        tickInterval: 5,
        tickPositions: undefined,
        threshold: { min: 0, normal: 10, max: 200 },
      },
      {
        unit: 'ms',
        name: 'Latency',
        description: '',
        format: 'decimal',
        id: 'backhaul_latency',
        tickInterval: 200,
        tickPositions: [0, 200, 400, 600, 800, 1000, 1050],
        threshold: { min: 0, normal: 800, max: 1000 },
      },
    ],
  },
  RESOURCES: {
    tnode: [
      {
        unit: '%',
        name: 'CPU',
        description: 'CPU usage percentage',
        id: 'cpu',
        format: 'number',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: 'MB',
        name: 'Memory',
        description: 'Used memory (MB)',
        id: 'memory',
        format: 'number',
        tickInterval: 20,
        tickPositions: undefined,
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: 'MB',
        name: 'Disk',
        description: 'Used disk space (MB)',
        id: 'disk',
        tickInterval: 20,
        format: 'number',
        tickPositions: undefined,
        threshold: { min: 0, normal: 80, max: 100 },
      },
    ],
    anode: [
      {
        unit: '%',
        name: 'CPU',
        description: 'CPU usage percentage',
        id: 'cpu',
        format: 'number',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: 'MB',
        name: 'Memory',
        description: 'Used memory (MB)',
        id: 'memory',
        format: 'number',
        tickInterval: 20,
        tickPositions: undefined,
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: 'MB',
        name: 'Disk',
        description: 'Used disk space (MB)',
        id: 'disk',
        tickInterval: 20,
        format: 'number',
        tickPositions: undefined,
        threshold: { min: 0, normal: 80, max: 100 },
      },
    ],
    cnode: [
      {
        unit: '%',
        name: 'CPU',
        description: 'CPU usage percentage',
        id: 'cpu',
        format: 'number',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: 'MB',
        name: 'Memory',
        description: 'Used memory (MB)',
        id: 'memory',
        format: 'number',
        tickInterval: 20,
        tickPositions: undefined,
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: 'MB',
        name: 'Disk',
        description: 'Used disk space (MB)',
        id: 'disk',
        tickInterval: 20,
        format: 'number',
        tickPositions: undefined,
        threshold: { min: 0, normal: 80, max: 100 },
      },
    ],
  },
  RADIO: {
    tnode: [
      {
        unit: 'dBm',
        id: 'power',
        description: 'Transmit power',
        name: 'TX Power',
        tickInterval: 10,
        format: 'decimal',
        tickPositions: [0, 10, 20, 30, 40, 41],
        threshold: { min: 0, normal: 31, max: 34 },
      },
    ],
    anode: [
      {
        unit: 'dBm',
        id: 'rx_power',
        description: 'Receive power',
        name: 'RX Power',
        tickInterval: 10,
        format: 'decimal',
        tickPositions: [0, 10, 20, 30, 40, 41],
        threshold: { min: 0, normal: 31, max: 34 },
      },
      {
        unit: 'dBm',
        id: 'pa_power',
        description: 'PA power',
        name: 'PA Power',
        tickInterval: 10,
        format: 'decimal',
        tickPositions: [0, 10, 20, 30, 40, 41],
        threshold: { min: 0, normal: 31, max: 34 },
      },
      {
        unit: 'dBm',
        id: 'tx_power',
        description: 'TX power',
        name: 'TX Power',
        tickInterval: 10,
        format: 'decimal',
        tickPositions: [0, 10, 20, 30, 40, 41],
        threshold: { min: 0, normal: 31, max: 34 },
      },
    ],
    cnode: [],
    hnode: [],
  },
};
