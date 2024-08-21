'use client';
import SiteMapComponent from '@/components/SiteMapComponent';
import { LField } from '@/components/Welcome';
import colors from '@/theme/colors';
import {
  Button,
  Divider,
  Paper,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { FormikProvider, FormikValues, useFormik } from 'formik';
import { useRouter, useSearchParams } from 'next/navigation';
import { useState } from 'react';
import * as Yup from 'yup';

const validationSchema = Yup.object({
  name: Yup.string()
    .required('Site name is required')
    .matches(
      /^[a-z0-9-]*$/,
      'Site name must be lowercase alphanumeric and should not contain spaces, "-" are allowed.',
    ),
});

interface ISiteName {
  params: {
    id: string;
  };
}

const SiteName = ({ params }: ISiteName) => {
  const { id } = params;
  const router = useRouter();
  const searchParams = useSearchParams();
  const qpLat = searchParams.get('lat') ?? '';
  const qpLng = searchParams.get('lng') ?? '';
  const address = searchParams.get('address') ?? '';
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
  const [latlng] = useState<[number, number]>([
    parseFloat(qpLat),
    parseFloat(qpLng),
  ]);

  const handleBack = () => router.back();

  const handleSubmit = (values: FormikValues) => {
    router.push(
      `/configure/node/${id}/site/${values.name}?${searchParams.toString()}`,
    );
  };

  return (
    <Paper elevation={0} sx={{ px: 4, py: 2 }}>
      <Typography variant="h6" fontWeight={500}>
        Name site -{' '}
        <span style={{ color: colors.black70, fontWeight: 400 }}>
          <i>optional</i> (5/6)
        </span>
      </Typography>
      <FormikProvider value={formik}>
        <form onSubmit={formik.handleSubmit}>
          <Stack direction="column" mt={3} mb={3} spacing={2}>
            <Typography variant="body1">
              Please name your recently created site for ease of reference.
            </Typography>

            <SiteMapComponent
              posix={[latlng[0], latlng[1]]}
              address={address}
              height={'128px'}
            />

            <LField label="Site Location" value={address} />
            <Divider sx={{ marginBottom: '8px !important' }} />

            <TextField
              fullWidth
              id={'name'}
              name={'name'}
              size="medium"
              label={'Site name'}
              placeholder="site-name"
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
              Next
            </Button>
          </Stack>
        </form>
      </FormikProvider>
    </Paper>
  );
};

export default SiteName;
