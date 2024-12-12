/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Paper, Typography, Stack, Box, Skeleton } from '@mui/material';
import colors from '@/theme/colors';

interface DataUsageProps {
  dataUsagePaid?: number;
  subscriberCount?: number;
  loading?: boolean;
}

const DataUsageComponent: React.FC<DataUsageProps> = ({
  dataUsagePaid,
  subscriberCount,
  loading = false,
}) => {
  const DataRow = ({ label, value, valuePrefix = '', isLoading }: any) => (
    <Stack
      direction="row"
      justifyContent="space-between"
      alignItems="center"
      sx={{
        width: '100%',
        mb: 1,
      }}
    >
      {isLoading ? (
        <>
          <Skeleton width="40%" height={20} />
          <Skeleton width="20%" height={24} />
        </>
      ) : (
        <>
          <Typography variant="body2" sx={{ color: colors.black54 }}>
            {label}
          </Typography>
          <Typography variant="h6" sx={{ fontWeight: 600 }}>
            {valuePrefix}
            {value}
          </Typography>
        </>
      )}
    </Stack>
  );

  return (
    <Paper
      elevation={2}
      sx={{
        p: 4,
        mt: 2,
        borderRadius: '10px',
        bgcolor: colors.white,
      }}
    >
      {loading ? (
        <Skeleton width="60%" height={32} sx={{ mb: 2 }} />
      ) : (
        <Typography variant="h6" sx={{ mb: 2 }}>
          Feature usage
        </Typography>
      )}

      <Box sx={{ width: '100%' }}>
        <DataRow
          label="Data usage paid for by subscribers"
          value={dataUsagePaid}
          valuePrefix="$ "
          isLoading={loading}
        />
        <DataRow
          label="Total subscribers"
          value={subscriberCount}
          isLoading={loading}
        />
      </Box>
    </Paper>
  );
};

export default DataUsageComponent;
