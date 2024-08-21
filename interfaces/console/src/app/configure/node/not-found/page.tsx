'use client';
import { Button, Paper, Stack, Typography } from '@mui/material';
import { useRouter } from 'next/navigation';

const NodeNotFoundPage = () => {
  const router = useRouter();
  return (
    <Paper elevation={0} sx={{ px: 4, py: 2 }}>
      <Stack direction={'column'} spacing={2}>
        <Typography variant="h6">No new node found!</Typography>
        <Typography variant="body1">
          Please check that your node is On. If it's On you'll get notification
          when it get online and ready to configure.
        </Typography>
        <Button
          variant="contained"
          sx={{ width: 'fit-content', alignSelf: 'flex-end' }}
          onClick={() => router.push('/')}
        >
          Back to Home
        </Button>
      </Stack>
    </Paper>
  );
};

export default NodeNotFoundPage;
