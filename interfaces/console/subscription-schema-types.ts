export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export enum Graphs_Type {
  Battery = 'BATTERY',
  Controller = 'CONTROLLER',
  DataUsage = 'DATA_USAGE',
  Home = 'HOME',
  MainBackhaul = 'MAIN_BACKHAUL',
  NetworkBackhaul = 'NETWORK_BACKHAUL',
  NetworkCellular = 'NETWORK_CELLULAR',
  NodeHealth = 'NODE_HEALTH',
  Radio = 'RADIO',
  Resources = 'RESOURCES',
  Site = 'SITE',
  Solar = 'SOLAR',
  Subscribers = 'SUBSCRIBERS',
  Switch = 'SWITCH'
}

export type GetMetricBySiteInput = {
  from: Scalars['Float']['input'];
  nodeIds?: InputMaybe<Array<Scalars['String']['input']>>;
  orgName: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
  step?: Scalars['Float']['input'];
  to: Scalars['Float']['input'];
  type: Graphs_Type;
  userId: Scalars['String']['input'];
  withSubscription?: Scalars['Boolean']['input'];
};

export type GetMetricByTabInput = {
  from: Scalars['Float']['input'];
  networkId?: InputMaybe<Scalars['String']['input']>;
  nodeId?: InputMaybe<Scalars['String']['input']>;
  orgName: Scalars['String']['input'];
  siteId?: InputMaybe<Scalars['String']['input']>;
  step?: Scalars['Float']['input'];
  to: Scalars['Float']['input'];
  type: Graphs_Type;
  userId: Scalars['String']['input'];
  withSubscription?: Scalars['Boolean']['input'];
};

export type GetMetricsSiteStatInput = {
  from: Scalars['Float']['input'];
  nodeIds?: InputMaybe<Array<Scalars['String']['input']>>;
  operation?: InputMaybe<Scalars['String']['input']>;
  orgName: Scalars['String']['input'];
  siteIds?: InputMaybe<Array<Scalars['String']['input']>>;
  step?: Scalars['Float']['input'];
  to: Scalars['Float']['input'];
  type: Stats_Type;
  userId?: InputMaybe<Scalars['String']['input']>;
  withSubscription?: Scalars['Boolean']['input'];
};

export type GetMetricsStatInput = {
  from: Scalars['Float']['input'];
  networkId?: InputMaybe<Scalars['String']['input']>;
  nodeId?: InputMaybe<Scalars['String']['input']>;
  operation?: InputMaybe<Scalars['String']['input']>;
  orgName: Scalars['String']['input'];
  siteId?: InputMaybe<Scalars['String']['input']>;
  step?: Scalars['Float']['input'];
  to?: InputMaybe<Scalars['Float']['input']>;
  type: Stats_Type;
  userId?: InputMaybe<Scalars['String']['input']>;
  withSubscription?: Scalars['Boolean']['input'];
};

export type LatestMetricSubRes = {
  __typename?: 'LatestMetricSubRes';
  dataPlanId?: Maybe<Scalars['String']['output']>;
  format?: Maybe<Scalars['String']['output']>;
  msg: Scalars['String']['output'];
  networkId?: Maybe<Scalars['String']['output']>;
  nodeId: Scalars['String']['output'];
  packageId?: Maybe<Scalars['String']['output']>;
  siteId: Scalars['String']['output'];
  success: Scalars['Boolean']['output'];
  threshold?: Maybe<MetricThreshold>;
  tickInterval?: Maybe<Scalars['Float']['output']>;
  tickPositions?: Maybe<Array<Scalars['Float']['output']>>;
  type: Scalars['String']['output'];
  unit?: Maybe<Scalars['String']['output']>;
  value: Array<Scalars['Float']['output']>;
};

export type MetricRes = {
  __typename?: 'MetricRes';
  dataPlanId?: Maybe<Scalars['String']['output']>;
  format?: Maybe<Scalars['String']['output']>;
  msg: Scalars['String']['output'];
  networkId?: Maybe<Scalars['String']['output']>;
  nodeId?: Maybe<Scalars['String']['output']>;
  packageId?: Maybe<Scalars['String']['output']>;
  siteId?: Maybe<Scalars['String']['output']>;
  success: Scalars['Boolean']['output'];
  threshold?: Maybe<MetricThreshold>;
  tickInterval?: Maybe<Scalars['Float']['output']>;
  tickPositions?: Maybe<Array<Scalars['Float']['output']>>;
  type: Scalars['String']['output'];
  unit?: Maybe<Scalars['String']['output']>;
  values: Array<Array<Scalars['Float']['output']>>;
};

export type MetricStateRes = {
  __typename?: 'MetricStateRes';
  dataPlanId?: Maybe<Scalars['String']['output']>;
  format?: Maybe<Scalars['String']['output']>;
  msg: Scalars['String']['output'];
  networkId?: Maybe<Scalars['String']['output']>;
  nodeId: Scalars['String']['output'];
  packageId?: Maybe<Scalars['String']['output']>;
  siteId?: Maybe<Scalars['String']['output']>;
  success: Scalars['Boolean']['output'];
  threshold?: Maybe<MetricThreshold>;
  tickInterval?: Maybe<Scalars['Float']['output']>;
  tickPositions?: Maybe<Array<Scalars['Float']['output']>>;
  type: Scalars['String']['output'];
  unit?: Maybe<Scalars['String']['output']>;
  value: Scalars['Float']['output'];
};

