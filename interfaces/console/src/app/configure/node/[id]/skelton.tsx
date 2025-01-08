import { Skeleton, Stack, Typography } from '@mui/material';

const LoadingSkelton = () => {
  return (
    <Stack spacing={2}>
      <Stack direction={'row'}>
        <Typography variant="h4" fontWeight={500}>
          <Skeleton width="60%" />
        </Typography>
      </Stack>

      <Stack mt={3} mb={3} direction={'column'} spacing={3}>
        <Typography variant={'body1'} fontWeight={400}>
          <Skeleton width="100%" />
          <Skeleton width="100%" />
          <Skeleton width="80%" />
        </Typography>

        <Skeleton variant="rounded" width={'100%'} height={128} />

        <Skeleton width="40%" height={32} />
        <Skeleton width="60%" height={32} />
      </Stack>

      <Stack
        direction={'row'}
        pt={{ xs: 4, md: 6 }}
        justifyContent={'space-between'}
      >
        <Skeleton variant="text" width="20%" />
        <Skeleton variant="rectangular" width="30%" height={36} />
      </Stack>
    </Stack>
  );
};

export default LoadingSkelton;
