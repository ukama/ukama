import { RoundedCard } from '@/styles/global';
import { Divider, Stack, Typography } from '@mui/material';
import LoadingWrapper from '../LoadingWrapper';
type CurrentBillProps = {
  amount: string;
  loading: boolean;
};

const CurrentBill = ({ amount, loading }: CurrentBillProps) => {
  return (
    <LoadingWrapper height={164} isLoading={loading}>
      <RoundedCard radius={'4px'}>
        <Stack direction="column" spacing={1.4} alignItems="flex-start">
          <Stack direction="row" width={'100%'} justifyContent="space-between">
            <Typography variant="body2" fontWeight={600}>
              Billing Month
            </Typography>
            <Typography variant="caption">06/14/2022 - 07/14/2022</Typography>
          </Stack>

          <Typography variant="caption">
            Detailed bill breakdown available below.
          </Typography>
          <Divider sx={{ width: '100%' }} />
          <Typography variant="h4" sx={{ m: '18px 0px' }}>
            {amount}
          </Typography>
        </Stack>
      </RoundedCard>
    </LoadingWrapper>
  );
};
export default CurrentBill;
