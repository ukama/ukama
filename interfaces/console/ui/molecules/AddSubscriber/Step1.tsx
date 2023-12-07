import { globalUseStyles } from '@/styles/global';
import colors from '@/styles/theme/colors';
import {
  Button,
  Grid,
  Stack,
  Switch,
  TextField,
  Typography,
} from '@mui/material';
import { Form, Formik } from 'formik';
import React, { useState } from 'react';
import * as Yup from 'yup';
import SimTypeRadio from './simTypeComp';

const validationSchema = Yup.object({
  name: Yup.string().required('Name is required'),
  email: Yup.string().email('Invalid email').required('Email is required'),
  phone: Yup.string().optional(),
});

interface SubscriberDialogProps {
  onClose: () => void;
  handleSimInstallation: Function;
  eSimCount: number | undefined;
  pSimCount: number | undefined;
}

interface SubscriberFormValues {
  name: string;
  email: string;
  phone: string;
}

const Step1: React.FC<SubscriberDialogProps> = ({
  onClose,
  handleSimInstallation,
  eSimCount,
  pSimCount,
}) => {
  const initialValues: SubscriberFormValues = {
    name: '',
    email: '',
    phone: '',
  };
  const gclasses = globalUseStyles();

  const [selectedSimType, setSelectedSimType] = useState<string>('eSim');
  const [roamingStatus, setRoamingStatus] = useState<boolean>(false);
  const handleSimTypeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedSimType(event.target.value);
  };

  return (
    <>
      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema}
        onSubmit={async (values) => {
          handleSimInstallation({ ...values, selectedSimType, roamingStatus });
        }}
      >
        {({
          values,
          errors,
          touched,
          handleChange,
          handleSubmit,
          handleBlur,
        }) => (
          <Form onSubmit={handleSubmit}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <TextField
                  required
                  fullWidth
                  label={'NAME'}
                  name={'name'}
                  InputLabelProps={{
                    shrink: true,
                  }}
                  onChange={handleChange}
                  value={values.name}
                  onBlur={handleBlur}
                  helperText={touched.name && errors.name}
                  error={touched.name && Boolean(errors.name)}
                  id={'name'}
                  spellCheck={false}
                  InputProps={{
                    classes: {
                      input: gclasses.inputFieldStyle,
                    },
                  }}
                />
              </Grid>
              <Grid item xs={6}>
                <SimTypeRadio
                  simType="eSim"
                  label="eSIM"
                  count={eSimCount}
                  selectedSimType={selectedSimType}
                  handleSimTypeChange={handleSimTypeChange}
                />
              </Grid>
              <Grid item xs={6}>
                <SimTypeRadio
                  simType="pSim"
                  label="pSIM"
                  count={pSimCount}
                  selectedSimType={selectedSimType}
                  handleSimTypeChange={handleSimTypeChange}
                />
              </Grid>

              <Grid item xs={12}>
                <TextField
                  required
                  fullWidth
                  label={'EMAIL'}
                  name={'email'}
                  InputLabelProps={{
                    shrink: true,
                  }}
                  onBlur={handleBlur}
                  onChange={handleChange}
                  value={values.email}
                  helperText={touched.email && errors.email}
                  error={touched.email && Boolean(errors.email)}
                  id={'email'}
                  spellCheck={false}
                  InputProps={{
                    classes: {
                      input: gclasses.inputFieldStyle,
                    },
                  }}
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label={'PHONE (OPTIONAL)'}
                  name={'phone'}
                  InputLabelProps={{
                    shrink: true,
                  }}
                  onBlur={handleBlur}
                  onChange={handleChange}
                  value={values.phone}
                  helperText={touched.phone && errors.phone}
                  error={touched.phone && Boolean(errors.phone)}
                  id={'phone'}
                  spellCheck={false}
                  InputProps={{
                    classes: {
                      input: gclasses.inputFieldStyle,
                    },
                  }}
                />
              </Grid>
              <Grid item xs={12}>
                <Stack direction="column" spacing={1}>
                  <Typography variant="caption" sx={{ color: colors.black38 }}>
                    ALLOW ROAMING
                  </Typography>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <Typography variant="body1">
                      Roaming allows subscriber to use data outside of their
                      home Ukama network. Insert billing information.
                    </Typography>
                    <Switch
                      size="small"
                      value={roamingStatus}
                      checked={roamingStatus}
                      onChange={(e) => setRoamingStatus(e.target.checked)}
                    />
                  </Stack>
                </Stack>

                <Stack
                  direction="row"
                  justifyContent={'flex-end'}
                  mt={1}
                  spacing={1}
                  sx={{ mb: 2 }}
                >
                  <Stack direction="row" spacing={3}>
                    <Button
                      variant="text"
                      onClick={() => {
                        onClose();
                      }}
                    >
                      {' CANCEL'}
                    </Button>

                    <Button
                      variant="contained"
                      type="submit"
                      sx={{
                        disabled: {
                          color: colors.primaryLight,
                        },
                      }}
                    >
                      <Typography variant="body1">NEXT</Typography>
                    </Button>
                  </Stack>
                </Stack>
              </Grid>
            </Grid>
          </Form>
        )}
      </Formik>
    </>
  );
};

export default Step1;
