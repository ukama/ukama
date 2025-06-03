/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import {
  Notification_Scope,
  Notification_Type,
  SubscribersResDto,
} from '@/client/graphql/generated';
import {
  Graphs_Type,
  LatestMetricSubRes,
} from '@/client/graphql/generated/subscriptions';
import { EmotionCache } from '@emotion/react';
import { AppProps } from 'next/app';

export type MenuItemType = {
  Icon?: any;
  id: number;
  title: string;
  route: string;
  color?: string;
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
export interface FormValues {
  switch: string;
  power: string;
  backhaul: string;
  access: string;
  siteName: string;
  network: string;
  latitude: number;
  longitude: number;
  location: string;
}

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
export type StatusType = 'uptime' | 'battery' | 'signal';
export interface DialogStates {
  addSubscriber: boolean;
  subscriberDetails: boolean;
  topupData: boolean;
  confirmation: boolean;
  simDeleteConfirmation: boolean;
}

export interface SubscriberData {
  details: any;
  simList: any[];
  packageHistories: any[];
  topUpSubscriberName: string;
  dataUsageForSim: string;
}

export interface OperationData {
  deletedSubscriber: {
    id: string;
    name: string;
  };
  simToDelete: {
    id: string;
    iccid: string;
    isLastSim: boolean;
  };
}

export interface UIState {
  search: string;
  subscribers: SubscribersResDto;
}

export const INITIAL_SUBSCRIBER_DIALOG_STATES: DialogStates = {
  addSubscriber: false,
  subscriberDetails: false,
  topupData: false,
  confirmation: false,
  simDeleteConfirmation: false,
};

export const INITIAL_SUBSCRIBER_DATA: SubscriberData = {
  details: undefined,
  simList: [],
  packageHistories: [],
  topUpSubscriberName: '',
  dataUsageForSim: '',
};

export const INITIAL_OPERATION_DATA: OperationData = {
  deletedSubscriber: { id: '', name: '' },
  simToDelete: { id: '', iccid: '', isLastSim: false },
};

export const INITIAL_SUBSCRIBER_UI_STATE: UIState = {
  search: '',
  subscribers: { subscribers: [] },
};

export interface StyleOutput {
  color: string;
  icon: JSX.Element;
}
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
  options?: {
    isSortable?: boolean;
  };
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
export interface SubscriberDetailsType {
  name: string;
  email: string;
  simIccid: string;
  plan: string;
}
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
  country: string;
  currency: string;
};

export type TSnackbarMessage = {
  id: string;
  message: string;
  type: string;
  show: boolean;
};

export type TSiteForm = {
  switch: string;
  power: string;
  access: string;
  backhaul: string;
  address: string;
  spectrum: string;
  siteName: string;
  latitude: number;
  longitude: number;
  network: string;
};

export type TCoordinates = {
  lat: number | null;
  lng: number | null;
};

export type TEnv = {
  APP_URL: string;
  SIM_TYPE: string;
  METRIC_URL: string;
  API_GW_URL: string;
  AUTH_APP_URL: string;
  MAP_BOX_TOKEN: string;
  METRIC_WEBSOCKET_URL: string;
  STRIPE_PK: string;
};

interface NotificationSubscription {
  createdAt: string;
  description: string;
  id: string;
  isRead: boolean;
  scope: Notification_Scope;
  title: string;
  type: Notification_Type;
  eventKey: string;
  resourceId: string;
  redirect: {
    action: string;
    title: string;
  };
}

interface Data {
  notificationSubscription: NotificationSubscription;
}

export interface TNotificationResDto {
  data: Data;
}

interface MetricsSubData {
  getMetricByTabSub: LatestMetricSubRes;
  getMetricStatSub: LatestMetricSubRes;
  getSiteMetricByTabSub: LatestMetricSubRes;
  getSiteMetricStatSub: LatestMetricSubRes;
  getMetricBySiteSub: LatestMetricSubRes;
}
export interface ActiveView {
  graphType: Graphs_Type;
  kpi: KPIType;
}

export interface TMetricResDto {
  data: MetricsSubData;
}
export type KPIType =
  | 'node'
  | 'solar'
  | 'controller'
  | 'battery'
  | 'backhaul'
  | 'switch';

export type TNodePoolData = {
  id: string;
  type: string;
  site: string;
  state: string;
  network: string;
  createdAt: string;
  connectivity: string;
};

export enum DataUnitType {
  Bytes = 'Bytes',
  KiloBytes = 'KiloBytes',
  MegaBytes = 'MegaBytes',
  GigaBytes = 'GigaBytes',
}

export enum DurationType {
  Day = 1,
  Month = 30,
}

export type CreatePlanType = {
  id: string;
  name: string;
  country: string;
  currency: string;
  duration: string;
  dataUnit: DataUnitType;
  amount: number | undefined;
  dataVolume: number | undefined;
};
