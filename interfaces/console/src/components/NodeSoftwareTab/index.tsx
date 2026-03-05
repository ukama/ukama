/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NodeApp } from '@/client/graphql/generated';
import { NodeAppsColumns } from '@/constants/tableColumns';
import colors from '@/theme/colors';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  Grid2,
  Paper,
  Stack,
  Typography,
} from '@mui/material';
import SimpleDataTable from '../SimpleDataTable';
interface INodeRadioTab {
  loading: boolean;
  nodeApps: NodeApp[];
}

const NodeSoftwareTab = ({ loading, nodeApps }: INodeRadioTab) => {
  return (
    <Grid2 container spacing={2} sx={{ overflowY: 'scroll' }}>
      <Grid2 size={12} sx={{ gridRowStart: 1, gridRowEnd: 6 }}>
        <Paper
          sx={{
            p: 3,
          }}
        >
          <Typography variant="h6" sx={{ marginBottom: 3 }}>
            Change Logs
          </Typography>
          <Box sx={{ overflow: 'scroll', height: '100%', pb: 4, pr: 3 }}>
            <SimpleDataTable
              height={'400'}
              dataset={nodeApps}
              columns={NodeAppsColumns}
            />
          </Box>
        </Paper>
      </Grid2>

      <Grid2 size={12}>
        <Paper
          sx={{
            p: 3,
          }}
        >
          <Typography variant="h6" sx={{ mb: 4 }}>
            Node Apps
          </Typography>
          <Box sx={{ overflow: 'scroll', height: '100%', pb: 8, pr: 3 }}>
            <Grid2 container spacing={3}>
              {nodeApps?.map(({ id, title, cpu, memory, version }: any) => (
                <Grid2 size={4} key={id}>
                  <Card variant="outlined">
                    <CardContent>
                      <Stack
                        spacing={1}
                        direction="row"
                        sx={{ alignItems: 'center' }}
                      >
                        <CheckCircleIcon
                          htmlColor={colors.green}
                          fontSize="medium"
                        />
                        <Typography variant="h5" textTransform={'capitalize'}>
                          {title}
                        </Typography>
                      </Stack>
                      <Typography
                        variant="body2"
                        color="text.secondary"
                        gutterBottom
                      >
                        Version: {version}
                      </Typography>
                      <Stack direction="row" spacing={1 / 2} mt={'12px'}>
                        <Typography variant="body2">CPU:</Typography>
                        <Typography
                          variant="body2"
                          sx={{ color: colors.darkBlue }}
                        >
                          {parseFloat(cpu).toFixed(2)} %
                        </Typography>
                      </Stack>
                      <Stack direction="row" spacing={1 / 2}>
                        <Typography variant="body2">MEMORY:</Typography>
                        <Typography
                          variant="body2"
                          sx={{ color: colors.darkBlue }}
                        >
                          {parseFloat(memory).toFixed(2)} KB
                        </Typography>
                      </Stack>
                    </CardContent>
                    <CardActions sx={{ ml: 1 }}>
                      <Button>VIEW MORE</Button>
                    </CardActions>
                  </Card>
                </Grid2>
              ))}
            </Grid2>
          </Box>
        </Paper>
      </Grid2>
    </Grid2>
  );
};

export default NodeSoftwareTab;
