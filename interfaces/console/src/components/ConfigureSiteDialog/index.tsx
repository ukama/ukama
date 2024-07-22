import React, { ChangeEvent, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stepper,
  Step,
  StepLabel,
  TextField,
  MenuItem,
  DialogContentText,
  IconButton,
  Stack,
} from '@mui/material';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';
import { globalUseStyles } from '@/styles/global';
import dynamic from 'next/dynamic';
import CloseIcon from '@mui/icons-material/Close';
import { NetworkDto } from '@/client/graphql/generated';
const SiteMapComponent = dynamic(() => import('../SiteMapComponent'), {
  loading: () => <p>Site map is loading</p>,
  ssr: false,
});
interface FormValues {
  switch: string;
  power: string;
  backhaul: string;
  access: string;
  siteName: string;
  selectedNetwork: string;
}

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

// Validation schema for both steps
const validationSchema = [
  Yup.object().shape({
    switch: Yup.string().required('Switch is required'),
    power: Yup.string().required('Power is required'),
    backhaul: Yup.string().required('Backhaul is required'),
    access: Yup.string().required('Access is required'),
  }),
  Yup.object().shape({
    siteName: Yup.string().required('Site Name is required'),
    selectedNetwork: Yup.string().required('Network is required'),
  }),
];

interface StepperDialogProps {
  open: boolean;
  onClose: () => void;
  components: Component[];
  networks: NetworkDto[];
  handleSiteConfiguration: (data: any) => void;
}

interface Coordinates {
  lat: number | null;
  lng: number | null;
}

const ConfigureSiteDialog: React.FC<StepperDialogProps> = ({
  open,
  onClose,
  components,
  networks,
  handleSiteConfiguration,
}) => {
  const [activeStep, setActiveStep] = useState(0);
  const [location, setLocation] = useState('');
  const gclasses = globalUseStyles();
  const handleNext = () => setActiveStep((prevStep) => prevStep + 1);
  const handleBack = () => setActiveStep((prevStep) => prevStep - 1);
  const [coordinates, setCoordinates] = useState<Coordinates>({
    lat: null,
    lng: null,
  });

  const initialValues: FormValues = {
    switch: '',
    power: '',
    backhaul: '',
    access: '',
    siteName: '',
    selectedNetwork: '',
  };

  const handleSubmit = (values: FormValues) => {
    handleSiteConfiguration({ ...values, coordinates, location });
    // onClose();
  };

  const handleCoordnatedChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    const numValue = value === '' ? null : parseFloat(value);
    setCoordinates((prev) => ({ ...prev, [name]: numValue }));
  };

  const steps = [
    'Select Switch, Power, Backhaul, and spectrum band',
    'Enter your site details',
  ];

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

  return (
    <Dialog
      open={open}
      onClose={onClose}
      sx={{
        '& .MuiDialog-paper': {
          width: '60%',
          maxWidth: '40%',
        },
      }}
    >
      <DialogTitle>
        Configure site installation ({activeStep + 1}/2)
        <IconButton
          aria-label="close"
          onClick={onClose}
          sx={{ position: 'absolute', right: 8, top: 8 }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          {`You have successfully installed your site, and need to configure
              it. Please note that if your power or backhaul choice is "other",
              it can't be monitored within Ukama's Console.`}
        </DialogContentText>
      </DialogContent>

      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema[activeStep]}
        onSubmit={handleSubmit}
      >
        {({ values, errors, touched, isValid, handleChange }) => (
          <Form>
            <DialogContent>
              <Stepper activeStep={activeStep} sx={{ mb: 4 }}>
                {steps.map((label) => (
                  <Step key={label}>
                    <StepLabel>{label}</StepLabel>
                  </Step>
                ))}
              </Stepper>
              {activeStep === 0 && (
                <>
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
                      <MenuItem
                        key={component.id}
                        value={component.description}
                      >
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
                      <MenuItem
                        key={component.id}
                        value={component.description}
                      >
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
                      <MenuItem
                        key={component.id}
                        value={component.description}
                      >
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
                    error={touched.access && Boolean(errors.access)}
                    helperText={touched.access && errors.access}
                  >
                    {accessComponents.map((component) => (
                      <MenuItem
                        key={component.id}
                        value={component.description}
                      >
                        {component.description}
                      </MenuItem>
                    ))}
                  </Field>
                </>
              )}
              {activeStep === 1 && (
                <>
                  {coordinates.lat !== null && coordinates.lng !== null && (
                    <SiteMapComponent
                      posix={[coordinates.lat, coordinates.lng]}
                      onAddressChange={(address: string) => {
                        setLocation(address);
                      }}
                    />
                  )}
                  {location}
                  <Field
                    as={TextField}
                    fullWidth
                    margin="normal"
                    name="siteName"
                    label="Site Name"
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
                  />
                  <Field
                    as={TextField}
                    fullWidth
                    select
                    required
                    name="selectedNetwork"
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
                    error={
                      touched.selectedNetwork && Boolean(errors.selectedNetwork)
                    }
                    helperText={
                      touched.selectedNetwork && errors.selectedNetwork
                    }
                  >
                    <MenuItem value="" disabled>
                      Choose a network to add your site to
                    </MenuItem>
                    {networks.map((network) => (
                      <MenuItem key={network.id} value={network.name}>
                        {network.name}
                      </MenuItem>
                    ))}
                  </Field>
                  <Stack direction="column" spacing={2} sx={{ mt: 2 }}>
                    <TextField
                      label="Latitude"
                      name="lat"
                      required
                      value={
                        coordinates.lat === null
                          ? ''
                          : coordinates.lat.toString()
                      }
                      onChange={handleCoordnatedChange}
                      fullWidth
                      type="number"
                      InputLabelProps={{
                        shrink: true,
                      }}
                      InputProps={{
                        classes: {
                          input: gclasses.inputFieldStyle,
                        },
                      }}
                    />
                    <TextField
                      label="Longitude"
                      name="lng"
                      required
                      value={
                        coordinates.lng === null
                          ? ''
                          : coordinates.lng.toString()
                      }
                      onChange={handleCoordnatedChange}
                      fullWidth
                      type="number"
                      InputLabelProps={{
                        shrink: true,
                      }}
                      InputProps={{
                        classes: {
                          input: gclasses.inputFieldStyle,
                        },
                      }}
                    />
                  </Stack>
                </>
              )}
            </DialogContent>
            <DialogActions>
              <Button onClick={onClose}>Cancel</Button>
              {activeStep > 0 && <Button onClick={handleBack}>Back</Button>}
              {activeStep < steps.length - 1 ? (
                <Button onClick={handleNext} disabled={!isValid}>
                  Next
                </Button>
              ) : (
                <Button type="submit" disabled={!isValid}>
                  Submit
                </Button>
              )}
            </DialogActions>
          </Form>
        )}
      </Formik>
    </Dialog>
  );
};

export default ConfigureSiteDialog;
