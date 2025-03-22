import React from 'react';
import { Box, Card, CardContent, Typography } from '@mui/material';
import colors from '@/theme/colors';

interface SiteOverviewProps {
  uptimeSeconds?: number;
  daysRange?: number;
  loading?: boolean;
}

const SiteOverview: React.FC<SiteOverviewProps> = ({
  uptimeSeconds = 0,
  daysRange = 90,
  loading = false,
}) => {
  const calculateUptimePercentage = (
    uptimeSeconds: number,
    days: number,
  ): number => {
    const totalSeconds = days * 24 * 60 * 60;
    const percentage = (uptimeSeconds / totalSeconds) * 100;
    return Math.min(Math.max(0, percentage), 100);
  };

  const actualUptimePercentage = calculateUptimePercentage(
    uptimeSeconds,
    daysRange,
  );

  const recentPeriodBars = Array(30).fill(actualUptimePercentage);
  const pastPeriodBars = Array(30).fill(actualUptimePercentage);

  const renderBar = (value: number, index: number) => {
    const heightPercentage = value;
    const barHeight = (heightPercentage / 100) * 75;

    return (
      <Box
        key={index}
        sx={{
          height: 75,
          width: 8,
          mx: 0.25,
          position: 'relative',
          borderRadius: 1,
          bgcolor: colors.gray,
        }}
      >
        <Box
          sx={{
            position: 'absolute',
            bottom: 0,
            width: '100%',
            height: `${barHeight}px`,
            bgcolor: colors.green,
            borderRadius: 1,
          }}
        />
      </Box>
    );
  };

  if (loading) {
    return (
      <Card
        sx={{
          borderRadius: 2,
          boxShadow: '0px 2px 6px rgba(0, 0, 0, 0.05)',
          height: '100%',
        }}
      ></Card>
    );
  }

  return (
    <Card
      sx={{
        borderRadius: 2,
        boxShadow: '0px 2px 6px rgba(0, 0, 0, 0.05)',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <CardContent sx={{ padding: 4, flexGrow: 1 }}>
        <Typography variant="h6" sx={{ mb: 3 }}>
          Site overview
        </Typography>

        <Typography
          variant="body2"
          sx={{
            mt: 2,
            mb: 4,
          }}
        >
          {actualUptimePercentage.toFixed(0)}% uptime over {daysRange} days
        </Typography>

        <Box sx={{ position: 'relative', mb: 3 }}>
          <Box
            sx={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'flex-end',
              height: 75,
              mb: 1,
            }}
          >
            {recentPeriodBars.map(renderBar)}
          </Box>

          <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
            <Typography variant="body2" color="text.secondary">
              30 days ago
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Today
            </Typography>
          </Box>
        </Box>

        <Box sx={{ position: 'relative', mt: 4 }}>
          <Box
            sx={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'flex-end',
              height: 75,
              mb: 1,
            }}
          >
            {pastPeriodBars.map(renderBar)}
          </Box>

          <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
            <Typography variant="body2" color="text.secondary">
              {daysRange} days ago
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {Math.floor(daysRange / 3)} days ago
            </Typography>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default SiteOverview;
