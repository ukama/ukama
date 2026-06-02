import { Skeleton, Stack, Typography } from '@mui/material';

const NetworkSkelton = () => {
  return (
    <Stack direction="column" height="100%" spacing={{ xs: 4, md: 6 }}>
      <Stack spacing={3}>
        <Typography variant="h4">
          <Skeleton width={200} />
        </Typography>
        <Stack direction="column" spacing={4}>
          <Typography variant="body1">
            <Skeleton width="100%" variant="text" />
            <Skeleton width="40%" variant="text" />
          </Typography>
          <Skeleton width="100%" height={42} variant="rectangular" />
        </Stack>
      </Stack>
      <Stack
        width="100%"
        direction="row"
        alignItems="center"
        height="fit-content"
        justifyContent="flex-end"
      >
        <Skeleton
          width={150}
          height={40}
          variant="rectangular"
          sx={{ borderRadius: '4px' }}
        />
      </Stack>
    </Stack>
  );
};

export default NetworkSkelton;
