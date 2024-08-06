import { NetworkDto } from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import { AddSiteValidationSchema } from '@/helpers/formValidators';
import { globalUseStyles } from '@/styles/global';
import colors from '@/theme/colors';
import { TSiteForm } from '@/types';
import { isValidLatLng } from '@/utils';
import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Dialog,
  DialogContent,
  DialogContentText,
  DialogTitle,
  MenuItem,
  Stack,
  Step,
  StepLabel,
  Stepper,
  IconButton,
  TextField,
  Typography,
} from '@mui/material';
import { Field, Form, Formik } from 'formik';
import dynamic from 'next/dynamic';
import React, { useEffect, useState } from 'react';
import { useFetchAddress } from '@/utils/useFetchAddress';

const SiteMapComponent = dynamic(() => import('../SiteMapComponent'), {
  loading: () => <p>Site map is loading</p>,
  ssr: false,
});

interface Component {
  id: string;
  inventory_id: string;
  category: string;
  type: string;
  user_id: string;
  description: string;
  datasheet_url: string;
  images_url: string;
  part_number: string;
  manufacturer: string;
  managed: string;
  warranty: number;
  specification: string;
}

interface IConfigureSiteDialog {
  open: boolean;
  site: TSiteForm;
  onClose: () => void;
  addSiteLoading: boolean;
  networks: NetworkDto[];
  components: Component[];
  handleSiteConfiguration: (values: TSiteForm) => void;
}

const steps = [
  'Select Switch, Power, Backhaul, and spectrum band',
  'Enter your site details',
];

