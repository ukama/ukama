/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { globalUseStyles } from '@/styles/global';
import colors from '@/theme/colors';
import { TAddSubscriberData } from '@/types';
import {
  Button,
  FormControl,
  FormControlLabel,
  Grid,
  Radio,
  RadioGroup,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { Form, Formik } from 'formik';
import React from 'react';
import * as Yup from 'yup';

const validationSchema = Yup.object({
  name: Yup.string().required('Name is required'),
  email: Yup.string().email('Invalid email').required('Email is required'),
  phone: Yup.string().optional(),
  simTYpe: Yup.string().default('eSim'),
  roamingStatus: Yup.boolean().default(false),
});

interface SubscriberDialogProps {
  onClose: () => void;
  handleStep1Submit: Function;
  eSimCount: number | undefined;
  pSimCount: number | undefined;
  formData: TAddSubscriberData;
}

const Step1: React.FC<SubscriberDialogProps> = ({
  onClose,
  handleStep1Submit,
  eSimCount,
  pSimCount,
  formData,
}) => {
  const initialValues: TAddSubscriberData = {
    name: formData.name,
    email: formData.email,
    phone: formData.phone,
    simType: formData.simType,
    roamingStatus: formData.roamingStatus,
    iccid: formData.iccid,
    plan: formData.plan,
  };
  const gclasses = globalUseStyles();

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={async (values) => {
        handleStep1Submit(values);
      }}
    >
      {({
        values,
        errors,
        touched,
        handleChange,
        handleSubmit,
        handleBlur,
      }) => (
        <Form onSubmit={handleSubmit}>
          <Grid container rowSpacing={2}>
            <Grid item xs={12}>
              <TextField
                required
                fullWidth
                label={'NAME'}
                name={'name'}
                InputLabelProps={{
                  shrink: true,
                }}
                onChange={handleChange}
                value={values.name}
                onBlur={handleBlur}
                helperText={touched.name && errors.name}
                error={touched.name && Boolean(errors.name)}
                id={'name'}
                spellCheck={false}
                InputProps={{
                  classes: {
                    input: gclasses.inputFieldStyle,
                  },
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <FormControl>
                <RadioGroup
                  row
                  id="simType"
                  name="simType"
                  aria-labelledby="simType-rg"
                  value={values.simType}
                  onChange={handleChange}
                  sx={{ px: 1.25 }}
                >
                  <FormControlLabel
                    value="eSim"
                    control={<Radio />}
                    label={`eSIM (${eSimCount || 0} left)`}
                    sx={{
                      pl: 1,
                      pr: 5,
                      py: 1.5,
                      borderRadius: '4px',
                      border: `1px solid ${colors.black10}`,
                    }}
                  />
                  <FormControlLabel
                    value="pSim"
                    control={<Radio />}
                    label={`pSIM (${pSimCount || 0} left)`}
                    sx={{
                      pl: 1,
                      pr: 5,
                      py: 1.5,
                      borderRadius: '4px',
                      border: `1px solid ${colors.black10}`,
                    }}
                  />
                </RadioGroup>
              </FormControl>
            </Grid>

            <Grid item xs={12}>
              <TextField
                required
                fullWidth
                label={'EMAIL'}
                name={'email'}
                InputLabelProps={{
                  shrink: true,
                }}
                onBlur={handleBlur}
                onChange={handleChange}
                value={values.email}
                helperText={touched.email && errors.email}
                error={touched.email && Boolean(errors.email)}
                id={'email'}
                spellCheck={false}
                InputProps={{
                  classes: {
                    input: gclasses.inputFieldStyle,
                  },
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label={'PHONE (OPTIONAL)'}
                name={'phone'}
                InputLabelProps={{
                  shrink: true,
                }}
                onBlur={handleBlur}
                onChange={handleChange}
                value={values.phone}
                helperText={touched.phone && errors.phone}
                error={touched.phone && Boolean(errors.phone)}
                id={'phone'}
                spellCheck={false}
                InputProps={{
                  classes: {
                    input: gclasses.inputFieldStyle,
                  },
                }}
              />
            </Grid>
            <Grid item xs={12}>
              {/* <Stack direction="column" spacing={1}>
                <Typography variant="caption" sx={{ color: colors.black38 }}>
                  ALLOW ROAMING
                </Typography>
                <Stack direction="row" spacing={1} alignItems="center">
                  <Typography variant="body1">
                    Roaming allows subscriber to use data outside of their home
                    Ukama network. Insert billing information.
                  </Typography>
                  <Switch
                    id="roamingStatus"
                    name="roamingStatus"
                    size="small"
                    value={values.roamingStatus}
                    checked={values.roamingStatus}
                    onChange={handleChange}
                  />
                </Stack>
              </Stack> */}

              <Stack
                direction="row"
                justifyContent={'flex-end'}
                mt={1}
                spacing={1}
                sx={{ mb: 2 }}
              >
                <Stack direction="row" spacing={3}>
                  <Button variant="text" onClick={onClose}>
                    {'CANCEL'}
                  </Button>

                  <Button
                    variant="contained"
                    type="submit"
                    sx={{
                      disabled: {
                        color: colors.primaryLight,
                      },
                    }}
                  >
                    <Typography variant="body1">NEXT</Typography>
                  </Button>
                </Stack>
              </Stack>
            </Grid>
          </Grid>
        </Form>
      )}
    </Formik>
  );
};

export default Step1;