export type MetricThreshold = {
  __typename?: 'MetricThreshold';
  max: Scalars['Float']['output'];
  min: Scalars['Float']['output'];
  normal: Scalars['Float']['output'];
};

export type MetricsRes = {
  __typename?: 'MetricsRes';
  metrics: Array<MetricRes>;
};

export type MetricsStateRes = {
  __typename?: 'MetricsStateRes';
  metrics: Array<MetricStateRes>;
};

export enum Notification_Scope {
  ScopeInvalid = 'SCOPE_INVALID',
  ScopeNetwork = 'SCOPE_NETWORK',
  ScopeNetworks = 'SCOPE_NETWORKS',
  ScopeNode = 'SCOPE_NODE',
  ScopeOrg = 'SCOPE_ORG',
  ScopeOwner = 'SCOPE_OWNER',
  ScopeSite = 'SCOPE_SITE',
  ScopeSites = 'SCOPE_SITES',
  ScopeSubscriber = 'SCOPE_SUBSCRIBER',
  ScopeSubscribers = 'SCOPE_SUBSCRIBERS',
  ScopeUser = 'SCOPE_USER',
  ScopeUsers = 'SCOPE_USERS'
}

export enum Notification_Type {
  TypeActionableCritical = 'TYPE_ACTIONABLE_CRITICAL',
  TypeActionableError = 'TYPE_ACTIONABLE_ERROR',
  TypeActionableInfo = 'TYPE_ACTIONABLE_INFO',
  TypeActionableWarning = 'TYPE_ACTIONABLE_WARNING',
  TypeCritical = 'TYPE_CRITICAL',
  TypeError = 'TYPE_ERROR',
  TypeInfo = 'TYPE_INFO',
  TypeInvalid = 'TYPE_INVALID',
  TypeWarning = 'TYPE_WARNING'
}

export type NotificationRedirect = {
  __typename?: 'NotificationRedirect';
  action: Scalars['String']['output'];
  title: Scalars['String']['output'];
};

export type NotificationsRes = {
  __typename?: 'NotificationsRes';
  notifications: Array<NotificationsResDto>;
};

export type NotificationsResDto = {
  __typename?: 'NotificationsResDto';
  createdAt: Scalars['String']['output'];
  description: Scalars['String']['output'];
  eventKey: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isRead: Scalars['Boolean']['output'];
  redirect: NotificationRedirect;
  resourceId: Scalars['String']['output'];
  scope: Notification_Scope;
  title: Scalars['String']['output'];
  type: Notification_Type;
};

export type Query = {
  __typename?: 'Query';
  getMetricBySite: MetricsRes;
  getMetricByTab: MetricsRes;
  getMetricsStat: MetricsStateRes;
  getNotifications: NotificationsRes;
  getSiteStat: MetricsStateRes;
};


export type QueryGetMetricBySiteArgs = {
  data: GetMetricBySiteInput;
};


export type QueryGetMetricByTabArgs = {
  data: GetMetricByTabInput;
};


export type QueryGetMetricsStatArgs = {
  data: GetMetricsStatInput;
};


export type QueryGetNotificationsArgs = {
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  orgName: Scalars['String']['input'];
  role: Scalars['String']['input'];
  startTimestamp: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};


export type QueryGetSiteStatArgs = {
  data: GetMetricsSiteStatInput;
};

export enum Stats_Type {
  AllNode = 'ALL_NODE',
  Battery = 'BATTERY',
  DataUsage = 'DATA_USAGE',
  Home = 'HOME',
  MainBackhaul = 'MAIN_BACKHAUL',
  Network = 'NETWORK',
  Overview = 'OVERVIEW',
  Radio = 'RADIO',
  Resources = 'RESOURCES',
  Site = 'SITE'
}

export type SubMetricsStatInput = {
  from: Scalars['Float']['input'];
  networkId?: InputMaybe<Scalars['String']['input']>;
  nodeId?: InputMaybe<Scalars['String']['input']>;
  orgName: Scalars['String']['input'];
  type: Stats_Type;
  userId: Scalars['String']['input'];
};

export type SubSiteMetricByTabInput = {
  from: Scalars['Float']['input'];
  orgName: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
  type: Graphs_Type;
  userId: Scalars['String']['input'];
};

export type SubSiteMetricsStatInput = {
  from: Scalars['Float']['input'];
  nodeIds?: InputMaybe<Array<Scalars['String']['input']>>;
  orgName: Scalars['String']['input'];
  siteIds?: InputMaybe<Array<Scalars['String']['input']>>;
  type: Stats_Type;
  userId: Scalars['String']['input'];
};

export type Subscription = {
  __typename?: 'Subscription';
  getMetricStatSub: LatestMetricSubRes;
  getSiteMetricByTabSub: LatestMetricSubRes;
  getSiteMetricStatSub: LatestMetricSubRes;
  notificationSubscription: NotificationsResDto;
};


export type SubscriptionGetMetricStatSubArgs = {
  data: SubMetricsStatInput;
};


export type SubscriptionGetSiteMetricByTabSubArgs = {
  data: SubSiteMetricByTabInput;
};


export type SubscriptionGetSiteMetricStatSubArgs = {
  data: SubSiteMetricsStatInput;
};


export type SubscriptionNotificationSubscriptionArgs = {
  networkId: Scalars['String']['input'];
  orgId: Scalars['String']['input'];
  orgName: Scalars['String']['input'];
  role: Scalars['String']['input'];
  startTimestamp: Scalars['String']['input'];
  subscriberId: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};
