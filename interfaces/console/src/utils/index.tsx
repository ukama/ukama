/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import {
  Invitation_Status,
  Node,
  NodeConnectivityEnum,
  Nodes,
  NodeTypeEnum,
  Role_Type,
  SitesResDto,
} from '@/client/graphql/generated';
import {
  Graphs_Type,
  MetricRes,
  MetricsRes,
  MetricsStateRes,
  SiteMetricsStateRes,
} from '@/client/graphql/generated/subscriptions';
import {
  INSTALLATION_FLOW,
  KPI_PLACEHOLDER_VALUE,
  NODE_ACTIONS_ENUM,
  ONBOARDING_FLOW,
} from '@/constants';
import colors from '@/theme/colors';
import { StatusType, StyleOutput, TNodeSiteTree } from '@/types';
import Battery50Icon from '@mui/icons-material/Battery50';
import BatteryAlertIcon from '@mui/icons-material/BatteryAlert';
import BatteryChargingFullIcon from '@mui/icons-material/BatteryChargingFull';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import RouterIcon from '@mui/icons-material/Router';
import SignalCellular1BarIcon from '@mui/icons-material/SignalCellular1Bar';
import SignalCellular2BarIcon from '@mui/icons-material/SignalCellular2Bar';
import SignalCellularAltIcon from '@mui/icons-material/SignalCellularAlt';
import SignalCellularConnectedNoInternet4BarIcon from '@mui/icons-material/SignalCellularConnectedNoInternet4Bar';
import SignalCellularOffIcon from '@mui/icons-material/SignalCellularOff';
import { Skeleton, Stack, Typography } from '@mui/material';
import { formatDistance } from 'date-fns';
import { DashStyleValue } from 'highcharts';
import { LatLngTuple } from 'leaflet';

type TConfigureStep = {
  totalStep: number;
  currentStep: number;
};

const getTitleFromPath = (path: string, id: string) => {
  if (id) {
    return (
      <Stack direction="row" alignItems="center" spacing={0.6}>
        <Typography variant="h5" sx={{ color: colors.black38 }}>
          {path.startsWith('/console/sites')
            ? 'Site'
            : path.startsWith('/console/nodes')
              ? 'Nodes'
              : ''}
        </Typography>
        <ChevronRightIcon sx={{ color: colors.black38 }} />
        <Typography variant="h5" sx={{ color: colors.black }}>
          {id}
        </Typography>
      </Stack>
    );
  }

  switch (path) {
    case '/console/home':
      return 'Home';
    case '/settings':
      return 'Settings';
    case '/console/sites':
      if (id) return `Site -> ${id}`;
      return 'Sites';
    case '/console/nodes':
      return 'Nodes';
    case '/console/subscribers':
      return 'Subscribers';
    // case '/site_planning':
    //   return 'Site Planning';
    case '/manage':
      return 'Manage';
    case '/onboarding':
      return 'OnBoarding';
    case '/unauthorized':
      return 'Unauthorized';
    case '/ping':
      return 'Ping';
    default:
      return <Skeleton variant="text" width={100} />;
  }
};

const hexToRGB = (hex: string, alpha: number): string => {
  const h = '0123456789ABCDEF';
  const r = h.indexOf(hex[1]) * 16 + h.indexOf(hex[2]);
  const g = h.indexOf(hex[3]) * 16 + h.indexOf(hex[4]);
  const b = h.indexOf(hex[5]) * 16 + h.indexOf(hex[6]);
  if (alpha) {
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
  }

  return `rgba(${r}, ${g}, ${b})`;
};

const getGraphFilterByType = (type: string) => {
  switch (type) {
    case 'DAY':
      return {
        to: Math.round(Date.now() / 1000),
        from: Math.round(Date.now() / 1000) - 86400,
      };
    case 'WEEK':
      return {
        to: Math.round(Date.now() / 1000),
        from: Math.round(Date.now() / 1000) - 604800,
      };
    case 'MONTH':
      return {
        to: Math.round(Date.now() / 1000),
        from: Math.round(Date.now() / 1000) - 2628002,
      };
  }
};

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

