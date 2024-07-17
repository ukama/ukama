import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  IconButton,
  Stepper,
  Step,
  StepLabel,
} from '@mui/material';
import { Formik, Form, FormikErrors, FormikTouched } from 'formik';
import * as Yup from 'yup';
import CloseIcon from '@mui/icons-material/Close';
import { SITE_CONFIG_STEPS } from '@/constants';
import { STEPPER_FORM_SCHEMA } from '@/helpers/formValidators';
import SiteStepForm from '@/components/SiteStepForm';

interface SiteConfigurationStepperDialogProps {
  open: boolean;
  handleClose: () => void;
  handleFormDataSubmit: (formData: any) => void;
  components: any[];
}

interface FormValues {
  switch: string;
  power: string;
  backhaul: string;
  spectrumBand: string;
  location: string;
  siteName: string;
  network: string;
  latitude: number;
  longitude: number;
}

const SiteConfigurationStepperDialog: React.FC<
  SiteConfigurationStepperDialogProps
> = ({ open, handleClose, handleFormDataSubmit, components }) => {
  const [activeStep, setActiveStep] = useState(0);
  const [lat, setLat] = useState<number>(0);
  const [lng, setLng] = useState<number>(0);
  const [location, setLocation] = useState('');

  useEffect(() => {
    console.log('Active step changed:', activeStep);
  }, [activeStep]);

  const initialValues: FormValues = {
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

  const handleNext = async (values: FormValues, helpers: any) => {
    console.log('handleNext called', values);
    try {
      await STEPPER_FORM_SCHEMA[activeStep].validate(values, {
        abortEarly: false,
      });
      console.log('Validation passed');
      if (activeStep === SITE_CONFIG_STEPS.length - 1) {
        handleFormDataSubmit({ ...values, location: location });
        handleClose();
      } else {
        setActiveStep((prevActiveStep) => prevActiveStep + 1);
        console.log('Active step updated', activeStep + 1);
      }
    } catch (error) {
      console.error('Validation error:', error);
      helpers.setSubmitting(false);
      if (error instanceof Yup.ValidationError) {
        const errorMessages: { [key: string]: string } = {};
        error.inner.forEach((err) => {
          if (err.path) {
            errorMessages[err.path] = err.message;
          }
        });
        helpers.setErrors(errorMessages);
      }
    }
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleAddressChange = (address: string) => {
    setLocation(address);
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="sm"
      fullWidth
      aria-labelledby="site-config-dialog-title"
    >
      <DialogTitle>
        <Typography variant="h6" color="initial">
          Site Configuration
        </Typography>
        <IconButton
          aria-label="close"
          onClick={handleClose}
          sx={{ position: 'absolute', right: 8, top: 8 }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <Formik
        initialValues={initialValues}
        validationSchema={STEPPER_FORM_SCHEMA[activeStep]}
        onSubmit={(values, helpers) => handleNext(values, helpers)}
      >
        {({ handleSubmit, errors, touched }) => (
          <Form onSubmit={handleSubmit}>
            <DialogContent>
              <Stepper activeStep={activeStep} alternativeLabel>
                {SITE_CONFIG_STEPS.map((label) => (
                  <Step key={label}>
                    <StepLabel>{label}</StepLabel>
                  </Step>
                ))}
              </Stepper>
              <SiteStepForm
                step={activeStep}
                lat={lat}
                lng={lng}
                location={location}
                handleLatChange={handleLatChange}
                handleLngChange={handleLngChange}
                handleAddressChange={handleAddressChange}
                components={components}
              />
              {Object.keys(errors).map((key) => {
                const touchedKey = key as keyof FormikTouched<FormValues>;
                const errorKey = key as keyof FormikErrors<FormValues>;
                return touched[touchedKey] && errors[errorKey] ? (
                  <Typography key={key} color="error">
                    {errors[errorKey]}
                  </Typography>
                ) : null;
              })}
            </DialogContent>
            <DialogActions>
              {activeStep > 0 && (
                <Button
                  onClick={handleBack}
                  variant="contained"
                  color="primary"
                >
                  Back
                </Button>
              )}
              <Button type="submit" variant="contained" color="primary">
                {activeStep === SITE_CONFIG_STEPS.length - 1
                  ? 'Finish'
                  : 'Next'}
              </Button>
            </DialogActions>
          </Form>
        )}
      </Formik>
    </Dialog>
  );
};

export default SiteConfigurationStepperDialog;
