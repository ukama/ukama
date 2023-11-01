import { colors } from '@/styles/theme';
import { Skeleton, Stack } from '@mui/material';

const GraphLoading = () => {
  return (
    <Stack
      spacing={1}
      direction={'row'}
      alignItems="flex-end"
      justifyContent="center"
    >
      <Skeleton
        variant="rectangular"
        width={10}
        height={24}
        sx={{ backgroundColor: colors.vulcan70 }}
      />
      <Skeleton
        variant="rectangular"
        width={10}
        height={34}
        sx={{ backgroundColor: colors.vulcan70 }}
      />
      <Skeleton
        variant="rectangular"
        width={10}
        height={44}
        sx={{ backgroundColor: colors.vulcan70 }}
      />
      <Skeleton
        variant="rectangular"
        width={10}
        height={40}
        sx={{ backgroundColor: colors.vulcan70 }}
      />
    </Stack>
  );
};

export default GraphLoading;
