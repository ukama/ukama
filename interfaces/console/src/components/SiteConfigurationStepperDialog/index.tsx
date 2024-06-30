import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Stack,
  IconButton,
  Stepper,
  Step,
  StepLabel,
} from '@mui/material';
import { Formik, Form } from 'formik';
import CloseIcon from '@mui/icons-material/Close';
import { SITE_CONFIG_STEPS } from '@/constants';
import { STEPPER_FORM_SCHEMA } from '@/helpers/formValidators';
import SiteStepForm from '@/components/SiteStepForm';

interface SiteConfigurationStepperDialogProps {
  open: boolean;
  handleClose: () => void;
  handleFormDataSubmit: (formData: any) => void;
}

const SiteConfigurationStepperDialog: React.FC<
  SiteConfigurationStepperDialogProps
> = ({ open, handleClose, handleFormDataSubmit }) => {
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
      await STEPPER_FORM_SCHEMA[activeStep].validate(values, {
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
        {({ handleSubmit }) => (
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
              />
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
