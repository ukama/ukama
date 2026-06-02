/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { LinearProgress, Stack, Typography } from '@mui/material';
import React, { useEffect } from 'react';

interface IInstallSiteLoading {
  title: string;
  subtitle: string;
  duration: number;
  description: string;
  handleBack: () => void;
  onCompleted: () => void;
}

const InstallSiteLoading = ({
  title,
  subtitle,
  duration,
  description,
  onCompleted,
}: IInstallSiteLoading) => {
  const [progress, setProgress] = React.useState(0);
  const [remainingTime, setRemainingTime] = React.useState(duration);

  useEffect(() => {
    if (duration !== remainingTime) {
      setProgress(0);
      setRemainingTime(duration);
    }
  }, [duration]);

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
    <Stack width={'100%'} direction="column" spacing={2}>
      <Typography variant="h4" fontWeight={500}>
        {title}
      </Typography>
      <Stack direction={'column'} spacing={1.5}>
        {subtitle && (
          <Typography variant="h6" sx={{ fontSize: '16px !important' }}>
            {subtitle}
          </Typography>
        )}
        <LinearProgress
          value={progress}
          variant="determinate"
          sx={{ height: '12px', borderRadius: '4px' }}
        />
        <Typography variant="body2">
          {description
            ? description
            : `About ${Math.ceil(remainingTime)} seconds remaining`}
        </Typography>
      </Stack>
    </Stack>
  );
};

export default InstallSiteLoading;
