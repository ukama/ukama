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
} from '@/client/graphql/generated';
import { CHECK_SITE_FLOW, NETWORK_FLOW } from '@/constants';
import { useAppContext } from '@/context';
import { setQueryParam } from '@/utils';
import { Button, Stack, TextField, Typography } from '@mui/material';
import { Formik } from 'formik';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useState } from 'react';
import * as Yup from 'yup';
import NetworkSkelton from './skelton';

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
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [loading, setLoading] = useState(true);
  const { setSnackbarMessage, setNetwork } = useAppContext();

  useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (data.getNetworks.networks.length > 0) {
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

  const [addNetwork] = useAddNetworkMutation({
    onCompleted: (data) => {
      if (data.addNetwork.id) {
        setNetwork({
          id: data.addNetwork.id,
          name: data.addNetwork.name,
        });
        const p = setQueryParam(
          'networkid',
          data.addNetwork.id,
          searchParams.toString(),
          pathname,
        );
        p.set('flow', NETWORK_FLOW);
        router.push(`/configure/check?${p.toString()}`);
      }
    },
    onError: (error) => {
      setLoading(false);
      setSnackbarMessage({
        id: 'add-networks-error',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  const handleAddNetwork = (values: any) => {
    setLoading(true);
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

  if (loading) return <NetworkSkelton />;

  return (
    <Stack direction="column" height={'100%'} spacing={3}>
      <Typography variant="h4" fontWeight={500}>
        Name your network
      </Typography>
      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema}
        onSubmit={(values) => {
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
        }) => (
          <form onSubmit={handleSubmit}>
            <Stack
              spacing={{ xs: 4, md: 6 }}
              direction={'column'}
              height={'100%'}
              alignItems="center"
            >
              <Stack spacing={4}>
                <Typography variant="body1">
                  Please name your first network. A network is made up of one or
                  more sites of Ukama hardware, allowing you to connect to the
                  cellular internet.
                </Typography>
                <TextField
                  fullWidth
                  id={'name'}
                  name={'name'}
                  size="medium"
                  value={values.name}
                  onBlur={handleBlur}
                  label={'Network name'}
                  onChange={handleChange}
                  placeholder="network-name"
                  InputLabelProps={{
                    shrink: true,
                  }}
                  sx={{
                    '.MuiOutlinedInput-input': {
                      height: '32px',
                      fontSize: '16px !important',
                    },
                  }}
                  helperText={touched.name && errors.name}
                  error={touched.name && Boolean(errors.name)}
                />
              </Stack>

              <Stack
                width="100%"
                direction={'row'}
                alignItems="center"
                justifyContent={'flex-end'}
              >
                <Button type="submit" variant="contained">
                  NAME NETWORK
                </Button>
              </Stack>
            </Stack>
          </form>
        )}
      </Formik>
    </Stack>
  );
};

export default Network;
