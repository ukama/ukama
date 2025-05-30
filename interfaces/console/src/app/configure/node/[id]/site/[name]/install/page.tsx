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
} from '@/client/graphql/generated';
import { ONBOARDING_FLOW } from '@/constants';
import { useAppContext } from '@/context';
import { SiteConfigureSchema } from '@/helpers/formValidators';
import colors from '@/theme/colors';
import { setQueryParam } from '@/utils';
import {
  AlertColor,
  Autocomplete,
  Button,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { formatISO } from 'date-fns';
import { FormikProvider, useFormik } from 'formik';
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
  const searchParams = useSearchParams();
  const [loading, setLoading] = useState<boolean>(false);
  const qpLat = searchParams.get('lat') ?? '';
  const qpLng = searchParams.get('lng') ?? '';
  const flow = searchParams.get('flow') ?? 'onb';
  const qpPower = searchParams.get('power') ?? '';
  const qpSwitch = searchParams.get('switch') ?? '';
  const qpAddress = searchParams.get('address') ?? '';
  const networkId = searchParams.get('networkid') ?? '';
  const qpbackhaul = searchParams.get('backhaul') ?? '';
  const { setSnackbarMessage } = useAppContext();

  const formik = useFormik({
    initialValues: {
      power: qpPower ?? '',
      switch: qpSwitch ?? '',
      backhaul: qpbackhaul ?? '',
    },
    validateOnChange: true,
    onSubmit: () => handleSubmit(),
    validationSchema: SiteConfigureSchema,
  });

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
      const p = setQueryParam('access', id, searchParams.toString(), pathname);
      p.set('name', name);
      if (flow === ONBOARDING_FLOW) {
        router.push(`/configure/sims?${p.toString()}`);
      } else {
        router.push(`/configure/complete?${p.toString()}`);
      }
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

  const { data: componentsData, loading: componentsLoading } =
    useGetComponentsByUserIdQuery({
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
            formik.setFieldValue('switch', switchRecords[0].id);
          }
          if (powerRecords.length === 1) {
            formik.setFieldValue('power', powerRecords[0].id);
          }
          if (backhaulRecords.length === 1) {
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
    if (formik.values.switch !== '') {
      setQueryParam(
        'switch',
        formik.values.switch,
        searchParams.toString(),
        pathname,
      );
    }

    if (formik.values.power !== '') {
      setQueryParam(
        'power',
        formik.values.power,
        searchParams.toString(),
        pathname,
      );
    }

    if (formik.values.backhaul !== '') {
      setQueryParam(
        'backhaul',
        formik.values.backhaul,
        searchParams.toString(),
        pathname,
      );
    }
  }, [formik.values]);

  useEffect(() => {
    if (addSiteLoading) setLoading(true);
  }, [addSiteLoading]);

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

    // TODO: Choosing first spectrum component. Need to revisit if spectrum have multiple components
    const spectrumId =
      spectrumComponentsData?.getComponentsByUserId.components[0].id;

    if (!accessId || !spectrumId) {
      setSnackbarMessage({
        id: 'add-site-error',
        message: 'Node or Spectrum components not found',
        type: 'error' as AlertColor,
        show: true,
      });
      return;
    }

    if (formik.isValid && networkId) {
      addSiteCall(accessId, spectrumId, networkId);
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
    router.back();
  };

  if (loading) return <LoadingSkeleton />;

  return (
    <Stack spacing={2} overflow={'scroll'}>
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

            <Autocomplete
              fullWidth
              loading={componentsLoading}
              options={
                (componentsData &&
                  componentsData?.getComponentsByUserId.components.filter(
                    (comp) => comp.category === 'SWITCH',
                  )) ??
                []
              }
              getOptionLabel={(option) => option.description}
              getOptionKey={(option) => option.id}
              value={
                componentsData?.getComponentsByUserId.components.find(
                  (comp) => comp.id === formik.values.switch,
                ) || null
              }
              onChange={(_, v: any) => {
                setQueryParam(
                  'switch',
                  v?.id || '',
                  searchParams.toString(),
                  pathname,
                );
                formik.setFieldValue('switch', v?.id || '');
              }}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="SWITCH"
                  placeholder="Select a switch"
                  value={formik.values.switch}
                  slotProps={{ inputLabel: { shrink: true } }}
                  error={formik.touched.switch && Boolean(formik.errors.switch)}
                  helperText={formik.touched.switch && formik.errors.switch}
                />
              )}
            />
            <Autocomplete
              fullWidth
              loading={componentsLoading}
              options={
                (componentsData &&
                  componentsData?.getComponentsByUserId.components.filter(
                    (comp) => comp.category === 'POWER',
                  )) ??
                []
              }
              getOptionLabel={(option) => option.description}
              getOptionKey={(option) => option.id}
              value={
                componentsData?.getComponentsByUserId.components.find(
                  (comp) => comp.id === formik.values.power,
                ) || null
              }
              onChange={(_, v: any) => {
                setQueryParam(
                  'power',
                  v?.id || '',
                  searchParams.toString(),
                  pathname,
                );
                formik.setFieldValue('power', v?.id || '');
              }}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="POWER"
                  placeholder="Select a power (Solar / Battery / Charge controller)"
                  value={formik.values.switch}
                  slotProps={{ inputLabel: { shrink: true } }}
                  error={formik.touched.power && Boolean(formik.errors.power)}
                  helperText={formik.touched.power && formik.errors.power}
                />
              )}
            />
            <Autocomplete
              fullWidth
              loading={componentsLoading}
              options={
                (componentsData &&
                  componentsData?.getComponentsByUserId.components.filter(
                    (comp) => comp.category === 'BACKHAUL',
                  )) ??
                []
              }
              getOptionLabel={(option) => option.description}
              getOptionKey={(option) => option.id}
              value={
                componentsData?.getComponentsByUserId.components.find(
                  (comp) => comp.id === formik.values.backhaul,
                ) || null
              }
              onChange={(_, v: any) => {
                setQueryParam(
                  'backhaul',
                  v?.id || '',
                  searchParams.toString(),
                  pathname,
                );
                formik.setFieldValue('backhaul', v?.id || '');
              }}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="BACKHAUL"
                  placeholder="Select a backhaul"
                  value={formik.values.switch}
                  slotProps={{ inputLabel: { shrink: true } }}
                  error={
                    formik.touched.backhaul && Boolean(formik.errors.backhaul)
                  }
                  helperText={formik.touched.backhaul && formik.errors.backhaul}
                />
              )}
            />
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
