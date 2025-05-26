/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { DATA_DURATION, DATA_UNIT } from '@/constants';
import { DataPlanSchema } from '@/helpers/formValidators';
import { CreatePlanType } from '@/types';
import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Dialog,
  DialogContent,
  DialogTitle,
  FormControl,
  IconButton,
  InputAdornment,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  TextField,
} from '@mui/material';
import Grid from '@mui/material/Grid2';
import { Formik } from 'formik';
import { Dispatch, SetStateAction } from 'react';

interface IDataPlanDialog {
  data: CreatePlanType;
  setData: Dispatch<SetStateAction<CreatePlanType>>;
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
        <Formik
          initialValues={dataplan}
          validationSchema={DataPlanSchema}
          onSubmit={(values) => {
            setDataplan({
              ...dataplan,
              name: values.name,
              amount: values.amount,
              dataUnit: values.dataUnit,
              duration: values.duration,
              dataVolume: values.dataVolume,
            });
            handleSuccessAction(action, {
              ...dataplan,
              name: values.name,
              amount: values.amount,
              dataUnit: values.dataUnit,
              duration: values.duration,
              dataVolume: values.dataVolume,
            });
          }}
        >
          {({
            values,
            errors,
            touched,
            handleBlur,
            handleSubmit,
            setFieldValue,
          }) => (
            <form onSubmit={handleSubmit}>
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
                    id="name"
                    fullWidth
                    required
                    label="DATA PLAN NAME"
                    value={values.name}
                    InputLabelProps={{
                      shrink: true,
                    }}
                    helperText={touched.name && errors.name}
                    error={touched.name && Boolean(errors.name)}
                    onChange={(e) => setFieldValue('name', e.target.value)}
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
                      <TextField
                        id="price"
                        fullWidth
                        required
                        label="PRICE"
                        onBlur={handleBlur}
                        InputProps={{
                          startAdornment: (
                            <InputAdornment position="start">
                              {currencySymbol}
                            </InputAdornment>
                          ),
                        }}
                        value={values.amount || ''}
                        onChange={(e) =>
                          setFieldValue(
                            'amount',
                            parseInt(e.target.value) || '',
                          )
                        }
                        error={touched.amount && Boolean(errors.amount)}
                        helperText={touched.amount && errors.amount}
                      />
                    </Grid>
                    <Grid size={{ xs: 6 }}>
                      <TextField
                        fullWidth
                        required
                        label="DATA LIMIT"
                        value={values.dataVolume || ''}
                        id="dataVolume"
                        InputLabelProps={{
                          shrink: true,
                        }}
                        onBlur={handleBlur}
                        onChange={(e) =>
                          setFieldValue(
                            'dataVolume',
                            parseInt(e.target.value) || '',
                          )
                        }
                        helperText={touched.dataVolume && errors.dataVolume}
                        error={touched.dataVolume && Boolean(errors.dataVolume)}
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
                        <InputLabel id={'unit-label'} shrink>
                          UNIT*
                        </InputLabel>
                        <Select
                          notched
                          required
                          label="UNIT"
                          onBlur={handleBlur}
                          value={values.dataUnit}
                          id={'unit'}
                          labelId="unit-label"
                          onChange={(e) =>
                            setFieldValue('dataUnit', e.target.value)
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
                      <TextField
                        select
                        fullWidth
                        required
                        id="duration"
                        label="DURATION"
                        onBlur={handleBlur}
                        InputLabelProps={{
                          shrink: true,
                        }}
                        value={values.duration}
                        onChange={(e) =>
                          setFieldValue('duration', e.target.value)
                        }
                        helperText={touched.duration && errors.duration}
                        error={touched.duration && Boolean(errors.duration)}
                      >
                        {DATA_DURATION.map(({ id, label, value }) => (
                          <MenuItem key={id} value={value}>
                            {label}
                          </MenuItem>
                        ))}
                      </TextField>
                    </Grid>
                  </Grid>
                )}

                <Grid size={{ xs: 12 }} mt={2}>
                  <Stack
                    direction={'row'}
                    justifyContent={'flex-end'}
                    spacing={2}
                  >
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
                      <Button type="submit" variant="contained">
                        {labelSuccessBtn}
                      </Button>
                    )}
                  </Stack>
                </Grid>
              </Grid>
            </form>
          )}
        </Formik>
      </DialogContent>
    </Dialog>
  );
};

export default DataPlanDialog;
