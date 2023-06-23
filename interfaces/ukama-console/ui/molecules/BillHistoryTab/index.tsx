import { HistoryBillingColumns } from '@/constants/tableColumns';
import { NoBillYet } from '@/public/svg';
import { RoundedCard } from '@/styles/global';
import colors from '@/styles/theme/colors';
import { Box, Stack, Typography } from '@mui/material';
import LoadingWrapper from '../LoadingWrapper';
import SimpleDataTable from '../SimpleDataTable';
import TableHeader from '../TableHeader';

const BillHistoryTab = () => {
  return (
    <LoadingWrapper
      height={'100%'}
      isLoading={false}
      cstyle={{
        overflow: 'auto',
        backgroundColor: false ? colors.white : 'transparent',
      }}
    >
      <RoundedCard radius="4px">
        <TableHeader title={'Billing history'} showSecondaryButton={false} />
        {true ? (
          <SimpleDataTable
            columns={HistoryBillingColumns}
            dataset={[
              {
                id: '1',
                date: 'July 11 2022',
                usage: '3 GB',
                total: '$ 100',
                pdf: 'www.google.com',
              },
            ]}
          />
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
