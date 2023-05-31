import { PageContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import { LoadingWrapper } from '@/ui/components';
import ConstructionSharpIcon from '@mui/icons-material/ConstructionSharp';
import { Stack, Typography } from '@mui/material';

export default function Page() {
  return (
    <LoadingWrapper
      radius="small"
      width={'100%'}
      isLoading={false}
      cstyle={{
        backgroundColor: false ? colors.white : 'transparent',
      }}
    >
      <PageContainer>
        <Stack
          height={'100%'}
          direction={'column'}
          alignItems={'center'}
          justifyContent={'center'}
        >
          <ConstructionSharpIcon />
          <Typography variant="subtitle2">Under developement.</Typography>
        </Stack>
      </PageContainer>
    </LoadingWrapper>
  );
}
