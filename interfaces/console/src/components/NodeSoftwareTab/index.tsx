/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Software, SoftwareStatusEnum } from '@/client/graphql/generated';
import { HorizontalContainerJustify } from '@/styles/global';
import colors from '@/theme/colors';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  Grid2,
  IconButton,
  Paper,
  Stack,
  Tooltip,
  Typography,
} from '@mui/material';
interface INodeRadioTab {
  loading: boolean;
  nodeApps: Software[];
  handleUpdateAvailable: (
    name: string,
    desiredVersion: string,
    nodeId: string,
  ) => void;
}

const NodeSoftwareTab = ({
  loading,
  nodeApps,
  handleUpdateAvailable,
}: INodeRadioTab) => {
  return (
    <Grid2 container spacing={2} sx={{ overflowY: 'scroll' }}>
      {/* <Grid2 size={12} sx={{ gridRowStart: 1, gridRowEnd: 6 }}>
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
      </Grid2> */}

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
              {nodeApps?.map(
                ({
                  id,
                  name,
                  currentVersion,
                  desiredVersion,
                  nodeId,
                  updatedAt,
                  status,
                  changeLog,
                }: Software) => (
                  <Grid2 size={3} key={id}>
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
                          <Typography
                            variant="h5"
                            fontWeight={400}
                            textTransform={'capitalize'}
                          >
                            {name}
                          </Typography>
                          {status == SoftwareStatusEnum.UpdateAvailable && (
                            <Tooltip
                              arrow
                              placement="right"
                              title={`Update Available: ${
                                changeLog?.[changeLog?.length - 1] ?? ''
                              }`}
                            >
                              <IconButton
                                color="info"
                                sx={{
                                  '&:hover svg path': {
                                    fill: 'inherit',
                                  },
                                }}
                              >
                                <InfoOutlinedIcon
                                  sx={{
                                    width: '16px',
                                    height: '16px',
                                  }}
                                />
                              </IconButton>
                            </Tooltip>
                          )}
                        </Stack>
                        <Typography
                          variant="body2"
                          color="text.secondary"
                          gutterBottom
                        >
                          Version: {currentVersion}
                        </Typography>
                        <Stack direction="row" spacing={1 / 2} mt={'12px'}>
                          <Typography variant="body2">CPU:</Typography>
                          <Typography
                            variant="body2"
                            sx={{ color: colors.darkBlue }}
                          >
                            {12.2} %
                          </Typography>
                        </Stack>
                        <Stack direction="row" spacing={1 / 2}>
                          <Typography variant="body2">MEMORY:</Typography>
                          <Typography
                            variant="body2"
                            sx={{ color: colors.darkBlue }}
                          >
                            {12.2} KB
                          </Typography>
                        </Stack>
                      </CardContent>
                      <CardActions sx={{ p: 2 }}>
                        <HorizontalContainerJustify>
                          <Button sx={{ p: 0 }}>View More</Button>
                          {status === SoftwareStatusEnum.UpdateAvailable && (
                            <Button
                              sx={{ p: 0, color: colors.green }}
                              onClick={() =>
                                handleUpdateAvailable(
                                  name,
                                  desiredVersion,
                                  nodeId,
                                )
                              }
                            >
                              Update Available
                            </Button>
                          )}
                        </HorizontalContainerJustify>
                      </CardActions>
                    </Card>
                  </Grid2>
                ),
              )}
            </Grid2>
          </Box>
        </Paper>
      </Grid2>
    </Grid2>
  );
};

export default NodeSoftwareTab;
