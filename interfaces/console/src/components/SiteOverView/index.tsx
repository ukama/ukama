/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Card, CardContent, Typography, Stack } from '@mui/material';

interface UptimeData {
  value: number;
  status: 'up' | 'down' | 'unknown';
}

interface SiteOverviewProps {
  uptimePercentage?: number;
  daysRange?: number;
  recentUptimeData?: UptimeData[];
  pastUptimeData?: UptimeData[];
  loading?: boolean;
}

const SiteOverview: React.FC<SiteOverviewProps> = ({
  uptimePercentage = 99,
  daysRange = 90,
  recentUptimeData,
  pastUptimeData,
  loading = false,
}) => {
  const generateMockData = (count: number): UptimeData[] => {
    return Array(count)
      .fill(0)
      .map(() => ({
        value: Math.random() * 100,
        status: Math.random() > 0.01 ? 'up' : 'down',
      }));
  };

  const recent = recentUptimeData || generateMockData(30);
  const past = pastUptimeData || generateMockData(30);

  const renderUptimeBar = (data: UptimeData, index: number) => (
    <Box
      key={index}
      sx={{
        height: 75,
        width: 8,
        backgroundColor:
          data.status === 'up'
            ? '#E5E5E5'
            : data.status === 'down'
              ? '#FF6B6B'
              : '#ADADAD',
        borderRadius: 1,
        mx: 0.25,
      }}
    />
  );

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
          {uptimePercentage}% uptime over {daysRange} days
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
            {recent.map(renderUptimeBar)}
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
            {past.map(renderUptimeBar)}
          </Box>

          <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
            <Typography variant="body2" color="text.secondary">
              90 days ago
            </Typography>
            <Typography variant="body2" color="text.secondary">
              60 days ago
            </Typography>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default SiteOverview;
