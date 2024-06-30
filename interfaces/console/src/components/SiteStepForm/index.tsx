import React from 'react';
import {
  Box,
  Typography,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Stack,
} from '@mui/material';
import { Field, ErrorMessage } from 'formik';
import { SITE_CONFIG_STEPS } from '@/constants';
import colors from '@/theme/colors';
import CustomTextField from '@/components/CustomTextField';
import dynamic from 'next/dynamic';

const SiteMapComponent = dynamic(() => import('../SiteMapComponent'), {
  loading: () => <p>Site map is loading</p>,
  ssr: false,
});

interface StepContentProps {
  step: number;
  lat: number;
  lng: number;
  location: string;
  handleLatChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleLngChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleAddressChange: (address: string) => void;
}

const SiteStepForm: React.FC<StepContentProps> = ({
  step,
  lat,
  lng,
  location,
  handleLatChange,
  handleLngChange,
  handleAddressChange,
}) => {
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
              it. Please note that if your power or backhaul choice is “other”,
              it can’t be monitored within Ukama’s Console.
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
            <Field as={Select} labelId="power-label" name="power" label="Power">
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
          <ErrorMessage name="longitude" component="div">
            {(msg) => <div style={{ color: colors.red }}>{msg}</div>}
          </ErrorMessage>
          <ErrorMessage name="latitude" component="div">
            {(msg) => <div style={{ color: colors.red }}>{msg}</div>}
          </ErrorMessage>
          <CustomTextField label="Site Name" name="siteName" />
          <ErrorMessage name="siteName" component="div">
            {(msg) => <div style={{ color: colors.red }}>{msg}</div>}
          </ErrorMessage>
          <CustomTextField label="Network" name="network" />
          <ErrorMessage name="network" component="div">
            {(msg) => <div style={{ color: colors.red }}>{msg}</div>}
          </ErrorMessage>
        </Box>
      );
    default:
      return <div>Not Found</div>;
  }
};

export default SiteStepForm;
