/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  Component_Type,
  useAddSiteMutation,
  useGetComponentsByUserIdQuery,
  useGetNetworksQuery,
} from '@/client/graphql/generated';
import { INSTALLATION_FLOW } from '@/constants';
import { useAppContext } from '@/context';
import { SiteConfigureSchema } from '@/helpers/formValidators';
import { globalUseStyles } from '@/styles/global';
import colors from '@/theme/colors';
import {
  AlertColor,
  Button,
  MenuItem,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { formatISO } from 'date-fns';
import { Field, FormikProvider, useFormik } from 'formik';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import LoadingSkeleton from './skelton';

interface IPage {
  params: {
    id: string;
    name: string;
  };
}

const SiteConfigure = ({ params }: IPage) => {
  const { id, name } = params;
  const router = useRouter();
  const pathname = usePathname();
  const gclasses = globalUseStyles();
  const searchParams = useSearchParams();
  const [loading, setLoading] = useState<boolean>(false);
  const qpLat = searchParams.get('lat') ?? '';
  const qpLng = searchParams.get('lng') ?? '';
  const flow = searchParams.get('flow') ?? 'onb';
  const qpPower = searchParams.get('power') ?? '';
  const qpSwitch = searchParams.get('switch') ?? '';
  const qpAddress = searchParams.get('address') ?? '';
  const step = parseInt(searchParams.get('step') ?? '1');
  const qpbackhaul = searchParams.get('backhaul') ?? '';
  const { setSnackbarMessage, network } = useAppContext();
  const formik = useFormik({
    initialValues: {
      power: qpPower,
      switch: qpSwitch,
      backhaul: qpbackhaul,
    },
    validateOnChange: true,
    onSubmit: (values) => {
      handleSubmit();
    },
    validationSchema: SiteConfigureSchema,
  });

  const { data: networkData } = useGetNetworksQuery();

  const { data: accessComponentsData } = useGetComponentsByUserIdQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      data: {
        category: Component_Type.Access,
      },
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'components-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: spectrumComponentsData } = useGetComponentsByUserIdQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      data: {
        category: Component_Type.Spectrum,
      },
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'components-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [addSite, { loading: addSiteLoading }] = useAddSiteMutation({
    onCompleted: () => {
      setSnackbarMessage({
        id: 'add-site-success',
        message: 'Site added successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
      router.push(
        `/configure/sims?step=${flow !== INSTALLATION_FLOW ? 4 : 5}&flow=${flow}`,
      );
    },
    onError: (error) => {
      setLoading(false);
      setSnackbarMessage({
        id: 'add-site-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: componentsData } = useGetComponentsByUserIdQuery({
    fetchPolicy: 'cache-first',
    variables: {
      data: {
        category: Component_Type.All,
      },
    },
    onCompleted: (data) => {
      if (data.getComponentsByUserId.components.length > 0) {
        const switchRecords = data.getComponentsByUserId.components.filter(
          (comp) => comp.category === 'SWITCH',
        );

        const powerRecords = data.getComponentsByUserId.components.filter(
          (comp) => comp.category === 'POWER',
        );

        const backhaulRecords = data.getComponentsByUserId.components.filter(
          (comp) => comp.category === 'BACKHAUL',
        );

        if (switchRecords.length === 1) {
          setQueryParam('switch', switchRecords[0].id);
          formik.setFieldValue('switch', switchRecords[0].id);
        }
        if (powerRecords.length === 1) {
          setQueryParam('power', switchRecords[0].id);
          formik.setFieldValue('power', powerRecords[0].id);
        }
        if (backhaulRecords.length === 1) {
          setQueryParam('backhaul', switchRecords[0].id);
          formik.setFieldValue('backhaul', backhaulRecords[0].id);
        }
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'components-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  useEffect(() => {
    if (addSiteLoading) setLoading(true);
  }, [addSiteLoading]);

  const setQueryParam = (key: string, value: string) => {
    const p = new URLSearchParams(searchParams.toString());
    p.set(key, value);
    window.history.replaceState({}, '', `${pathname}?${p.toString()}`);
    return p;
  };

  const handleSubmit = () => {
    if (
      qpAddress === '' ||
      qpPower === '' ||
      qpSwitch === '' ||
      qpbackhaul === ''
    ) {
      setSnackbarMessage({
        id: 'add-site-error',
        message: 'Require data is missing. Please complete the previous steps',
        type: 'error' as AlertColor,
        show: true,
      });
      return;
    }

    const accessId =
      accessComponentsData?.getComponentsByUserId.components.find(
        (component) => component.partNumber === id,
      )?.id;

    const spectrumId =
      spectrumComponentsData?.getComponentsByUserId.components[0].id;

    if (!accessId || !spectrumId) {
      setSnackbarMessage({
        id: 'add-site-error',
        message: 'Access or Spectrum components not found',
        type: 'error' as AlertColor,
        show: true,
      });
      return;
    }

    if (
      formik.isValid &&
      networkData &&
      networkData?.getNetworks.networks.length > 0
    ) {
      addSiteCall(
        accessId,
        spectrumId,
        networkData?.getNetworks.networks[0].id,
      );
    }
  };

  const addSiteCall = (
    accessId: string,
    spectrumId: string,
    networkId: string,
  ) => {
    addSite({
      variables: {
        data: {
          name: name,
          power_id: qpPower,
          access_id: accessId,
          switch_id: qpSwitch,
          location: qpAddress,
          backhaul_id: qpbackhaul,
          spectrum_id: spectrumId,
          latitude: parseFloat(qpLat),
          longitude: parseFloat(qpLng),
          install_date: formatISO(new Date()),
          network_id: networkId,
        },
      },
    });
  };

  const handleBack = () => {
    setQueryParam('step', (step - 1).toString());
    router.back();
  };

  if (loading) return <LoadingSkeleton />;

  return (
    <Stack spacing={2}>
      <Typography variant="h4" fontWeight={500}>
        Configure site settings
      </Typography>

      <FormikProvider value={formik}>
        <form onSubmit={formik.handleSubmit}>
          <Stack direction="column" spacing={2.5}>
            <Typography variant="body1" color={colors.vulcan}>
              If you did not install your site yourself, or used all the default
              options, click “Next”.
              <br />
              <br />
              If you used your own custom options for the switch, power, or
              backhaul, please select “Other” in the corresponding dropdown.
              Please note that we currently cannot track real time KPIs for
              custom options.
              <br />
              <br />
            </Typography>

            <Field
              as={TextField}
              select
              required
              fullWidth
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
              onChange={(e: any) => {
                setQueryParam('switch', e.target.value);
                formik.handleChange(e);
              }}
              error={formik.touched.switch && Boolean(formik.errors.switch)}
              helperText={formik.touched.switch && formik.errors.switch}
            >
              {componentsData?.getComponentsByUserId.components
                .filter((comp) => comp.category === 'SWITCH')
                .map((component) => (
                  <MenuItem key={component.id} value={component.id}>
                    {component.description}
                  </MenuItem>
                ))}
            </Field>
            <Field
              as={TextField}
              select
              required
              fullWidth
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
              onChange={(e: any) => {
                setQueryParam('power', e.target.value);
                formik.handleChange(e);
              }}
              error={formik.touched.power && Boolean(formik.errors.power)}
              helperText={formik.touched.power && formik.errors.power}
            >
              {componentsData?.getComponentsByUserId.components
                .filter((comp) => comp.category === 'POWER')
                .map((component) => (
                  <MenuItem key={component.id} value={component.id}>
                    {component.description}
                  </MenuItem>
                ))}
            </Field>
            <Field
              select
              required
              fullWidth
              as={TextField}
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
              onChange={(e: any) => {
                setQueryParam('backhaul', e.target.value);
                formik.handleChange(e);
              }}
              error={formik.touched.backhaul && Boolean(formik.errors.backhaul)}
              helperText={formik.touched.backhaul && formik.errors.backhaul}
            >
              {componentsData?.getComponentsByUserId.components
                .filter((comp) => comp.category === 'BACKHAUL')
                .map((component) => (
                  <MenuItem key={component.id} value={component.id}>
                    {component.description}
                  </MenuItem>
                ))}
            </Field>
          </Stack>
          <Stack
            mt={{ xs: 4, md: 6 }}
            direction={'row'}
            justifyContent={'space-between'}
          >
            <Button
              variant="text"
              onClick={handleBack}
              sx={{ color: colors.black70, p: 0 }}
            >
              Back
            </Button>
            <Button type="submit" variant="contained">
              Configure site
            </Button>
          </Stack>
        </form>
      </FormikProvider>
    </Stack>
  );
};

export default SiteConfigure;
