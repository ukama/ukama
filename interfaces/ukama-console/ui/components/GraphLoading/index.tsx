import { Skeleton, Stack } from '@mui/material';

const GraphLoading = () => {
  return (
    <Stack
      spacing={1}
      direction={'row'}
      alignItems="flex-end"
      justifyContent="center"
    >
      <Skeleton variant="rectangular" width={10} height={24} />
      <Skeleton variant="rectangular" width={10} height={34} />
      <Skeleton variant="rectangular" width={10} height={44} />
      <Skeleton variant="rectangular" width={10} height={40} />
    </Stack>
  );
};

export default GraphLoading;