const formatBytes = (bytes = 0): string => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const dm = 3;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ` ${sizes[i]}`;
};

/**
 * Converts bytes to gigabytes (GB) and returns a string.
 * @param bytes Number of bytes to convert.
 * @param decimals Number of decimal places to show (default is 2).
 * @returns String representation in GB.
 */
const formatBytesToGB = (bytes = 0, decimals = 2): string => {
  if (bytes === 0) return '0';
  const gb = bytes / (1024 * 1024 * 1024);
  return gb.toFixed(decimals);
};

const formatBytesToGigabit = (bytes = 0, decimals = 2): string => {
  if (bytes === 0) return '0';
  // Convert bytes to bits (multiply by 8) then to gigabits
  const gigabit = (bytes * 8) / (1024 * 1024 * 1024);
  return gigabit.toFixed(decimals);
};

const getDataUsageSymbol = (dataUnit: string): string => {
  switch (dataUnit) {
    case 'GigaBytes':
      return 'GB';
    case 'MegaBytes':
      return 'MB';
    case 'KiloBytes':
      return 'KB';
    default:
      return 'MB';
  }
};

const getPortInfo: Record<string, { number: number; desc: string }> = {
  solar: { number: 2, desc: 'Solar Controller' },
  backhaul: { number: 3, desc: 'Backhaul' },
  node: { number: 1, desc: 'Node' },
};
const getDataPlanUsage = (
  duration: string,
  currency: string,
  amount: string,
  dataVolume: string,
  dataUnit: string,
): string => {
  const symbol = currency === 'Dollar' ? '$' : currency;
  return `${symbol} ${amount} / ${dataVolume} ${getDataUsageSymbol(
    dataUnit,
  )} / ${duration}`;
};

const fileToBase64 = (file: File): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => {
      const base64String = reader.result as string;
      resolve(base64String.split(',')[1]);
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
};

const getUnixTime = (): number => {
  return Math.floor(Date.now() / 1000);
};

const getDuration = (number: number): string => {
  return number > 1 ? `${number} Days` : `${number} Day`;
};

