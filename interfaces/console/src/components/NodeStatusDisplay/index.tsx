/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Paper, Stack, Typography, Button } from '@mui/material';
import FiberManualRecordIcon from '@mui/icons-material/FiberManualRecord';
import { colors } from '@/theme';
import { duration } from '@/utils';
import { useRouter } from 'next/navigation';

interface NodeStatusDisplayProps {
  nodeIds: string[];
  nodeUpTime?: number;
}

const NodeStatusDisplay: React.FC<NodeStatusDisplayProps> = ({
  nodeIds = [],
  nodeUpTime = 0,
}) => {
  const router = useRouter();

  return (
    <Paper
      sx={{
        p: 4,
        borderRadius: 2,
        height: {
          xs: 'calc(100vh - 480px)',
          md: 'calc(100vh - 328px)',
        },
        overflow: 'auto',
        background: colors.gray,
      }}
    >
      <Stack spacing={4}>
        {nodeIds.map((nodeId, index) => (
          <Paper
            key={nodeId || index}
            sx={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'flex-start',
              mb: 2,
              p: 2,
            }}
          >
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              <FiberManualRecordIcon
                sx={{
                  color: colors.green,
                  mr: 2,
                  fontSize: 24,
                }}
              />
              <Typography variant="h6" fontWeight="500">
                {nodeId} is online and well
              </Typography>
            </Box>

            <Typography variant="body1" sx={{ ml: 4, mb: 3 }}>
              Node health has been up for {duration(nodeUpTime)}
            </Typography>

            <Button
              variant="text"
              sx={{
                ml: 4,
                fontWeight: 'bold',
                fontSize: '14px',
                color: colors.primaryMain,
              }}
              onClick={() => {
                router.push(`/console/nodes/${nodeId}`);
              }}
            >
              VIEW NODE
            </Button>
          </Paper>
        ))}
      </Stack>
    </Paper>
  );
};

export default NodeStatusDisplay;
