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

const CurrentBillTab = ({
  data,
  loading,
  planName,
  totalAmount,
  currentBill,
}: ICurrentBillTab) => {
  return (
    <Grid container item spacing={2}>
      <Grid xs={12} md={6} item>
        <CurrentBill loading={loading} amount={currentBill} />
      </Grid>
      <Grid xs={12} md={6} item>
        {/* Will be part of phase 2
         <PaymentCard
          selectedPM={''}
          onChangePM={() => {}}
          title={'Payment settings'}
          paymentMethodData={[]}
          onAddPaymentMethod={() => {}}
        /> */}
        <NotificationContainer />
      </Grid>

      <Grid xs={12} item>
        <RoundedCard radius="4px">
          <TableHeader
            title={'Billing breakdown'}
            showSecondaryButton={false}
          />
          {data.length > 0 ? (
            <Stack alignItems={'flex-end'}>
              <SimpleDataTable columns={CurrentBillColumns} dataset={data} />
              <Stack
                mt={2}
                width={'270px'}
                direction={'row'}
                justifyContent={'space-between'}
              >
                <Typography variant="h6" textAlign={'end'}>
                  Total
                </Typography>
                <Typography variant="h6" textAlign={'end'}>
                  {totalAmount}
                </Typography>
              </Stack>
            </Stack>
          ) : (
            <Stack
              direction="column"
              spacing={2}
              justifyItems={'center'}
              alignItems={'center'}
            >
              <NoBillYet />
              <Typography variant="body1">Nothing in your bill yet!</Typography>
            </Stack>
          )}
        </RoundedCard>
      </Grid>
    </Grid>
  );
};

export default CurrentBillTab;
