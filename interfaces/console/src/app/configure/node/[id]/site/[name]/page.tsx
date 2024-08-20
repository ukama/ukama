'use client';

import { useAddSiteMutation } from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import colors from '@/theme/colors';
import {
  AlertColor,
  Button,
  Paper,
  Stack,
  SvgIcon,
  TextField,
  Typography,
} from '@mui/material';
import { formatISO } from 'date-fns';
import { FormikProvider, FormikValues, useFormik } from 'formik';
import { useRouter, useSearchParams } from 'next/navigation';
import * as Yup from 'yup';
import NetworkInfo from '../../../../../../../public/svg/NetworkInfo';

const validationSchema = Yup.object({
  name: Yup.string()
    .required('Network name is required')
    .matches(
      /^[a-z0-9-]*$/,
      'Network name must be lowercase alphanumeric and should not contain spaces, "-" are allowed.',
    ),
});

interface IPage {
  params: {
    id: string;
    name: string;
  };
}

const Page = ({ params }: IPage) => {
  const { id, name } = params;
  const { setSnackbarMessage } = useAppContext();
  const router = useRouter();
  const searchParams = useSearchParams();
  const qpLat = searchParams.get('lat') ?? '';
  const qpLng = searchParams.get('lng') ?? '';
  const qpPower = searchParams.get('power') || '';
  const qpSwitch = searchParams.get('switch') || '';
  const qpbackhaul = searchParams.get('backhaul') || '';
  const formik = useFormik({
    initialValues: {
      name: '',
    },
    validateOnChange: true,
    onSubmit: (values) => {
      handleSubmit(values);
    },
    validationSchema: validationSchema,
  });

  const [addSite, { loading: addSiteLoading }] = useAddSiteMutation({
    onCompleted: () => {
      setSnackbarMessage({
        id: 'add-site-success',
        message: 'Site added successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'add-site-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const handleBack = () => router.back();

  const handleSubmit = (values: FormikValues) => {
    console.log(
      id,
      name,
      qpLat,
      qpLng,
      qpPower,
      qpSwitch,
      qpbackhaul,
      values.name,
    );

    addSite({
      variables: {
        data: {
          name: name,
          power_id: qpPower,
          location: values.address,
          access_id: id,
          switch_id: qpSwitch,
          latitude: parseFloat(qpLat),
          network_id: values.network,
          longitude: parseFloat(qpLng),
          backhaul_id: qpbackhaul,
          spectrum_id: '',
          install_date: formatISO(new Date()),
        },
      },
    });
  };

  return (
    <Paper elevation={0} sx={{ px: 4, py: 2 }}>
      <Typography variant="h6" fontWeight={500}>
        Name network -{' '}
        <span style={{ color: colors.black70, fontWeight: 400 }}>
          <i>optional</i> (6/6)
        </span>
      </Typography>
      <FormikProvider value={formik}>
        <form onSubmit={formik.handleSubmit}>
          <Stack
            direction="column"
            mt={3}
            mb={3}
            spacing={3}
            alignItems={'center'}
          >
            <Typography variant="body1">
              You have successfully created your first network, and can always
              add more sites to it later! Please name it for your ease of
              reference.
            </Typography>

            <SvgIcon sx={{ width: 240, height: 176, mt: 2, mb: 4 }}>
              {NetworkInfo}
            </SvgIcon>

            <TextField
              fullWidth
              id={'name'}
              name={'name'}
              size="medium"
              label={'Network name'}
              placeholder="network-name"
              onBlur={formik.handleBlur}
              value={formik.values.name}
              onChange={formik.handleChange}
              InputLabelProps={{
                shrink: true,
              }}
              error={formik.touched.name && Boolean(formik.errors.name)}
              helperText={formik.touched.name && formik.errors.name}
            />
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
              Finish setup
            </Button>
          </Stack>
        </form>
      </FormikProvider>
    </Paper>
  );
};

export default Page;
