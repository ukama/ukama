import { Skeleton, Stack, Typography } from '@mui/material';

const LoadingSkeleton = () => (
  <Stack>
    <Typography variant="h4" fontWeight={500} mb={2}>
      <Skeleton width="60%" />
    </Typography>

    <Stack direction="column" spacing={2.5}>
      <Typography variant="body1" color="textSecondary">
        <Skeleton width="100%" />
        <br />
        <Skeleton width="100%" />
        <Skeleton width="80%" />
        <br />
      </Typography>

      <Skeleton variant="rectangular" width="100%" height={56} />
      <Skeleton variant="rectangular" width="100%" height={56} />
      <Skeleton variant="rectangular" width="100%" height={56} />
    </Stack>

    <Stack
      mt={{ xs: 4, md: 6 }}
      direction={'row'}
      justifyContent={'space-between'}
    >
      <Skeleton variant="text" width="20%" />
      <Skeleton variant="rectangular" width="30%" height={36} />
    </Stack>
  </Stack>
);

export default LoadingSkeleton;
