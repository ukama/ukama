/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import BarChartIcon from '@mui/icons-material/BarChart';
import { Typography } from '@mui/material';
import Grid from '@mui/material/Grid2';
import { Variant } from '@mui/material/styles/createTypography';
import React from 'react';
import EmptyView from '../EmptyView';
import GraphLoading from '../GraphLoading';
import TimeFilter from '../TimeFilter';

interface IGraphTitleWrapper {
  title?: string;
  hasData?: boolean;
  loading?: boolean;
  variant?: Variant;
  showFilter?: boolean;
  children: React.ReactNode;
  handleFilterChange?: (value: string) => void;
}

const GraphTitleWrapper = ({
  children,
  title = '',
  hasData = false,
  loading = true,
  showFilter = true,
  variant = 'subtitle1',
  handleFilterChange = undefined,
}: IGraphTitleWrapper) => {
  const GTChild = hasData ? (
    children
  ) : (
    <EmptyView size="large" title="No activity yet!" icon={BarChartIcon} />
  );
  return (
    <Grid container width="100%">
      {(title ?? showFilter) && (
        <Grid container width="100%" mb={2}>
          {title && (
            <Grid size={{ xs: 6 }}>
              <Typography variant={variant} fontWeight={600}>
                {title}
              </Typography>
            </Grid>
          )}
          {hasData && showFilter && (
            <Grid size={{ xs: 6 }} display="flex" justifyContent="flex-end">
              <TimeFilter
                handleFilterSelect={(v: string) => handleFilterChange?.(v)}
              />
            </Grid>
          )}
        </Grid>
      )}
      <Grid
        container
        width={'100%'}
        height={'400px'}
        alignItems={'center'}
        justifyContent="center"
      >
        {loading ? <GraphLoading /> : GTChild}
      </Grid>
    </Grid>
  );
};

export default GraphTitleWrapper;
