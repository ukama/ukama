import { Button, Paper, Stack, Typography } from '@mui/material';
import { useRouter } from 'next/navigation';

const SiteSuccess = () => {
  const router = useRouter();
  return (
    <Paper elevation={0} sx={{ px: 4, py: 2 }}>
      <Stack direction={'column'} spacing={2}>
        <Typography variant="h6">Site successfully created</Typography>
        <Button
          variant="contained"
          onClick={() => router.push('/configure/success')}
        >
          Back to Home
        </Button>
      </Stack>
    </Paper>
  );
};

export default SiteSuccess;
