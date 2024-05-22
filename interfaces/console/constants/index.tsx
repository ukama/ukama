import { ColumnsWithOptions, MenuItemType } from '@/types';
import AccountTreeIcon from '@mui/icons-material/AccountTree';
import AppsIcon from '@mui/icons-material/Apps';
import BillingIcon from '@mui/icons-material/CreditCard';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import HomeIcon from '@mui/icons-material/Home';
import {
  default as PeopleAltIcon,
  default as UsersIcon,
} from '@mui/icons-material/PeopleAlt';
import PermDataSettingIcon from '@mui/icons-material/PermDataSetting';
import RouterIcon from '@mui/icons-material/Router';
import SimCardIcon from '@mui/icons-material/SimCard';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';
import { DataTableWithOptionColumns } from './tableColumns';

const DRAWER_WIDTH = 200;
const APP_VERSION = 'v0.0.1';
const COPY_RIGHTS = 'Copyright © Ukama Inc.';
const IPFY_URL = 'https://api.ipify.org/?format=json';
const IP_API_BASE_URL = 'https://ipapi.co';
const SETTING_MENU = [
  { id: 'personal-settings', name: 'My Account' },
  { id: 'network-settings', name: 'Network' },
  { id: 'billing', name: 'Billing' },
  { id: 'appearance', name: 'Appearance' },
];
const NODE_ACTIONS_BUTTONS = [
  {
    id: 'node-on-off',
    name: 'Turn Node Off',
  },
  {
    id: 'node-restart',
    name: 'Restart',
  },
  {
    id: 'node-rf-off',
    name: 'Turn RF Off',
  },
];

const TABLE_EXPORT_OPTIONS = [
  { id: 1, label: 'EXPORT', value: 'EXPORT' },
  { id: 2, label: 'WORD', value: 'WORD' },
  { id: 3, label: 'PDF', value: 'PDF' },
];
export const MEMBER_ROLES = [
  { id: 1, label: 'Administrator', value: 'ADMIN' },
  { id: 2, label: 'Vendor', value: 'VENDOR' },
  { id: 3, label: 'Member', value: 'MEMBER' },
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
  { id: 1, label: 'BYTE', value: 'Bytes' },
  { id: 2, label: 'KB', value: 'KiloBytes' },
  { id: 3, label: 'MB', value: 'MegaBytes' },
  { id: 4, label: 'GB', value: 'KiloBGigaBytesytes' },
];
export const DATA_DURATION = [
  { id: 1, label: 'Day', value: '1' },
  { id: 2, label: 'Month', value: '30' },
];
export const SIM_TYPES = [
  { id: 1, label: 'Unknown', value: 'unknown' },
  { id: 2, label: 'Test', value: 'test' },
  { id: 3, label: 'Operator Data', value: 'operator_data' },
  { id: 4, label: 'Ukama Data', value: 'ukama_data' },
];
const NODE_TYPE = [
  {
    id: 1,
    label: 'Home',
    value: 'HOME',
  },
  {
    id: 3,
    label: 'Tower',
    value: 'TOWER',
  },
  {
    id: 2,
    label: 'Amplifier',
    value: 'AMPLIFIER',
  },
];
const LANGUAGE_OPTIONS = [
  { id: 1, label: '🇺🇸  English, US', value: 'en' },
  { id: 2, label: '🇫🇷  French, France', value: 'fr' },
];

const SIDEBAR_MENU1 = [
  { id: '1', title: 'Home', Icon: HomeIcon, route: '/' },
  { id: '2', title: 'Nodes', Icon: RouterIcon, route: '/nodes' },
  { id: '3', title: 'Users', Icon: UsersIcon, route: '/users' },
  { id: '4', title: 'Billing', Icon: BillingIcon, route: '/billing' },
];
const STATS_OPTIONS = [
  { id: 1, label: 'Connected', value: 'Connected' },
  { id: 2, label: 'Device uptime', value: 'Device uptime' },
  { id: 3, label: 'Throughput', value: 'Throughput' },
];
const STATS_PERIOD = [
  { id: '1', label: 'DAY' },
  { id: '2', label: 'WEEK' },
  { id: '3', label: 'MONTH' },
  { id: '4', label: 'LIVE' },
];

