/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Role_Type, Sim_Types } from '@/client/graphql/generated';
import { colors } from '@/theme';
import { ColumnsWithOptions, MenuItemType } from '@/types';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';

export const STAT_STEP_29 = 29;
export const METRIC_RANGE_3600 = 3600;
export const METRIC_RANGE_10800 = 10800;
export const NETWORK_FLOW = 'net';
export const ONBOARDING_FLOW = 'onb';
export const INSTALLATION_FLOW = 'ins';
export const CHECK_SITE_FLOW = 'chk';
export const DRAWER_WIDTH = 200;
export const APP_VERSION = 'v0.0.1';
export const COPY_RIGHTS = 'Copyright © Ukama Inc.';
export const IPFY_URL = 'https://api.ipify.org/?format=json';
export const IP_API_BASE_URL = 'https://ipapi.co';
export const REFRESH_INTERVAL = 30000;

export const SETTING_MENU = [
  { id: 'personal-settings', name: 'My Account' },
  { id: 'network-settings', name: 'Network' },
  { id: 'billing', name: 'Billing' },
  { id: 'appearance', name: 'Appearance' },
];

export const NODE_ACTIONS_ENUM = {
  NODE_RESTART: 'node-restart',
  NODE_RF_OFF: 'node-rf-off',
  NODE_OFF: 'node-off',
  NODE_RF_ON: 'node-rf-on',
};

export const NODE_ACTIONS_BUTTONS = [
  {
    id: NODE_ACTIONS_ENUM.NODE_RESTART,
    name: 'Restart',
    consent: 'Are you sure you want to restart node?',
  },
  { id: NODE_ACTIONS_ENUM.NODE_RF_OFF, name: 'Turn RF Off' },
  { id: NODE_ACTIONS_ENUM.NODE_RF_ON, name: 'Turn RF On' },
  // { id: 'node-on-off', name: 'Turn Node Off' },
];

export const SITE_CONFIG_STEPS = [
  'Configure Site Installation (1/2)',
  'Configure Site Installation (2/2)',
];

export const MEMBER_ROLES = [
  { id: 1, label: 'Administrator', value: Role_Type.RoleAdmin },
  { id: 2, label: 'Network owner', value: Role_Type.RoleNetworkOwner },
];

export const SITE_PLANNING_AP_OPTIONS = [
  { id: 1, label: 'Tower Node + 1 Amplifier Unit', value: 'ONE_TO_ONE' },
  { id: 2, label: 'Tower Node + 2 Amplifier Units', value: 'ONE_TO_TWO' },
];

export const SOLAR_UPTIME_OPTIONS = [
  { id: 1, label: '95%', value: 95 },
  { id: 2, label: '98%', value: 98 },
  { id: 3, label: '99%', value: 99 },
];

export const DATA_UNIT = [
  { id: 1, label: 'Bytes', value: 'Bytes' },
  { id: 2, label: 'KB', value: 'KiloBytes' },
  { id: 3, label: 'MB', value: 'MegaBytes' },
  { id: 4, label: 'GB', value: 'GigaBytes' },
];

export const DATA_DURATION = [
  { id: 1, label: 'Day', value: '1' },
  { id: 2, label: 'Month', value: '30' },
];

export const SIM_TYPES = [
  { id: 1, label: 'Test', value: Sim_Types.Test },
  { id: 2, label: 'Operator Data', value: Sim_Types.OperatorData },
];

export const LANGUAGE_OPTIONS = [
  { id: 1, label: '🇺🇸  English, US', value: 'en' },
  { id: 2, label: '🇫🇷  French, France', value: 'fr' },
];

export const MONTH_FILTER = [
  { id: 1, label: 'January', value: 'JANUARY' },
  { id: 2, label: 'February', value: 'FEBRUARY' },
  { id: 3, label: 'March', value: 'MARCH' },
  { id: 4, label: 'April', value: 'APRIL' },
  { id: 5, label: 'May', value: 'MAY' },
  { id: 6, label: 'June', value: 'JUNE' },
  { id: 7, label: 'July', value: 'JULY' },
  { id: 8, label: 'August', value: 'AUGUST' },
  { id: 9, label: 'September', value: 'SEPTEMBER' },
  { id: 10, label: 'October', value: 'OCTOBER' },
  { id: 11, label: 'November', value: 'NOVEMBER' },
  { id: 12, label: 'December', value: 'DECEMBER' },
];

export const TIME_FILTER = [
  { id: 1, label: 'Today', value: 'TODAY' },
  { id: 2, label: 'This week', value: 'WEEK' },
  { id: 3, label: 'Month', value: 'MONTH' },
  { id: 4, label: 'Total', value: 'TOTAL' },
];

