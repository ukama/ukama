/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Card, CardContent, Typography, Tooltip } from '@mui/material';
import colors from '@/theme/colors';

interface NodeUptime {
  id: string;
  uptimeSeconds: number;
}

interface SiteOverviewProps {
  siteUptimeSeconds?: number;
  nodeUptimes?: NodeUptime[];
  daysRange?: number;
  loading?: boolean;
}

const SiteOverview: React.FC<SiteOverviewProps> = ({
  siteUptimeSeconds = 0,
  nodeUptimes = [],
  daysRange = 90,
  loading = false,
}) => {
  const totalPossibleSeconds = daysRange * 24 * 60 * 60;

  const siteUptimePercentage = (siteUptimeSeconds / totalPossibleSeconds) * 100;

  let averageNodeUptimePercentage = 0;
  if (nodeUptimes.length > 0) {
    const totalNodeUptimeSeconds = nodeUptimes.reduce(
      (sum, node) => sum + node.uptimeSeconds,
      0,
    );

    const averageNodeUptimeSeconds =
      totalNodeUptimeSeconds / nodeUptimes.length;

    averageNodeUptimePercentage =
      (averageNodeUptimeSeconds / totalPossibleSeconds) * 100;
  }

  let overallUptimePercentage = siteUptimePercentage;
  if (nodeUptimes.length > 0) {
    overallUptimePercentage =
      (siteUptimePercentage + averageNodeUptimePercentage) / 2;
  }

  overallUptimePercentage = Math.min(Math.max(0, overallUptimePercentage), 100);

  const nodePercentages = nodeUptimes.map((node) => ({
    id: node.id,
    percentage: (node.uptimeSeconds / totalPossibleSeconds) * 100,
  }));

  const recentPeriodBars = Array(30).fill(overallUptimePercentage);
  const pastPeriodBars = Array(30).fill(overallUptimePercentage);

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
            bgcolor: colors.lightGreen,
            borderRadius: 1,
          }}
        />
      </Box>
    );
  };

  const createTooltipText = () => {
    let text = `Site Uptime: ${siteUptimePercentage.toFixed(1)}%`;

    if (nodeUptimes.length > 0) {
      text += `\nAverage Node Uptime: ${averageNodeUptimePercentage.toFixed(1)}%`;
      text += '\n';

      nodePercentages.forEach((node) => {
        text += `\nNode ${node.id}: ${node.percentage.toFixed(1)}%`;
      });
    }

    return text;
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

        <Tooltip
          title={createTooltipText()}
          placement="top"
          componentsProps={{
            tooltip: {
              sx: {
                whiteSpace: 'pre-line',
                maxWidth: 'none',
              },
            },
          }}
        >
          <Typography
            variant="body2"
            sx={{
              mt: 2,
              mb: 4,
            }}
          >
            {overallUptimePercentage.toFixed(1)}% uptime over {daysRange} days
          </Typography>
        </Tooltip>

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
