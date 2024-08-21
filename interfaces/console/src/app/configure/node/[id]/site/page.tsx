'use client';
import {
  Component_Type,
  useGetComponentsByUserIdQuery,
} from '@/client/graphql/generated';
import { LField } from '@/components/Welcome';
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
  const qpPower = searchParams.get('power') || '';
  const qpSwitch = searchParams.get('switch') || '';
  const qpbackhaul = searchParams.get('backhaul') || '';

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

  const handleSubmit = () => {
    if (formik.isValid) {
      router.push(`/configure/node/${id}/site/name?${searchParams.toString()}`);
    }
  };

  const setQueryParam = (key: string, value: string) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set(key, value);
    window.history.pushState(null, '', `?${params.toString()}`);
  };

  const handleBack = () => router.back();

  return (
    <Paper elevation={0} sx={{ px: 4, py: 2 }}>
      <Typography variant="h6" fontWeight={500}>
        Configure site installation -{' '}
        <span style={{ color: colors.black70, fontWeight: 400 }}>
          <i>optional</i> (4/6)
        </span>
      </Typography>
      <FormikProvider value={formik}>
        <form onSubmit={formik.handleSubmit}>
          <Stack direction="column" mt={3} mb={3} spacing={2}>
            <Typography variant="body1" color={colors.vulcan}>
              You have successfully created your site, and now need to configure
              some settings. If the node or site location details are wrong,
              please check on your installation.
              <br />
              <br />
              If you did not install your site yourself, or used all the default
              options, click “Next”.
            </Typography>

            <LField
              label="Node"
              value={
                'Tower node #12381293891283192 + Amplifier unit #18238192398128931'
              }
            />
            <LField
              label="SITE LOCATION"
              value={
                '10349 Monstera Hills Road, San Ramon, CA 94611, United States'
              }
            />
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
