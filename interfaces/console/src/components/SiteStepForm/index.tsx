import React, { useEffect } from 'react';
import {
  Box,
  Typography,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Stack,
} from '@mui/material';
import { Field, ErrorMessage, FormikProps } from 'formik';
import colors from '@/theme/colors';
import CustomTextField from '@/components/CustomTextField';
import dynamic from 'next/dynamic';

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

interface FormValues {
  switch: string;
  power: string;
  backhaul: string;
  access: string;
  siteName: string;
  network: string;
  latitude: number;
  longitude: number;
}

interface StepContentProps {
  step: number;
  lat: number;
  lng: number;
  location: string;
  handleLatChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleLngChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleAddressChange: (address: string) => void;
  components: Component[];
  formik: FormikProps<FormValues>;
  getComponentInfos: (components: FormikProps<FormValues>) => void;
}

const SiteStepForm: React.FC<StepContentProps> = ({
  step,
  lat,
  lng,
  location,
  handleLatChange,
  handleLngChange,
  handleAddressChange,
  components,
  formik,
  getComponentInfos,
}) => {
  const getComponentsByType = (type: string) => {
    return components.filter(
      (component) => component.type.toLowerCase() === type.toLowerCase(),
    );
  };

  useEffect(() => {
    getComponentInfos(formik);
  }, [formik]);

  switch (step) {
    case 0:
      return (
        <Box
          component="form"
          style={{ display: 'flex', flexDirection: 'column', gap: 16 }}
        >
          <Box sx={{ mt: 2, mb: 2 }}>
            <Typography>
              {`You have successfully installed your site, and need to configure
              it. Please note that if your power or backhaul choice is "other",
              it can't be monitored within Ukama's Console.`}
            </Typography>
          </Box>
          <FormControl fullWidth>
            <InputLabel id="switch-label">Switch</InputLabel>
            <Field
              as={Select}
              labelId="switch-label"
              name="switch"
              label="Switch"
              value={formik.values.switch}
              onChange={formik.handleChange}
            >
              {getComponentsByType('switch').map((component) => (
                <MenuItem key={component.id} value={component.description}>
                  {component.description}
                </MenuItem>
              ))}
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
              value={formik.values.power}
              onChange={formik.handleChange}
            >
              {getComponentsByType('power').map((component) => (
                <MenuItem key={component.id} value={component.description}>
                  {component.description}
                </MenuItem>
              ))}
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
              value={formik.values.backhaul}
              onChange={formik.handleChange}
            >
              {getComponentsByType('backhaul').map((component) => (
                <MenuItem key={component.id} value={component.description}>
                  {component.description}
                </MenuItem>
              ))}
            </Field>
            <ErrorMessage name="backhaul" component="div">
              {(msg) => <div style={{ color: colors.red }}>{msg}</div>}
            </ErrorMessage>
          </FormControl>
          <FormControl fullWidth>
            <InputLabel id="access-label">Access</InputLabel>
            <Field
              as={Select}
              labelId="access-label"
              name="access"
              label="Access"
              value={formik.values.access}
              onChange={formik.handleChange}
            >
              {getComponentsByType('access').map((component) => (
                <MenuItem key={component.id} value={component.description}>
                  {component.description}
                </MenuItem>
              ))}
            </Field>
            <ErrorMessage name="access" component="div">
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
              Please name your site for your ease of reference, and assign it to
              a network.
            </Typography>
          </Box>
          <SiteMapComponent
            posix={[lat, lng]}
            onAddressChange={handleAddressChange}
          />
          <Box>
            <Stack direction="column" spacing={1} justifyItems={'center'}>
              <Typography variant="body2" sx={{ color: `${colors.darkGray}` }}>
                LOCATION
              </Typography>
              <Typography variant="body2" color="initial">
                {location || 'Fetching site location...'}
              </Typography>
            </Stack>
          </Box>
          <CustomTextField label="Location" name="location" />
          <CustomTextField label="Latitude" name="latitude" />
          <CustomTextField label="Longitude" name="longitude" />
          <CustomTextField label="Site Name" name="siteName" />
        </Box>
      );
    default:
      return <div>Not Found</div>;
  }
};

export default SiteStepForm;