const structureNodeSiteDate = (nodes: Nodes, sites: SitesResDto) => {
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

export const getMetricValue = (key: string, metrics: MetricsRes) => {
  const metric = metrics.metrics.find((item: MetricRes) => item.type === key);
  return metric?.values ?? [];
};

export const isMetricValue = (key: string, metrics: MetricsRes) => {
  const metric = metrics.metrics.find((item: MetricRes) => item.type === key);
  return (metric && metric.values.length > 1) ?? false;
};

const getSimValuefromSimType = (simType: string) => {
  switch (simType) {
    case 'operator_data':
      return 'Operator Data';
    case 'ukama_data':
      return 'Ukama Data';
    case 'test':
      return 'Test';
    default:
      return 'Unknown';
  }
};

const getInvitationStatusColor = (status: string, isExpired: boolean) => {
  if (isExpired) {
    return (
      <Typography variant="body2" color={colors.red}>
        Expired
      </Typography>
    );
  }

  switch (status) {
    case Invitation_Status.InviteAccepted:
      return (
        <Typography variant="body2" color={colors.green}>
          Accepted
        </Typography>
      );
    case Invitation_Status.InviteDeclined:
      return (
        <Typography variant="body2" color={colors.red}>
          Declined
        </Typography>
      );
    case Invitation_Status.InvitePending:
      return (
        <Typography variant="body2" color={colors.yellow}>
          Pending
        </Typography>
      );
  }
};

const provideStatusColor = (status: Invitation_Status) => {
  switch (status) {
    case Invitation_Status.InvitePending:
      return colors.blueGray;
    case Invitation_Status.InviteAccepted:
      return 'success';
    case Invitation_Status.InviteDeclined:
      return 'warning';
    default:
      return 'info';
  }
};

const formatTime = (isoString: string) => {
  const date = new Date(isoString);
  const day = date.getDate().toString().padStart(2, '0');
  const month = (date.getMonth() + 1).toString();
  const hours = date.getHours();
  const period = hours >= 12 ? 'PM' : 'AM';
  const formattedHours = (hours % 12).toString();
  return `${month}/${day} ${formattedHours}${period}`;
};

const roleEnumToString = (role: Role_Type): string => {
  switch (role) {
    case Role_Type.RoleOwner:
      return 'OWNER';
    case Role_Type.RoleAdmin:
      return 'ADMIN';
    case Role_Type.RoleNetworkOwner:
      return 'NETWORK OWNER';
    case Role_Type.RoleVendor:
      return 'VENDOR';
    case Role_Type.RoleUser:
      return 'USER';
    default:
      return 'Invalid';
  }
};

const inviteStatusEnumToString = (status: Invitation_Status): string => {
  switch (status) {
    case Invitation_Status.InviteAccepted:
      return 'ACCEPTED';
    case Invitation_Status.InviteDeclined:
      return 'DECLINED';
    case Invitation_Status.InvitePending:
      return 'PENDING';
  }
};

function isValidLatitude(latitude: number) {
  if (typeof latitude !== 'number' || latitude < -90 || latitude > 90) {
    return false;
  }
  return true;
}

function isValidLongitude(longitude: number) {
  if (typeof longitude !== 'number' || longitude < -180 || longitude > 180) {
    return false;
  }
  return true;
}

const isValidLatLng = (position: LatLngTuple): boolean => {
  const [latitude, longitude] = position || [0, 0];
  return (
    latitude !== 0 &&
    longitude !== 0 &&
    !isNaN(latitude) &&
    !isNaN(longitude) &&
    isValidLatitude(latitude) &&
    isValidLongitude(longitude)
  );
};

const ConfigureStep = (path: string, flow: string): TConfigureStep => {
  switch (flow) {
    case ONBOARDING_FLOW:
      if (path.includes('configure/network'))
        return { currentStep: 1, totalStep: 6 };
      else if (path.includes('check')) return { currentStep: 1, totalStep: 6 };
      else if (path.includes('sims')) return { currentStep: 6, totalStep: 6 };
      else if (path.includes('install'))
        return { currentStep: 5, totalStep: 6 };
      else if (path.includes('site/name'))
        return { currentStep: 4, totalStep: 6 };
      else if (path.includes('configure/node'))
        return { currentStep: 3, totalStep: 6 };
      else if (path.includes('configure'))
        return { currentStep: 2, totalStep: 6 };
    case INSTALLATION_FLOW:
      if (path.includes('check')) return { currentStep: 1, totalStep: 4 };
      else if (path.includes('install'))
        return { currentStep: 4, totalStep: 4 };
      else if (path.includes('site/name'))
        return { currentStep: 3, totalStep: 4 };
      else if (path.includes('configure/node'))
        return { currentStep: 2, totalStep: 4 };
    default:
      return { currentStep: 1, totalStep: 5 };
  }
};

const NodeEnumToString = (type: NodeTypeEnum): string => {
  switch (type) {
    case NodeTypeEnum.Tnode:
      return 'Tower Node';
    case NodeTypeEnum.Anode:
      return 'Amplifier Node';
    case NodeTypeEnum.Hnode:
      return 'Home Node';
    default:
      return 'Unknown';
  }
};

const getKPIStatValue = (
  id: string,
  loading: boolean,
  statsData: MetricsStateRes | SiteMetricsStateRes,
): string => {
  if (loading || !statsData?.metrics) return KPI_PLACEHOLDER_VALUE;
  const stat = statsData.metrics.find((item) => item.type === id);
  return stat?.value?.toString() ?? KPI_PLACEHOLDER_VALUE;
};

const base64ToBlob = (base64: string, contentType = ''): Blob => {
  const byteCharacters = atob(base64.split(',')[1] || base64);
  const byteArrays = [];

  for (let offset = 0; offset < byteCharacters.length; offset += 512) {
    const slice = byteCharacters.slice(offset, offset + 512);
    const byteNumbers = new Array(slice.length);

    for (let i = 0; i < slice.length; i++) {
      byteNumbers[i] = slice.charCodeAt(i);
    }

    byteArrays.push(new Uint8Array(byteNumbers));
  }

  return new Blob(byteArrays, { type: contentType });
};

export const duration = (s: number) =>
  formatDistance(0, s * 1000, { includeSeconds: true });

const findNullZones = (data: any) => {
  const zones = [];
  let inNullZone = false;
  let start = null;

  for (let i = 0; i < data.length; i++) {
    const [x, y] = data[i];

    if (y === null) {
      if (!inNullZone) {
        start = x;
        inNullZone = true;
      }
    } else {
      if (inNullZone) {
        zones.push({ value: start });
        zones.push({
          value: data[i - 1][0],
          color: colors.black38,
          dashStyle: 'dash' as DashStyleValue,
        });
        inNullZone = false;
      }
    }
  }

  if (inNullZone) {
    zones.push({ value: start });
    zones.push({
      value: data[data.length - 1][0],
      color: colors.black38,
      dashStyle: 'dash' as DashStyleValue,
    });
  }

  return zones;
};

const kpiToGraphType: Record<string, Graphs_Type> = {
  solar: Graphs_Type.Solar,
  battery: Graphs_Type.Battery,
  controller: Graphs_Type.Controller,
  main_backhaul: Graphs_Type.MainBackhaul,
  backhaul: Graphs_Type.MainBackhaul,
  switch: Graphs_Type.Switch,
  node: Graphs_Type.Site,
};
const graphTypeToSection: Record<Graphs_Type | string, string> = {
  [Graphs_Type.Solar]: 'SOLAR',
  [Graphs_Type.Battery]: 'BATTERY',
  [Graphs_Type.Controller]: 'CONTROLLER',
  [Graphs_Type.MainBackhaul]: 'MAIN_BACKHAUL',
  [Graphs_Type.Switch]: 'SWITCH',
  [Graphs_Type.Site]: 'SITE',
};

export const generatePlotLines = (values: number[] | undefined): any[] => {
  if (!values) {
    return [];
  }
  if (values.length < 3 || values.length > 7) {
    throw new Error('invalid length');
  }

  return values.slice(1).map((value, index, arr) => ({
    value,
    color:
      index === 0
        ? colors.dullGrey
        : index === arr.length - 2
          ? colors.dullRed
          : index === arr.length - 1
            ? colors.white
            : colors.dullGreen,
    width: 2,
    zIndex: 4,
    dashStyle: 'Dash',
  }));
};

export const formatKPIValue = (value: string, type: string): any => {
  switch (type) {
    case 'number':
      return Math.floor(parseFloat(value));
    case 'decimal':
      return parseFloat(value).toFixed(2);
    default:
      return value.toString();
  }
};

const getConnectionStyles = (connectionStatus: string) => {
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

const getBatteryStyles = (batteryStatus: string) => {
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

const getSignalStyles = (signalStrength: string) => {
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

const getStatusStyles = (type: StatusType, value: number): StyleOutput => {
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

const setQueryParam = (
  key: string,
  value: string,
  params: string,
  pathname: string,
): URLSearchParams => {
  const p = new URLSearchParams(params);
  p.set(key, value);
  window.history.replaceState({}, '', `${pathname}?${p.toString()}`);
  return p;
};
const getSectionFromKPI = (kpi: string) => {
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

const getNodeActionDescriptionByProgress = (
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

export {
  base64ToBlob,
  ConfigureStep,
  fileToBase64,
  findNullZones,
  formatBytes,
  formatBytesToGB,
  formatBytesToGigabit,
  formatTime,
  getBatteryStyles,
  getConnectionStyles,
  getDataPlanUsage,
  getDuration,
  getGraphFilterByType,
  getInvitationStatusColor,
  getKPIStatValue,
  getNodeActionDescriptionByProgress,
  getPortInfo,
  getSectionFromKPI,
  getSignalStyles,
  getSimValuefromSimType,
  getStatusStyles,
  getTitleFromPath,
  getUnixTime,
  graphTypeToSection,
  hexToRGB,
  inviteStatusEnumToString,
  isValidLatLng,
  kpiToGraphType,
  NodeEnumToString,
  provideStatusColor,
  roleEnumToString,
  setQueryParam,
  structureNodeSiteDate,
};
