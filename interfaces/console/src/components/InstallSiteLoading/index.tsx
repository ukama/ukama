import colors from '@/theme/colors';
import GradientWrapper from '@/wrappers/gradiantWrapper';
import { LinearProgress, Paper, Stack, Typography } from '@mui/material';
import React from 'react';

interface IInstallSiteLoading {
  duration: number;
  onCompleted: () => void;
}

const InstallSiteLoading = ({ duration, onCompleted }: IInstallSiteLoading) => {
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
    <GradientWrapper>
      <Paper elevation={0} sx={{ px: 4, py: 2 }}>
        <Typography variant="h6" fontWeight={500}>
          Install site -{' '}
          <span style={{ color: colors.black70, fontWeight: 400 }}>
            <i>optional</i> (2/6)
          </span>
        </Typography>
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
      </Paper>
    </GradientWrapper>
  );
};

export default InstallSiteLoading;
