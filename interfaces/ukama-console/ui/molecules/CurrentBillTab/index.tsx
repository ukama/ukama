import { CurrentBillColumns } from '@/constants/tableColumns';
import { RoundedCard } from '@/styles/global';
import { Grid, Stack, Typography } from '@mui/material';
import CurrentBill from '../CurrentBill';
import PaymentCard from '../PaymentCard';
import SimpleDataTable from '../SimpleDataTable';
import TableHeader from '../TableHeader';

export default function CurrentBillTab() {
  return (
    <Grid container item spacing={2}>
      <Grid xs={12} md={5} item>
        <CurrentBill
          amount={`$ 0.00`}
          billMonth={''}
          dueDate={''}
          loading={false}
        />
      </Grid>
      <Grid xs={12} md={7} item>
        <RoundedCard>
          <PaymentCard
            selectedPM={''}
            onChangePM={() => {}}
            title={'Payment settings'}
            paymentMethodData={[]}
            onAddPaymentMethod={() => {}}
          />
        </RoundedCard>
      </Grid>
      <Grid xs={12} item>
        <RoundedCard>
          <TableHeader
            title={'Billing breakdown'}
            showSecondaryButton={false}
          />
          {true ? (
            <SimpleDataTable
              columns={CurrentBillColumns}
              dataset={[]}
              // totalAmount={0}
            />
          ) : (
            <Stack
              direction="column"
              spacing={2}
              justifyItems={'center'}
              alignItems={'center'}
            >
              {/* <NoBillYet /> */}
              <Typography variant="body1">Nothing in your bill yet!</Typography>
            </Stack>
          )}
        </RoundedCard>
      </Grid>
    </Grid>
  );
}