export const SUBSCRIBER_TABLE_COLUMNS: ColumnsWithOptions[] = [
  { id: 'name', label: 'Name', minWidth: 160 },
  { id: 'email', label: 'Email', minWidth: 180 },
  { id: 'dataPlan', label: 'Data Plan', minWidth: 140 },
  { id: 'dataUsage', label: 'Data Usage', minWidth: 140 },
  { id: 'actions', label: 'Actions', align: 'right', minWidth: 80 },
];

export const BILLING_TABLE_COLUMNS: ColumnsWithOptions[] = [
  { id: 'billing', label: 'Billing period', minWidth: 100 },
  { id: 'posted', label: 'Posted', minWidth: 100 },
  { id: 'description', label: 'Description', minWidth: 100 },
  { id: 'payment', label: 'Payment', minWidth: 100 },
  { id: 'pdf', label: 'Pdf', minWidth: 100 },
];

export const SUBSCRIBER_TABLE_MENU: MenuItemType[] = [
  { id: 1, Icon: null, title: 'Edit subscriber', route: 'edit-sub' },
  { id: 2, Icon: null, title: 'Top up data', route: 'top-up-data' },
];

export const BILLING_HISTORY_TABLE_MENU: MenuItemType[] = [
  { id: 1, Icon: null, title: 'View Details', route: 'View Details' },
  { id: 2, Icon: null, title: 'Top up data', route: 'top-up-data' },
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
export const SITE_STATUS = {
  ONLINE: 'Online',
  OFFLINE: 'Offline',
  WARNING: 'Warning',
};

export const INVITATION_TABLE_COLUMN: ColumnsWithOptions[] = [
  { id: 'name', label: 'Name', minWidth: 120 },
  { id: 'role', label: 'Role', minWidth: 144 },
  { id: 'status', label: 'Status', minWidth: 144 },
  { id: 'email', label: 'Email', minWidth: 120 },
  { id: 'actions', label: 'Action', minWidth: 72 },
];

export const INVITATION_TABLE_MENU: MenuItemType[] = [
  { id: 1, Icon: null, title: 'Accept invite', route: 'accept-invite' },
  { id: 2, Icon: null, title: 'Reject invite', route: 'reject-invite' },
];

export const MEMBER_TABLE_COLUMN: ColumnsWithOptions[] = [
  { id: 'name', label: 'Name', minWidth: 160 },
  { id: 'role', label: 'Role', minWidth: 180 },
  { id: 'email', label: 'Email', minWidth: 140 },
  { id: 'actions', label: 'Actions', align: 'right', minWidth: 80 },
];

export const MEMBER_TABLE_MENU: MenuItemType[] = [
  {
    id: 1,
    Icon: null,
    title: 'Deactivate/Activate member',
    route: 'member-status-update',
  },
  { id: 2, Icon: null, title: 'Remove member', route: 'remove-member' },
];

export const MANAGE_SIM_POOL_COLUMN: ColumnsWithOptions[] = [
  { id: 'iccid', label: 'ICCID', minWidth: 160 },
  { id: 'isPhysical', label: 'Type', minWidth: 180 },
  {
    id: 'isAllocated',
    label: 'Status',
    minWidth: 140,
    options: {
      isSortable: true,
    },
  },
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

export const PAYMENT_METHODS = ['Stripe'];

export const BASIC_MENU_ACTIONS: MenuItemType[] = [
  { id: 1, Icon: EditIcon, title: 'Edit', route: 'edit' },
  { id: 2, Icon: DeleteIcon, title: 'Delete', route: 'delete' },
  { id: 3, Icon: UpdateIcon, title: 'Update available', route: 'update' },
];

export const ROAMING_SELECT = [
  { id: 1, value: 'all', text: 'CHANGE ROAMING FOR INDIVIDUAL SIMS' },
  { id: 2, value: 'esim1', text: 'ESIM 1' },
  { id: 3, value: 'esim2', text: 'ESIM 2' },
];

export const NodePageTabs = [
  { id: 'node-tab-0', label: 'Overview', value: 0 },
  { id: 'node-tab-1', label: 'Network', value: 1 },
  { id: 'node-tab-2', label: 'Resources', value: 2 },
  { id: 'node-tab-3', label: 'Radio', value: 3 },
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

export { NodeApps } from './stubData';

export const KPI_PLACEHOLDER_VALUE = '-';

export interface SiteKpiConfig {
  id: string;
  name: string;
  unit: string;
  description: string;
  tickInterval?: number;
  tickPositions?: number[];
  threshold?: {
    min: number;
    normal: number;
    max: number;
  } | null;
  show?: boolean;
  format: string;
}

export interface SectionData {
  [key: string]: SiteKpiConfig[];
}

export const SITE_KPIS = {
  SOLAR: {
    metrics: [
      {
        unit: 'W',
        show: true,
        name: 'Solar panel power',
        id: 'solar_panel_power',
        description: 'Solar power generation',
        tickInterval: 100,
        min: 0,
        max: 600,
        threshold: {
          min: 100,
          normal: 300,
          max: 500,
        },
        format: 'number',
        tickPositions: [0, 100, 200, 300, 400, 500, 600],
      },
      {
        unit: 'V',
        show: true,
        name: 'Solar panel voltage',
        id: 'solar_panel_voltage',
        description: 'Solar panel voltage',
        tickInterval: 10,
        min: 0,
        max: 120,
        threshold: {
          min: 50,
          normal: 75,
          max: 100,
        },
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 120],
      },
      {
        unit: 'A',
        show: true,
        name: 'Solar panel current',
        id: 'solar_panel_current',
        description: 'Solar panel current',
        tickInterval: 2,
        min: 0,
        max: 12,
        threshold: {
          min: 2,
          normal: 5,
          max: 8,
        },
        format: 'number',
        tickPositions: [0, 2, 4, 6, 8, 10, 12],
      },
    ],
  },
  BATTERY: {
    metrics: [
      {
        unit: '%',
        show: true,
        name: 'Battery charge',
        id: 'battery_charge_percentage',
        description: 'Battery charge percentage',
        tickInterval: 20,
        min: 0,
        max: 120,
        threshold: {
          min: 20,
          normal: 70,
          max: 100,
        },
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 120],
      },
    ],
  },
  CONTROLLER: {
    metrics: [
      {
        unit: 'W',
        show: true,
        name: 'Solar panel power',
        id: 'solar_panel_power',
        description: 'Solar power generation',
        tickInterval: 100,
        min: 0,
        max: 600,
        threshold: {
          min: 100,
          normal: 300,
          max: 500,
        },
        format: 'number',
        tickPositions: [0, 100, 200, 300, 400, 500, 600],
      },
      {
        unit: 'V',
        show: true,
        name: 'Solar panel voltage',
        id: 'solar_panel_voltage',
        description: 'Solar panel voltage',
        tickInterval: 10,
        min: 0,
        max: 120,
        threshold: {
          min: 50,
          normal: 75,
          max: 100,
        },
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 120],
      },
      {
        unit: 'A',
        show: true,
        name: 'Solar panel current',
        id: 'solar_panel_current',
        description: 'Solar panel current',
        tickInterval: 2,
        min: 0,
        max: 12,
        threshold: {
          min: 2,
          normal: 5,
          max: 8,
        },
        format: 'number',
        tickPositions: [0, 2, 4, 6, 8, 10, 12],
      },
      {
        unit: '%',
        show: true,
        name: 'Battery charge',
        id: 'battery_charge_percentage',
        description: 'Battery charge percentage',
        tickInterval: 20,
        min: 0,
        max: 120,
        threshold: {
          min: 20,
          normal: 70,
          max: 100,
        },
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 120],
      },
    ],
  },
  MAIN_BACKHAUL: {
    metrics: [
      {
        unit: 'ms',
        show: true,
        name: 'Backhaul latency',
        id: 'main_backhaul_latency',
        description: 'Main backhaul latency',
        tickInterval: 500,
        min: 0,
        max: 3000,
        threshold: {
          min: 0,
          normal: 500,
          max: 1000,
        },
        format: 'number',
        lowerIsBetter: true,
        tickPositions: [0, 500, 1000, 1500, 2000, 2500, 3000],
      },
      {
        unit: 'Mbps',
        show: true,
        name: 'Backhaul speed',
        id: 'backhaul_speed',
        description: 'Backhaul network speed',
        tickInterval: 10,
        min: 0,
        max: 120,
        threshold: {
          min: 20,
          normal: 50,
          max: 100,
        },
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 120],
      },
    ],
  },
  SWITCH: {
    metrics: [
      {
        unit: 'Mbps',
        show: true,
        name: 'Backhaul switch port speed',
        id: 'backhaul_switch_port_speed',
        description: 'Backhaul switch port speed',
        tickInterval: 100,
        min: 0,
        max: 1000,
        threshold: {
          min: 50,
          normal: 200,
          max: 1000,
        },
        format: 'number',
      },
      {
        unit: '',
        show: true,
        name: 'Node switch port status',
        id: 'node_switch_port_status',
        description: 'Node switch port status (1 = up, 0 = down)',
        tickInterval: 1,
        min: 0,
        max: 1,
        threshold: {
          min: 0,
          normal: 1,
          max: 1,
        },
        format: 'number',
      },
      {
        unit: '',
        show: true,
        name: 'Solar switch port status',
        id: 'solar_switch_port_status',
        description: 'Solar controller switch port status (1 = up, 0 = down)',
        tickInterval: 1,
        min: 0,
        max: 1,
        threshold: {
          min: 0,
          normal: 1,
          max: 1,
        },
        format: 'number',
      },
      {
        unit: '',
        show: true,
        name: 'Backhaul switch port status',
        id: 'backhaul_switch_port_status',
        description: 'Backhaul switch port status (1 = up, 0 = down)',
        tickInterval: 1,
        min: 0,
        max: 1,
        threshold: {
          min: 0,
          normal: 1,
          max: 1,
        },
        format: 'number',
      },
      {
        unit: 'W',
        show: true,
        name: 'Backhaul switch port power',
        id: 'backhaul_switch_port_power',
        description: 'Backhaul switch port power consumption',
        tickInterval: 0.5,
        min: 0,
        max: 7,
        threshold: {
          min: 5,
          normal: 6,
          max: 7,
        },
        format: 'number',
      },

      {
        unit: 'Mbps',
        show: true,
        name: 'Solar switch port speed',
        id: 'solar_switch_port_speed',
        description: 'Solar controller switch port speed',
        tickInterval: 100,
        min: 0,
        max: 1000,
        threshold: {
          min: 0,
          normal: 200,
          max: 1000,
        },
        format: 'number',
      },
      {
        unit: 'W',
        show: true,
        name: 'Solar switch port power',
        id: 'solar_switch_port_power',
        description: 'Solar controller switch port power consumption',
        tickInterval: 0.5,
        min: 0,
        max: 7,
        threshold: {
          min: 5,
          normal: 6,
          max: 7,
        },
        format: 'number',
      },

      {
        unit: 'Mbps',
        show: true,
        name: 'Node switch port speed',
        id: 'node_switch_port_speed',
        description: 'Node switch port speed',
        tickInterval: 100,
        min: 0,
        max: 1000,
        threshold: {
          min: 0,
          normal: 100,
          max: 1000,
        },
        format: 'number',
      },
      {
        unit: 'W',
        show: true,
        name: 'Node switch port power',
        id: 'node_switch_port_power',
        description: 'Node switch port power consumption',
        tickInterval: 10,
        min: 0,
        max: 120,
        threshold: {
          min: 0,
          normal: 75,
          max: 110,
        },
        format: 'number',
      },
    ],
  },
  SITE: {
    stats: [
      {
        unit: 's',
        show: true,
        name: 'Site uptime',
        id: 'site_uptime_seconds',
        description: 'Cumulative site operational time',
        format: 'number',
      },
      {
        unit: '%',
        show: true,
        name: 'Site uptime',
        id: 'site_uptime_percentage',
        description: 'Site uptime percentage',
        format: 'number',
      },
      {
        unit: '',
        show: true,
        name: 'Node uptime',
        id: 'unit_uptime',
        description: 'Node uptime',
        format: 'number',
      },
    ],
  },
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
        unit: '',
        show: true,
        format: 'number',
        name: 'Nde uptime',
        id: 'unit_uptime',
        description: 'Node uptime',
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
        id: 'node_load',
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
        id: 'unit_health',
        description: 'Hardware health index (Temperature)',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
    ],
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
        name: 'Load',
        id: 'hwd_load',
        unit: '%',
        format: 'number',
        description: 'Hardware focus  NLI',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: '%',
        name: 'Memory',
        description: '',
        id: 'memory_usage',
        format: 'number',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: '%',
        name: 'CPU',
        description: '',
        id: 'cpu_usage',
        tickInterval: 20,
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
      {
        unit: '%',
        name: 'Disk',
        description: '',
        id: 'disk_usage',
        tickInterval: 20,
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 101],
        threshold: { min: 0, normal: 80, max: 100 },
      },
    ],
  },
  RADIO: {
    tnode: [
      {
        unit: 'dBm',
        id: 'txpower',
        description: '',
        name: 'TX Power',
        tickInterval: 10,
        format: 'decimal',
        tickPositions: [0, 10, 20, 30, 40, 41],
        threshold: { min: 0, normal: 31, max: 34 },
      },
    ],
  },
};

export const TIME_FILTER_OPTIONS = [
  { id: '1', label: 'LIVE' },
  { id: '2', label: 'ZOOM' },
];

export const SITE_KPI_TYPES = {
  SITE_UPTIME: 'site_uptime_seconds',
  BATTERY_CHARGE_PERCENTAGE: 'battery_charge_percentage',
  BACKHAUL_SPEED: 'backhaul_speed',
  SITE_UPTIME_PERCENTAGE: 'site_uptime_percentage',
  NODE_UPTIME: 'unit_uptime',
};
