/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import {
  Invitation_Status,
  NodeTypeEnum,
  Role_Type,
} from '@/client/graphql/generated';
import {
  Graphs_Type,
  MetricRes,
  MetricsRes,
} from '@/client/graphql/generated/subscriptions';
import { INSTALLATION_FLOW, ONBOARDING_FLOW } from '@/constants';
import colors from '@/theme/colors';
import { TNodeSiteTree } from '@/types';
import ArrowForwardIosIcon from '@mui/icons-material/ArrowForwardIos';
import { Skeleton, Stack, Typography } from '@mui/material';
import { LatLngTuple } from 'leaflet';

type TConfigureStep = {
  totalStep: number;
  currentStep: number;
};

const getTitleFromPath = (path: string, id: string) => {
  if (path.startsWith('/console/sites') && id) {
    return (
      <Stack direction="row" alignItems="center" spacing={1}>
        <Typography variant="h5" sx={{ color: colors.black38 }}>
          Site
        </Typography>
        <ArrowForwardIosIcon sx={{ color: colors.black38 }} />
        <Typography variant="h5" sx={{ color: colors.black }}>
          {id}
        </Typography>
      </Stack>
    );
  }

  if (path.startsWith('/console/nodes') && id) {
    return (
      <Stack direction="row" alignItems="center" spacing={1}>
        <Typography variant="h5" sx={{ color: colors.black38 }}>
          Nodes
        </Typography>
        <ArrowForwardIosIcon sx={{ color: colors.black38 }} />
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

const getTabByIndex = (index: number) => {
  switch (index) {
    case 0:
      return 'Graphs_Tab.Overview';
    case 1:
      return 'Graphs_Tab.Network';
    case 2:
      return 'Graphs_Tab.Resources';
    case 3:
      return 'Graphs_Tab.Radio';
    case 4:
      return 'Graphs_Tab.Home';
    default:
      return 'Graphs_Tab.Overview';
  }
};

const getTitleByKey = (key: string) => {
  switch (key) {
    case 'uptime_trx':
      return 'Uptime TRX';
    case 'temperaturetrx':
      return 'Temp. (TRX)';
    case 'temperaturerfe':
      return 'Temp. (RFE)';
    case 'subscribersactive':
      return 'Active';
    case 'subscribersattached':
      return 'Attached';
    case 'temperaturectl':
      return 'Temp. (CTL)';
    case 'temperaturecom':
      return 'Temp. (COM)';
    case 'rrc':
      return 'RRC CNX success';
    case 'rlc':
      return 'RLS  drop rate';
    case 'erab':
      return 'ERAB drop rate';
    case 'throughputuplink':
      return 'Throughput (U/L)';
    case 'throughputdownlink':
      return 'Throughput (D/L)';
    case 'cputrxusage':
      return 'CPU-TRX';
    case 'memorytrxused':
      return 'Memory-TRX';
    case 'disktrxused':
      return 'DISK-TRX';
    case 'cpuctlused':
      return 'CPU-CTL';
    case 'diskctlused':
      return 'DISK-CTL';
    case 'memoryctlused':
      return 'Memory-CTL';
    case 'powerlevel':
      return 'Power';
    case 'cpucomusage':
      return 'CPU-COM';
    case 'diskcomused':
      return 'DISK-COM';
    case 'memorycomused':
      return 'Memory-COM';
    case 'txpower':
      return 'TX Power';
    case 'rxpower':
      return 'RX Power';
    case 'papower':
      return 'PA Power';
    default:
      return '';
  }
};

export const getNodeTabTypeByIndex = (index: number) => {
  switch (index) {
    case 0:
      return Graphs_Type.NodeHealth;
    case 1:
      return Graphs_Type.Network;
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

const formatBytesToMB = (bytes = 0): string => {
  if (bytes === 0) return '0';
  return (bytes / (1024 * 1024)).toFixed(2);
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

const structureNodeSiteDate = (data: any) => {
  let count = 1;
  const t: TNodeSiteTree[] = [];

  data.forEach((node: any) => {
    if (node.type === NodeTypeEnum.Tnode) {
      t.push({
        id: node.site?.siteId ?? '',
        name: `Site ${count++}`,
        nodeId: node.id,
        nodeType: node.type,
        nodeName: node.name,
      });
    }
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

export {
  ConfigureStep,
  fileToBase64,
  formatBytes,
  formatBytesToMB,
  formatTime,
  getDataPlanUsage,
  getDuration,
  getGraphFilterByType,
  getInvitationStatusColor,
  getSimValuefromSimType,
  getTitleFromPath,
  getUnixTime,
  hexToRGB,
  inviteStatusEnumToString,
  isValidLatLng,
  NodeEnumToString,
  provideStatusColor,
  roleEnumToString,
  structureNodeSiteDate,
};
