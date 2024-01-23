import { Grid, Typography, Stack } from '@mui/material';
import colors from '@/styles/theme/colors';

interface SubscriberOverviewProps {
  revenuePerYear: number;
}
const SubscriberOverView: React.FC<SubscriberOverviewProps> = ({
  revenuePerYear,
}) => {
  return (
    <Grid container spacing={0}>
      <Stack direction="row" spacing={1} alignItems={'center'}>
        <Typography variant="h6">Yearly overview </Typography>
        <Typography variant="body1" sx={{ color: colors.black38 }}>
          (${revenuePerYear})
        </Typography>
      </Stack>
    </Grid>
  );
};
export default SubscriberOverView;
