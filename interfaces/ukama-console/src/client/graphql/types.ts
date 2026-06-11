export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type ActivityItemDto = {
  __typename?: 'ActivityItemDto';
  description?: Maybe<Scalars['String']['output']>;
  occurredAt?: Maybe<Scalars['String']['output']>;
  routingKey?: Maybe<Scalars['String']['output']>;
};

export type AddMemberInputDto = {
  role: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type AddNetworkInputDto = {
  budget?: InputMaybe<Scalars['Float']['input']>;
  countries?: InputMaybe<Array<Scalars['String']['input']>>;
  isDefault?: Scalars['Boolean']['input'];
  name: Scalars['String']['input'];
  networks?: InputMaybe<Array<Scalars['String']['input']>>;
};

export type AddNodeInput = {
  id: Scalars['String']['input'];
  name: Scalars['String']['input'];
};

export type AddNodeToSiteInput = {
  networkId: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
};

export type AddPackagSimResDto = {
  __typename?: 'AddPackagSimResDto';
  error?: Maybe<Scalars['String']['output']>;
  packageId?: Maybe<Scalars['String']['output']>;
  success?: Maybe<Scalars['Boolean']['output']>;
};

export type AddPackageInputDto = {
  amount: Scalars['Float']['input'];
  country: Scalars['String']['input'];
  currency: Scalars['String']['input'];
  dataUnit: Scalars['String']['input'];
  dataVolume: Scalars['Int']['input'];
  duration: Scalars['Int']['input'];
  name: Scalars['String']['input'];
};

export type AddPackagesSimResDto = {
  __typename?: 'AddPackagesSimResDto';
  packages: Array<AddPackagSimResDto>;
};

export type AddPackagesToSimInputDto = {
  packages: Array<PackagesToSimInputDto>;
  sim_id: Scalars['String']['input'];
};

export type AddSiteInputDto = {
  access_id: Scalars['String']['input'];
  backhaul_id: Scalars['String']['input'];
  install_date: Scalars['String']['input'];
  latitude: Scalars['String']['input'];
  location: Scalars['String']['input'];
  longitude: Scalars['String']['input'];
  name: Scalars['String']['input'];
  network_id: Scalars['String']['input'];
  power_id: Scalars['String']['input'];
  spectrum_id: Scalars['String']['input'];
  switch_id: Scalars['String']['input'];
};

export type AlarmRowDto = {
  __typename?: 'AlarmRowDto';
  alarmId?: Maybe<Scalars['String']['output']>;
  closedAt?: Maybe<Scalars['String']['output']>;
  customersAffected: Scalars['Int']['output'];
  description?: Maybe<Scalars['String']['output']>;
  openedAt?: Maybe<Scalars['String']['output']>;
  recommendedAction?: Maybe<Scalars['String']['output']>;
  resourceId?: Maybe<Scalars['String']['output']>;
  resourceType?: Maybe<Scalars['String']['output']>;
  revenueAtRisk: Scalars['Float']['output'];
  severity?: Maybe<Scalars['String']['output']>;
  state?: Maybe<Scalars['String']['output']>;
};

export type AlertsSection = {
  __typename?: 'AlertsSection';
  error?: Maybe<SectionError>;
  notifications?: Maybe<Array<NotificationsDto>>;
};

export type AllocateSimApiDto = {
  __typename?: 'AllocateSimAPIDto';
  allocated_at: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi?: Maybe<Scalars['String']['output']>;
  is_physical: Scalars['Boolean']['output'];
  msisdn: Scalars['String']['output'];
  network_id: Scalars['String']['output'];
  package: SimAllocatePackageDto;
  status: Scalars['String']['output'];
  subscriber_id: Scalars['String']['output'];
  sync_status: Scalars['String']['output'];
  traffic_policy: Scalars['Float']['output'];
  type: Scalars['String']['output'];
};

export type AllocateSimInputDto = {
  iccid?: InputMaybe<Scalars['String']['input']>;
  network_id: Scalars['String']['input'];
  package_id: Scalars['String']['input'];
  sim_type: Scalars['String']['input'];
  subscriber_id: Scalars['String']['input'];
  traffic_policy: Scalars['Float']['input'];
};

export type AnalyticsNodeInput = {
  from?: InputMaybe<Scalars['String']['input']>;
  nodeId: Scalars['String']['input'];
  period?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
  to?: InputMaybe<Scalars['String']['input']>;
};

export type AnalyticsSiteInput = {
  from?: InputMaybe<Scalars['String']['input']>;
  period?: InputMaybe<Scalars['String']['input']>;
  siteId: Scalars['String']['input'];
  timezone?: InputMaybe<Scalars['String']['input']>;
  to?: InputMaybe<Scalars['String']['input']>;
};

export type AnalyticsWindowInput = {
  from?: InputMaybe<Scalars['String']['input']>;
  metric?: InputMaybe<Scalars['String']['input']>;
  networkId?: InputMaybe<Scalars['String']['input']>;
  nodeId?: InputMaybe<Scalars['String']['input']>;
  page?: InputMaybe<Scalars['Int']['input']>;
  pageSize?: InputMaybe<Scalars['Int']['input']>;
  period?: InputMaybe<Scalars['String']['input']>;
  query?: InputMaybe<Scalars['String']['input']>;
  severity?: InputMaybe<Scalars['String']['input']>;
  siteId?: InputMaybe<Scalars['String']['input']>;
  state?: InputMaybe<Scalars['String']['input']>;
  status?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
  to?: InputMaybe<Scalars['String']['input']>;
};

export type App = {
  __typename?: 'App';
  name: Scalars['String']['output'];
  resource?: Maybe<AppResource>;
  status: Scalars['String']['output'];
  tag: Scalars['String']['output'];
  version: Scalars['String']['output'];
};

export type AppChangeLog = {
  __typename?: 'AppChangeLog';
  date: Scalars['Float']['output'];
  version: Scalars['String']['output'];
};

export type AppChangeLogs = {
  __typename?: 'AppChangeLogs';
  logs: Array<AppChangeLog>;
  type: NodeTypeEnum;
};

export type AppResource = {
  __typename?: 'AppResource';
  cpuPercent: Scalars['Float']['output'];
  diskReadBytes: Scalars['Float']['output'];
  diskWriteBytes: Scalars['Float']['output'];
  memoryRssKb: Scalars['Float']['output'];
};

export type Apps = {
  __typename?: 'Apps';
  apps: Array<App>;
};

export type AttachNodeInput = {
  anodel: Scalars['String']['input'];
  anoder: Scalars['String']['input'];
  parentNode: Scalars['String']['input'];
};

export type AttachedNodes = {
  __typename?: 'AttachedNodes';
  id: Scalars['String']['output'];
  latitude: Scalars['String']['output'];
  longitude: Scalars['String']['output'];
  name: Scalars['String']['output'];
  site: NodeSite;
  status: NodeStatus;
  type: NodeTypeEnum;
};

export type BalanceSection = {
  __typename?: 'BalanceSection';
  error?: Maybe<SectionError>;
  latestUnpaidPeriod?: Maybe<Scalars['String']['output']>;
  outstandingCount?: Maybe<Scalars['Int']['output']>;
};

export type BillingSummaryDto = {
  __typename?: 'BillingSummaryDto';
  invoices: Array<InvoiceRowDto>;
  kpis: Array<KpiDto>;
  lastInvoiceDate?: Maybe<Scalars['String']['output']>;
  meta?: Maybe<MetaDto>;
};

export type BusinessHomeDto = {
  __typename?: 'BusinessHomeDto';
  kpis: Array<KpiDto>;
  recentActivity: Array<ActivityItemDto>;
  sites: Array<SiteSummaryDto>;
  topPackages: Array<NamedValueDto>;
};

export type BusinessSiteDto = {
  __typename?: 'BusinessSiteDto';
  kpis: Array<KpiDto>;
  revenueTrend?: Maybe<TimeSeriesDto>;
  site?: Maybe<BusinessSiteRowDto>;
};

export type BusinessSiteRowDto = {
  __typename?: 'BusinessSiteRowDto';
  customers: Scalars['Int']['output'];
  dataUsed: Scalars['Float']['output'];
  issue?: Maybe<Scalars['String']['output']>;
  latitude: Scalars['Float']['output'];
  longitude: Scalars['Float']['output'];
  name?: Maybe<Scalars['String']['output']>;
  revenue: Scalars['Float']['output'];
  revenueToday: Scalars['Float']['output'];
  siteId: Scalars['String']['output'];
  status?: Maybe<Scalars['String']['output']>;
  topPackage?: Maybe<Scalars['String']['output']>;
  uptime: Scalars['Float']['output'];
};

export type BusinessSitesDto = {
  __typename?: 'BusinessSitesDto';
  meta?: Maybe<MetaDto>;
  sites: Array<BusinessSiteRowDto>;
};

export type CBooleanResponse = {
  __typename?: 'CBooleanResponse';
  success: Scalars['Boolean']['output'];
};

export enum Component_Type {
  Access = 'access',
  All = 'all',
  Backhaul = 'backhaul',
  Power = 'power',
  Spectrum = 'spectrum',
  Switch = 'switch'
}

export type CategoryCountDto = {
  __typename?: 'CategoryCountDto';
  category: Scalars['String']['output'];
  count: Scalars['Int']['output'];
};

export type CommerceView = {
  __typename?: 'CommerceView';
  balance: BalanceSection;
  invoices: InvoicesSection;
  networkId?: Maybe<Scalars['String']['output']>;
  plans: PlanStatsSection;
  revenue: RevenueSection;
};


export type CommerceViewInvoicesArgs = {
  limit?: Scalars['Int']['input'];
};

export type Component = {
  __typename?: 'Component';
  componentId?: Maybe<Scalars['String']['output']>;
  componentName?: Maybe<Scalars['String']['output']>;
  elementType: Scalars['String']['output'];
};

export type ComponentDto = {
  __typename?: 'ComponentDto';
  category: Scalars['String']['output'];
  datasheetUrl: Scalars['String']['output'];
  description: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imageUrl: Scalars['String']['output'];
  inventoryId: Scalars['String']['output'];
  managed: Scalars['String']['output'];
  manufacturer: Scalars['String']['output'];
  partNumber: Scalars['String']['output'];
  specification: Scalars['String']['output'];
  type: Scalars['String']['output'];
  userId: Scalars['String']['output'];
  warranty: Scalars['Float']['output'];
};

export type ComponentStatsSection = {
  __typename?: 'ComponentStatsSection';
  byCategory?: Maybe<Array<CategoryCountDto>>;
  error?: Maybe<SectionError>;
  total?: Maybe<Scalars['Int']['output']>;
};

export type ComponentTypeInputDto = {
  category: Component_Type;
};

export type ComponentsResDto = {
  __typename?: 'ComponentsResDto';
  components: Array<ComponentDto>;
};

export type CountriesRes = {
  __typename?: 'CountriesRes';
  countries: Array<CountryDto>;
};

export type CountryDto = {
  __typename?: 'CountryDto';
  code: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type CreateInvitationInputDto = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  role: Role_Type;
};

export type CurrencyRes = {
  __typename?: 'CurrencyRes';
  code: Scalars['String']['output'];
  image: Scalars['String']['output'];
  symbol: Scalars['String']['output'];
};

export type CustomerByIdInput = {
  customerId: Scalars['String']['input'];
  from?: InputMaybe<Scalars['String']['input']>;
  period?: InputMaybe<Scalars['String']['input']>;
  timezone?: InputMaybe<Scalars['String']['input']>;
  to?: InputMaybe<Scalars['String']['input']>;
};

export type CustomerDetailDto = {
  __typename?: 'CustomerDetailDto';
  customer?: Maybe<CustomerRowDto>;
  kpis: Array<KpiDto>;
  packageHistory: Array<PackageIntervalDto>;
};

export type CustomerDto = {
  __typename?: 'CustomerDto';
  addressLine1?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  email?: Maybe<Scalars['String']['output']>;
  externalId: Scalars['String']['output'];
  legalName?: Maybe<Scalars['String']['output']>;
  legalNumber?: Maybe<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  phone?: Maybe<Scalars['String']['output']>;
  timezone?: Maybe<Scalars['String']['output']>;
  vatRate: Scalars['Float']['output'];
};

export type CustomerListDto = {
  __typename?: 'CustomerListDto';
  customers: Array<CustomerRowDto>;
  meta?: Maybe<MetaDto>;
};

export type CustomerOverviewDto = {
  __typename?: 'CustomerOverviewDto';
  kpis: Array<KpiDto>;
};

export type CustomerRowDto = {
  __typename?: 'CustomerRowDto';
  customerId: Scalars['String']['output'];
  dataUsage: Scalars['Float']['output'];
  email?: Maybe<Scalars['String']['output']>;
  lastSeen?: Maybe<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
  packageName?: Maybe<Scalars['String']['output']>;
  packageStatus?: Maybe<Scalars['String']['output']>;
  simIccid?: Maybe<Scalars['String']['output']>;
  simStatus?: Maybe<Scalars['String']['output']>;
  siteId?: Maybe<Scalars['String']['output']>;
  siteName?: Maybe<Scalars['String']['output']>;
  status?: Maybe<Scalars['String']['output']>;
};

export type CustomerSimsDto = {
  __typename?: 'CustomerSimsDto';
  meta?: Maybe<MetaDto>;
  sims: Array<SimRowDto>;
};

export type CustomerSupportDto = {
  __typename?: 'CustomerSupportDto';
  customer?: Maybe<CustomerRowDto>;
  escalationNeeded: Scalars['Boolean']['output'];
  likelyIssue?: Maybe<Scalars['String']['output']>;
  recentActivity: Array<ActivityItemDto>;
  recommendedAction?: Maybe<Scalars['String']['output']>;
  signals: Array<SupportSignalDto>;
};

export type DataPlan = {
  __typename?: 'DataPlan';
  elementType: Scalars['String']['output'];
  planId: Scalars['String']['output'];
  planName: Scalars['String']['output'];
};

export type DefaultMarkupHistoryDto = {
  __typename?: 'DefaultMarkupHistoryDto';
  Markup: Scalars['Float']['output'];
  createdAt: Scalars['String']['output'];
  deletedAt: Scalars['String']['output'];
};

export type DefaultMarkupHistoryResDto = {
  __typename?: 'DefaultMarkupHistoryResDto';
  markupRates?: Maybe<Array<DefaultMarkupHistoryDto>>;
};

export type DefaultMarkupInputDto = {
  markup: Scalars['Float']['input'];
};

export type DefaultMarkupResDto = {
  __typename?: 'DefaultMarkupResDto';
  markup: Scalars['Float']['output'];
};

export type DeleteInvitationResDto = {
  __typename?: 'DeleteInvitationResDto';
  id: Scalars['String']['output'];
};

export type DeleteNode = {
  __typename?: 'DeleteNode';
  id: Scalars['String']['output'];
};

export type DeleteSimInputDto = {
  simId: Scalars['String']['input'];
};

export type DeleteSimResDto = {
  __typename?: 'DeleteSimResDto';
  simId?: Maybe<Scalars['String']['output']>;
};

export type EventRowDto = {
  __typename?: 'EventRowDto';
  description?: Maybe<Scalars['String']['output']>;
  occurredAt?: Maybe<Scalars['String']['output']>;
  resourceId?: Maybe<Scalars['String']['output']>;
  resourceType?: Maybe<Scalars['String']['output']>;
  routingKey?: Maybe<Scalars['String']['output']>;
};

export type FeeDto = {
  __typename?: 'FeeDto';
  eventsCount: Scalars['String']['output'];
  item: ItemResDto;
  taxesAmountCents: Scalars['String']['output'];
  taxesPreciseAmount: Scalars['String']['output'];
  totalAmountCents: Scalars['String']['output'];
  totalAmountCurrency: Scalars['String']['output'];
  units: Scalars['Float']['output'];
};

export type GapSection = {
  __typename?: 'GapSection';
  error?: Maybe<SectionError>;
};

export type GetAppsInputDto = {
  appName?: InputMaybe<Scalars['String']['input']>;
  nodeId: Scalars['String']['input'];
};

export type GetHealthReportInputDto = {
  id: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  timeframe: Timeframe_Filter;
  timestamp: Scalars['String']['input'];
};

export type GetNodeLatestMetricInput = {
  nodeId: Scalars['String']['input'];
  type: Scalars['String']['input'];
};

export type GetNodesByStateInput = {
  connectivity: NodeConnectivityEnum;
  state: NodeStateEnum;
};

export type GetOperationInputDto = {
  id: Scalars['String']['input'];
};

export type GetPackagesForSimInputDto = {
  sim_id: Scalars['String']['input'];
};

export type GetPaymentsInputDto = {
  paymentMethod?: InputMaybe<Scalars['String']['input']>;
  status?: InputMaybe<Scalars['String']['input']>;
  type?: InputMaybe<Scalars['String']['input']>;
};

export type GetPdfReportUrlDto = {
  __typename?: 'GetPdfReportUrlDto';
  contentType: Scalars['String']['output'];
  downloadUrl: Scalars['String']['output'];
  filename: Scalars['String']['output'];
  id: Scalars['String']['output'];
};

export type GetReportDto = {
  __typename?: 'GetReportDto';
  report: ReportDto;
};

export type GetReportsDto = {
  __typename?: 'GetReportsDto';
  reports: Array<ReportDto>;
};

export type GetReportsInputDto = {
  count?: InputMaybe<Scalars['Float']['input']>;
  isPaid?: InputMaybe<Scalars['Boolean']['input']>;
  networkId?: InputMaybe<Scalars['String']['input']>;
  ownerId?: InputMaybe<Scalars['String']['input']>;
  ownerType: Scalars['String']['input'];
  report_type: Scalars['String']['input'];
  sort?: InputMaybe<Scalars['Boolean']['input']>;
};

export type GetResourceLockInputDto = {
  resourceKey: Scalars['String']['input'];
};

export type GetSimBySubscriberInputDto = {
  subscriberId: Scalars['String']['input'];
};

export type GetSimInputDto = {
  simId: Scalars['String']['input'];
};

export type GetSimPackagesDtoApi = {
  __typename?: 'GetSimPackagesDtoAPI';
  packages: Array<SimToPackagesDto>;
  sim_id: Scalars['String']['output'];
};

export type GetSimsInput = {
  status: Sim_Status;
  type: Sim_Types;
};

export type GetSoftwaresInput = {
  name: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  status: SoftwareStatusEnum;
};

export type HealthCappInfo = {
  __typename?: 'HealthCappInfo';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  resources: Array<HealthResourceInfo>;
  space: Scalars['String']['output'];
  status: Scalars['String']['output'];
  tag: Scalars['String']['output'];
};

export type HealthInfo = {
  __typename?: 'HealthInfo';
  capps: Array<HealthCappInfo>;
  id: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  system: Array<HealthSystemInfo>;
  timestamp: Scalars['String']['output'];
};

export type HealthResourceInfo = {
  __typename?: 'HealthResourceInfo';
  cappId: Scalars['String']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type HealthSection = {
  __typename?: 'HealthSection';
  error?: Maybe<SectionError>;
  health?: Maybe<HealthInfo>;
};

export type HealthSystemInfo = {
  __typename?: 'HealthSystemInfo';
  healthId: Scalars['String']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export enum Invitation_Status {
  InviteAccepted = 'INVITE_ACCEPTED',
  InviteDeclined = 'INVITE_DECLINED',
  InvitePending = 'INVITE_PENDING'
}

export type IdResponse = {
  __typename?: 'IdResponse';
  uuid: Scalars['String']['output'];
};

export type InventoryReadinessDto = {
  __typename?: 'InventoryReadinessDto';
  kpis: Array<KpiDto>;
};

export type InventoryView = {
  __typename?: 'InventoryView';
  components: ComponentStatsSection;
  orgName: Scalars['String']['output'];
  simStock: SimPoolStatsSection;
  unassignedNodes: NodesSection;
};

export type InvitationDto = {
  __typename?: 'InvitationDto';
  email: Scalars['String']['output'];
  expireAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  link: Scalars['String']['output'];
  name: Scalars['String']['output'];
  role: Scalars['String']['output'];
  status: Invitation_Status;
  userId: Scalars['String']['output'];
};

export type InvitationsResDto = {
  __typename?: 'InvitationsResDto';
  invitations: Array<InvitationDto>;
};

export type InvoiceRowDto = {
  __typename?: 'InvoiceRowDto';
  amount: Scalars['Float']['output'];
  generatedAt?: Maybe<Scalars['String']['output']>;
  invoiceId?: Maybe<Scalars['String']['output']>;
  status?: Maybe<Scalars['String']['output']>;
};

export type InvoicesSection = {
  __typename?: 'InvoicesSection';
  error?: Maybe<SectionError>;
  reports?: Maybe<Array<ReportDto>>;
};

export type ItemResDto = {
  __typename?: 'ItemResDto';
  code: Scalars['String']['output'];
  name: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type KpiDto = {
  __typename?: 'KpiDto';
  asOf?: Maybe<Scalars['String']['output']>;
  delta?: Maybe<Scalars['Float']['output']>;
  deltaPeriod?: Maybe<Scalars['String']['output']>;
  formatted?: Maybe<Scalars['String']['output']>;
  key: Scalars['String']['output'];
  stale?: Maybe<Scalars['Boolean']['output']>;
  value: Scalars['Float']['output'];
};

export type KpiEntryDto = {
  __typename?: 'KpiEntryDto';
  format?: Maybe<Scalars['String']['output']>;
  key: Scalars['String']['output'];
  label?: Maybe<Scalars['String']['output']>;
  success: Scalars['Boolean']['output'];
  threshold?: Maybe<MetricThreshold>;
  timestamp: Scalars['Float']['output'];
  unit?: Maybe<Scalars['String']['output']>;
  value: Scalars['Float']['output'];
};

export type KpisSection = {
  __typename?: 'KpisSection';
  error?: Maybe<SectionError>;
  metrics?: Maybe<Array<KpiEntryDto>>;
};

export type ListSimsInput = {
  networkId: Scalars['String']['input'];
  status: Scalars['String']['input'];
};

export type MemberDto = {
  __typename?: 'MemberDto';
  email: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  memberId: Scalars['String']['output'];
  memberSince?: Maybe<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  role: Scalars['String']['output'];
  userId: Scalars['String']['output'];
};

export type Members = {
  __typename?: 'Members';
  activeMembers: Scalars['String']['output'];
  inactiveMembers: Scalars['String']['output'];
  totalMembers: Scalars['String']['output'];
};

export type MembersResDto = {
  __typename?: 'MembersResDto';
  members: Array<MemberDto>;
};

export type MembersView = {
  __typename?: 'MembersView';
  orgName: Scalars['String']['output'];
  team: TeamSection;
};

export type MetaDto = {
  __typename?: 'MetaDto';
  count: Scalars['Int']['output'];
  page: Scalars['Int']['output'];
  pages: Scalars['Int']['output'];
  size: Scalars['Int']['output'];
};

export type MetricInfoDto = {
  __typename?: 'MetricInfoDto';
  lastSampleAt?: Maybe<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
  stale: Scalars['Boolean']['output'];
  unit?: Maybe<Scalars['String']['output']>;
};

export type MetricPanelDto = {
  __typename?: 'MetricPanelDto';
  alarms: Array<AlarmRowDto>;
  kpis: Array<KpiDto>;
  series: Array<TimeSeriesDto>;
};

export type MetricRes = {
  __typename?: 'MetricRes';
  dataPlanId?: Maybe<Scalars['String']['output']>;
  format?: Maybe<Scalars['String']['output']>;
  label?: Maybe<Scalars['String']['output']>;
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

export type MetricThreshold = {
  __typename?: 'MetricThreshold';
  max: Scalars['Float']['output'];
  min: Scalars['Float']['output'];
  normal: Scalars['Float']['output'];
};

export type MetricsRangeInput = {
  from: Scalars['Int']['input'];
  keys: Array<Scalars['String']['input']>;
  nodeId?: InputMaybe<Scalars['String']['input']>;
  operation?: InputMaybe<Scalars['String']['input']>;
  to?: InputMaybe<Scalars['Int']['input']>;
};

export type MetricsRes = {
  __typename?: 'MetricsRes';
  metrics: Array<MetricRes>;
};

export type Mutation = {
  __typename?: 'Mutation';
  addMember: MemberDto;
  addNetwork: NetworkDto;
  addNode: Node;
  addNodeToSite: CBooleanResponse;
  addPackage: PackageDto;
  addPackagesToSim: AddPackagesSimResDto;
  addSite: SiteDto;
  addSubscriber: SubscriberDto;
  allocateSim: AllocateSimApiDto;
  attachNode: CBooleanResponse;
  createInvitation: InvitationDto;
  defaultMarkup: CBooleanResponse;
  deleteInvitation: DeleteInvitationResDto;
  deleteNodeFromOrg: DeleteNode;
  deletePackage: IdResponse;
  deleteSim: DeleteSimResDto;
  deleteSubscriber: CBooleanResponse;
  detachhNode: CBooleanResponse;
  processPayment: ProcessPaymentDto;
  rebuildRollups: RebuildRollupsResultDto;
  refreshAnalytics: RefreshResultDto;
  releaseNodeFromSite: CBooleanResponse;
  removeMember: CBooleanResponse;
  removePackageForSim: RemovePackageFromSimResDto;
  restartNode: CBooleanResponse;
  restartNodes: CBooleanResponse;
  restartSite: CBooleanResponse;
  setDefaultNetwork: CBooleanResponse;
  toggleInternetSwitch: CBooleanResponse;
  toggleRFStatus: CBooleanResponse;
  toggleService: CBooleanResponse;
  toggleSimStatus: SimStatusResDto;
  updateFirstVisit: UserFirstVisitResDto;
  updateInvitation: UpdateInvitationResDto;
  updateMember: CBooleanResponse;
  updateNode: Node;
  updateNodeState: Node;
  updateNotification: UpdateNotificationResDto;
  updatePackage: PackageDto;
  updatePayment: PaymentDto;
  updateSite: SiteDto;
  updateSoftware: StringResponse;
  updateSubscriber: CBooleanResponse;
  uploadSims: UploadSimsResDto;
};


export type MutationAddMemberArgs = {
  data: AddMemberInputDto;
};


export type MutationAddNetworkArgs = {
  data: AddNetworkInputDto;
};


export type MutationAddNodeArgs = {
  data: AddNodeInput;
};


export type MutationAddNodeToSiteArgs = {
  data: AddNodeToSiteInput;
};


export type MutationAddPackageArgs = {
  data: AddPackageInputDto;
};


export type MutationAddPackagesToSimArgs = {
  data: AddPackagesToSimInputDto;
};


export type MutationAddSiteArgs = {
  data: AddSiteInputDto;
};


export type MutationAddSubscriberArgs = {
  data: SubscriberInputDto;
};


export type MutationAllocateSimArgs = {
  data: AllocateSimInputDto;
};


export type MutationAttachNodeArgs = {
  data: AttachNodeInput;
};


export type MutationCreateInvitationArgs = {
  data: CreateInvitationInputDto;
};


export type MutationDefaultMarkupArgs = {
  data: DefaultMarkupInputDto;
};


export type MutationDeleteInvitationArgs = {
  id: Scalars['String']['input'];
};


export type MutationDeleteNodeFromOrgArgs = {
  data: NodeInput;
};


export type MutationDeletePackageArgs = {
  packageId: Scalars['String']['input'];
};


export type MutationDeleteSimArgs = {
  data: DeleteSimInputDto;
};


export type MutationDeleteSubscriberArgs = {
  subscriberId: Scalars['String']['input'];
};


export type MutationDetachhNodeArgs = {
  data: NodeInput;
};


export type MutationProcessPaymentArgs = {
  data: ProcessPaymentInputDto;
};


export type MutationRebuildRollupsArgs = {
  data: RebuildRollupsInput;
};


export type MutationRefreshAnalyticsArgs = {
  data: RefreshInput;
};


export type MutationReleaseNodeFromSiteArgs = {
  data: NodeInput;
};


export type MutationRemoveMemberArgs = {
  id: Scalars['String']['input'];
};


export type MutationRemovePackageForSimArgs = {
  data: RemovePackageFormSimInputDto;
};


export type MutationRestartNodeArgs = {
  data: RestartNodeInputDto;
};


export type MutationRestartNodesArgs = {
  data: RestartNodesInputDto;
};


export type MutationRestartSiteArgs = {
  data: RestartSiteInputDto;
};


export type MutationSetDefaultNetworkArgs = {
  data: SetDefaultNetworkInputDto;
};


export type MutationToggleInternetSwitchArgs = {
  data: ToggleInternetSwitchInputDto;
};


export type MutationToggleRfStatusArgs = {
  data: ToggleRfStatusInputDto;
};


export type MutationToggleServiceArgs = {
  data: ToggleRfStatusInputDto;
};


export type MutationToggleSimStatusArgs = {
  data: ToggleSimStatusInputDto;
};


export type MutationUpdateFirstVisitArgs = {
  data: UserFirstVisitInputDto;
};


export type MutationUpdateInvitationArgs = {
  data: UpdateInvitationInputDto;
};


export type MutationUpdateMemberArgs = {
  data: UpdateMemberInputDto;
  memberId: Scalars['String']['input'];
};


export type MutationUpdateNodeArgs = {
  data: UpdateNodeInput;
};


export type MutationUpdateNodeStateArgs = {
  data: UpdateNodeStateInput;
};


export type MutationUpdateNotificationArgs = {
  id: Scalars['String']['input'];
  isRead: Scalars['Boolean']['input'];
};


export type MutationUpdatePackageArgs = {
  data: UpdatePackageInputDto;
  packageId: Scalars['String']['input'];
};


export type MutationUpdatePaymentArgs = {
  data: UpdatePaymentInputDto;
};


export type MutationUpdateSiteArgs = {
  data: UpdateSiteInputDto;
  siteId: Scalars['String']['input'];
};


export type MutationUpdateSoftwareArgs = {
  data: UpdateSoftwareInputDto;
};


export type MutationUpdateSubscriberArgs = {
  data: UpdateSubscriberInputDto;
  subscriberId: Scalars['String']['input'];
};


export type MutationUploadSimsArgs = {
  data: UploadSimsInputDto;
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

export type NamedValueDto = {
  __typename?: 'NamedValueDto';
  id?: Maybe<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
  value: Scalars['Float']['output'];
};

export type Network = {
  __typename?: 'Network';
  elementType: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  networkName: Scalars['String']['output'];
  sites: Array<Site>;
  subscribers?: Maybe<Subscribers>;
};

export type NetworkAlarmsDto = {
  __typename?: 'NetworkAlarmsDto';
  alarms: Array<AlarmRowDto>;
  kpis: Array<KpiDto>;
  meta?: Maybe<MetaDto>;
};

export type NetworkDto = {
  __typename?: 'NetworkDto';
  budget: Scalars['Float']['output'];
  countries: Array<Scalars['String']['output']>;
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  isDefault: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  networks: Array<Scalars['String']['output']>;
  overdraft: Scalars['Float']['output'];
  paymentLinks: Scalars['Boolean']['output'];
  trafficPolicy: Scalars['Float']['output'];
};

export type NetworkEventsDto = {
  __typename?: 'NetworkEventsDto';
  events: Array<EventRowDto>;
  meta?: Maybe<MetaDto>;
};

export type NetworkMetricsDto = {
  __typename?: 'NetworkMetricsDto';
  metrics: Array<MetricInfoDto>;
  series: Array<TimeSeriesDto>;
};

export type NetworkNodeDto = {
  __typename?: 'NetworkNodeDto';
  kpis: Array<KpiDto>;
  node?: Maybe<NodeRowDto>;
  recentEvents: Array<EventRowDto>;
  series: Array<TimeSeriesDto>;
};

export type NetworkNodesDto = {
  __typename?: 'NetworkNodesDto';
  kpis: Array<KpiDto>;
  meta?: Maybe<MetaDto>;
  nodes: Array<NodeRowDto>;
};

export type NetworkOverview = {
  __typename?: 'NetworkOverview';
  kpis: KpisSection;
  latestAlerts: AlertsSection;
  network: NetworkSection;
  networkId: Scalars['String']['output'];
  nodeStats: NodeStatsSection;
  siteStats: SitesSection;
  subscriberStats: SubscriberStatsSection;
};


export type NetworkOverviewLatestAlertsArgs = {
  limit?: Scalars['Int']['input'];
};

export type NetworkOverviewDto = {
  __typename?: 'NetworkOverviewDto';
  kpis: Array<KpiDto>;
  networkStatus?: Maybe<Scalars['String']['output']>;
  recentEvents: Array<EventRowDto>;
};

export type NetworkSection = {
  __typename?: 'NetworkSection';
  error?: Maybe<SectionError>;
  network?: Maybe<NetworkDto>;
};

export type NetworkSiteDto = {
  __typename?: 'NetworkSiteDto';
  alarms: Array<AlarmRowDto>;
  kpis: Array<KpiDto>;
  series: Array<TimeSeriesDto>;
  site?: Maybe<SiteRowDto>;
};

export type NetworkSitesDto = {
  __typename?: 'NetworkSitesDto';
  kpis: Array<KpiDto>;
  meta?: Maybe<MetaDto>;
  sites: Array<SiteRowDto>;
};

export type NetworkStats = {
  __typename?: 'NetworkStats';
  activeSubscriber: Scalars['Float']['output'];
  averageSignalStrength: Scalars['Float']['output'];
  averageThroughput: Scalars['Float']['output'];
};

export type NetworkSupportSearchDto = {
  __typename?: 'NetworkSupportSearchDto';
  results: Array<SupportResultDto>;
};

export type NetworkTopologyDto = {
  __typename?: 'NetworkTopologyDto';
  sites: Array<TopologySiteDto>;
};

export type NetworksResDto = {
  __typename?: 'NetworksResDto';
  networks: Array<NetworkDto>;
};

export type Node = {
  __typename?: 'Node';
  attached: Array<AttachedNodes>;
  id: Scalars['String']['output'];
  latitude: Scalars['String']['output'];
  longitude: Scalars['String']['output'];
  name: Scalars['String']['output'];
  site: NodeSite;
  status: NodeStatus;
  type: NodeTypeEnum;
};

export type NodeApp = {
  __typename?: 'NodeApp';
  cpu: Scalars['String']['output'];
  date: Scalars['Float']['output'];
  memory: Scalars['String']['output'];
  name: Scalars['String']['output'];
  notes: Scalars['String']['output'];
  version: Scalars['String']['output'];
};

export type NodeApps = {
  __typename?: 'NodeApps';
  apps: Array<NodeApp>;
  type: NodeTypeEnum;
};

export type NodeAppsChangeLogInput = {
  type: NodeTypeEnum;
};

/** Node connectivity enums */
export enum NodeConnectivityEnum {
  Offline = 'Offline',
  Online = 'Online',
  Unknown = 'Unknown'
}

export type NodeInput = {
  id: Scalars['String']['input'];
};

export type NodeLatestMetric = {
  __typename?: 'NodeLatestMetric';
  msg: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  success: Scalars['Boolean']['output'];
  type: Scalars['String']['output'];
  value: Array<Scalars['Float']['output']>;
};

export type NodePoolDto = {
  __typename?: 'NodePoolDto';
  kpis: Array<KpiDto>;
  nodes: Array<NodeRowDto>;
};

export type NodeRowDto = {
  __typename?: 'NodeRowDto';
  configuringDurationSeconds: Scalars['Float']['output'];
  lastTelemetry?: Maybe<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
  noTelemetryWarning: Scalars['Boolean']['output'];
  nodeId: Scalars['String']['output'];
  siteId?: Maybe<Scalars['String']['output']>;
  siteName?: Maybe<Scalars['String']['output']>;
  status?: Maybe<Scalars['String']['output']>;
  type?: Maybe<Scalars['String']['output']>;
  uptime: Scalars['Float']['output'];
};

export type NodeSection = {
  __typename?: 'NodeSection';
  error?: Maybe<SectionError>;
  node?: Maybe<Node>;
};

export type NodeSite = {
  __typename?: 'NodeSite';
  addedAt?: Maybe<Scalars['String']['output']>;
  networkId?: Maybe<Scalars['String']['output']>;
  nodeId?: Maybe<Scalars['String']['output']>;
  siteId?: Maybe<Scalars['String']['output']>;
};

/** Node state enums */
export enum NodeStateEnum {
  Configured = 'Configured',
  Faulty = 'Faulty',
  Operational = 'Operational',
  Unknown = 'Unknown'
}

export type NodeStateRes = {
  __typename?: 'NodeStateRes';
  createdAt: Scalars['String']['output'];
  currentState: NodeStateEnum;
  id: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  previousState?: Maybe<NodeStateEnum>;
  previousStateId?: Maybe<Scalars['String']['output']>;
};

export type NodeStateSection = {
  __typename?: 'NodeStateSection';
  error?: Maybe<SectionError>;
  stateHistory?: Maybe<NodeStateRes>;
};

export type NodeStatsSection = {
  __typename?: 'NodeStatsSection';
  error?: Maybe<SectionError>;
  offline?: Maybe<Scalars['Int']['output']>;
  online?: Maybe<Scalars['Int']['output']>;
  total?: Maybe<Scalars['Int']['output']>;
};

export type NodeStatus = {
  __typename?: 'NodeStatus';
  connectivity: Scalars['String']['output'];
  state: Scalars['String']['output'];
};

/** Node type enums */
export enum NodeTypeEnum {
  Anode = 'anode',
  Cnode = 'cnode',
  Hnode = 'hnode',
  Tnode = 'tnode'
}

export type NodeView = {
  __typename?: 'NodeView';
  health: HealthSection;
  kpis: KpisSection;
  node: NodeSection;
  nodeId: Scalars['String']['output'];
  radioStatus: GapSection;
  siblings: NodesSection;
  site: SiteSection;
  software: SoftwareSection;
  stateHistory: NodeStateSection;
};

export type Nodes = {
  __typename?: 'Nodes';
  nodes: Array<Node>;
};

export type NodesFilterInput = {
  connectivity?: InputMaybe<Scalars['String']['input']>;
  id?: InputMaybe<Scalars['String']['input']>;
  networkId?: InputMaybe<Scalars['String']['input']>;
  siteId?: InputMaybe<Scalars['String']['input']>;
  state?: InputMaybe<Scalars['String']['input']>;
  type?: InputMaybe<Scalars['String']['input']>;
};

export type NodesSection = {
  __typename?: 'NodesSection';
  error?: Maybe<SectionError>;
  nodes?: Maybe<Array<Node>>;
};

export type NodesView = {
  __typename?: 'NodesView';
  health: GapSection;
  networkId?: Maybe<Scalars['String']['output']>;
  nodes: NodesSection;
};

export type NotificationRedirectDto = {
  __typename?: 'NotificationRedirectDto';
  action: Scalars['String']['output'];
  title: Scalars['String']['output'];
};

export type NotificationResDto = {
  __typename?: 'NotificationResDto';
  createdAt: Scalars['String']['output'];
  description: Scalars['String']['output'];
  id: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  orgId: Scalars['String']['output'];
  resourceId: Scalars['String']['output'];
  scope: Notification_Scope;
  subscriberId: Scalars['String']['output'];
  title: Scalars['String']['output'];
  type: Notification_Type;
  userId: Scalars['String']['output'];
};

export type NotificationsDto = {
  __typename?: 'NotificationsDto';
  createdAt: Scalars['String']['output'];
  description: Scalars['String']['output'];
  eventKey: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isRead: Scalars['Boolean']['output'];
  redirect?: Maybe<NotificationRedirectDto>;
  resourceId: Scalars['String']['output'];
  scope: Notification_Scope;
  title: Scalars['String']['output'];
  type: Notification_Type;
};

export type NotificationsResDto = {
  __typename?: 'NotificationsResDto';
  notifications: Array<NotificationsDto>;
};

export type OnboardingStatusDto = {
  __typename?: 'OnboardingStatusDto';
  hasNetwork: Scalars['Boolean']['output'];
  hasNode: Scalars['Boolean']['output'];
  hasSite: Scalars['Boolean']['output'];
  networkId?: Maybe<Scalars['String']['output']>;
  networkName?: Maybe<Scalars['String']['output']>;
};

export type OperationDto = {
  __typename?: 'OperationDto';
  createdAt?: Maybe<Scalars['String']['output']>;
  error?: Maybe<Scalars['String']['output']>;
  fencingToken: Scalars['Float']['output'];
  id: Scalars['String']['output'];
  idempotencyKey?: Maybe<Scalars['String']['output']>;
  leaseExpiresAt?: Maybe<Scalars['String']['output']>;
  requestedBy?: Maybe<Scalars['String']['output']>;
  resourceKey: Scalars['String']['output'];
  startedAt?: Maybe<Scalars['String']['output']>;
  status: Scalars['String']['output'];
  system: Scalars['String']['output'];
  terminalAt?: Maybe<Scalars['String']['output']>;
  type: Scalars['String']['output'];
};

export type Org = {
  __typename?: 'Org';
  country: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  dataplans: Array<DataPlan>;
  elementType: Scalars['String']['output'];
  members?: Maybe<Members>;
  networks: Array<Network>;
  orgId: Scalars['String']['output'];
  orgName: Scalars['String']['output'];
  ownerEmail: Scalars['String']['output'];
  ownerId: Scalars['String']['output'];
  ownerName: Scalars['String']['output'];
  sims?: Maybe<Sims>;
};

export type OrgDto = {
  __typename?: 'OrgDto';
  certificate: Scalars['String']['output'];
  country: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  owner: Scalars['String']['output'];
};

export type OrgTreeRes = {
  __typename?: 'OrgTreeRes';
  org: Org;
};

export type OrgsResDto = {
  __typename?: 'OrgsResDto';
  memberOf: Array<OrgDto>;
  ownerOf: Array<OrgDto>;
  user: Scalars['String']['output'];
};

export type PackageDto = {
  __typename?: 'PackageDto';
  active: Scalars['Boolean']['output'];
  amount: Scalars['Float']['output'];
  apn: Scalars['String']['output'];
  country: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  dataUnit: Scalars['String']['output'];
  dataVolume: Scalars['Float']['output'];
  deletedAt: Scalars['String']['output'];
  dlbr: Scalars['String']['output'];
  duration: Scalars['Float']['output'];
  flatrate: Scalars['Boolean']['output'];
  from: Scalars['String']['output'];
  markup: PackageMarkupApiDto;
  messageUnit: Scalars['String']['output'];
  name: Scalars['String']['output'];
  ownerId: Scalars['String']['output'];
  provider: Scalars['String']['output'];
  rate: PackageRateApiDto;
  simType: Scalars['String']['output'];
  smsVolume: Scalars['Float']['output'];
  to: Scalars['String']['output'];
  type: Scalars['String']['output'];
  ulbr: Scalars['String']['output'];
  updatedAt: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
  voiceUnit: Scalars['String']['output'];
  voiceVolume: Scalars['Float']['output'];
};

export type PackageIntervalDto = {
  __typename?: 'PackageIntervalDto';
  endAt?: Maybe<Scalars['String']['output']>;
  packageId?: Maybe<Scalars['String']['output']>;
  packageName?: Maybe<Scalars['String']['output']>;
  startAt?: Maybe<Scalars['String']['output']>;
  state?: Maybe<Scalars['String']['output']>;
};

export type PackageMarkupApiDto = {
  __typename?: 'PackageMarkupAPIDto';
  baserate: Scalars['String']['output'];
  markup: Scalars['Float']['output'];
};

export type PackagePerformanceDto = {
  __typename?: 'PackagePerformanceDto';
  kpis: Array<KpiDto>;
  meta?: Maybe<MetaDto>;
  packages: Array<PackageRowDto>;
  revenueMix: Array<NamedValueDto>;
};

export type PackageRateApiDto = {
  __typename?: 'PackageRateAPIDto';
  amount: Scalars['Float']['output'];
  data: Scalars['Float']['output'];
  sms_mo: Scalars['String']['output'];
  sms_mt: Scalars['Float']['output'];
};

export type PackageRowDto = {
  __typename?: 'PackageRowDto';
  activeSubscribers: Scalars['Int']['output'];
  dataQuota?: Maybe<Scalars['String']['output']>;
  dataUsed: Scalars['Float']['output'];
  name?: Maybe<Scalars['String']['output']>;
  packageId: Scalars['String']['output'];
  price: Scalars['Float']['output'];
  revenue: Scalars['Float']['output'];
  revenueSharePct?: Maybe<Scalars['Float']['output']>;
  soldCount: Scalars['Int']['output'];
  status?: Maybe<Scalars['String']['output']>;
  validity?: Maybe<Scalars['String']['output']>;
};

export type PackagesResDto = {
  __typename?: 'PackagesResDto';
  packages: Array<PackageDto>;
};

export type PackagesToSimInputDto = {
  package_id: Scalars['String']['input'];
  start_date: Scalars['String']['input'];
};

export type PaymentDto = {
  __typename?: 'PaymentDto';
  amount: Scalars['String']['output'];
  correspondent: Scalars['String']['output'];
  country: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  currency: Scalars['String']['output'];
  depositedAmount: Scalars['String']['output'];
  description: Scalars['String']['output'];
  extra: Scalars['String']['output'];
  failureReason: Scalars['String']['output'];
  id: Scalars['String']['output'];
  itemId: Scalars['String']['output'];
  itemType: Scalars['String']['output'];
  paidAt: Scalars['String']['output'];
  payerEmail: Scalars['String']['output'];
  payerName: Scalars['String']['output'];
  payerPhone: Scalars['String']['output'];
  paymentMethod: Scalars['String']['output'];
  status: Scalars['String']['output'];
};

export type PaymentsDto = {
  __typename?: 'PaymentsDto';
  payments: Array<PaymentDto>;
};

export type PlanNameDto = {
  __typename?: 'PlanNameDto';
  name: Scalars['String']['output'];
  packageId: Scalars['String']['output'];
};

export type PlanStatsDto = {
  __typename?: 'PlanStatsDto';
  active: Scalars['Boolean']['output'];
  amount: Scalars['Float']['output'];
  attachCount?: Maybe<Scalars['Int']['output']>;
  currency: Scalars['String']['output'];
  name: Scalars['String']['output'];
  packageId: Scalars['String']['output'];
  revenue: Scalars['Float']['output'];
  revenueSharePct: Scalars['Int']['output'];
};

export type PlanStatsSection = {
  __typename?: 'PlanStatsSection';
  arpu?: Maybe<Scalars['Float']['output']>;
  error?: Maybe<SectionError>;
  mrr?: Maybe<Scalars['Float']['output']>;
  plans?: Maybe<Array<PlanStatsDto>>;
};

export type PlansSection = {
  __typename?: 'PlansSection';
  error?: Maybe<SectionError>;
  plans?: Maybe<Array<PlanNameDto>>;
};

export type PointDto = {
  __typename?: 'PointDto';
  time?: Maybe<Scalars['String']['output']>;
  value: Scalars['Float']['output'];
};

export type PoolSimsSection = {
  __typename?: 'PoolSimsSection';
  error?: Maybe<SectionError>;
  sims?: Maybe<Array<SimPoolResDto>>;
};

export type ProcessPaymentDto = {
  __typename?: 'ProcessPaymentDto';
  payment: PaymentDto;
};

export type ProcessPaymentInputDto = {
  correspondent?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['String']['input'];
  token: Scalars['String']['input'];
};

export type Query = {
  __typename?: 'Query';
  commerceView: CommerceView;
  example?: Maybe<Scalars['String']['output']>;
  getAlarms: NetworkAlarmsDto;
  getApps?: Maybe<Apps>;
  getAppsChangeLog: AppChangeLogs;
  getBackhaul: MetricPanelDto;
  getBillingSummary: BillingSummaryDto;
  getBusinessHome: BusinessHomeDto;
  getBusinessSite: BusinessSiteDto;
  getBusinessSites: BusinessSitesDto;
  getComponentById: ComponentDto;
  getComponentsByUserId: ComponentsResDto;
  getCountries: CountriesRes;
  getCurrencySymbol: CurrencyRes;
  getCustomer: CustomerDetailDto;
  getCustomerOverview: CustomerOverviewDto;
  getCustomerSims: CustomerSimsDto;
  getCustomerSupport: CustomerSupportDto;
  getDataUsage: SimDataUsage;
  getDataUsages: SimDataUsages;
  getDefaultMarkup: DefaultMarkupResDto;
  getDefaultMarkupHistory: DefaultMarkupHistoryResDto;
  getEvents: NetworkEventsDto;
  getGeneratedPdfReport: GetPdfReportUrlDto;
  getHealthReport: HealthInfo;
  getInventoryReadiness: InventoryReadinessDto;
  getInvitation: InvitationDto;
  getInvitations: InvitationsResDto;
  getInvitationsByEmail: InvitationsResDto;
  getMember: MemberDto;
  getMemberByUserId: MemberDto;
  getMembers: MembersResDto;
  getMetrics: NetworkMetricsDto;
  getNetwork: NetworkDto;
  getNetworkNode: NetworkNodeDto;
  getNetworkNodes: NetworkNodesDto;
  getNetworkOverview: NetworkOverviewDto;
  getNetworkSite: NetworkSiteDto;
  getNetworkSites: NetworkSitesDto;
  getNetworkStats: NetworkStats;
  getNetworks: NetworksResDto;
  getNode: Node;
  getNodeApps: NodeApps;
  getNodeLatestMetric: NodeLatestMetric;
  getNodePool: NodePoolDto;
  getNodeState: NodeStateRes;
  getNodes: Nodes;
  getNodesByNetwork: Nodes;
  getNodesByState: Nodes;
  getNodesForSite: Nodes;
  getNodesLocation: Nodes;
  getNotification: NotificationResDto;
  getNotifications: NotificationsResDto;
  getOperation?: Maybe<OperationDto>;
  getOrg: OrgDto;
  getOrgTree: OrgTreeRes;
  getOrgs: OrgsResDto;
  getPackage: PackageDto;
  getPackagePerformance: PackagePerformanceDto;
  getPackages: PackagesResDto;
  getPackagesForSim: GetSimPackagesDtoApi;
  getPayment: PaymentDto;
  getPayments: PaymentsDto;
  getPower: MetricPanelDto;
  getRadio: MetricPanelDto;
  getRefreshState: RefreshStateDto;
  getReport: GetReportDto;
  getReportPdf: GetReportDto;
  getReports: GetReportsDto;
  getResourceLock: ResourceLockDto;
  getSalesOverview: SalesOverviewDto;
  getSim: SimDto;
  getSimPool: SimPoolDto;
  getSimPoolStats: SimPoolStatsDto;
  getSims: SimsResDto;
  getSimsByNetwork: SubscriberSimsResDto;
  getSimsBySubscriber: SubscriberToSimsDto;
  getSimsFromPool: SimsPoolResDto;
  getSite: SiteDto;
  getSites: SitesResDto;
  getSoftwares: Softwares;
  getSubscriber: SubscriberDto;
  getSubscriberMetricsByNetwork: SubscriberMetricsByNetworkDto;
  getSubscribersByNetwork: SubscribersResDto;
  getTimezones: TimezoneRes;
  getToken: TokenResDto;
  getTopology: NetworkTopologyDto;
  getUser: UserResDto;
  inventoryView: InventoryView;
  listCustomers: CustomerListDto;
  membersView: MembersView;
  metricsRange: MetricsRes;
  networkOverview: NetworkOverview;
  nodeView: NodeView;
  nodesView: NodesView;
  onboardingStatus: OnboardingStatusDto;
  searchCustomers: CustomerListDto;
  simPoolView: SimPoolView;
  siteView: SiteView;
  sitesView: SitesView;
  subscriberView: SubscriberView;
  subscribersView: SubscribersView;
  supportSearch: NetworkSupportSearchDto;
  whoami: WhoamiDto;
};


export type QueryCommerceViewArgs = {
  networkId?: InputMaybe<Scalars['String']['input']>;
};


export type QueryGetAlarmsArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetAppsArgs = {
  data: GetAppsInputDto;
};


export type QueryGetAppsChangeLogArgs = {
  data: NodeAppsChangeLogInput;
};


export type QueryGetBackhaulArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetBillingSummaryArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetBusinessHomeArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetBusinessSiteArgs = {
  data: AnalyticsSiteInput;
};


export type QueryGetBusinessSitesArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetComponentByIdArgs = {
  componentId: Scalars['String']['input'];
};


export type QueryGetComponentsByUserIdArgs = {
  data: ComponentTypeInputDto;
};


export type QueryGetCurrencySymbolArgs = {
  code: Scalars['String']['input'];
};


export type QueryGetCustomerArgs = {
  data: CustomerByIdInput;
};


export type QueryGetCustomerOverviewArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetCustomerSimsArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetCustomerSupportArgs = {
  data: CustomerByIdInput;
};


export type QueryGetDataUsageArgs = {
  data: SimUsageInputDto;
};


export type QueryGetDataUsagesArgs = {
  data: SimUsagesInputDto;
};


export type QueryGetEventsArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetGeneratedPdfReportArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetHealthReportArgs = {
  data: GetHealthReportInputDto;
};


export type QueryGetInventoryReadinessArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetInvitationArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetInvitationsByEmailArgs = {
  email: Scalars['String']['input'];
};


export type QueryGetMemberArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetMemberByUserIdArgs = {
  userId: Scalars['String']['input'];
};


export type QueryGetMetricsArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNetworkNodeArgs = {
  data: AnalyticsNodeInput;
};


export type QueryGetNetworkNodesArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetNetworkOverviewArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetNetworkSiteArgs = {
  data: AnalyticsSiteInput;
};


export type QueryGetNetworkSitesArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetNetworkStatsArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNodeArgs = {
  data: NodeInput;
};


export type QueryGetNodeAppsArgs = {
  data: NodeAppsChangeLogInput;
};


export type QueryGetNodeLatestMetricArgs = {
  data: GetNodeLatestMetricInput;
};


export type QueryGetNodePoolArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetNodeStateArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetNodesArgs = {
  data: NodesFilterInput;
};


export type QueryGetNodesByNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetNodesByStateArgs = {
  data: GetNodesByStateInput;
};


export type QueryGetNodesForSiteArgs = {
  siteId: Scalars['String']['input'];
};


export type QueryGetNotificationArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetOperationArgs = {
  data: GetOperationInputDto;
};


export type QueryGetPackageArgs = {
  packageId: Scalars['String']['input'];
};


export type QueryGetPackagePerformanceArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetPackagesForSimArgs = {
  data: GetPackagesForSimInputDto;
};


export type QueryGetPaymentArgs = {
  paymentId: Scalars['String']['input'];
};


export type QueryGetPaymentsArgs = {
  data: GetPaymentsInputDto;
};


export type QueryGetPowerArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetRadioArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetReportArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetReportPdfArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetReportsArgs = {
  data: GetReportsInputDto;
};


export type QueryGetResourceLockArgs = {
  data: GetResourceLockInputDto;
};


export type QueryGetSalesOverviewArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetSimArgs = {
  data: GetSimInputDto;
};


export type QueryGetSimPoolArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetSimPoolStatsArgs = {
  data: GetSimsInput;
};


export type QueryGetSimsArgs = {
  data: ListSimsInput;
};


export type QueryGetSimsByNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetSimsBySubscriberArgs = {
  data: GetSimBySubscriberInputDto;
};


export type QueryGetSimsFromPoolArgs = {
  data: GetSimsInput;
};


export type QueryGetSiteArgs = {
  siteId: Scalars['String']['input'];
};


export type QueryGetSitesArgs = {
  data: SitesInputDto;
};


export type QueryGetSoftwaresArgs = {
  data: GetSoftwaresInput;
};


export type QueryGetSubscriberArgs = {
  subscriberId: Scalars['String']['input'];
};


export type QueryGetSubscriberMetricsByNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetSubscribersByNetworkArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryGetTokenArgs = {
  paymentId: Scalars['String']['input'];
};


export type QueryGetTopologyArgs = {
  data: AnalyticsWindowInput;
};


export type QueryGetUserArgs = {
  userId: Scalars['String']['input'];
};


export type QueryListCustomersArgs = {
  data: AnalyticsWindowInput;
};


export type QueryMetricsRangeArgs = {
  data: MetricsRangeInput;
};


export type QueryNetworkOverviewArgs = {
  networkId: Scalars['String']['input'];
};


export type QueryNodeViewArgs = {
  nodeId: Scalars['String']['input'];
};


export type QueryNodesViewArgs = {
  networkId?: InputMaybe<Scalars['String']['input']>;
};


export type QuerySearchCustomersArgs = {
  data: AnalyticsWindowInput;
};


export type QuerySimPoolViewArgs = {
  simType: Scalars['String']['input'];
};


export type QuerySiteViewArgs = {
  siteId: Scalars['String']['input'];
};


export type QuerySitesViewArgs = {
  networkId: Scalars['String']['input'];
};


export type QuerySubscriberViewArgs = {
  subscriberId: Scalars['String']['input'];
};


export type QuerySubscribersViewArgs = {
  networkId: Scalars['String']['input'];
};


export type QuerySupportSearchArgs = {
  data: AnalyticsWindowInput;
};

export enum Role_Type {
  RoleAdmin = 'ROLE_ADMIN',
  RoleInvalid = 'ROLE_INVALID',
  RoleNetworkOwner = 'ROLE_NETWORK_OWNER',
  RoleOwner = 'ROLE_OWNER',
  RoleUser = 'ROLE_USER',
  RoleVendor = 'ROLE_VENDOR'
}

export type RawReportDto = {
  __typename?: 'RawReportDto';
  currency: Scalars['String']['output'];
  customer: CustomerDto;
  fees: Array<FeeDto>;
  feesAmountCents: Scalars['String']['output'];
  fileUrl: Scalars['String']['output'];
  invoiceType: Scalars['String']['output'];
  issuingDate: Scalars['String']['output'];
  paymentDueDate: Scalars['String']['output'];
  paymentOverdue: Scalars['Boolean']['output'];
  paymentStatus: Scalars['String']['output'];
  status: Scalars['String']['output'];
  subTotalExcludingTaxesAmountCents: Scalars['String']['output'];
  subTotalIncludingTaxesAmountCents: Scalars['String']['output'];
  subscriptions: Array<SubscriptionDto>;
  taxesAmountCents: Scalars['String']['output'];
  totalAmountCents: Scalars['String']['output'];
  vatAmountCents: Scalars['String']['output'];
  vatAmountCurrency?: Maybe<Scalars['String']['output']>;
};

export type RebuildRollupsInput = {
  family: Scalars['String']['input'];
  from?: InputMaybe<Scalars['String']['input']>;
  to?: InputMaybe<Scalars['String']['input']>;
};

export type RebuildRollupsResultDto = {
  __typename?: 'RebuildRollupsResultDto';
  rollups: Array<RollupStateDto>;
};

export type RefreshInput = {
  source: Scalars['String']['input'];
};

export type RefreshResultDto = {
  __typename?: 'RefreshResultDto';
  states: Array<SourceStateDto>;
};

export type RefreshStateDto = {
  __typename?: 'RefreshStateDto';
  rollups: Array<RollupStateDto>;
  states: Array<SourceStateDto>;
};

export type RemovePackageFormSimInputDto = {
  packageId: Scalars['String']['input'];
  simId: Scalars['String']['input'];
};

export type RemovePackageFromSimResDto = {
  __typename?: 'RemovePackageFromSimResDto';
  packageId?: Maybe<Scalars['String']['output']>;
};

export type ReportDto = {
  __typename?: 'ReportDto';
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isPaid: Scalars['Boolean']['output'];
  networkId: Scalars['String']['output'];
  ownerId: Scalars['String']['output'];
  ownerType: Scalars['String']['output'];
  period: Scalars['String']['output'];
  rawReport: RawReportDto;
  type: Scalars['String']['output'];
};

export type ResourceLockDto = {
  __typename?: 'ResourceLockDto';
  locked: Scalars['Boolean']['output'];
  operation?: Maybe<OperationDto>;
};

export type RestartNodeInputDto = {
  nodeId: Scalars['String']['input'];
};

export type RestartNodesInputDto = {
  networkId: Scalars['String']['input'];
  nodeIds: Array<Scalars['String']['input']>;
};

export type RestartSiteInputDto = {
  networkId: Scalars['String']['input'];
  siteId: Scalars['String']['input'];
};

export type RevenueSection = {
  __typename?: 'RevenueSection';
  error?: Maybe<SectionError>;
  momPct?: Maybe<Scalars['Int']['output']>;
  monthPaid?: Maybe<Scalars['Float']['output']>;
  prevMonthPaid?: Maybe<Scalars['Float']['output']>;
  totalPaid?: Maybe<Scalars['Float']['output']>;
  totalPending?: Maybe<Scalars['Float']['output']>;
};

export type RollupStateDto = {
  __typename?: 'RollupStateDto';
  dirty: Scalars['Boolean']['output'];
  rollup?: Maybe<Scalars['String']['output']>;
  watermark?: Maybe<Scalars['String']['output']>;
};

export enum Sim_Status {
  All = 'ALL',
  Assigned = 'ASSIGNED',
  Unassigned = 'UNASSIGNED'
}

export enum Sim_Types {
  OperatorData = 'operator_data',
  Test = 'test',
  UkamaData = 'ukama_data',
  Unknown = 'unknown'
}

export type SalesOverviewDto = {
  __typename?: 'SalesOverviewDto';
  kpis: Array<KpiDto>;
  revenueByPackage: Array<NamedValueDto>;
  revenueBySite: Array<NamedValueDto>;
  revenueTrend?: Maybe<TimeSeriesDto>;
};

/** Typed failure of one section of a composite query. The section's data field resolves to null and a SectionError describes why, so the UI can distinguish 'failed' from 'genuinely empty'. */
export type SectionError = {
  __typename?: 'SectionError';
  code: SectionErrorCode;
  message: Scalars['String']['output'];
  section: Scalars['String']['output'];
};

/** Machine-readable failure code for a composite query section. UI branches on this code; `message` is for display/logs only. */
export enum SectionErrorCode {
  Forbidden = 'FORBIDDEN',
  Internal = 'INTERNAL',
  NotFound = 'NOT_FOUND',
  NotImplemented = 'NOT_IMPLEMENTED',
  UpstreamError = 'UPSTREAM_ERROR',
  UpstreamTimeout = 'UPSTREAM_TIMEOUT'
}

export type SetDefaultNetworkInputDto = {
  id: Scalars['String']['input'];
};

export type SimAllocatePackageDto = {
  __typename?: 'SimAllocatePackageDto';
  endDate?: Maybe<Scalars['String']['output']>;
  id?: Maybe<Scalars['String']['output']>;
  isActive?: Maybe<Scalars['Boolean']['output']>;
  packageId?: Maybe<Scalars['String']['output']>;
  startDate?: Maybe<Scalars['String']['output']>;
};

export type SimBatchDto = {
  __typename?: 'SimBatchDto';
  assigned: Scalars['Int']['output'];
  assignedPercent: Scalars['Float']['output'];
  batchId: Scalars['String']['output'];
  quantity: Scalars['Int']['output'];
  uploadedAt?: Maybe<Scalars['String']['output']>;
};

export type SimDataUsage = {
  __typename?: 'SimDataUsage';
  simId: Scalars['String']['output'];
  usage: Scalars['String']['output'];
};

export type SimDataUsages = {
  __typename?: 'SimDataUsages';
  usages: Array<SimDataUsage>;
};

export type SimDto = {
  __typename?: 'SimDto';
  activationsCount: Scalars['String']['output'];
  allocatedAt: Scalars['String']['output'];
  deactivationsCount: Scalars['String']['output'];
  firstActivatedOn: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi: Scalars['String']['output'];
  isPhysical: Scalars['Boolean']['output'];
  lastActivatedOn: Scalars['String']['output'];
  msisdn: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  package?: Maybe<SimPackage>;
  status: Scalars['String']['output'];
  subscriberId: Scalars['String']['output'];
  syncStatus: Scalars['String']['output'];
  trafficPolicy: Scalars['Float']['output'];
  type: Scalars['String']['output'];
};

export type SimPackage = {
  __typename?: 'SimPackage';
  asExpired: Scalars['Boolean']['output'];
  defaultDuration: Scalars['String']['output'];
  endDate: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isActive: Scalars['Boolean']['output'];
  packageId: Scalars['String']['output'];
  startDate: Scalars['String']['output'];
};

export type SimPackageDto = {
  __typename?: 'SimPackageDto';
  created_at: Scalars['String']['output'];
  end_date: Scalars['String']['output'];
  id: Scalars['String']['output'];
  is_active: Scalars['Boolean']['output'];
  package_id: Scalars['String']['output'];
  start_date: Scalars['String']['output'];
  updated_at: Scalars['String']['output'];
};

export type SimPoolDto = {
  __typename?: 'SimPoolDto';
  batches: Array<SimBatchDto>;
  kpis: Array<KpiDto>;
};

export type SimPoolResDto = {
  __typename?: 'SimPoolResDto';
  activationCode: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  deletedAt: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  isAllocated: Scalars['Boolean']['output'];
  isFailed: Scalars['Boolean']['output'];
  isPhysical: Scalars['Boolean']['output'];
  msisdn: Scalars['String']['output'];
  qrCode: Scalars['String']['output'];
  simType: Scalars['String']['output'];
  smApAddress: Scalars['String']['output'];
  updatedAt: Scalars['String']['output'];
};

export type SimPoolStatsDto = {
  __typename?: 'SimPoolStatsDto';
  available: Scalars['Float']['output'];
  consumed: Scalars['Float']['output'];
  esim: Scalars['Float']['output'];
  failed: Scalars['Float']['output'];
  physical: Scalars['Float']['output'];
  total: Scalars['Float']['output'];
};

export type SimPoolStatsSection = {
  __typename?: 'SimPoolStatsSection';
  available?: Maybe<Scalars['Int']['output']>;
  consumed?: Maybe<Scalars['Int']['output']>;
  error?: Maybe<SectionError>;
  esim?: Maybe<Scalars['Int']['output']>;
  failed?: Maybe<Scalars['Int']['output']>;
  lowStock?: Maybe<Scalars['Boolean']['output']>;
  pctAssigned?: Maybe<Scalars['Int']['output']>;
  physical?: Maybe<Scalars['Int']['output']>;
  total?: Maybe<Scalars['Int']['output']>;
};

export type SimPoolView = {
  __typename?: 'SimPoolView';
  simType: Scalars['String']['output'];
  sims: PoolSimsSection;
  stats: SimPoolStatsSection;
};


export type SimPoolViewSimsArgs = {
  limit?: Scalars['Int']['input'];
};

export type SimRowDto = {
  __typename?: 'SimRowDto';
  allocatedAt?: Maybe<Scalars['String']['output']>;
  batchId?: Maybe<Scalars['String']['output']>;
  customerId?: Maybe<Scalars['String']['output']>;
  iccid?: Maybe<Scalars['String']['output']>;
  simId: Scalars['String']['output'];
  status?: Maybe<Scalars['String']['output']>;
};

export type SimStatusResDto = {
  __typename?: 'SimStatusResDto';
  simId?: Maybe<Scalars['String']['output']>;
};

export type SimToPackagesDto = {
  __typename?: 'SimToPackagesDto';
  end_date: Scalars['String']['output'];
  id: Scalars['String']['output'];
  is_active: Scalars['Boolean']['output'];
  package_id: Scalars['String']['output'];
  start_date: Scalars['String']['output'];
};

export type SimUsageInputDto = {
  iccid: Scalars['String']['input'];
  simId: Scalars['String']['input'];
  type: Scalars['String']['input'];
};

export type SimUsagesInputDto = {
  networkId: Scalars['String']['input'];
  type: Scalars['String']['input'];
};

export type Sims = {
  __typename?: 'Sims';
  availableSims: Scalars['String']['output'];
  consumed: Scalars['String']['output'];
  totalSims: Scalars['String']['output'];
};

export type SimsPoolResDto = {
  __typename?: 'SimsPoolResDto';
  sims: Array<SimPoolResDto>;
};

export type SimsResDto = {
  __typename?: 'SimsResDto';
  sims: Array<SimDto>;
};

export type Site = {
  __typename?: 'Site';
  components: Array<Component>;
  elementType: Scalars['String']['output'];
  siteId: Scalars['String']['output'];
  siteName: Scalars['String']['output'];
};

export type SiteComponentDto = {
  __typename?: 'SiteComponentDto';
  componentId?: Maybe<Scalars['String']['output']>;
  componentName?: Maybe<Scalars['String']['output']>;
  elementType: Scalars['String']['output'];
};

export type SiteComponentsSection = {
  __typename?: 'SiteComponentsSection';
  components?: Maybe<Array<SiteComponentDto>>;
  error?: Maybe<SectionError>;
};

export type SiteCustomersSection = {
  __typename?: 'SiteCustomersSection';
  count?: Maybe<Scalars['Int']['output']>;
  error?: Maybe<SectionError>;
};

export type SiteDto = {
  __typename?: 'SiteDto';
  accessId: Scalars['String']['output'];
  backhaulId: Scalars['String']['output'];
  createdAt: Scalars['String']['output'];
  id: Scalars['String']['output'];
  installDate: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  latitude: Scalars['String']['output'];
  location: Scalars['String']['output'];
  longitude: Scalars['String']['output'];
  name: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  powerId: Scalars['String']['output'];
  spectrumId: Scalars['String']['output'];
  switchId: Scalars['String']['output'];
};

export type SiteNodeCountDto = {
  __typename?: 'SiteNodeCountDto';
  offline: Scalars['Int']['output'];
  online: Scalars['Int']['output'];
  siteId: Scalars['String']['output'];
  total: Scalars['Int']['output'];
};

export type SiteNodeCountsSection = {
  __typename?: 'SiteNodeCountsSection';
  counts?: Maybe<Array<SiteNodeCountDto>>;
  error?: Maybe<SectionError>;
};

export type SiteRowDto = {
  __typename?: 'SiteRowDto';
  backhaulLatencyHigh: Scalars['Boolean']['output'];
  batteryCritical: Scalars['Boolean']['output'];
  customers: Scalars['Int']['output'];
  issueSummary?: Maybe<Scalars['String']['output']>;
  latitude: Scalars['Float']['output'];
  longitude: Scalars['Float']['output'];
  name?: Maybe<Scalars['String']['output']>;
  nodeCount: Scalars['Int']['output'];
  offlineDurationSeconds: Scalars['Float']['output'];
  siteId: Scalars['String']['output'];
  status?: Maybe<Scalars['String']['output']>;
  uptime: Scalars['Float']['output'];
};

export type SiteSection = {
  __typename?: 'SiteSection';
  error?: Maybe<SectionError>;
  site?: Maybe<SiteDto>;
};

export type SiteSummaryDto = {
  __typename?: 'SiteSummaryDto';
  customers: Scalars['Int']['output'];
  name?: Maybe<Scalars['String']['output']>;
  revenue: Scalars['Float']['output'];
  siteId: Scalars['String']['output'];
  status?: Maybe<Scalars['String']['output']>;
};

export type SiteView = {
  __typename?: 'SiteView';
  components: SiteComponentsSection;
  financials: GapSection;
  kpis: KpisSection;
  nodes: NodesSection;
  power: KpisSection;
  site: SiteSection;
  siteId: Scalars['String']['output'];
};

export type SitesInputDto = {
  networkId?: InputMaybe<Scalars['String']['input']>;
};

export type SitesResDto = {
  __typename?: 'SitesResDto';
  sites: Array<SiteDto>;
};

export type SitesSection = {
  __typename?: 'SitesSection';
  error?: Maybe<SectionError>;
  sites?: Maybe<Array<SiteDto>>;
};

export type SitesView = {
  __typename?: 'SitesView';
  customers: SiteCustomersSection;
  financials: GapSection;
  kpis: GapSection;
  networkId: Scalars['String']['output'];
  nodeCounts: SiteNodeCountsSection;
  sites: SitesSection;
};

export type Software = {
  __typename?: 'Software';
  changeLog: Array<Scalars['String']['output']>;
  createdAt: Scalars['String']['output'];
  currentVersion: Scalars['String']['output'];
  desiredVersion: Scalars['String']['output'];
  id: Scalars['String']['output'];
  metricsKeys: Array<Scalars['String']['output']>;
  name: Scalars['String']['output'];
  nodeId: Scalars['String']['output'];
  notes: Scalars['String']['output'];
  releaseDate: Scalars['String']['output'];
  space: Scalars['String']['output'];
  status: SoftwareStatusEnum;
  updatedAt: Scalars['String']['output'];
};

export type SoftwareSection = {
  __typename?: 'SoftwareSection';
  error?: Maybe<SectionError>;
  softwares?: Maybe<Softwares>;
};

/** Software status enums */
export enum SoftwareStatusEnum {
  Unknown = 'unknown',
  UpToDate = 'up_to_date',
  UpdateAvailable = 'update_available',
  UpdateFailed = 'update_failed',
  UpdateInProgress = 'update_in_progress'
}

export type Softwares = {
  __typename?: 'Softwares';
  software: Array<Software>;
};

export type SourceStateDto = {
  __typename?: 'SourceStateDto';
  detail?: Maybe<Scalars['String']['output']>;
  lastRunAt?: Maybe<Scalars['String']['output']>;
  lastSuccessAt?: Maybe<Scalars['String']['output']>;
  source?: Maybe<Scalars['String']['output']>;
  status?: Maybe<Scalars['String']['output']>;
};

export type StringResponse = {
  __typename?: 'StringResponse';
  message: Scalars['String']['output'];
};

export type SubscriberBillingSection = {
  __typename?: 'SubscriberBillingSection';
  error?: Maybe<SectionError>;
  payments?: Maybe<Array<PaymentDto>>;
};

export type SubscriberDto = {
  __typename?: 'SubscriberDto';
  address: Scalars['String']['output'];
  dob: Scalars['String']['output'];
  email: Scalars['String']['output'];
  gender: Scalars['String']['output'];
  idSerial: Scalars['String']['output'];
  name: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  phone: Scalars['String']['output'];
  proofOfIdentification: Scalars['String']['output'];
  sim?: Maybe<Array<SubscriberSimDto>>;
  uuid: Scalars['String']['output'];
};

export type SubscriberInputDto = {
  email: Scalars['String']['input'];
  name: Scalars['String']['input'];
  network_id: Scalars['String']['input'];
  phone?: InputMaybe<Scalars['String']['input']>;
};

export type SubscriberMetricsByNetworkDto = {
  __typename?: 'SubscriberMetricsByNetworkDto';
  active: Scalars['Float']['output'];
  inactive: Scalars['Float']['output'];
  terminated: Scalars['Float']['output'];
  total: Scalars['Float']['output'];
};

export type SubscriberPlansSection = {
  __typename?: 'SubscriberPlansSection';
  error?: Maybe<SectionError>;
  plans?: Maybe<Array<PlanNameDto>>;
};

export type SubscriberSection = {
  __typename?: 'SubscriberSection';
  error?: Maybe<SectionError>;
  subscriber?: Maybe<SubscriberDto>;
};

export type SubscriberSimDto = {
  __typename?: 'SubscriberSimDto';
  allocatedAt: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi: Scalars['String']['output'];
  isPhysical?: Maybe<Scalars['Boolean']['output']>;
  msisdn: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  package?: Maybe<SimPackageDto>;
  status: Scalars['String']['output'];
  subscriberId: Scalars['String']['output'];
  sync_status?: Maybe<Scalars['String']['output']>;
  type: Scalars['String']['output'];
};

export type SubscriberSimsDto = {
  __typename?: 'SubscriberSimsDto';
  allocatedAt: Scalars['String']['output'];
  iccid: Scalars['String']['output'];
  id: Scalars['String']['output'];
  imsi: Scalars['String']['output'];
  isPhysical: Scalars['Boolean']['output'];
  msisdn: Scalars['String']['output'];
  networkId: Scalars['String']['output'];
  status: Scalars['String']['output'];
  subscriberId: Scalars['String']['output'];
  syncStatus: Scalars['String']['output'];
  trafficPolicy: Scalars['Float']['output'];
  type: Scalars['String']['output'];
};

export type SubscriberSimsResDto = {
  __typename?: 'SubscriberSimsResDto';
  sims: Array<SubscriberSimDto>;
};

export type SubscriberStatsSection = {
  __typename?: 'SubscriberStatsSection';
  active?: Maybe<Scalars['Int']['output']>;
  error?: Maybe<SectionError>;
  inactive?: Maybe<Scalars['Int']['output']>;
  total?: Maybe<Scalars['Int']['output']>;
};

export type SubscriberToSimsDto = {
  __typename?: 'SubscriberToSimsDto';
  sims: Array<SubscriberSimsDto>;
  subscriberId: Scalars['String']['output'];
};

export type SubscriberView = {
  __typename?: 'SubscriberView';
  billing: SubscriberBillingSection;
  plans: SubscriberPlansSection;
  subscriber: SubscriberSection;
  subscriberId: Scalars['String']['output'];
  usage: GapSection;
};

export type Subscribers = {
  __typename?: 'Subscribers';
  activeSubscribers: Scalars['String']['output'];
  inactiveSubscribers: Scalars['String']['output'];
  totalSubscribers: Scalars['String']['output'];
};

export type SubscribersResDto = {
  __typename?: 'SubscribersResDto';
  subscribers: Array<SubscriberDto>;
};

export type SubscribersSection = {
  __typename?: 'SubscribersSection';
  error?: Maybe<SectionError>;
  subscribers?: Maybe<Array<SubscriberDto>>;
};

export type SubscribersView = {
  __typename?: 'SubscribersView';
  networkId: Scalars['String']['output'];
  plans: PlansSection;
  subscribers: SubscribersSection;
  usage: GapSection;
};

export type SubscriptionDto = {
  __typename?: 'SubscriptionDto';
  canceledAt?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['String']['output'];
  externalCustomerId: Scalars['String']['output'];
  externalId: Scalars['String']['output'];
  name?: Maybe<Scalars['String']['output']>;
  planCode: Scalars['String']['output'];
  startedAt: Scalars['String']['output'];
  status: Scalars['String']['output'];
  terminatedAt?: Maybe<Scalars['String']['output']>;
};

export type SupportResultDto = {
  __typename?: 'SupportResultDto';
  batteryPercent: Scalars['Float']['output'];
  customers: Scalars['Int']['output'];
  name?: Maybe<Scalars['String']['output']>;
  recommendation?: Maybe<Scalars['String']['output']>;
  resourceId?: Maybe<Scalars['String']['output']>;
  resourceType?: Maybe<Scalars['String']['output']>;
  signalDbm: Scalars['Float']['output'];
  status?: Maybe<Scalars['String']['output']>;
  statusSummary?: Maybe<Scalars['String']['output']>;
  uptime30d: Scalars['Float']['output'];
};

export type SupportSignalDto = {
  __typename?: 'SupportSignalDto';
  detail?: Maybe<Scalars['String']['output']>;
  key: Scalars['String']['output'];
  state?: Maybe<Scalars['String']['output']>;
};

export enum Timeframe_Filter {
  All = 'ALL',
  Latest = 'LATEST',
  Unknown = 'UNKNOWN'
}

export type TeamMemberDto = {
  __typename?: 'TeamMemberDto';
  email?: Maybe<Scalars['String']['output']>;
  id: Scalars['String']['output'];
  inviteExpiresAt?: Maybe<Scalars['String']['output']>;
  memberSince?: Maybe<Scalars['String']['output']>;
  name?: Maybe<Scalars['String']['output']>;
  role: Scalars['String']['output'];
  status: Scalars['String']['output'];
};

export type TeamSection = {
  __typename?: 'TeamSection';
  error?: Maybe<SectionError>;
  rows?: Maybe<Array<TeamMemberDto>>;
};

export type TimeSeriesDto = {
  __typename?: 'TimeSeriesDto';
  key: Scalars['String']['output'];
  points: Array<PointDto>;
};

export type TimezoneDto = {
  __typename?: 'TimezoneDto';
  abbr: Scalars['String']['output'];
  isdst: Scalars['Boolean']['output'];
  offset: Scalars['Float']['output'];
  text: Scalars['String']['output'];
  utc: Array<Scalars['String']['output']>;
  value: Scalars['String']['output'];
};

export type TimezoneRes = {
  __typename?: 'TimezoneRes';
  timezones: Array<TimezoneDto>;
};

export type ToggleInternetSwitchInputDto = {
  port: Scalars['Float']['input'];
  siteId: Scalars['String']['input'];
  status: Scalars['Boolean']['input'];
};

export type ToggleRfStatusInputDto = {
  nodeId: Scalars['String']['input'];
  status: Scalars['Boolean']['input'];
};

export type ToggleSimStatusInputDto = {
  sim_id: Scalars['String']['input'];
  status: Scalars['String']['input'];
};

export type TokenResDto = {
  __typename?: 'TokenResDto';
  token: Scalars['String']['output'];
};

export type TopologyNodeDto = {
  __typename?: 'TopologyNodeDto';
  name?: Maybe<Scalars['String']['output']>;
  nodeId?: Maybe<Scalars['String']['output']>;
  status?: Maybe<Scalars['String']['output']>;
  type?: Maybe<Scalars['String']['output']>;
};

export type TopologySiteDto = {
  __typename?: 'TopologySiteDto';
  latitude: Scalars['Float']['output'];
  longitude: Scalars['Float']['output'];
  name?: Maybe<Scalars['String']['output']>;
  nodes: Array<TopologyNodeDto>;
  siteId?: Maybe<Scalars['String']['output']>;
  status?: Maybe<Scalars['String']['output']>;
};

export type UpdateInvitationInputDto = {
  email: Scalars['String']['input'];
  id: Scalars['String']['input'];
  status: Invitation_Status;
};

export type UpdateInvitationResDto = {
  __typename?: 'UpdateInvitationResDto';
  id: Scalars['String']['output'];
};

export type UpdateMemberInputDto = {
  isDeactivated: Scalars['Boolean']['input'];
  role: Scalars['String']['input'];
};

export type UpdateNodeInput = {
  id: Scalars['String']['input'];
  name: Scalars['String']['input'];
};

export type UpdateNodeStateInput = {
  id: Scalars['String']['input'];
  state: NodeStateEnum;
};

export type UpdateNotificationResDto = {
  __typename?: 'UpdateNotificationResDto';
  id: Scalars['String']['output'];
};

export type UpdatePackageInputDto = {
  active: Scalars['Boolean']['input'];
  name: Scalars['String']['input'];
};

export type UpdatePaymentInputDto = {
  country?: InputMaybe<Scalars['String']['input']>;
  currency?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['String']['input'];
  payerEmail?: InputMaybe<Scalars['String']['input']>;
  payerName?: InputMaybe<Scalars['String']['input']>;
  paymentMethod?: InputMaybe<Scalars['String']['input']>;
};

export type UpdateSiteInputDto = {
  name: Scalars['String']['input'];
};

export type UpdateSoftwareInputDto = {
  name: Scalars['String']['input'];
  nodeId: Scalars['String']['input'];
  tag: Scalars['String']['input'];
};

export type UpdateSubscriberInputDto = {
  address?: InputMaybe<Scalars['String']['input']>;
  email?: InputMaybe<Scalars['String']['input']>;
  id_serial?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  phone?: InputMaybe<Scalars['String']['input']>;
  proof_of_identification?: InputMaybe<Scalars['String']['input']>;
};

export type UploadSimsInputDto = {
  data: Scalars['String']['input'];
  simType: Sim_Types;
};

export type UploadSimsResDto = {
  __typename?: 'UploadSimsResDto';
  iccid: Array<Scalars['String']['output']>;
};

export type UserFirstVisitInputDto = {
  email: Scalars['String']['input'];
  firstVisit: Scalars['Boolean']['input'];
  name: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type UserFirstVisitResDto = {
  __typename?: 'UserFirstVisitResDto';
  firstVisit: Scalars['Boolean']['output'];
};

export type UserResDto = {
  __typename?: 'UserResDto';
  authId: Scalars['String']['output'];
  email: Scalars['String']['output'];
  isDeactivated: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  phone: Scalars['String']['output'];
  registeredSince: Scalars['String']['output'];
  uuid: Scalars['String']['output'];
};

export type WhoamiDto = {
  __typename?: 'WhoamiDto';
  memberOf: Array<OrgDto>;
  ownerOf: Array<OrgDto>;
  user: UserResDto;
};
