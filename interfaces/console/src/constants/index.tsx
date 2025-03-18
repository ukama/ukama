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

export const SETTING_MENU = [
  { id: 'personal-settings', name: 'My Account' },
  { id: 'network-settings', name: 'Network' },
  { id: 'billing', name: 'Billing' },
  { id: 'appearance', name: 'Appearance' },
];

export const NODE_ACTIONS_BUTTONS = [
  {
    id: 'node-restart',
    name: 'Restart',
    consent: 'Are you sure you want to restart node?',
  },
  // { id: 'node-on-off', name: 'Turn Node Off' },
  // { id: 'node-rf-off', name: 'Turn RF Off' },
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
  { id: 'isAllocated', label: 'Status', minWidth: 140 },
];

export const MANAGE_NODE_POOL_COLUMN: ColumnsWithOptions[] = [
  { id: 'id', label: 'Node #', minWidth: 160 },
  { id: 'type', label: 'Type', minWidth: 180 },
  { id: 'connectivity', label: 'Connectivity', minWidth: 120 },
  { id: 'state', label: 'State', minWidth: 120 },
  { id: 'site', label: 'Site', minWidth: 180 },
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

export const NODE_KPIS = {
  HOME: {
    stats: [
      {
        unit: '',
        show: true,
        name: 'Network uptime',
        id: 'network_uptime',
        description: 'Network uptime',
        threshold: { min: 0, normal: 0, max: 2678400 },
      },
      {
        unit: '',
        show: true,
        name: 'Sales',
        id: 'network_sales',
        description: 'Network sales',
        threshold: { min: 0, normal: 10000, max: 50000 },
      },
      {
        unit: '',
        show: true,
        name: 'Data volume',
        id: 'network_data_volume',
        description: 'Network data volume',
        threshold: { min: 0, normal: 512000, max: 1024000 },
      },
      {
        unit: '',
        show: true,
        name: 'Network subscribers',
        id: 'network_active_ue',
        description: 'Network active subscribers',
        threshold: { min: 0, normal: 500, max: 10000 },
      },
    ],
  },
  NODE_UPTIME: {
    tnode: [
      {
        unit: '',
        show: true,
        name: 'Nde uptime',
        id: 'unit_uptime',
        description: 'Node uptime',
        threshold: { min: 0, normal: 0, max: 2678400 },
      },
    ],
  },
  HEALTH: {
    tnode: [
      {
        unit: '%',
        show: true,
        name: 'Load',
        id: 'node_load',
        description: 'Node Load index',
        threshold: { min: 50, normal: 75, max: 90 },
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100],
      },
      {
        unit: '°',
        show: true,
        name: 'Hardware',
        id: 'unit_health',
        description: 'Hardware health index (Temperature)',
        threshold: { min: 50, normal: 80, max: 100 },
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100],
      },
    ],
  },
  SUBSCRIBER: {
    tnode: [
      {
        unit: '',
        name: 'Active',
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
        unit: 'Mbps',
        name: 'Uplink',
        description: '',
        threshold: null,
        id: 'cellular_uplink',
        tickInterval: 300,
        tickPositions: undefined,
      },
      {
        unit: 'Mbps',
        description: '',
        threshold: null,
        name: 'Downlink',
        id: 'cellular_downlink',
        tickInterval: 300,
        tickPositions: undefined,
      },
    ],
  },
  NETWORK_BACKHAUL: {
    tnode: [
      {
        unit: 'Mbps',
        description: '',
        name: 'Uplink',
        threshold: null,
        id: 'backhaul_uplink',
        tickInterval: 300,
        tickPositions: undefined,
      },
      {
        unit: 'Mbps',
        description: '',
        name: 'Downlink',
        id: 'backhaul_downlink',
        tickInterval: 300,
        tickPositions: undefined,
        threshold: null,
      },
      {
        unit: 'ms',
        name: 'Latency',
        description: '',
        id: 'backhaul_latency',
        tickInterval: 50,
        tickPositions: [0, 50, 100, 150, 200],
        threshold: { min: 100, normal: 150, max: 200 },
      },
    ],
  },
  RESOURCES: {
    tnode: [
      {
        name: 'Load',
        id: 'hwd_load',
        unit: '%',
        description: 'Hardware focus  NLI',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100],
        threshold: { min: 50, normal: 70, max: 80 },
      },
      {
        unit: '%',
        name: 'Memory',
        description: '',
        id: 'memory_usage',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100],
        threshold: { min: 40, normal: 70, max: 80 },
      },
      {
        unit: '%',
        name: 'CPU',
        description: '',
        id: 'cpu_usage',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100],
        threshold: { min: 40, normal: 70, max: 80 },
      },
      {
        unit: '%',
        name: 'Disk',
        description: '',
        id: 'disk_usage',
        tickInterval: 20,
        tickPositions: [0, 20, 40, 60, 80, 100],
        threshold: { min: 40, normal: 70, max: 80 },
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
        tickInterval: 30,
        tickPositions: [0, 30, 60, 90, 120],
        threshold: { min: 30, normal: 60, max: 95 },
      },
    ],
  },
};

export const TIME_FILTER_OPTIONS = [
  { id: '1', label: 'LIVE' },
  { id: '2', label: 'ZOOM' },
];
