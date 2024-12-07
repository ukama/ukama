import { Skeleton, Stack, Typography } from '@mui/material';

const LoadingSkeleton = () => (
  <Stack>
    <Typography variant="h4" fontWeight={500} mb={2}>
      <Skeleton width="60%" />
    </Typography>

    <Stack direction={'column'} spacing={4}>
      <Typography variant={'body1'} fontWeight={400}>
        <Skeleton width="100%" />
        <Skeleton width="100%" />
        <Skeleton width="70%" />
      </Typography>

      <Skeleton variant="rectangular" width="100%" height={122} />
    </Stack>

    <Stack
      direction={'row'}
      mt={{ xs: 4, md: 6 }}
      justifyContent={'space-between'}
    >
      <Skeleton variant="text" width="20%" />
      <Skeleton variant="rectangular" width="30%" height={36} />
    </Stack>
  </Stack>
);

export default LoadingSkeleton;
