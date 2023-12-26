/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import { LinkStyle } from '@/styles/global';
import AutorenewIcon from '@mui/icons-material/Autorenew';
import SimpleDataTable from '@/ui/molecules/SimpleDataTable';
import { NODE_SETTINGS_TABLE_COLUMN } from '@/constants';

import {
  Button,
  FormControlLabel,
  Box,
  FormGroup,
  Typography,
  Divider,
  Grid,
  Checkbox,
  Paper,
  Stack,
} from '@mui/material';
import { useRouter } from 'next/router';
import { useState } from 'react';
interface Device {
  id: number;
  name: string;
}

export default function NodeSettings() {
  const [selectedDevices, setSelectedDevices] = useState<number[]>([]);

  const handleCheckboxChange = (deviceId: number) => {
    const isSelected = selectedDevices.includes(deviceId);
    setSelectedDevices((prevSelected) =>
      isSelected
        ? prevSelected.filter((id) => id !== deviceId)
        : [...prevSelected, deviceId],
    );
  };

  const devices: Device[] = [
    { id: 1, name: 'Device 1' },
    { id: 2, name: 'Device 2' },
    { id: 3, name: 'Device 3' },
  ];

  const handleAssignToNetwork = () => {
    // Implement your logic to assign selected nodes to the network
    console.log('Assigned to network:', selectedDevices);
  };

  return (
    <Paper
      sx={{
        py: 3,
        px: 4,
        width: '100%',
        borderRadius: '5px',
        height: 'calc(100vh - 200px)',
      }}
    >
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={false}
        cstyle={{
          overflow: 'auto',
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <Grid container spacing={2}>
          <Grid item container xs={12} spacing={2}>
            <Grid item xs={12} sm={4}>
              <Stack direction="column" spacing={2}>
                <Typography variant="h6">Assign network nodes</Typography>
                <Stack direction={'row'} spacing={1} alignItems={'center'}>
                  <AutorenewIcon sx={{ color: `${colors.black70}` }} />
                  <Typography
                    variant="body1"
                    sx={{ color: `${colors.black70}` }}
                  >
                    Saving
                  </Typography>
                </Stack>
              </Stack>
            </Grid>
            <Grid item xs={12} sm={7}>
              <Typography variant="body2">
                Check the corresponding boxes for nodes you want to add to this
                network. To move nodes around in between networks and more,
                <LinkStyle
                  underline="hover"
                  href={`node pool link`}
                  sx={{
                    typography: 'body1',
                  }}
                >
                  {`view node pool.`}
                </LinkStyle>
                If the node you purchased does not show up below,{' '}
                <LinkStyle
                  underline="hover"
                  href={`contact link`}
                  sx={{
                    typography: 'body1',
                  }}
                >
                  {`contact us.`}
                </LinkStyle>
              </Typography>
              <Box mt={2}>
                <FormGroup>
                  <FormControlLabel
                    control={
                      <Checkbox
                        checked={selectedDevices.length === devices.length}
                        indeterminate={
                          selectedDevices.length > 0 &&
                          selectedDevices.length < devices.length
                        }
                        onChange={() =>
                          selectedDevices.length === devices.length
                            ? setSelectedDevices([])
                            : setSelectedDevices(
                                devices.map((device) => device.id),
                              )
                        }
                      />
                    }
                    label={`Available nodes (${devices.length})`}
                  />
                  {devices.map((device) => (
                    <FormControlLabel
                      key={device.id}
                      control={
                        <Checkbox
                          checked={selectedDevices.includes(device.id)}
                          onChange={() => handleCheckboxChange(device.id)}
                        />
                      }
                      label={device.name}
                    />
                  ))}
                </FormGroup>
                <Box mt={2}>
                  <Button variant="contained" onClick={handleAssignToNetwork}>
                    Assign to Network
                  </Button>
                </Box>
              </Box>
            </Grid>
            <Grid item xs={12}>
              <Divider />
            </Grid>
          </Grid>

          <Grid item container xs={12} spacing={2}>
            <Grid item xs={12} md={4}>
              <Typography variant="h6">Current network nodes</Typography>
            </Grid>
            <Grid item md={8} xs={12} spacing={3} container>
              <SimpleDataTable
                dataKey="uuid"
                dataset={[]}
                columns={NODE_SETTINGS_TABLE_COLUMN}
              />
            </Grid>
          </Grid>
        </Grid>
      </LoadingWrapper>
    </Paper>
  );
}
