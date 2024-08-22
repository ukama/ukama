import colors from '@/theme/colors';
import { LinearProgress, Stack, Typography } from '@mui/material';
import React from 'react';

interface IInstallSiteLoading {
  duration: number;
  onCompleted: () => void;
  isShowStepTracker: boolean;
}

const InstallSiteLoading = ({
  duration,
  isShowStepTracker,
  onCompleted,
}: IInstallSiteLoading) => {
  const [progress, setProgress] = React.useState(0);
  const [remainingTime, setRemainingTime] = React.useState(duration);

  React.useEffect(() => {
    const startTime = Date.now();

    const timer = setInterval(() => {
      const elapsedTime = (Date.now() - startTime) / 1000;
      const newProgress = (elapsedTime / duration) * 100;
      setProgress(Math.min(newProgress, 100));

      const newRemainingTime = Math.max(duration - elapsedTime, 0);
      setRemainingTime(newRemainingTime);

      if (newProgress >= 100) {
        onCompleted();
        clearInterval(timer);
      }
    }, 1000);

    return () => {
      clearInterval(timer);
    };
  }, [duration]);

  return (
    <Stack direction="column">
      <Stack direction={'row'}>
        <Typography variant="h6"> {'Install site'}</Typography>
        <Typography
          variant="h6"
          fontWeight={400}
          sx={{
            color: colors.black70,
            display: isShowStepTracker ? 'flex' : 'none',
          }}
        >
          <i>&nbsp;- optional</i>&nbsp;(2/6)
        </Typography>
      </Stack>

      <Stack direction={'column'} mt={3} mb={3} spacing={1.5}>
        <Typography variant="body1" fontWeight={700}>
          Loading up your site...
        </Typography>
        <LinearProgress
          value={progress}
          variant="determinate"
          sx={{ height: '12px', borderRadius: '4px' }}
        />
        <Typography variant="body1">
          About {Math.ceil(remainingTime / 60)} minutes remaining
        </Typography>
      </Stack>
    </Stack>
  );
};

export default InstallSiteLoading;
