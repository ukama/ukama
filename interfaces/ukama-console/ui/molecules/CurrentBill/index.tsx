import { RoundedCard } from '@/styles/global';
import { Divider, Stack, Typography } from '@mui/material';
import { format } from 'date-fns';
import LoadingWrapper from '../LoadingWrapper';
type CurrentBillProps = {
  amount: string;
  plan: string;
  loading: boolean;
};

const CurrentBill = ({ amount, plan, loading }: CurrentBillProps) => {
  return (
    <LoadingWrapper height={194} isLoading={loading}>
      <RoundedCard radius={'4px'}>
        <Stack direction="column" spacing={1} alignItems="flex-start">
          <Typography variant="h6">
            {`${format(new Date(), 'MMMM')} bill`}
          </Typography>

          <Typography variant="caption">{plan}</Typography>
          <Divider sx={{ width: '100%' }} />
          <Typography variant="h3" sx={{ m: '18px 0px' }}>
            {amount}
          </Typography>
        </Stack>
      </RoundedCard>
    </LoadingWrapper>
  );
};
export default CurrentBill;
