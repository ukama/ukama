/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import {
  useAddNetworkMutation,
  useGetNetworksLazyQuery,
} from '@/client/graphql/generated';
import { CHECK_SITE_FLOW, NETWORK_FLOW, ONBOARDING_FLOW } from '@/constants';
import { useAppContext } from '@/context';
import { CenterContainer } from '@/styles/global';
import colors from '@/theme/colors';
import {
  Button,
  CircularProgress,
  FormControlLabel,
  Paper,
  Stack,
  Switch,
  TextField,
  Typography,
} from '@mui/material';
import { Formik } from 'formik';
import { useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import * as Yup from 'yup';
import NetworkInfo from '../../../../public/svg/NetworkInfo';

interface AddNetworkForm {
  name: string;
  isDefault: boolean;
  budget: number;
  countries: { name: string; code: string }[];
  networks: { id: string; name: string; isDefault: boolean }[];
}

const validationSchema = Yup.object({
  networks: Yup.array().optional().default([]),
  isDefault: Yup.boolean().default(false),
  countries: Yup.array().optional().default([]),
  name: Yup.string()
    .required('Network name is required')
    .matches(
      /^[a-z0-9-]*$/,
      'Network name must be lowercase alphanumeric and should not contain spaces, "-" are allowed.',
    ),
  budget: Yup.number().default(0),
});

const initialValues: AddNetworkForm = {
  name: '',
  budget: 0,
  isDefault: true,
  countries: [],
  networks: [],
};

const Network = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [loading, setLoading] = useState(true);
  const flow = searchParams.get('flow') ?? ONBOARDING_FLOW;
  const totalStep = flow === ONBOARDING_FLOW ? 5 : 4;
  const step = parseInt(searchParams.get('step') ?? '1');
  const { setSnackbarMessage, network, setNetwork } = useAppContext();

  const [getNetworks] = useGetNetworksLazyQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (data.getNetworks.networks.length > 0) {
        setNetwork({
          id: data.getNetworks.networks[0].id,
          name: data.getNetworks.networks[0].name,
        });
        router.push(`/configure/check?flow=${CHECK_SITE_FLOW}`);
      } else {
        setLoading(false);
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'networks-msg',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [addNetwork, { loading: addNetworkLoading }] = useAddNetworkMutation({
    onCompleted: (data) => {
      if (data.addNetwork.id) {
        setNetwork({
          id: data.addNetwork.id,
          name: data.addNetwork.name,
        });
        router.push(`/configure/check?flow=${NETWORK_FLOW}`);
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'add-networks-error',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  useEffect(() => {
    if (network.id) {
      router.push(`/configure/check?flow=${CHECK_SITE_FLOW}`);
    } else {
      getNetworks();
    }
  }, []);

  const handleAddNetwork = (values: any) => {
    addNetwork({
      variables: {
        data: {
          isDefault: true,
          name: values.name,
          budget: values.budget,
          networks: values.networks,
          countries: values.countries,
        },
      },
    });
  };

  return (
    <Paper elevation={0} sx={{ px: { xs: 2, md: 4 }, py: { xs: 1, md: 2 } }}>
      {loading ? (
        <CenterContainer>
          <CircularProgress />
        </CenterContainer>
      ) : (
        <Stack direction="column" spacing={1.5}>
          <Stack direction={'row'}>
            <Typography variant="h6"> {'Name your network'}</Typography>
            <Typography
              variant="h6"
              fontWeight={400}
              sx={{
                color: colors.black70,
              }}
            >
              {flow === ONBOARDING_FLOW && <i>&nbsp;</i>}&nbsp;
              {`(${step}/${totalStep})`}
            </Typography>
          </Stack>
          <Formik
            initialValues={initialValues}
            validationSchema={validationSchema}
            onSubmit={async (values) => {
              handleAddNetwork(values);
            }}
          >
            {({
              values,
              errors,
              touched,
              handleChange,
              handleSubmit,
              handleBlur,
              setFieldValue,
            }) => (
              <form onSubmit={handleSubmit}>
                <Stack spacing={3} direction={'column'} alignItems="center">
                  <Typography variant="body1">
                    Please name your first network. A network is made up of one
                    or more sites of Ukama hardware, allowing you to connect to
                    the cellular internet.
                  </Typography>
                  {NetworkInfo}
                  <TextField
                    fullWidth
                    name={'name'}
                    size="medium"
                    placeholder="network-name"
                    label={'Network name'}
                    InputLabelProps={{
                      shrink: true,
                    }}
                    onBlur={handleBlur}
                    onChange={handleChange}
                    value={values.name}
                    helperText={touched.name && errors.name}
                    error={touched.name && Boolean(errors.name)}
                    id={'name'}
                  />
                  <FormControlLabel
                    sx={{ display: 'none' }}
                    control={
                      <Switch
                        defaultChecked={false}
                        value={values.isDefault}
                        checked={values.isDefault}
                        onChange={() =>
                          setFieldValue('isDefault', !values.isDefault)
                        }
                      />
                    }
                    label="Make this network default"
                  />
                  <Stack
                    width="100%"
                    spacing={2}
                    direction={'row'}
                    alignItems="center"
                    justifyContent={'flex-end'}
                  >
                    <Button
                      type="submit"
                      variant="contained"
                      disabled={addNetworkLoading}
                    >
                      CREATE NETWORK
                    </Button>
                  </Stack>
                </Stack>
              </form>
            )}
          </Formik>
        </Stack>
      )}
    </Paper>
  );
};

export default Network;
