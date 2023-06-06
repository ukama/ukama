import { colors } from '@/styles/theme';
import { ArrowBack } from '@mui/icons-material';
import { IconButton, Stack, Typography } from '@mui/material';
import { useRouter } from 'next/router';

interface IBackButton {
  title: string;
}

const BackButton = ({ title }: IBackButton) => {
  const router = useRouter();
  return (
    <>
      <Stack
        direction={'row'}
        alignItems={'center'}
        spacing={1.5}
        sx={{
          ':hover': {
            p: {
              color: colors.primaryMain,
              cursor: 'pointer',
            },
            '.MuiButtonBase-root': {
              color: colors.primaryMain,
            },
          },
        }}
        onClick={() => router.back()}
      >
        <IconButton size="small" sx={{ p: 0 }}>
          <ArrowBack />
        </IconButton>
        <Typography variant={'body2'} fontWeight={500}>
          {title}
        </Typography>
      </Stack>
    </>
  );
};
export default BackButton;