const ConfigureSiteDialog: React.FC<IConfigureSiteDialog> = ({
  open,
  site,
  onClose,
  networks,
  components,
  addSiteLoading = false,
  handleSiteConfiguration,
}) => {
  const { setSnackbarMessage } = useAppContext();
  const gclasses = globalUseStyles();
  const [activeStep, setActiveStep] = useState(0);
  const [formValues, setFormValues] = useState(site);
  const { address, isLoading, error, fetchAddress } = useFetchAddress();
  const resetForm = () => {
    setFormValues({
      switch: '',
      power: '',
      backhaul: '',
      access: '',
      spectrum: '',
      siteName: '',
      network: '',
      latitude: 0,
      longitude: 0,
      address: '',
    });
    setActiveStep(0);
  };

  useEffect(() => {
    if (open) {
      resetForm();
      setFormValues(site);
    }
  }, [open, site]);
  const handleNext = () => setActiveStep((prevStep) => prevStep + 1);
  const handleBack = () => setActiveStep((prevStep) => prevStep - 1);

  const handleSubmit = (values: TSiteForm) => {
    if (address === 'Location not found') {
      setSnackbarMessage({
        id: 'error-fetching-address',
        type: 'error',
        show: true,
        message: 'Error fetching address from coordinates',
      });
    }

    handleSiteConfiguration({
      ...values,
      address: address,
    });
    resetForm();
    onClose();
  };

  const handleStepSubmit = (values: Partial<TSiteForm>) => {
    setFormValues((prev) => ({ ...prev, ...values }));
    handleNext();
  };

  const switchComponents = components.filter(
    (comp) => comp.category === 'SWITCH',
  );
  const powerComponents = components.filter(
    (comp) => comp.category === 'POWER',
  );
  const backhaulComponents = components.filter(
    (comp) => comp.category === 'BACKHAUL',
  );
  const accessComponents = components.filter(
    (comp) => comp.category === 'ACCESS',
  );

  useEffect(() => {
    if (error) {
      setSnackbarMessage({
        id: 'error-fetching-address',
        type: 'error',
        show: true,
        message: 'Error fetching address from coordinates',
      });
    }
  }, [error, setSnackbarMessage]);

  const handleFetchAddress = async (lat: number, lng: number) => {
    setSnackbarMessage({
      id: 'fetching-address',
      type: 'success',
      show: true,
      message: 'Fetching address with coordinates',
    });
    await fetchAddress(lat, lng);
  };
  const handleClose = () => {
    resetForm();
    onClose();
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      sx={{
        '& .MuiDialog-paper': {
          width: '60%',
          maxWidth: '40%',
        },
      }}
    >
      <DialogTitle>
        Configure site installation ({activeStep + 1}/2)
      </DialogTitle>
      <IconButton
        aria-label="close"
        onClick={handleClose}
        sx={{
          position: 'absolute',
          right: 8,
          top: 8,
          color: (theme) => theme.palette.grey[500],
        }}
      >
        <CloseIcon />
      </IconButton>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          {`You have successfully installed your site, and need to configure
              it. Please note that if your power or backhaul choice is "other",
              it can't be monitored within Ukama's Console.`}
        </DialogContentText>
      </DialogContent>

      <DialogContent>
        <Stepper activeStep={activeStep} sx={{ mb: 4 }}>
          {steps.map((label) => (
            <Step key={label}>
              <StepLabel>{label}</StepLabel>
            </Step>
          ))}
        </Stepper>
        {activeStep === 0 && (
          <Formik
            initialValues={formValues}
            onSubmit={handleStepSubmit}
            validationSchema={AddSiteValidationSchema[0]}
          >
            {({ errors, touched, isValid, handleReset }) => (
              <Form>
                <Stack>
                  <Field
                    as={TextField}
                    fullWidth
                    select
                    required
                    name="switch"
                    label="SWITCH"
                    margin="normal"
                    InputLabelProps={{
                      shrink: true,
                    }}
                    InputProps={{
                      classes: {
                        input: gclasses.inputFieldStyle,
                      },
                    }}
                    error={touched.switch && Boolean(errors.switch)}
                    helperText={touched.switch && errors.switch}
                  >
                    {switchComponents.map((component) => (
                      <MenuItem key={component.id} value={component.id}>
                        {component.description}
                      </MenuItem>
                    ))}
                  </Field>
                  <Field
                    as={TextField}
                    fullWidth
                    select
                    required
                    name="power"
                    label="POWER"
                    margin="normal"
                    InputLabelProps={{
                      shrink: true,
                    }}
                    InputProps={{
                      classes: {
                        input: gclasses.inputFieldStyle,
                      },
                    }}
                    error={touched.power && Boolean(errors.power)}
                    helperText={touched.power && errors.power}
                  >
                    {powerComponents.map((component) => (
                      <MenuItem key={component.id} value={component.id}>
                        {component.description}
                      </MenuItem>
                    ))}
                  </Field>
                  <Field
                    as={TextField}
                    fullWidth
                    select
                    required
                    name="backhaul"
                    label="BACKHAUL"
                    margin="normal"
                    InputLabelProps={{
                      shrink: true,
                    }}
                    InputProps={{
                      classes: {
                        input: gclasses.inputFieldStyle,
                      },
                    }}
                    error={touched.backhaul && Boolean(errors.backhaul)}
                    helperText={touched.backhaul && errors.backhaul}
                  >
                    {backhaulComponents.map((component) => (
                      <MenuItem key={component.id} value={component.id}>
                        {component.description}
                      </MenuItem>
                    ))}
                  </Field>
                  <Field
                    as={TextField}
                    fullWidth
                    select
                    required
                    name="access"
                    label="ACCESS"
                    margin="normal"
                    InputLabelProps={{
                      shrink: true,
                    }}
                    InputProps={{
                      classes: {
                        input: gclasses.inputFieldStyle,
                      },
                    }}
                    error={touched.access && Boolean(errors.access)}
                    helperText={touched.access && errors.access}
                  >
                    {accessComponents.map((component) => (
                      <MenuItem key={component.id} value={component.id}>
                        {component.description}
                      </MenuItem>
                    ))}
                  </Field>
                  <Field
                    as={TextField}
                    fullWidth
                    select
                    required
                    name="spectrum"
                    label="SPECTRUM BAND"
                    margin="normal"
                    InputLabelProps={{
                      shrink: true,
                    }}
                    InputProps={{
                      classes: {
                        input: gclasses.inputFieldStyle,
                      },
                    }}
                    error={touched.spectrum && Boolean(errors.spectrum)}
                    helperText={touched.spectrum && errors.spectrum}
                  >
                    {accessComponents.map((component) => (
                      <MenuItem key={component.id} value={component.id}>
                        {component.description}
                      </MenuItem>
                    ))}
                  </Field>
                </Stack>
                <Stack
                  direction="row"
                  spacing={1}
                  justifyItems={'center'}
                  justifyContent={'flex-end'}
                  sx={{ mt: 1 }}
                >
                  <Button
                    onClick={() => {
                      handleReset();
                      handleClose();
                    }}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    disabled={!isValid}
                  >
                    Next
                  </Button>
                </Stack>
              </Form>
            )}
          </Formik>
        )}
        {activeStep === 1 && (
          <Formik
            initialValues={formValues}
            onSubmit={handleSubmit}
            validationSchema={AddSiteValidationSchema[1]}
          >
            {({
              values,
              errors,
              touched,
              isValid,
              setFieldValue,
              validateField,
              handleReset,
            }) => {
              // New function to update formValues
              const updateFormValues = (field: string, value: any) => {
                setFormValues((prev) => ({ ...prev, [field]: value }));
                setFieldValue(field, value);
              };
              return (
                <Form>
                  <Stack spacing={2}>
                    {address && (
                      <SiteMapComponent
                        posix={[values.latitude, values.longitude]}
                        address={address}
                      />
                    )}
                    <Typography variant="body2" mt={1}>
                      {Boolean(errors.latitude) || Boolean(errors.longitude) ? (
                        <span style={{ color: colors.redMatt }}>
                          Please enter valid coordinates
                        </span>
                      ) : isLoading ? (
                        'Loading address...'
                      ) : error ? (
                        <span style={{ color: colors.redMatt }}>
                          Error fetching address. Please try again.
                        </span>
                      ) : (
                        address
                      )}
                    </Typography>
                    <Field
                      as={TextField}
                      fullWidth
                      required
                      margin="normal"
                      name="siteName"
                      label="Site name"
                      placeholder="site-name"
                      error={touched.siteName && Boolean(errors.siteName)}
                      helperText={touched.siteName && errors.siteName}
                      InputLabelProps={{
                        shrink: true,
                      }}
                      InputProps={{
                        classes: {
                          input: gclasses.inputFieldStyle,
                        },
                      }}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                        updateFormValues('siteName', e.target.value);
                      }}
                    />
                    <Field
                      as={TextField}
                      fullWidth
                      select
                      required
                      name="network"
                      label="Network"
                      margin="normal"
                      InputLabelProps={{
                        shrink: true,
                      }}
                      InputProps={{
                        classes: {
                          input: gclasses.inputFieldStyle,
                        },
                      }}
                      error={touched.network && Boolean(errors.network)}
                      helperText={touched.network && errors.network}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                        updateFormValues('network', e.target.value);
                      }}
                    >
                      <MenuItem value="" disabled>
                        Choose a network to add your site to
                      </MenuItem>
                      {networks.map((network) => (
                        <MenuItem key={network.id} value={network.id}>
                          {network.name}
                        </MenuItem>
                      ))}
                    </Field>
                    <TextField
                      required
                      fullWidth
                      type="number"
                      label="Latitude"
                      name="latitude"
                      value={values.latitude}
                      onBlur={() => {
                        validateField('latitude');
                        if (
                          !errors.latitude &&
                          !errors.longitude &&
                          isValidLatLng([values.latitude, values.longitude])
                        ) {
                          handleFetchAddress(values.latitude, values.longitude);
                          updateFormValues('address', address);
                        }
                      }}
                      InputLabelProps={{
                        shrink: true,
                      }}
                      InputProps={{
                        classes: {
                          input: gclasses.inputFieldStyle,
                        },
                      }}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                        const value = parseFloat(e.target.value);
                        updateFormValues('latitude', value);
                      }}
                      error={touched.latitude && Boolean(errors.latitude)}
                      helperText={touched.latitude && errors.latitude}
                    />
                    <TextField
                      required
                      fullWidth
                      type="number"
                      label="Longitude"
                      name="longitude"
                      value={values.longitude}
                      onBlur={() => {
                        validateField('longitude');
                        if (
                          !errors.latitude &&
                          !errors.longitude &&
                          isValidLatLng([values.latitude, values.longitude])
                        ) {
                          handleFetchAddress(values.latitude, values.longitude);
                          updateFormValues('address', address);
                        }
                      }}
                      InputLabelProps={{
                        shrink: true,
                      }}
                      InputProps={{
                        classes: {
                          input: gclasses.inputFieldStyle,
                        },
                      }}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                        const value = parseFloat(e.target.value);
                        updateFormValues('longitude', value);
                      }}
                      error={touched.longitude && Boolean(errors.longitude)}
                      helperText={touched.longitude && errors.longitude}
                    />
                  </Stack>
                  <Stack
                    direction="row"
                    spacing={1}
                    justifyItems={'center'}
                    justifyContent={'flex-end'}
                    sx={{ mt: 2 }}
                  >
                    <Button onClick={handleClose}>Cancel</Button>
                    <Button onClick={handleBack}>Back</Button>
                    <Button
                      type="submit"
                      variant="contained"
                      color="primary"
                      disabled={
                        !isValid || addSiteLoading || isLoading || !address
                      }
                    >
                      Submit
                    </Button>
                  </Stack>
                </Form>
              );
            }}
          </Formik>
        )}
      </DialogContent>
    </Dialog>
  );
};

export default ConfigureSiteDialog;
