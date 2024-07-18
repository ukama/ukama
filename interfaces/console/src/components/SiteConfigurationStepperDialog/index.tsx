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
import { Formik, Form, FormikErrors, FormikTouched, FormikProps } from 'formik';
import * as Yup from 'yup';
import CloseIcon from '@mui/icons-material/Close';
import { SITE_CONFIG_STEPS } from '@/constants';
import { STEPPER_FORM_SCHEMA } from '@/helpers/formValidators';
import SiteStepForm from '@/components/SiteStepForm';
import { FormValues } from '@/types';

interface SiteConfigurationStepperDialogProps {
  open: boolean;
  handleClose: () => void;
  handleFormDataSubmit: (formData: FormValues & { location: string }) => void;
  components: any[];
}

const SiteConfigurationStepperDialog: React.FC<
  SiteConfigurationStepperDialogProps
> = ({ open, handleClose, handleFormDataSubmit, components }) => {
  const [activeStep, setActiveStep] = useState(0);
  const [lat, setLat] = useState<number>(0);
  const [lng, setLng] = useState<number>(0);
  const [location, setLocation] = useState('');
  const [component, setComponent] = useState<any>();

  useEffect(() => {
    console.log('Active step changed:', activeStep);
  }, [activeStep]);

  const initialValues: FormValues = {
    switch: '',
    power: '',
    backhaul: '',
    access: '',
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
    console.log('handleNext called', values, activeStep);
    try {
      await STEPPER_FORM_SCHEMA[activeStep].validate(values, {
        abortEarly: false,
      });
      console.log('Validation passed');
      if (activeStep === SITE_CONFIG_STEPS.length - 1) {
        handleFormDataSubmit({ ...values, location });
        handleClose();
        console.log('DATA:', values);
      } else {
        setActiveStep((prevActiveStep) => {
          console.log(
            'Updating activeStep from',
            prevActiveStep,
            'to',
            prevActiveStep + 1,
          );
          return prevActiveStep + 1;
        });
      }
    } catch (error) {
      console.error('Validation error:', error);
      if (error instanceof Yup.ValidationError) {
        const errorMessages: { [key: string]: string } = {};
        error.inner.forEach((err) => {
          if (err.path) {
            errorMessages[err.path] = err.message;
          }
        });
        console.log('Validation errors:', errorMessages);
        helpers.setErrors(errorMessages);
      }
      helpers.setSubmitting(false);
    }
  };

  const getComponentInfos = (form: any) => {
    setComponent(form.values);
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
        {(formik: FormikProps<FormValues>) => (
          <Form onSubmit={formik.handleSubmit}>
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
                getComponentInfos={getComponentInfos}
                formik={formik}
                components={components}
              />
            </DialogContent>
            <DialogActions>
              <Button disabled={activeStep === 0} onClick={handleBack}>
                Back
              </Button>
              <Button type="submit">
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
