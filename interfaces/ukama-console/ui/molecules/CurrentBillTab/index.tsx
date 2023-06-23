import { CurrentBillColumns } from '@/constants/tableColumns';
import { NoBillYet } from '@/public/svg';
import { RoundedCard } from '@/styles/global';
import { Grid, Stack, Typography } from '@mui/material';
import CurrentBill from '../CurrentBill';
import PaymentCard from '../PaymentCard';
import SimpleDataTable from '../SimpleDataTable';
import TableHeader from '../TableHeader';

interface ICurrentBillTab {
  loading: boolean;
  data: any;
  totalAmount: string;
  currentBill: string;
  planName: string;
}

const CurrentBillTab = ({
  loading,
  data,
  totalAmount,
  currentBill,
  planName,
}: ICurrentBillTab) => {
  return (
    <Grid container item spacing={2}>
      <Grid xs={12} md={5} item>
        <CurrentBill loading={loading} amount={currentBill} plan={planName} />
      </Grid>
      <Grid xs={12} md={7} item>
        <PaymentCard
          selectedPM={''}
          onChangePM={() => {}}
          title={'Payment settings'}
          paymentMethodData={[]}
          onAddPaymentMethod={() => {}}
        />
      </Grid>
      <Grid xs={12} item>
        <RoundedCard>
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
