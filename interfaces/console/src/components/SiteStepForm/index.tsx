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
import { Formik, Form, Field, ErrorMessage, FormikProps } from 'formik';
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

interface StepContentProps {
  step: number;
  lat: number;
  lng: number;
  location: string;
  handleLatChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleLngChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleAddressChange: (address: string) => void;
  components: Component[];
}

interface FormValues {
  [key: string]: string;
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
}) => {
  const initialValues: FormValues = components.reduce(
    (acc: FormValues, component) => {
      acc[component.type.toLowerCase()] = '';
      return acc;
    },
    {},
  );

  const renderComponentField = (
    component: Component,
    formikProps: FormikProps<FormValues>,
  ) => {
    const fieldName = component.type.toLowerCase();
    return (
      <FormControl fullWidth key={component.id}>
        <InputLabel id={`${fieldName}-label`}>{component.category}</InputLabel>
        <Field
          as={Select}
          labelId={`${fieldName}-label`}
          name={fieldName}
          label={component.category}
          value={formikProps.values[fieldName] || ''}
          onChange={(e: React.ChangeEvent<{ value: unknown }>) =>
            formikProps.setFieldValue(fieldName, e.target.value as string)
          }
        >
          <MenuItem value={component.id}>{component.description}</MenuItem>
        </Field>
        <ErrorMessage name={fieldName} component="div">
          {(msg) => <div style={{ color: colors.red }}>{msg}</div>}
        </ErrorMessage>
      </FormControl>
    );
  };

  const renderStep0 = (formikProps: FormikProps<FormValues>) => (
    <>
      <Box sx={{ mt: 2, mb: 2 }}>
        <Typography>
          {`You have successfully installed your site, and need to configure it.
          Please note that if your power or backhaul choice is "other", it can't
          be monitored within Ukama's Console.`}
        </Typography>
      </Box>
      {components.map((component) =>
        renderComponentField(component, formikProps),
      )}
    </>
  );

  const renderStep1 = () => (
    <>
      <Box sx={{ mt: 2, mb: 2 }}>
        <Typography>
          Please name your site for your ease of reference, and assign it to a
          network.
        </Typography>
      </Box>
      <SiteMapComponent
        posix={[lat, lng]}
        onAddressChange={handleAddressChange}
      />
      <Box>
        <Stack direction="column" spacing={1} justifyItems={'center'}>
          <Typography variant="body2" sx={{ color: colors.darkGray }}>
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
    </>
  );

  return (
    <Formik<FormValues>
      initialValues={initialValues}
      onSubmit={(values: FormValues) => {
        console.log(values);
      }}
    >
      {(formikProps: FormikProps<FormValues>) => (
        <Form style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          {step === 0 ? renderStep0(formikProps) : renderStep1()}
        </Form>
      )}
    </Formik>
  );
};

export default SiteStepForm;
