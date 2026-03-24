/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { colors } from '@/theme';
import { duration } from '@/utils';
import FiberManualRecordIcon from '@mui/icons-material/FiberManualRecord';
import { Box, Button, Paper, Stack, Typography } from '@mui/material';
import { useRouter } from 'next/navigation';
import React from 'react';

interface NodeStatusDisplayProps {
  nodeIds: string[];
  nodeUptimes: Record<string, number>;
}

const NodeStatusDisplay: React.FC<NodeStatusDisplayProps> = ({
  nodeIds,
  nodeUptimes,
}) => {
  const router = useRouter();

  return (
    <Paper
      elevation={0}
      sx={{
        p: 2,
        height: 'fit-content',
        borderRadius: 2,
        background: colors.gray,
      }}
    >
      <Stack spacing={4}>
        {nodeIds.map((id) => {
          // const isNodeDown = uptime <= 0;
          return (
            <Paper
              key={id}
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
                    // color: isNodeDown ? 'red' : colors.green,
                    color: colors.green,
                    mr: 2,
                    fontSize: 24,
                  }}
                />

                <Typography variant="h6" fontWeight={500}>
                  {`${id} is online and well`}
                </Typography>
              </Box>

              <Typography
                variant="body1"
                sx={{
                  ml: 4,
                  mb: 3,
                  color: colors.black70,
                }}
              >
                {`Node health has been up for ${duration(nodeUptimes[id] ?? 0)}`}
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
                  router.push(`/console/nodes/${id}`);
                }}
              >
                VIEW NODE
              </Button>
            </Paper>
          );
        })}
      </Stack>
    </Paper>
  );
};

export default NodeStatusDisplay;