const SIDEBAR_MENU2 = [
  { id: '5', title: 'Integrations', Icon: AppsIcon, route: '/store' },
];
const MONTH_FILTER = [
  { id: 1, label: 'January ', value: 'JANUARY ' },
  { id: 2, label: 'February', value: 'FEBRUARY' },
  { id: 3, label: 'March', value: 'MARCH ' },
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

const TIME_FILTER = [
  { id: 1, label: 'Today', value: 'TODAY' },
  { id: 2, label: 'This week', value: 'WEEK' },
  { id: 3, label: 'Month', value: 'MONTH' },
  { id: 4, label: 'Total', value: 'TOTAL' },
];

export const SUBSCRIBER_TABLE_COLUMNS: ColumnsWithOptions[] = [
  { id: 'name', label: 'Name', minWidth: 160 },
  { id: 'network', label: 'network', minWidth: 180 },
  { id: 'dataUsage', label: 'Data Usage', minWidth: 140 },
  { id: 'dataPlan', label: 'Data Plan', minWidth: 140 },
  { id: 'actions', label: 'Actions', align: 'right', minWidth: 80 },
];
export const SUBSCRIBER_TABLE_MENU: MenuItemType[] = [
  {
    id: 1,
    Icon: null,
    title: 'Edit subscriber',
    route: 'edit-sub',
  },
  { id: 2, Icon: null, title: 'Top up data', route: 'top-up-data' },
  { id: 3, Icon: null, title: 'Pause service', route: 'pause-service' },
  { id: 4, Icon: null, title: 'Delete subscriber', route: 'delete-sub' },
];
export const NODE_TABLE_COLUMNS: ColumnsWithOptions[] = [
  { id: 'name', label: 'Name', minWidth: 160 },
  { id: 'type', label: 'Type', minWidth: 180 },
  { id: 'state', label: 'State', minWidth: 140 },
  { id: 'network', label: 'Network', minWidth: 140 },
  { id: 'actions', label: 'Actions', align: 'right', minWidth: 80 },
];
export const NODE_TABLE_MENU: MenuItemType[] = [
  {
    id: 1,
    Icon: null,
    title: 'Edit node',
    route: 'edit-node',
  },
  { id: 2, Icon: null, title: 'Detach node', route: 'detach-node' },
  { id: 3, Icon: null, title: 'Delete node', route: 'delete-node' },
];
export const MANAGE_MENU_LIST = [
  {
    id: 'manage-members',
    name: 'Manage members',
    path: '/manage/members',
    icon: PeopleAltIcon,
  },
  {
    id: 'manage-sim',
    name: 'Manage SIM pool',
    path: '/manage/sim',
    icon: SimCardIcon,
  },
  {
    id: 'manage-node',
    name: 'Manage node pool',
    path: '/manage/node',
    icon: AccountTreeIcon,
  },
  {
    id: 'manage-data-plan',
    name: 'Manage data plans',
    path: '/manage/data-plans',
    icon: PermDataSettingIcon,
  },
];
export const MANAGE_TABLE_COLUMN: ColumnsWithOptions[] = [
  { id: 'name', label: 'Name', minWidth: 160 },
  { id: 'role', label: 'role', minWidth: 180 },
  { id: 'email', label: 'Email', minWidth: 140 },
];
export const MANAGE_SIM_POOL_COLUMN: ColumnsWithOptions[] = [
  { id: 'iccid', label: 'ICCID', minWidth: 160 },
  { id: 'simType', label: 'Type', minWidth: 180 },
  { id: 'qrCode', label: 'QR Code', minWidth: 140 },
  { id: 'isAllocated', label: 'Status', minWidth: 140 },
];
export const MANAGE_NODE_POOL_COLUMN: ColumnsWithOptions[] = [
  { id: 'name', label: 'Node #', minWidth: 160 },
  { id: 'type', label: 'Type', minWidth: 180 },
  { id: 'network', label: 'Network', minWidth: 180 },
  { id: 'created_at', label: 'Date claimed', minWidth: 140 },
];

const BASIC_MENU_ACTIONS: MenuItemType[] = [
  { id: 1, Icon: EditIcon, title: 'Edit', route: 'edit' },
  {
    id: 2,
    Icon: DeleteIcon,
    title: 'Delete',
    route: 'delete',
  },
  {
    id: 3,
    Icon: UpdateIcon,
    title: 'Update available',
    route: 'update',
  },
];

const DEACTIVATE_EDIT_ACTION_MENU: MenuItemType[] = [
  {
    id: 1,
    Icon: DeleteIcon,
    title: 'Deactivate',
    route: 'deactivate',
  },
  { id: 2, Icon: EditIcon, title: 'Edit', route: 'edit' },
];

const UserActivation = {
  title: 'Add Users',
  subTitle: `Install Ukama eSIMs on all desired devices to start using the network. For remote installation, you can share digital instructions below.`,
};

const BillingTabs = [
  { id: 0, label: 'CURRENT BILL', value: '1' },
  { id: 1, label: 'BILLING HISTORY', value: '2' },
];

const DEFAULT_RESPONSE = {
  error: null,
  loading: false,
  response: null,
};

const ROAMING_SELECT = [
  {
    id: 1,
    value: 'all',
    text: 'CHANGE ROAMING FOR INDIVIDAL SIMS',
  },
  {
    id: 2,
    value: 'esim1',
    text: 'ESIM 1',
  },
  {
    id: 3,
    value: 'esim2',
    text: 'ESIM 2',
  },
];

const RF_NODES = [
  {
    id: 1,
    value: 'off',
    text: 'TURN OFF RF FOR INDIVIDAL NODES',
  },
  {
    id: 2,
    value: 'on',
    text: 'TURN ON RF FOR INDIVIDAL NODES',
  },
];

const TooltipsText = {
  TRX: 'TRX Tooltip text',
  COM: 'COM Tooltip text',
  TRX_ALERT: 'TRX ALERT tooltip text',
  COM_ALERT: 'COM ALERT tooltip text',
  ATTACHED: 'Attached Tooltip text',
  ACTIVE: 'Active Tooltip text',
  DL: 'DL Tooltip text',
  UL: 'UL Tooltip text',
  RRCCNX: 'RRCCNX Tooltip text',
  ERAB: 'ERAB Tooltip text',
  RLS: 'RLS Tooltip text',
  MTRX: 'MTRX Tooltip text',
  MCOM: 'MCOM Tooltip text',
  CPUTRX: 'CPUTRX Tooltip text',
  CPUCOM: 'CPUCOM Tooltip text',
  DISKTRX: 'DISKTRX Tooltip text',
  DISKCOM: 'DISKCOM Tooltip text',
  POWER: 'POWER Tooltip text',
  TXPOWER: 'TXPOWER Tooltip text',
  RXPOWER: 'RXPOWER Tooltip text',
  PAPOWER: 'PAPOWER Tooltip text',
};

const NodePageTabs = [
  { id: 'node-tab-0', label: 'Overview', value: 0 },
  { id: 'node-tab-1', label: 'Network', value: 1 },
  { id: 'node-tab-2', label: 'Resources', value: 2 },
  { id: 'node-tab-3', label: 'Radio', value: 3 },
  { id: 'node-tab-4', label: 'Software', value: 4 },
  { id: 'node-tab-5', label: 'Schematic', value: 5 },
];

const NodeResourcesTabConfigure: any = {
  hnode: [
    { name: 'MEMORY-TRX', show: true, id: 'memory_trx_used' },
    { name: 'none', show: false, id: 'memory_com_used' },
    { name: 'CPU-TRX', show: true, id: 'cpu_trx_usage' },
    { name: 'none', show: false, id: 'cpu_com_usage' },
    { name: 'DISK-TRX', show: true, id: 'disk_trx_used' },
    { name: 'none', show: false, id: 'disk_com_used' },
    { name: 'none', show: false, id: 'power_level' },
  ],
  anode: [
    { name: 'MEMORY-TRX', show: true, id: 'memory_ctl_used' },
    { name: 'none', show: false, id: 'memory_com_used' },
    { name: 'CPU-TRX', show: true, id: 'cpu_ctl_used' },
    { name: 'none', show: false, id: 'cpu_com_usage' },
    { name: 'DISK-TRX', show: true, id: 'disk_ctl_used' },
    { name: 'none', show: false, id: 'disk_com_used' },
    { name: 'none', show: false, id: 'power_level' },
  ],
  tnode: [
    { name: 'MEMORY-TRX', show: true, id: 'memory_trx_used' },
    { name: 'MEMORY-COM', show: true, id: 'memory_com_used' },
    { name: 'CPU-TRX', show: true, id: 'cpu_trx_usage' },
    { name: 'CPU-COM', show: true, id: 'cpu_com_usage' },
    { name: 'DISK-TRX', show: true, id: 'disk_trx_used' },
    { name: 'DISK-COM', show: true, id: 'disk_com_used' },
    { name: 'POWER', show: true, id: 'power_level' },
  ],
};
const HealtChartsConfigure: any = {
  hnode: [
    { name: 'Uptime', show: true, id: 'uptime_trx' },
    { name: 'Temp. (TRX)', show: true, id: 'temperature_trx' },
    { name: 'Temp. (RFE)', show: true, id: 'temperature_rfe' },
    { name: 'none', show: false, id: '' },
    { name: 'Attached ', show: true, id: 'subscribers_attached' },
    { name: 'Active ', show: true, id: 'subscribers_active' },
  ],
  anode: [
    { name: 'Temp. (CTL)', show: true, id: 'temperature_ctl' },
    { name: 'Temp. (RFE)', show: true, id: 'temperature_rfe' },
    { name: 'none', show: false, id: '' },
  ],
  tnode: [
    { name: 'Temp. (TRX)', show: true, id: 'temperature_trx' },
    { name: 'Temp. (COM)', show: true, id: 'temperature_com' },
    { name: 'Attached ', show: true, id: 'subscribers_attached' },
    { name: 'Active ', show: true, id: 'subscribers_active' },
  ],
};

const MASK_BY_TYPE = {
  hnode: '{uk- }######{ -hnode- }##{ - }####',
  anode: '{uk- }######{ -\\anode- }##{ - }}####',
  tnode: '{uk- }######{ -tnode- }##{ - }}####',
};

const MASK_PLACEHOLDERS = {
  hnode: 'uk- ______ -hnode- __ - ____',
  anode: 'uk- ______ -anode- __ - ____',
  tnode: 'uk- ______ -tnode- __ - ____',
};

const SIM_VARIANT = [
  { id: 2, title: 'eSIM', type: 'eSIM' },
  { id: 1, title: 'Physical SIM', type: 'Physical SIM' },
];

export { NodeAppLogs, NodeApps, SimCardData } from './stubData';

const NODE_IMAGES = {
  tnode:
    'https://ukama-site-assets.s3.amazonaws.com/images/ukama_tower_node.png',
  anode:
    'https://ukama-site-assets.s3.amazonaws.com/images/ukama_amplifier_node.png',
  hnode:
    'https://ukama-site-assets.s3.amazonaws.com/images/ukama_home_node.png',
};

const BILLING_PLANS = [
  {
    id: 1,
    label: 'Default roaming - $5/GB; prorated per MB',
    value: 'default',
  },
  { id: 2, label: 'Unlimited plan - ', value: 'unlimited' },
  { id: 3, label: 'Set plan - $', value: 'plan1' },
  {
    id: 4,
    label: 'No plan - disables roaming for all residents.',
    value: 'no_plan',
  },
];

const BillingDialogList = [
  {
    id: 0,
    title: 'Choose roaming plan',
    description:
      'Choose a roaming plan below, and it will apply to all residents that have roaming enabled. Your selection can always be changed later.',
  },
  {
    id: 1,
    title: 'Customize preferences',
    description: 'Monitor and budget data usage with these settings.',
  },
  {
    id: 2,
    title: 'Enter payment information',
    description: 'Enter your payment information',
  },
  {
    id: 3,
    title: 'Payment set up successfully ',
    description:
      'Your payment and preferences have been set up successfully! You can change your settings at any time.',
  },
];

export {
  APP_VERSION,
  BASIC_MENU_ACTIONS,
  BILLING_PLANS,
  BillingDialogList,
  BillingTabs,
  COPY_RIGHTS,
  DEACTIVATE_EDIT_ACTION_MENU,
  DEFAULT_RESPONSE,
  DRAWER_WIDTH,
  DataTableWithOptionColumns,
  HealtChartsConfigure,
  IPFY_URL,
  IP_API_BASE_URL,
  LANGUAGE_OPTIONS,
  MASK_BY_TYPE,
  MASK_PLACEHOLDERS,
  MONTH_FILTER,
  NODE_ACTIONS_BUTTONS,
  NODE_IMAGES,
  NODE_TYPE,
  NodePageTabs,
  NodeResourcesTabConfigure,
  RF_NODES,
  ROAMING_SELECT,
  SETTING_MENU,
  SIDEBAR_MENU1,
  SIDEBAR_MENU2,
  SIM_VARIANT,
  STATS_OPTIONS,
  STATS_PERIOD,
  TABLE_EXPORT_OPTIONS,
  TIME_FILTER,
  TooltipsText,
  UserActivation,
};
