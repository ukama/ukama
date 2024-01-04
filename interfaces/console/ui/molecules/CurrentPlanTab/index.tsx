/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { CurrentBillColumns } from '@/constants/tableColumns';
import { NoBillYet } from '@/public/svg';
import { RoundedCard } from '@/styles/global';
import { Grid, Stack, Typography } from '@mui/material';
import CurrentBill from '../CurrentBill';
import NotificationContainer from '../NotificationContainer';
import SimpleDataTable from '../SimpleDataTable';
import TableHeader from '../TableHeader';

interface ICurrentBillTab {
  data: any;
  loading: boolean;
  planName: string;
  totalAmount: string;
  currentBill: string;
}

const CurrentPlanTab = () => {
  return (
    <Grid container item spacing={2}>
      <Grid xs={12} item>
        <RoundedCard radius="4px">
          <Typography variant="body1" color="initial"></Typography>
          <TableHeader title={'Current plan'} showSecondaryButton={false} />
          <Grid container spacing={2} sx={{ py: 2 }}>
            <Grid item xs={6}>
              <Typography variant="body2" color="initial">
                Community bundle - free Console plan for basic network
                management needs.
              </Typography>
            </Grid>
            <Grid item xs={6} container justifyContent={'flex-end'}>
              <Typography variant="h6" color="initial">
                $ 20.30
              </Typography>
            </Grid>
          </Grid>
        </RoundedCard>
      </Grid>
    </Grid>
  );
};

export default CurrentPlanTab;
