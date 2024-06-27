/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Stepper,
  Step,
  StepLabel,
  TextField,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  Typography,
  Stack,
  IconButton,
} from '@mui/material';
import { Formik, Form, Field, ErrorMessage, useField } from 'formik';
import * as Yup from 'yup';
import { globalUseStyles } from '@/styles/global';
import colors from '@/theme/colors';
import { SITE_CONFIG_STEPS } from '@/constants';
import CloseIcon from '@mui/icons-material/Close';
import dynamic from 'next/dynamic';

interface SiteConfigurationStepperDialogProps {
  open: boolean;
  handleClose: () => void;
  handleFormDataSubmit: (formData: any) => void;
}

interface CustomTextFieldProps {
  label: string;
  name: string;
  onChange?: (event: any) => void;
}

const CustomTextField: React.FC<CustomTextFieldProps> = ({
  label,
  name,
  onChange,
}) => {
  const [field, meta] = useField(name);
  const gclasses = globalUseStyles();

  return (
    <TextField
      {...field}
      label={label}
      fullWidth
      onChange={(e) => {
        field.onChange(e);
        if (onChange) {
          onChange(e);
        }
      }}
      InputLabelProps={{
        shrink: true,
      }}
      helperText={meta.touched && meta.error ? meta.error : ''}
      error={meta.touched && Boolean(meta.error)}
      spellCheck={false}
      InputProps={{
        classes: {
          input: gclasses.inputFieldStyle,
        },
      }}
    />
  );
};

const validationSchema = [
  Yup.object().shape({
    switch: Yup.string().required('Switch is required'),
    power: Yup.string().required('Power is required'),
    backhaul: Yup.string().required('Backhaul is required'),
    spectrumBand: Yup.string().required('Spectrum Band is required'),
  }),
  Yup.object().shape({
    siteName: Yup.string().required('Site Name is required'),
    network: Yup.string().required('Network is required'),
    latitude: Yup.number()
      .required('Latitude is required')
      .min(-90, 'Latitude must be between -90 and 90')
      .max(90, 'Latitude must be between -90 and 90'),
    longitude: Yup.number()
      .required('Longitude is required')
      .min(-180, 'Longitude must be between -180 and 180')
      .max(180, 'Longitude must be between -180 and 180'),
  }),
];

const SiteConfigurationStepperDialog: React.FC<
  SiteConfigurationStepperDialogProps
