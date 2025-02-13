/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { duration } from '@/utils';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import SignalCellularOffIcon from '@mui/icons-material/SignalCellularOff';
import {
  Stack,
  Tooltip,
  Typography,
  useMediaQuery,
  useTheme,
} from '@mui/material';
import LoadingWrapper from '../LoadingWrapper';

const getIconByStatus = (status: string, info: string) => {
  switch (status) {
    case 'DOWN':
      return (
        <Tooltip title={info}>
          <CheckCircleIcon sx={{ color: colors.red }} />
        </Tooltip>
      );

    case 'ONLINE':
      return (
        <Tooltip title={info}>
          <CheckCircleIcon sx={{ color: colors.green }} />
        </Tooltip>
      );

    default:
      return (
        <Tooltip title={info}>
          <SignalCellularOffIcon />
        </Tooltip>
      );
  }
};

type NetworkStatusProps = {
  title: string;
  subtitle: number;
  loading?: boolean;
  tooltipInfo: string;
  availableNodes: number | undefined;
  statusType: string | undefined;
};

const NetworkStatus = ({
  title,
  subtitle,
  loading,
  tooltipInfo,
  availableNodes = undefined,
  statusType = 'onboarding',
}: NetworkStatusProps) => {
  const theme = useTheme();
  const matches = useMediaQuery(theme.breakpoints.down('md'));
  return (
    <LoadingWrapper
      isLoading={loading}
      height={loading ? '30px' : 'fit-content'}
      width={loading ? '40%' : 'fit-content'}
    >
      <Stack direction={'row'} alignItems="center" mt={{ xs: 0.5, md: 1.5 }}>
        {getIconByStatus(statusType, tooltipInfo)}
        <Typography
          variant={matches ? 'subtitle1' : 'h6'}
          sx={{ fontWeight: { xs: 400, md: 500 }, ml: 1 }}
        >
          {title} {subtitle === 0 ? '-' : duration(subtitle)}
        </Typography>
        {/* <Typography variant={'h6'} sx={{ fontWeight: { xs: 400, md: 500 } }}>
          {subtitle}
        </Typography> */}
        {availableNodes && (
          <Typography
            variant={'h6'}
            color="secondary"
            sx={{
              fontWeight: {
                xs: 400,
                md: 500,
              },
              whiteSpace: 'break-spaces',
              ml: { xs: '28px', md: '8px' },
            }}
          >
            {`${availableNodes} available nodes`}
          </Typography>
        )}
      </Stack>
    </LoadingWrapper>
  );
};

export default NetworkStatus;
