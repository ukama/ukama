import colors from '@/styles/theme/colors';
import { Stack, Typography } from '@mui/material';

const OnBoarding = () => {
  return (
    <Stack direction={'column'} spacing={4}>
      <Typography variant="h4" color={colors.primaryMain}>
        <b>Welcome to Ukama</b>
      </Typography>
    </Stack>
  );
};

export default OnBoarding;
