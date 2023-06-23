import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import ConstructionSharpIcon from '@mui/icons-material/ConstructionSharp';
import { Paper, Stack, Typography } from '@mui/material';

export default function Alerts() {
  return (
    <LoadingWrapper
      radius="small"
      width={'100%'}
      isLoading={false}
      cstyle={{
        backgroundColor: false ? colors.white : 'transparent',
      }}
    >
      <Paper
        sx={{
          py: 3,
          px: 4,
          width: '100%',
          borderRadius: '5px',
          height: 'calc(100vh - 200px)',
        }}
      >
        <Stack
          height={'100%'}
          direction={'column'}
          alignItems={'center'}
          justifyContent={'center'}
        >
          <ConstructionSharpIcon />
          <Typography variant="subtitle2">Alerts</Typography>
        </Stack>
      </Paper>
    </LoadingWrapper>
  );
}
