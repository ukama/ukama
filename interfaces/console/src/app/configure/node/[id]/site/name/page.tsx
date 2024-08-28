'use client';
import SiteMapComponent from '@/components/SiteMapComponent';
import { LField } from '@/components/Welcome';
import { SiteNameSchemaValidation } from '@/helpers/formValidators';
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
  const stepTracker = searchParams.get('step') ?? '1';
  const formik = useFormik({
    initialValues: {
      name: '',
    },
    validateOnChange: true,
    onSubmit: (values) => {
      handleSubmit(values);
    },
    validationSchema: SiteNameSchemaValidation,
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
      <Stack direction={'row'}>
        <Typography variant="h6">{'Name site'}</Typography>
        <Typography
          variant="h6"
          fontWeight={400}
          sx={{
            color: colors.black70,
            display: stepTracker !== '1' ? 'none' : 'flex',
          }}
        >
          <i>&nbsp;- optional</i>&nbsp;(5/6)
        </Typography>
      </Stack>

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
