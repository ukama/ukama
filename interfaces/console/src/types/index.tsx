/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { EmotionCache } from '@emotion/react';
import { AppProps } from 'next/app';

export type MenuItemType = {
  Icon?: any;
  id: number;
  title: string;
  route: string;
};
export type UserSettingsMenuType = {
  id: number;
  label: string;
};
export type StatsItemType = {
  id: number;
  label: string;
  value: string;
};

export type Site = {
  name: string;
  health: 'online' | 'offline';
  duration: string;
};

export type NodeAppDetailsTypes = {
  id: number;
  cpu: number;
  memory: number;
  nodeAppName: string;
};
export type StatsPeriodItemType = {
  id: string;
  label: string;
};
export type HeaderMenuItemType = {
  id: string;
  Icon: any;
  title: string;
};
export type SelectItemType = {
  id: number | string;
  label: string;
  value: string;
};
export type BillingType = {
  id?: number;
  value: string;
  label: string;
};
export type ExportOptionsType = {
  id?: number | string;
  value: string | number;
  label: string;
};
export type BillingTableHeaderOptionsType = {
  id?: number;
  label: string;
};
export type CurrentBillType = {
  id: number;
  name: string;
  rate: string;
  subTotal: number;
  dataUsage: string;
};
export type PaymentMethodType = {
  id?: number;
  card_experintionDetails: string;
};
export type SVGType = {
  color?: string;
  color2?: string;
  width?: string;
  height?: string;
};
export type SettingsMenuTypes = {
  id: number;
  title: string;
};
export type ColumnsWithOptions = {
  id: any;
  label: string;
  minWidth?: number;
  align?: 'right';
};
export type UserSearchFormType = {
  text: string;
};
export type SimActivateFormType = {
  email: string;
  phone: string;
  number: string;
  lastName: string;
  firstName: string;
};
export type UserActivateFormType = {
  nodeName: string;
  serialNumber: string;
  securityCode: string;
};

export type TVariant = 'small' | 'medium' | 'large';

export type TObject = { [key: string]: boolean | string | number };

export type Record = {
  [key: string]: string;
};

export type TSnackMessage = {
  id: string;
  message: string;
  type: string;
  show: boolean;
};

export interface MyAppProps extends AppProps {
  emotionCache?: EmotionCache;
}

export type TNodeSiteChild = {};

export type TNodeSiteTree = {
  id: string;
  name: string;
  nodeId: string;
  nodeType: string;
  nodeName: string;
};

export type TAddSubscriberData = {
  name: string;
  email: string;
  phone: string;
  simType: string;
  roamingStatus: boolean;
  iccid: string;
  plan: string;
};

export type TNetwork = {
  id: string;
  name: string;
};

export type TUser = {
  id: string;
  name: string;
  email: string;
  role: string;
  orgId: string;
  orgName: string;
};

export type TSnackbarMessage = {
  id: string;
  message: string;
  type: string;
  show: boolean;
};
