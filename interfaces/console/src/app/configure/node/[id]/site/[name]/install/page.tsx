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
  useGetComponentsByUserIdQuery,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import { SiteConfigureSchema } from '@/helpers/formValidators';
import { globalUseStyles } from '@/styles/global';
import colors from '@/theme/colors';
import {
  AlertColor,
  Button,
  Divider,
  MenuItem,
  Paper,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { Field, FormikProvider, useFormik } from 'formik';
import { useRouter, useSearchParams } from 'next/navigation';

interface ISiteConfigure {
  params: {
    id: string;
  };
}

const SiteConfigure = ({ params }: ISiteConfigure) => {
  const { id } = params;
  const router = useRouter();
  const searchParams = useSearchParams();
  const qpPower = searchParams.get('power') ?? '';
  const qpSwitch = searchParams.get('switch') ?? '';
  const flow = searchParams.get('flow') ?? 'onb';
  const step = parseInt(searchParams.get('step') ?? '1');
  const qpbackhaul = searchParams.get('backhaul') ?? '';

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
  const gclasses = globalUseStyles();
  const { setSnackbarMessage } = useAppContext();

  const { data: componentsData } = useGetComponentsByUserIdQuery({
    fetchPolicy: 'cache-and-network',
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

  const setQueryParam = (key: string, value: string) => {
    const p = new URLSearchParams(searchParams.toString());
    p.set(key, value);
    return p;
  };

  const handleSubmit = () => {
    if (formik.isValid) {
      const p = setQueryParam('step', (step + 1).toString());
      router.push(`/configure/node/${id}/site/name?${p.toString()}`);
    }
  };

  const handleBack = () => router.back();

  return (
    <Paper elevation={0} sx={{ px: { xs: 2, md: 4 }, py: { xs: 1, md: 2 } }}>
      <Stack direction={'row'}>
        <Typography variant="h6">{'Configure site installation'}</Typography>
        <Typography
          variant="h6"
          fontWeight={400}
          sx={{
            color: colors.black70,
          }}
        >
          {flow === 'onb' && <i>&nbsp;- optional</i>}&nbsp;({step}/
          {flow === 'onb' ? 6 : 4})
        </Typography>
      </Stack>

      <FormikProvider value={formik}>
        <form onSubmit={formik.handleSubmit}>
          <Stack direction="column" my={3} spacing={2.5}>
            <Typography variant="body1" color={colors.vulcan}>
              You have successfully created your site, and now need to configure
              some settings. If the node or site location details are wrong,
              please check on your installation.
              <br />
              <br />
              If you did not install your site yourself, or used all the default
              options, click “Next”.
            </Typography>

            <Divider sx={{ marginBottom: '8px !important' }} />
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
          <Stack mb={1} direction={'row'} justifyContent={'space-between'}>
            <Button
              variant="text"
              onClick={handleBack}
              sx={{ color: colors.black70, p: 0 }}
            >
              Back
            </Button>
            <Button type="submit" variant="contained">
              Next
            </Button>
          </Stack>
        </form>
      </FormikProvider>
    </Paper>
  );
};

export default SiteConfigure;
