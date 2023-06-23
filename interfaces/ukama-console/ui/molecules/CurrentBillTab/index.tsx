import { CurrentBillColumns } from '@/constants/tableColumns';
import { NoBillYet } from '@/public/svg';
import { RoundedCard } from '@/styles/global';
import { Grid, Stack, Typography } from '@mui/material';
import CurrentBill from '../CurrentBill';
import PaymentCard from '../PaymentCard';
import SimpleDataTable from '../SimpleDataTable';
import TableHeader from '../TableHeader';

const CurrentBillTab = () => {
  return (
    <Grid container item spacing={2}>
      <Grid xs={12} md={5} item>
        <CurrentBill
          loading={false}
          amount={`$ 20.00`}
          plan={'Console cloud plan - [community/empowerment]'}
        />
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
          {true ? (
            <Stack alignItems={'flex-end'}>
              <SimpleDataTable
                columns={CurrentBillColumns}
                dataset={[
                  {
                    id: '1',
                    name: 'Tryphena Nelson',
                    rate: '3 GB',
                    subtotal: '200',
                  },
                ]}
              />
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
                  $20.00
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
