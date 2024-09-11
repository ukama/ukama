/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { PackageDto, SimDto } from '@/client/graphql/generated';
import { TAddSubscriberData } from '@/types';
import {
  Button,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Select,
  Stack,
  Typography,
} from '@mui/material';
import { makeStyles } from '@mui/styles';
import React from 'react';

const useStyles = makeStyles(() => ({
  selectStyle: () => ({
    width: '100%',
    height: '48px',
  }),
  formControl: {
    width: '100%',
    height: '48px',
  },
}));

interface SubscriberDialogProps {
  sims: SimDto[];
  onClose: () => void;
  goBack: () => void;
  packages: PackageDto[];
  setFormData: Function;
  formData: TAddSubscriberData;
  handleSubmitButton: () => void;
}

const Step2: React.FC<SubscriberDialogProps> = React.memo(
  ({
    sims,
    goBack,
    onClose,
    packages,
    formData,
    setFormData,
    handleSubmitButton,
  }) => {
    const classes = useStyles();
    return (
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <FormControl variant="outlined" className={classes.formControl}>
            <InputLabel
              shrink
              variant="outlined"
              required
              htmlFor="outlined-age-always-notched"
            >
              {formData.simType == 'pSim' ? `pSIM ICCID` : `eSIM ICCID`}
            </InputLabel>

            <Select
              required
              variant="outlined"
              value={formData.iccid}
              onChange={(e) =>
                setFormData({ ...formData, iccid: e.target.value })
              }
              sx={{
                '& legend': { width: '93px' },
              }}
              input={
                <OutlinedInput
                  notched
                  fullWidth
                  label="ICCID"
                  name={'iccid'}
                  id="outlined-age-always-notched"
                />
              }
              MenuProps={{
                disablePortal: false,
                PaperProps: {
                  sx: {
                    boxShadow:
                      '0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)',
                    borderRadius: '4px',
                  },
                },
              }}
              className={classes.selectStyle}
            >
              {sims.length === 0 ? (
                <MenuItem
                  disabled
                  value={''}
                  sx={{
                    m: 0,
                    p: '6px 16px',
                  }}
                >
                  <Typography variant="body1">
                    {
                      'Sim pool is empty available. Please upload sims to simpool first.'
                    }
                  </Typography>
                </MenuItem>
              ) : (
                sims.map((sim) => (
                  <MenuItem
                    key={sim.id}
                    value={sim.iccid}
                    sx={{
                      m: 0,
                      p: '6px 16px',
                    }}
                  >
                    <Typography variant="body1">{sim.iccid}</Typography>
                  </MenuItem>
                ))
              )}
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={12}>
          <FormControl variant="outlined" className={classes.formControl}>
            <InputLabel
              shrink
              variant="outlined"
              required
              htmlFor="outlined-age-always-notched"
            >
              DATA PLAN
            </InputLabel>

            <Select
              required
              variant="outlined"
              value={formData.plan}
              onChange={(e) => {
                setFormData({ ...formData, plan: e.target.value });
              }}
              sx={{
                '& legend': { width: '93px' },
              }}
              input={
                <OutlinedInput
                  fullWidth
                  notched
                  label="Plan"
                  name={'plan'}
                  id="outlined-age-always-notched"
                />
              }
              MenuProps={{
                disablePortal: false,
                PaperProps: {
                  sx: {
                    boxShadow:
                      '0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)',
                    borderRadius: '4px',
                  },
                },
              }}
              className={classes.selectStyle}
            >
              {packages.length === 0 ? (
                <MenuItem
                  disabled
                  value={''}
                  sx={{
                    m: 0,
                    p: '6px 16px',
                  }}
                >
                  <Typography variant="body1">
                    {'No packages available. Please add packages first.'}
                  </Typography>
                </MenuItem>
              ) : (
                packages.map((pkg) => (
                  <MenuItem
                    key={pkg.uuid}
                    value={pkg.uuid}
                    sx={{
                      m: 0,
                      p: '6px 16px',
                    }}
                  >
                    <Typography variant="body1">
                      {`${pkg.name} - ${pkg.currency} ${pkg.amount}/${pkg.dataVolume} ${pkg.dataUnit}`}
                    </Typography>
                  </MenuItem>
                ))
              )}
            </Select>
          </FormControl>
        </Grid>

        <Grid item xs={12}>
          <Stack
            direction="row"
            justifyContent="space-between"
            mt={1}
            sx={{ mb: 2 }}
          >
            <Button
              variant="text"
              onClick={() => {
                goBack();
              }}
            >
              {'Go Back'}
            </Button>

            <Stack direction="row" spacing={3}>
              <Button
                variant="text"
                onClick={() => {
                  onClose();
                }}
              >
                {'CANCEL'}
              </Button>

              <Button
                type="submit"
                variant="contained"
                onClick={handleSubmitButton}
                disabled={!formData.iccid || !formData.plan}
              >
                <Typography variant="body1"> ADD SUBSCRIBER</Typography>
              </Button>
            </Stack>
          </Stack>
        </Grid>
      </Grid>
    );
  },
);
Step2.displayName = 'Step2';
export default Step2;
