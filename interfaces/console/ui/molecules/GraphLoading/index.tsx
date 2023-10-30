/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/styles/theme';
import { Skeleton, Stack } from '@mui/material';

const GraphLoading = () => {
  return (
    <Stack
      spacing={1}
      direction={'row'}
      alignItems="flex-end"
      justifyContent="center"
    >
      <Skeleton
        variant="rectangular"
        width={10}
        height={24}
        sx={{ backgroundColor: colors.vulcan70 }}
      />
      <Skeleton
        variant="rectangular"
        width={10}
        height={34}
        sx={{ backgroundColor: colors.vulcan70 }}
      />
      <Skeleton
        variant="rectangular"
        width={10}
        height={44}
        sx={{ backgroundColor: colors.vulcan70 }}
      />
      <Skeleton
        variant="rectangular"
        width={10}
        height={40}
        sx={{ backgroundColor: colors.vulcan70 }}
      />
    </Stack>
  );
};

export default GraphLoading;
