/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import SiteMapComponent from '@/components/SiteMapComponent';
import { LField } from '@/components/Welcome';
import { SiteNameSchemaValidation } from '@/helpers/formValidators';
import colors from '@/theme/colors';
import { Button, Divider, Stack, TextField, Typography } from '@mui/material';
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

  const handleBack = () => {
    router.back();
  };

  const handleSubmit = (values: FormikValues) => {
    router.push(
      `/configure/node/${id}/site/${values.name}/install?${searchParams.toString()}`,
    );
  };

  return (
    <Stack spacing={2} overflow={'scroll'}>
      <Typography variant="h4" fontWeight={500}>
        Name site
      </Typography>

      <FormikProvider value={formik}>
        <form onSubmit={formik.handleSubmit}>
          <Stack direction="column">
            <Typography variant="body1" mb={4}>
              Please name your recently created site for ease of reference.
            </Typography>

            <Stack spacing={1} mb={3}>
              <SiteMapComponent
                posix={[latlng[0], latlng[1]]}
                address={address}
                height={'128px'}
              />

              <LField label="Site Location" value={address} />
            </Stack>

            <Divider sx={{ marginBottom: '16px !important' }} />

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

          <Stack
            direction={'row'}
            mt={{ xs: 4, md: 6 }}
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
              Next
            </Button>
          </Stack>
        </form>
      </FormikProvider>
    </Stack>
  );
};

export default SiteName;
