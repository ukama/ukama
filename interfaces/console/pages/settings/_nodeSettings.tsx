/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { LANGUAGE_OPTIONS, TimeZones } from '@/constants';
import { colors } from '@/styles/theme';
import { ExportOptionsType } from '@/types';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import SettingsArrowIcon from '@mui/icons-material/CallMade';
import { LinkStyle } from '@/styles/global';
import AutorenewIcon from '@mui/icons-material/Autorenew';
import FormControlCheckboxes from '@/ui/components/FormControlCheckboxes';

import {
  Button,
  FormControlLabel,
  Box,
  FormGroup,
  Divider,
  Grid,
  MenuItem,
  Checkbox,
  Paper,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { useRouter } from 'next/router';
import { useCallback, useState } from 'react';
const defaultTimeZone = 'Pacific Standard Time';
// localStorage['timeZone']? localStorage['timeZone']:'Pacific Standard Time'

export default function NodeSettings() {
  const router = useRouter();
  const [language, setLanguage] = useState('en');
  const [timezone, setTimezone] = useState(defaultTimeZone);
  const [alertList, setAlertList] = useState<Object>({});
  const [selectedDevices, setSelectedDevices] = useState<any>([]);

  const handleLanguageChange = (event: any) => {
    setLanguage(event.target.value);
    // localStorage.setItem('i18nextLng', event.target.value);
  };

  const handleTimezoneChange = (event: any) => {
    setTimezone(event.target.value);
    // localStorage.setItem('timeZone', event.target.value);
  };

  const handleAccountSettings = () => {
    router.push(
      `${process.env.NEXT_PUBLIC_AUTH_APP_URL}/auth/userAccountSettings`,
    );
  };
  const devices = [
    { id: 1, name: 'Device 1' },
    { id: 2, name: 'Device 2' },
    { id: 3, name: 'Device 3' },
    // Add more devices as needed
  ];

  const handleCheckboxChange = (deviceId: any) => {
    const isSelected = selectedDevices.includes(deviceId);
    setSelectedDevices((prevSelected: any[]) =>
      isSelected
        ? prevSelected.filter((id) => id !== deviceId)
        : [...prevSelected, deviceId]
    );
  };

  const handleSelectAll = () => {
    const allDeviceIds = devices.map((device) => device.id);
    setSelectedDevices(allDeviceIds);
  };

  const handleDeselectAll = () => {
    setSelectedDevices([]);
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
            <AutorenewIcon sx={{color:`${colors.black70}`}}/>
            <Typography variant="body1" sx={{color:`${colors.black70}`}} >Saving</Typography>

            </Stack>
            </Stack>
            </Grid>
            <Grid item xs={12} sm={7}>
             <Typography variant="body2" >
             Check the corresponding boxes for nodes you want to add to this network. To move nodes around in between networks and more, 
             <LinkStyle
              underline="hover"
              href={`node pool link`}
              sx={{
                typography: 'body1',
              }}
            >
              {`view node pool.`}
            </LinkStyle>
             If the node you purchased does not show up below, <LinkStyle
              underline="hover"
              href={`contact link`}
              sx={{
                typography: 'body1',
              }}
            >
              {`contact us.`}
            </LinkStyle> 
             </Typography>
             <Box>
      <Typography variant="h6">Devices to Update</Typography>
      <FormGroup>
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
        <Button variant="contained" onClick={handleSelectAll}>
          Select All
        </Button>
        <Button variant="contained" onClick={handleDeselectAll} sx={{ marginLeft: 2 }}>
          Deselect All
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
              <Typography variant="h6">Language & Region</Typography>
            </Grid>
            <Grid item md={8} xs={12} spacing={3} container>
              <Grid item xs={12} sm={12} md={8}>
                <TextField
                  select
                  id="language"
                  label="LANGUAGE"
                  value={language}
                  sx={{ width: '100%' }}
                  onChange={handleLanguageChange}
                >
                  {LANGUAGE_OPTIONS.map(
                    ({ value, label }: ExportOptionsType) => (
                      <MenuItem key={value} value={value}>
                        <Typography variant="body1">{label}</Typography>
                      </MenuItem>
                    ),
                  )}
                </TextField>
              </Grid>
              <Grid item xs={12} sm={12} md={8}>
                <TextField
                  select
                  id="timezone"
                  label="TIME ZONE"
                  value={timezone}
                  onChange={handleTimezoneChange}
                  sx={{ width: '100%' }}
                >
                  {TimeZones.map(({ value, text }: any) => (
                    <MenuItem key={value} value={value}>
                      <Typography variant="body1">{text}</Typography>
                    </MenuItem>
                  ))}
                </TextField>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </LoadingWrapper>
    </Paper>
  );
}
