/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Card, Grid } from '@mui/material';

import SiteFlowDiagram from '../../../public/svg/sitecomps';
import LineChart from '../LineChart';

const SiteComponents: React.FC = () => {
  return (
    <Box>
      <Card
        sx={{
          borderRadius: 2,
          boxShadow: '0px 2px 6px rgba(0, 0, 0, 0.05)',
          width: '100%',
        }}
      >
        <Grid container>
          <Grid item xs={6} sx={{ p: 2 }}>
            <SiteFlowDiagram />
          </Grid>
          <Grid item xs={6} sx={{ p: 2 }}>
            <LineChart topic={''} initData={undefined} from={0} />
          </Grid>
        </Grid>
      </Card>
    </Box>
  );
};

export default SiteComponents;
