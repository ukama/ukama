/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

export * from './nodes';
export * from './sites';
export * from './onboarding';

import { Role_Type, Sim_Types } from '@/client/graphql/generated';
import { ColumnsWithOptions, MenuItemType } from '@/types';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';

export const DRAWER_WIDTH = 200;
export const APP_VERSION = 'v0.0.1';
export const COPY_RIGHTS = 'Copyright © Ukama Inc.';
export const REFRESH_INTERVAL = 30000;
export const KPI_PLACEHOLDER_VALUE = '-';

export const SETTING_MENU = [
  { id: 'personal-settings', name: 'My Account' },
  { id: 'network-settings', name: 'Network' },
  { id: 'billing', name: 'Billing' },
  { id: 'appearance', name: 'Appearance' },
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

export { NodeApps } from './stubData';
