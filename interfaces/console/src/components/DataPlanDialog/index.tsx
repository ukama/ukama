/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { DATA_DURATION, DATA_UNIT } from '@/constants';
import { CreatePlanType } from '@/types';
import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  IconButton,
  InputAdornment,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Select,
  Stack,
  TextField,
} from '@mui/material';
import Grid from '@mui/material/Grid2';
import { Dispatch, SetStateAction, useMemo } from 'react';

interface IDataPlanDialog {
  data: CreatePlanType;
  setData: Dispatch<SetStateAction<any>>;
  title: string;
  action: string;
  isOpen: boolean;
  currencySymbol?: string;
  handleCloseAction: any;
  labelSuccessBtn?: string;
  handleSuccessAction: any;
  labelNegativeBtn?: string;
}

const DataPlanDialog = ({
  title,
  isOpen,
  action,
  currencySymbol,
  data: dataplan,
  setData: setDataplan,
  labelSuccessBtn,
  labelNegativeBtn,
  handleCloseAction,
  handleSuccessAction,
}: IDataPlanDialog) => {
  // Check if all required fields are filled
  const isFormValid = useMemo(() => {
    if (action === 'update') {
      return !!dataplan.name;
    }

    return !!(
      dataplan.name &&
      dataplan.amount &&
      dataplan.dataVolume &&
      dataplan.dataUnit &&
      dataplan.duration
    );
  }, [dataplan, action]);

  return (
    <Dialog
      fullWidth
      open={isOpen}
      maxWidth="sm"
      onClose={handleCloseAction}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <DialogTitle>{title}</DialogTitle>
        <IconButton
          onClick={handleCloseAction}
          sx={{ position: 'relative', right: 8 }}
        >
          <CloseIcon />
        </IconButton>
      </Stack>

      <DialogContent>
        <Grid
          container
          rowSpacing={2}
          gridAutoRows={2}
          columnSpacing={2}
          gridAutoColumns={1}
          alignItems={'center'}
          justifyContent={'center'}
        >
          <Grid size={{ xs: 12 }}>
            <TextField
              fullWidth
              required
              label="DATA PLAN NAME"
              value={dataplan.name}
              id={'data-plan-name'}
              InputLabelProps={{
                shrink: true,
              }}
              onChange={(e) =>
                setDataplan({ ...dataplan, name: e.target.value })
              }
            />
          </Grid>
          {action !== 'update' && (
            <Grid
              container
              size={{ xs: 12, sm: 6 }}
              columnSpacing={1}
              rowSpacing={2}
            >
              <Grid size={{ xs: 6 }}>
                <FormControl fullWidth>
                  <InputLabel htmlFor="package-price">PRICE</InputLabel>
                  <OutlinedInput
                    id="package-price"
                    fullWidth
                    required
                    type="number"
                    label="PRICE"
                    startAdornment={
                      <InputAdornment position="start">
                        {currencySymbol}
                      </InputAdornment>
                    }
                    value={dataplan.amount}
                    onChange={(e) =>
                      setDataplan({
                        ...dataplan,
                        amount: parseInt(e.target.value),
                      })
                    }
                  />
                </FormControl>
              </Grid>
              <Grid size={{ xs: 6 }}>
                <TextField
                  fullWidth
                  required
                  type="number"
                  label="DATA LIMIT"
                  value={dataplan.dataVolume}
                  id={'data-plan-limit'}
                  InputLabelProps={{
                    shrink: true,
                  }}
                  onChange={(e) =>
                    setDataplan({
                      ...dataplan,
                      dataVolume: parseInt(e.target.value),
                    })
                  }
                />
              </Grid>
            </Grid>
          )}
          {action !== 'update' && (
            <Grid
              container
              size={{ xs: 12, sm: 6 }}
              columnSpacing={1}
              rowSpacing={2}
            >
              <Grid size={{ xs: 5 }}>
                <FormControl fullWidth>
                  <InputLabel id={'data-plan-unit-label'} shrink>
                    UNIT*
                  </InputLabel>
                  <Select
                    notched
                    required
                    label="UNIT"
                    value={dataplan.dataUnit}
                    id={'data-plan-unit'}
                    labelId="data-plan-unit-label"
                    onChange={(e) =>
                      setDataplan({
                        ...dataplan,
                        dataUnit: e.target.value,
                      })
                    }
                  >
                    {DATA_UNIT.map(({ id, label, value }) => (
                      <MenuItem key={id} value={value}>
                        {label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>

              <Grid size={{ xs: 7 }}>
                <FormControl fullWidth>
                  <InputLabel
                    id={'data-plan-unit-price-label'}
                    shrink
                    sx={{
                      '& legend': {
                        letterSpacing: 0.6,
                      },
                    }}
                  >
                    DURATION*
                  </InputLabel>
                  <Select
                    notched
                    required
                    label="UNIT"
                    value={dataplan.duration}
                    id={'data-plan-unit'}
                    labelId="data-plan-unit-price-label"
                    sx={{
                      '& legend': {
                        letterSpacing: 4,
                      },
                    }}
                    onChange={(e) =>
                      setDataplan({
                        ...dataplan,
                        duration: parseInt(e.target.value as string),
                      })
                    }
                  >
                    {DATA_DURATION.map(({ id, label, value }) => (
                      <MenuItem key={id} value={value}>
                        {label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          )}
        </Grid>
      </DialogContent>

      <DialogActions>
        <Stack direction={'row'} alignItems="center" spacing={2}>
          {labelNegativeBtn && (
            <Button
              variant="text"
              color={'primary'}
              onClick={handleCloseAction}
            >
              {labelNegativeBtn}
            </Button>
          )}
          {labelSuccessBtn && (
            <Button
              variant="contained"
              onClick={() => handleSuccessAction(action)}
              disabled={!isFormValid}
            >
              {labelSuccessBtn}
            </Button>
          )}
        </Stack>
      </DialogActions>
    </Dialog>
  );
};

export default DataPlanDialog;
