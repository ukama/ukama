import React, { useEffect, useState } from 'react';
import { Box, Typography, TextField, MenuItem, Stack } from '@mui/material';
import { Field, ErrorMessage, FormikProps } from 'formik';
import colors from '@/theme/colors';
import CustomTextField from '@/components/CustomTextField';
import dynamic from 'next/dynamic';
import { globalUseStyles } from '@/styles/global';

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

export interface FormValues {
  switch: string;
  power: string;
  backhaul: string;
  access: string;
  siteName: string;
  network: string;
  latitude: number;
  longitude: number;
  location: string;
}

interface StepContentProps {
  step: number;
  handleAddressChange: (address: string) => void;
  components: Component[];
  formik: FormikProps<FormValues>;
  onSiteInfoChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onNameChange: (name: any) => void;
}

const SiteStepForm: React.FC<StepContentProps> = ({
  step,
  handleAddressChange,
  onSiteInfoChange,
  onNameChange,
  components,
  formik,
}) => {
  const getComponentsByType = (type: string) => {
    return components.filter(
      (component) => component.type.toLowerCase() === type.toLowerCase(),
    );
  };

  const gclasses = globalUseStyles();

  switch (step) {
    case 0:
      return (
        <Box style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          <Box sx={{ mt: 2, mb: 2 }}>
            <Typography>
              {`You have successfully installed your site, and need to configure
              it. Please note that if your power or backhaul choice is "other",
              it can't be monitored within Ukama's Console.`}
            </Typography>
          </Box>
          <TextField
            fullWidth
            label={'SWITCH'}
            select
            required
            name="switch"
            InputLabelProps={{
              shrink: true,
            }}
            onBlur={formik.handleBlur}
            onChange={formik.handleChange}
            value={formik.values.switch}
            helperText={formik.touched.switch && formik.errors.switch}
            error={formik.touched.switch && Boolean(formik.errors.switch)}
            id={'switch'}
            spellCheck={false}
            InputProps={{
              classes: {
                input: gclasses.inputFieldStyle,
              },
            }}
          >
            {getComponentsByType('switch').map((component) => (
              <MenuItem key={component.id} value={component.description}>
                {component.description}
              </MenuItem>
            ))}
          </TextField>

          <TextField
            fullWidth
            label={'POWER'}
            select
            required
            name="power"
            InputLabelProps={{
              shrink: true,
            }}
            onBlur={formik.handleBlur}
            onChange={formik.handleChange}
            value={formik.values.power}
            helperText={formik.touched.power && formik.errors.power}
            error={formik.touched.power && Boolean(formik.errors.power)}
            id={'power'}
            spellCheck={false}
            InputProps={{
              classes: {
                input: gclasses.inputFieldStyle,
              },
            }}
          >
            {getComponentsByType('power').map((component) => (
              <MenuItem key={component.id} value={component.description}>
                {component.description}
              </MenuItem>
            ))}
          </TextField>

          <TextField
            fullWidth
            label={'BACKHAUL'}
            select
            required
            name="backhaul"
            InputLabelProps={{
              shrink: true,
            }}
            onBlur={formik.handleBlur}
            onChange={formik.handleChange}
            value={formik.values.backhaul}
            helperText={formik.touched.backhaul && formik.errors.backhaul}
            error={formik.touched.backhaul && Boolean(formik.errors.backhaul)}
            id={'backhaul'}
            spellCheck={false}
            InputProps={{
              classes: {
                input: gclasses.inputFieldStyle,
              },
            }}
          >
            {getComponentsByType('backhaul').map((component) => (
              <MenuItem key={component.id} value={component.description}>
                {component.description}
              </MenuItem>
            ))}
          </TextField>

          <TextField
            fullWidth
            label={'SPECTRUM BAND'}
            select
            required
            name="access"
            InputLabelProps={{
              shrink: true,
            }}
            onBlur={formik.handleBlur}
            onChange={formik.handleChange}
            value={formik.values.access}
            helperText={formik.touched.access && formik.errors.access}
            error={formik.touched.access && Boolean(formik.errors.access)}
            id={'access'}
            spellCheck={false}
            InputProps={{
              classes: {
                input: gclasses.inputFieldStyle,
              },
            }}
          >
            {getComponentsByType('access').map((component) => (
              <MenuItem key={component.id} value={component.description}>
                {component.description}
              </MenuItem>
            ))}
          </TextField>
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
            posix={[formik.values.latitude, formik.values.longitude]}
            onAddressChange={(address) => {
              formik.setFieldValue('location', address);
              handleAddressChange(address);
            }}
          />
          <Box>
            <Stack direction="column" spacing={1} justifyItems={'center'}>
              <Typography variant="body2" sx={{ color: `${colors.darkGray}` }}>
                LOCATION
              </Typography>
              <Typography variant="body2" color="initial">
                {formik.values.location || 'Fetching site location...'}
              </Typography>
            </Stack>
          </Box>
          <Field
            as={TextField}
            fullWidth
            label="Location"
            name="location"
            InputLabelProps={{ shrink: true }}
            onChange={formik.handleChange}
            value={formik.values.location}
            id="location"
            spellCheck={false}
            InputProps={{ classes: { input: gclasses.inputFieldStyle } }}
          />
          <Field
            as={TextField}
            fullWidth
            label="Latitude"
            name="latitude"
            InputLabelProps={{ shrink: true }}
            onChange={formik.handleChange}
            value={formik.values.latitude}
            id="latitude"
            spellCheck={false}
            InputProps={{ classes: { input: gclasses.inputFieldStyle } }}
          />
          <Field
            as={TextField}
            fullWidth
            label="Longitude"
            name="longitude"
            InputLabelProps={{ shrink: true }}
            onChange={formik.handleChange}
            value={formik.values.longitude}
            id="longitude"
            spellCheck={false}
            InputProps={{ classes: { input: gclasses.inputFieldStyle } }}
          />
          <Field
            as={TextField}
            fullWidth
            label="Site Name"
            name="siteName"
            InputLabelProps={{ shrink: true }}
            onChange={formik.handleChange}
            value={formik.values.siteName}
            id="siteName"
            spellCheck={false}
            InputProps={{ classes: { input: gclasses.inputFieldStyle } }}
          />
        </Box>
      );
    default:
      return <div>Not Found</div>;
  }
};

export default SiteStepForm;
