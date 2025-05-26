/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { colors } from '@/theme';
import Skeleton from '@mui/material/Skeleton';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import React from 'react';

interface Props {
  packageName: string | undefined;
  currentSite: string | undefined;
  bundle: string | undefined;
}

const DataPlanComponent: React.FC<Props> = ({ packageName, bundle }) => {
  return (
    <Stack direction="column" spacing={2}>
      <Stack direction="row" spacing={2}>
        <Typography variant="body1" sx={{ color: colors.black }}>
          Data plan:
        </Typography>
        <Typography variant="subtitle1" fontWeight={600} color="text.primary">
          {packageName && packageName.length ? (
            packageName
          ) : (
            <Skeleton
              variant="rectangular"
              width={120}
              height={24}
              sx={{ backgroundColor: colors.black10 }}
            />
          )}
        </Typography>
      </Stack>
      {/* 
      TODO: Need more discussion
        <Stack direction="row" spacing={2}>
        <Typography variant="body1" sx={{ color: colors.black }}>
          Current site:
        </Typography>
        <Typography variant="subtitle1" sx={{ color: colors.black }}>
          {currentSite ?? ''}
        </Typography>
      </Stack> */}
      <Stack direction="row" spacing={2}>
        <Typography variant="body1" sx={{ color: colors.black }}>
          Month usage:
        </Typography>
        <Typography variant="subtitle1" sx={{ color: colors.black }}>
          {bundle && bundle.length ? (
            bundle
          ) : (
            <Skeleton
              variant="rectangular"
              width={120}
              height={24}
              sx={{ backgroundColor: colors.black10 }}
            />
          )}
        </Typography>
      </Stack>
    </Stack>
  );
};

export default DataPlanComponent;
