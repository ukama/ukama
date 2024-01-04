/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/styles/theme';
import EditableTextField from '@/ui/molecules/EditableTextField';
import {
  Box,
  Button,
  Switch,
  Grid,
  Divider,
  InputAdornment,
  Paper,
  TextField,
  Stack,
  Typography,
} from '@mui/material';
import { useState } from 'react';

interface INetworkSetting {
  name: string;
  handleSubmit: Function;
  handleDeleteNetwork: Function;
}

const NetworkSetting = ({
  name,
  handleSubmit,
  handleDeleteNetwork,
}: INetworkSetting) => {
  const [value, setValue] = useState('Democratic Republic of the Congo');
  const [isEditing, setIsEditing] = useState(false);
  const [textValue, setTextValue] = useState(name);

  const handleEditClick = () => {
    setIsEditing(true);
  };

  const handleSaveClick = () => {
    setIsEditing(false);
    // Perform save action or update state as needed
  };

  const handleChange = (event: any) => {
    setTextValue(event.target.value);
  };

  return (
    <Paper
      sx={{
        py: 3,
        px: 4,
        width: '100%',
        borderRadius: '4px',
        position: 'relative',
        height: 'calc(100vh - 200px)',
      }}
    >
      <Grid container spacing={2} pb={5}>
        <Grid item container spacing={2}>
          <Grid item xs={12} md={4}>
            <Typography variant="h6">Network details</Typography>
          </Grid>
          <Grid item xs={12} md={4} spacing={2}>
            <Stack direction={'column'} spacing={3}>
              <Typography
                variant="body2"
                sx={{
                  mb: '18px',
                  lineHeight: '19px',
                }}
              >
                You can edit this again at any point.
              </Typography>

              <TextField
                value={textValue}
                onChange={handleChange}
                disabled={!isEditing}
                required
                label="NETWORK NAME"
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="end">
                      {isEditing ? (
                        <Button
                          onClick={handleSaveClick}
                          variant="text"
                          sx={{ color: colors.primaryMain }}
                        >
                          Save
                        </Button>
                      ) : (
                        <Button
                          onClick={handleEditClick}
                          variant="text"
                          sx={{ color: colors.primaryMain }}
                        >
                          Edit
                        </Button>
                      )}
                    </InputAdornment>
                  ),
                }}
              />
              <EditableTextField
                type="text"
                value={value}
                isEditable={true}
                label="NETWORK COUNTRY"
                handleOnChange={(e: any) => setValue(e.target.value)}
              />
              <Button
                variant="contained"
                sx={{ width: '60%', bgcolor: colors.red }}
                onClick={() => handleDeleteNetwork()}
              >
                Delete Network
              </Button>
            </Stack>
          </Grid>
          <Grid item xs={12}>
            <Divider />
          </Grid>
          <Grid item xs={12} md={4}>
            <Typography variant="h6">Roaming options</Typography>
          </Grid>
          <Grid item xs={12} md={4} spacing={2}>
            <Stack direction={'column'} spacing={3}>
              <Typography
                variant="body2"
                sx={{
                  mb: '18px',
                  lineHeight: '19px',
                }}
              >
                Roaming is when you use cellular data outside your Ukama
                network, for [INSERT RATE HERE]. Roaming is by default an
                available option for all, and can be adjusted based on
                individual subscriber, unless you chose to disable roaming for
                all.
              </Typography>
              <Stack direction="row" spacing={1} alignItems={'center'}>
                <Switch disabled />
                <Typography variant="body1" color="initial">
                  Disable roaming for all
                </Typography>
              </Stack>
            </Stack>
          </Grid>
        </Grid>
        {/* <Divider sx={{ width: '100%' }} />
        <Grid item container spacing={2}>
          <Grid item xs={12} md={4}>
            <Typography variant="h6">Roaming Options</Typography>
          </Grid>
          <Grid item container xs={12} md={8}>
            <Grid item xs={12} sm={10} md={8}>
              <Typography
                variant="body2"
                sx={{
                  mb: '18px',
                  lineHeight: '19px',
                }}
              >
                Explanation of roaming & its rates. Your temporary eSIM has
                roaming enabled by default, and cannot be disabled.
              </Typography>
            </Grid>
            <Grid item xs={12} sm={10} md={8}>
              <FormControlLabel
                control={
                  <Switch
                    checked={networkSettings.roamingOption}
                    onChange={(e: any) => {
                      setNetworkSettings({
                        roamingOption: e.target.checked,
                      });
                      localStorage.setItem('roamingOption', e.target.checked);
                    }}
                  />
                }
                label="Enable roaming for all"
                sx={{ typography: 'body1' }}
              />
            </Grid>
            <Grid item xs={12} sm={10} md={8}>
              <TextField
                select
                id="eSims"
                InputProps={{
                  disabled: !networkSettings.roamingOption,
                  disableUnderline: true,
                }}
                value={esim}
                variant={'standard'}
                sx={{ mt: '18px' }}
                onChange={handleTimezoneChange}
              >
                {ROAMING_SELECT.map(({ value, text }: any) => (
                  <MenuItem key={value} value={value}>
                    <Typography variant="body2" sx={{ fontWeight: 500 }}>
                      {text}
                    </Typography>
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
          </Grid>
        </Grid> */}
      </Grid>
    </Paper>
  );
};
export default NetworkSetting;
