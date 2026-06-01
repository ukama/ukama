/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Invitation_Status, Role_Type } from '@/client/graphql/generated';
import {
  INSTALLATION_FLOW,
  ONBOARDING_FLOW,
} from '@/constants';
import colors from '@/theme/colors';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { Skeleton, Stack, Typography } from '@mui/material';
import { LatLngTuple } from 'leaflet';

type TConfigureStep = {
  totalStep: number;
  currentStep: number;
};

export const getTitleFromPath = (path: string, id: string) => {
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

export const formatBytes = (bytes = 0): string => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const dm = 3;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ` ${sizes[i]}`;
};

export const formatBytesToGB = (bytes = 0, decimals = 2): string => {
  if (bytes === 0) return '0';
  const gb = bytes / (1024 * 1024 * 1024);
  return gb.toFixed(decimals);
};

export const formatBytesToGigabit = (bytes = 0, decimals = 2): string => {
  if (bytes === 0) return '0';
  const gigabit = (bytes * 8) / (1024 * 1024 * 1024);
  return gigabit.toFixed(decimals);
};

export const getDataUsageSymbol = (dataUnit: string): string => {
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

export const getPortInfo: Record<string, { number: number; desc: string }> = {
  solar: { number: 2, desc: 'Solar Controller' },
  backhaul: { number: 3, desc: 'Backhaul' },
  node: { number: 1, desc: 'Node' },
};

export const getDataPlanUsage = (
  duration: string,
  currency: string,
  amount: string,
  dataVolume: string,
  dataUnit: string,
): string => {
  const symbol = currency === 'Dollar' ? '$' : currency;
  return `${symbol} ${amount} / ${dataVolume} ${getDataUsageSymbol(dataUnit)} / ${duration}`;
};

export const fileToBase64 = (file: File): Promise<string> => {
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

export const getUnixTime = (): number => {
  return Math.floor(Date.now() / 1000);
};

export const getDuration = (number: number): string => {
  return number > 1 ? `${number} Days` : `${number} Day`;
};

export const getSimValuefromSimType = (simType: string) => {
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

export const getInvitationStatusColor = (status: string, isExpired: boolean) => {
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

export const provideStatusColor = (status: Invitation_Status) => {
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

export const formatTime = (isoString: string) => {
  const date = new Date(isoString);
  const day = date.getDate().toString().padStart(2, '0');
  const month = (date.getMonth() + 1).toString();
  const hours = date.getHours();
  const period = hours >= 12 ? 'PM' : 'AM';
  const formattedHours = (hours % 12).toString();
  return `${month}/${day} ${formattedHours}${period}`;
};

export const roleEnumToString = (role: Role_Type): string => {
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

export const inviteStatusEnumToString = (status: Invitation_Status): string => {
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

export const isValidLatLng = (position: LatLngTuple): boolean => {
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

export const ConfigureStep = (path: string, flow: string): TConfigureStep => {
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

export const base64ToBlob = (base64: string, contentType = ''): Blob => {
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

export const duration = (s: number | string | null | undefined) => {
  const normalizedInput =
    typeof s === 'number' && Number.isFinite(s)
      ? s
      : Number(
          String(s ?? '')
            .trim()
            .replaceAll(',', ''),
        );
  const totalSeconds = Number.isFinite(normalizedInput)
    ? Math.max(0, Math.floor(normalizedInput))
    : 0;
  const units: Array<[label: string, value: number]> = [
    ['day', Math.floor(totalSeconds / 86400)],
    ['hour', Math.floor((totalSeconds % 86400) / 3600)],
    ['minute', Math.floor((totalSeconds % 3600) / 60)],
    ['second', totalSeconds % 60],
  ];

  const parts = units
    .filter(([label, value]) => value > 0 || label === 'second')
    .map(([label, value]) => `${value} ${label}${value === 1 ? '' : 's'}`)
    .filter(
      (part, index, list) => part !== '0 seconds' || index === list.length - 1,
    );

  return parts.join(' ');
};

export const stringToBoolean = (value: string): boolean => {
  return (
    value === 'true' ||
    value === '1' ||
    value === 'on' ||
    value === 'yes' ||
    value === 'enabled'
  );
};

export const setQueryParam = (
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

export const getMapStyleURL = (style: string, token: string) => {
  switch (style) {
    case 'terrain':
      return `https://api.mapbox.com/styles/v1/mapbox/outdoors-v11/tiles/256/{z}/{x}/{y}@2x?access_token=${token}`;
    case 'satellite':
      return `https://api.mapbox.com/styles/v1/mapbox/satellite-v9/tiles/256/{z}/{x}/{y}@2x?access_token=${token}`;
    case 'streets':
      return `https://api.mapbox.com/styles/v1/mapbox/streets-v11/tiles/256/{z}/{x}/{y}@2x?access_token=${token}`;
    case 'light':
      return `https://api.mapbox.com/styles/v1/mapbox/light-v10/tiles/256/{z}/{x}/{y}@2x?access_token=${token}`;
    case 'dark':
      return `https://api.mapbox.com/styles/v1/mapbox/dark-v10/tiles/256/{z}/{x}/{y}@2x?access_token=${token}`;
  }
};
