import { user } from '@/app-recoil';
import { NETWORK_NAME_SCHEMA_VALIDATOR } from '@/helpers/formValidators';
import { globalUseStyles } from '@/styles/global';
import {
  Box,
  Button,
  Grid,
  Paper,
  Radio,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { Formik } from 'formik';
import React, { useEffect, useState } from 'react';
import { useRecoilValue } from 'recoil';
import * as Yup from 'yup';
const eSimFormSchema = Yup.object(NETWORK_NAME_SCHEMA_VALIDATOR);

interface INetworkTypes {
  nextStep: Function;
  networkData: any;
}

const NetworkSetup = ({ nextStep, networkData }: INetworkTypes) => {
  const [networkType, setNetworkType] = useState('personal');
  const gclasses = globalUseStyles();
  // const setNetworkNames = useSetRecoilState();
  // const getNetworkName = useRecoilValue();
  const getUser: any = useRecoilValue(user);
  useEffect(() => {
    // setNetworkNames('');
  }, []);
  const handleSimTypeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setNetworkType(event.target.value);
  };
  const handleNetworksetup = (value: any) => {
    networkData(value);
    nextStep();
    // setNetworkNames(value.name);
  };
  const backToSignUp = () => {
    typeof window !== 'undefined' &&
      window.location.replace(
        `${process.env.NEXT_PUBLIC_REACT_AUTH_APP_URL}/logout?goTo=signUp&name=${getUser?.name}&email=${getUser?.email}`,
      );
  };

  return (
    <Box sx={{ pb: 2 }}>
      <Formik
        initialValues={{ name:  '' }}
        validationSchema={eSimFormSchema}
        onSubmit={async (values) =>
          handleNetworksetup({ ...values, networkType })
        }
      >
        {({
          values,
          errors,
          touched,
          handleChange,
          handleSubmit,
          handleBlur,
        }) => (
          <form onSubmit={handleSubmit}>
            <Stack direction="column" spacing={3} sx={{ mb: 2 }}>
              <Typography variant="h6">
                What kind of network are you setting up?
              </Typography>
              <Typography variant="body2">
                Get a customized Console for your specialized needs, depending
                on what type of network you’re setting up.
              </Typography>
            </Stack>
            <Grid container spacing={1}>
              <Grid item xs={6}>
                <Paper variant="outlined" sx={{}}>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <Radio
                      checked={networkType === 'personal'}
                      onChange={handleSimTypeChange}
                      value="personal"
                      name="personal"
                      inputProps={{
                        'aria-label': 'personal',
                      }}
                    />
                    <Typography variant="body1">Personal network</Typography>
                  </Stack>
                </Paper>
              </Grid>
              <Grid item xs={6}>
                <Paper variant="outlined" sx={{}}>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <Radio
                      checked={networkType === 'community'}
                      onChange={handleSimTypeChange}
                      value="community"
                      name="community"
                      inputProps={{
                        'aria-label': 'community',
                      }}
                    />
                    <Typography variant="body1">Community network</Typography>
                  </Stack>
                </Paper>
              </Grid>
              <Grid item xs={12} sx={{ mt: 2, mb: 2 }}>
                <TextField
                  fullWidth
                  id="name"
                  name="name"
                  label="NETWORK NAME"
                  onBlur={handleBlur}
                  onChange={handleChange}
                  value={values.name}
                  sx={{ mb: 1 / 2 }}
                  InputLabelProps={{ shrink: true }}
                  InputProps={{
                    classes: {
                      input: gclasses.inputFieldStyle,
                    },
                  }}
                  // helperText={(touched.name && errors.name) || <></>}
                  error={touched.name && Boolean(errors.name)}
                />
              </Grid>
            </Grid>

            <Stack direction="row" justifyContent="space-between">
              <Button variant="text" onClick={backToSignUp}>
                BACK TO SIGN UP
              </Button>

              <Button variant="contained" type="submit">
                NEXT
              </Button>
            </Stack>
          </form>
        )}
      </Formik>
    </Box>
  );
};
export default NetworkSetup;
