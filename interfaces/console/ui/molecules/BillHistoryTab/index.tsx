/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { HistoryBillingColumns } from '@/constants/tableColumns';
import { NoBillYet } from '@/public/svg';
import { RoundedCard } from '@/styles/global';
import colors from '@/styles/theme/colors';
import { Box, Stack, Typography } from '@mui/material';
import LoadingWrapper from '../LoadingWrapper';
import SimpleDataTable from '../SimpleDataTable';
import TableHeader from '../TableHeader';

interface IBillHistoryTab {
  loading: boolean;
  data: any;
}

const BillHistoryTab = ({ loading, data }: IBillHistoryTab) => {
  return (
    <LoadingWrapper
      height={'100%'}
      isLoading={loading}
      cstyle={{
        overflow: 'auto',
        backgroundColor: loading ? colors.white : 'transparent',
      }}
    >
      <RoundedCard radius="4px">
        <TableHeader title={'Billing history'} showSecondaryButton={false} />
        {data.length > 0 ? (
          <SimpleDataTable columns={HistoryBillingColumns} dataset={data} />
        ) : (
          <Box
            display="flex"
            justifyContent="center"
            alignItems="center"
            minHeight="60vh"
          >
            <Stack direction="column" spacing={2}>
              <NoBillYet color={colors.silver} color2={colors.white} />
              <Typography variant="body1">No bill History yet!</Typography>
            </Stack>
          </Box>
        )}
      </RoundedCard>
    </LoadingWrapper>
  );
};

export default BillHistoryTab;
