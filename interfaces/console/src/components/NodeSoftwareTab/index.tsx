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
import UpdateIcon from '@mui/icons-material/Update';
import {
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
import { useEffect, useState } from 'react';
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
  nodeApps,
  handleUpdateAvailable,
}: INodeRadioTab) => {
  const [isUpdating, setIsUpdating] = useState<{
    [key: string]: boolean;
  }>({});

  useEffect(() => {
    nodeApps.forEach((app) => {
      if (app.status == SoftwareStatusEnum.UpdateInProgress) {
        setIsUpdating({
          ...isUpdating,
          [app.id]: true,
        });
      }
    });
  }, [nodeApps]);

  const statausByTitle = (status: string) => {
    switch (status) {
      case 'update_available':
        return 'Update Available';
      case 'update_in_progress':
        return 'Update In Progress';
      case 'update_failed':
        return 'Update Failed';
      case 'up_to_date':
        return 'Update Completed';
      default:
        return 'Unknown Status';
    }
  };
  return (
    <Paper
      sx={{
        p: 3,
        overflow: 'auto',
        height: { xs: 'calc(100vh - 480px)', md: 'calc(100vh - 328px)' },
      }}
    >
      <Typography variant="h6" sx={{ mb: 4 }}>
        Node Apps
      </Typography>
      <Grid2 container spacing={3}>
        {nodeApps?.map(
          ({
            id,
            name,
            currentVersion,
            desiredVersion,
            nodeId,
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
                        title={`${statausByTitle(status)}: ${
                          changeLog?.[changeLog?.length - 1] ?? ''
                        }`}
                      >
                        <IconButton
                          color={
                            status == SoftwareStatusEnum.UpdateAvailable
                              ? 'info'
                              : status == SoftwareStatusEnum.UpdateInProgress ||
                                  isUpdating[id]
                                ? 'warning'
                                : 'error'
                          }
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
                  <Stack direction="row" spacing={1 / 2} mt={1.5}>
                    <Typography variant="body2">CPU:</Typography>
                    <Typography variant="body2" sx={{ color: colors.darkBlue }}>
                      {12.2} %
                    </Typography>
                  </Stack>
                  <Stack direction="row" spacing={1 / 2}>
                    <Typography variant="body2">MEMORY:</Typography>
                    <Typography variant="body2" sx={{ color: colors.darkBlue }}>
                      {12.2} KB
                    </Typography>
                  </Stack>
                </CardContent>
                <CardActions sx={{ pb: 2, pt: 0, px: 2 }}>
                  <HorizontalContainerJustify>
                    <Button sx={{ p: 0 }}>View More</Button>
                    {(status == SoftwareStatusEnum.UpdateInProgress ||
                      isUpdating[id]) && (
                      <UpdateIcon
                        htmlColor={colors.yellow}
                        sx={{ width: '24px', height: '24px' }}
                      />
                    )}
                    {status == SoftwareStatusEnum.UpdateAvailable && (
                      <Button
                        sx={{ p: 0, color: colors.green }}
                        onClick={() => {
                          setIsUpdating({ ...isUpdating, [id]: true });
                          handleUpdateAvailable(name, desiredVersion, nodeId);
                          setIsUpdating({ ...isUpdating, [id]: false });
                        }}
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
    </Paper>
  );
};

export default NodeSoftwareTab;