> = ({ open, handleClose, handleFormDataSubmit }) => {
  const SiteMapComponent = dynamic(() => import('../SiteMapComponent'), {
    loading: () => <p>Site map is loading</p>,
    ssr: false,
  });
  const [activeStep, setActiveStep] = useState(0);
  const [lat, setLat] = useState<number>(0);
  const [lng, setLng] = useState<number>(0);
  const [location, setLocation] = useState('');

  const initialValues = {
    switch: '',
    power: '',
    backhaul: '',
    spectrumBand: '',
    location: '',
    siteName: '',
    network: '',
    latitude: 0,
    longitude: 0,
  };

  const handleLatChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = Number(e.target.value);
    setLat(value);
  };

  const handleLngChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = Number(e.target.value);
    setLng(value);
  };

  const handleNext = async (values: any, helpers: any) => {
    try {
      await validationSchema[activeStep].validate(values, {
        abortEarly: false,
      });
      if (activeStep === SITE_CONFIG_STEPS.length - 1) {
        handleFormDataSubmit({ ...values, location: location });
        handleClose();
      } else {
        setActiveStep((prevActiveStep) => prevActiveStep + 1);
      }
    } catch (error) {
      helpers.setSubmitting(false);
    }
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleAddressChange = (address: string) => {
    setLocation(address);
  };
  const stepContent = (step: number) => {
    switch (step) {
      case 0:
        return (
          <Box
            component="form"
            style={{ display: 'flex', flexDirection: 'column', gap: 16 }}
          >
            <Box sx={{ mt: 2, mb: 2 }}>
              <Typography>
                You have successfully installed your site, and need to configure
                it. Please note that if your power or backhaul choice is
                “other”, it can’t be monitored within Ukama’s Console.
              </Typography>
            </Box>

            <FormControl fullWidth>
              <InputLabel id="switch-label">Switch</InputLabel>
              <Field
                as={Select}
                labelId="switch-label"
                name="switch"
                label="Switch"
              >
                <MenuItem value="8 port switch">8 port switch</MenuItem>
                <MenuItem value="16 port switch">16 port switch</MenuItem>
              </Field>

              <ErrorMessage name="switch" component="div">
                {(msg) => <div style={{ color: colors.red }}>{msg}</div>}
              </ErrorMessage>
            </FormControl>

            <FormControl fullWidth>
              <InputLabel id="power-label">Power</InputLabel>
              <Field
                as={Select}
                labelId="power-label"
                name="power"
                label="Power"
              >
                <MenuItem value="Battery">Battery</MenuItem>
                <MenuItem value="AC Power">AC Power</MenuItem>
              </Field>

              <ErrorMessage name="power" component="div">
                {(msg) => <div style={{ color: colors.red }}>{msg}</div>}
              </ErrorMessage>
            </FormControl>

            <FormControl fullWidth>
              <InputLabel id="backhaul-label">Backhaul</InputLabel>
              <Field
                as={Select}
                labelId="backhaul-label"
                name="backhaul"
                label="Backhaul"
              >
                <MenuItem value="ViaSAT">ViaSAT</MenuItem>
                <MenuItem value="Other">Other</MenuItem>
              </Field>

              <ErrorMessage name="backhaul" component="div">
                {(msg) => <div style={{ color: colors.red }}>{msg}</div>}
              </ErrorMessage>
            </FormControl>

            <FormControl fullWidth>
              <InputLabel id="spectrumBand-label">Spectrum Band</InputLabel>
              <Field
                as={Select}
                labelId="spectrumBand-label"
                name="spectrumBand"
                label="Spectrum Band"
              >
                <MenuItem value="Band 40">Band 40</MenuItem>
                <MenuItem value="Band 41">Band 41</MenuItem>
              </Field>

              <ErrorMessage name="spectrumBand" component="div">
                {(msg) => <div style={{ color: colors.red }}>{msg}</div>}
              </ErrorMessage>
            </FormControl>
          </Box>
        );
      case 1:
        return (
          <Box
            component="form"
            style={{ display: 'flex', flexDirection: 'column', gap: 16 }}
          >
            <Box sx={{ mt: 2, mb: 2 }}>
              <Typography>
                Please name your site for your ease of reference, and assign it
                to a network.
              </Typography>
            </Box>
            <SiteMapComponent
              posix={[lat, lng]}
              onAddressChange={handleAddressChange}
            />

            <Box>
              <Stack direction="column" spacing={1} justifyItems={'center'}>
                <Typography
                  variant="body2"
                  sx={{ color: `${colors.darkGray}` }}
                >
                  LOCATION
                </Typography>
                <Typography variant="body2" color="initial">
                  {location || 'Fetching site location...'}
                </Typography>
              </Stack>
            </Box>
            <CustomTextField
              label="Longitude"
              name="longitude"
              onChange={handleLngChange}
            />
            <CustomTextField
              label="Latitude"
              name="latitude"
              onChange={handleLatChange}
            />
            <CustomTextField label="Site Name" name="siteName" />
            <CustomTextField label="Network" name="network" />
          </Box>
        );
      default:
        return 'Unknown step';
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>
        Site Configuration
        <IconButton
          aria-label="close"
          onClick={handleClose}
          style={{ position: 'absolute', right: '8px', top: '8px' }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema[activeStep]}
        onSubmit={handleNext}
      >
        {({ handleSubmit, isSubmitting }) => (
          <Form onSubmit={handleSubmit}>
            <DialogContent>
              <Stepper activeStep={activeStep}>
                {SITE_CONFIG_STEPS.map((label) => (
                  <Step key={label}>
                    <StepLabel>{label}</StepLabel>
                  </Step>
                ))}
              </Stepper>
              {stepContent(activeStep)}
            </DialogContent>
            <DialogActions>
              <Button
                onClick={() => {
                  handleClose(), setActiveStep(0);
                }}
                color="secondary"
              >
                Cancel
              </Button>
              {activeStep === SITE_CONFIG_STEPS.length - 1 ? (
                <Button type="submit" color="primary" disabled={isSubmitting}>
                  Save
                </Button>
              ) : (
                <>
                  <Button
                    onClick={handleBack}
                    color="primary"
                    disabled={activeStep === 0 || isSubmitting}
                  >
                    Back
                  </Button>
                  <Button type="submit" color="primary" disabled={isSubmitting}>
                    Next
                  </Button>
                </>
              )}
            </DialogActions>
          </Form>
        )}
      </Formik>
    </Dialog>
  );
};

export default SiteConfigurationStepperDialog;
