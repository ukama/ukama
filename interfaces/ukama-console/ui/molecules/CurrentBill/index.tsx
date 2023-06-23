import { RoundedCard } from '@/styles/global';
import colors from '@/styles/theme/colors';
import PictureAsPdfIcon from '@mui/icons-material/PictureAsPdf';
import { Divider, Stack, Typography } from '@mui/material';
import { format } from 'date-fns';
import LoadingWrapper from '../LoadingWrapper';
type CurrentBillProps = {
  amount: string;
  billMonth: string;
  dueDate: string;
  loading: boolean;
};

const CurrentBill = ({
  amount,
  billMonth,
  dueDate,
  loading,
}: CurrentBillProps) => {
  return (
    <RoundedCard radius={'4px'}>
      <LoadingWrapper height={200} isLoading={loading}>
        <Stack direction="row" spacing={1} alignItems="center">
          <Typography variant="h6">
            {`${format(new Date(), 'MMMM')} bill`}
          </Typography>
          <PictureAsPdfIcon sx={{ color: colors.primaryMain }} />
        </Stack>

        <Typography variant="body2">{`${dueDate} - ${billMonth}`}</Typography>

        <Divider />
        <Typography variant="h3" sx={{ m: '18px 0px' }}>
          {amount}
        </Typography>
      </LoadingWrapper>
    </RoundedCard>
  );
};
export default CurrentBill;
