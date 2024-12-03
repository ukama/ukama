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
  useGetNetworksQuery,
  useGetSitesLazyQuery,
} from '@/client/graphql/generated';
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
import { useState } from 'react';
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
  const [loading, setLoading] = useState(false);
  const { setSnackbarMessage } = useAppContext();
  const flow = searchParams.get('flow') ?? 'onb';
  const totalStep = flow === 'onb' ? 5 : 4;
  const step = parseInt(searchParams.get('step') ?? '1');

  const [getSites, { data: sites }] = useGetSitesLazyQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (data.getSites.sites.length > 0) {
        // TODO: CHECK IF ANY SITE IS AVAILABLE FOR CONFIGURE & REDIRECT TO SITE CONFIGURE
        router.push(`/console/home`);
      } else {
        router.push(`/configure?step=2&flow=${flow}`);
      }
    },
  });

  const { data: networksData, loading: networksLoading } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (data.getNetworks.networks.length >= 1) {
        getSites({
          variables: {
            networkId: data.getNetworks.networks[0].id,
          },
        });
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
        getSites({
          variables: {
            networkId: data.addNetwork.id,
          },
        });

        setSnackbarMessage({
          id: 'add-networks-success',
          message: 'Network added successfully',
          type: 'success',
          show: true,
        });
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
              {flow === 'onb' && <i>&nbsp;</i>}&nbsp;
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
