/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { NetworkDto } from '@/client/graphql/generated';
import { AddSiteValidationSchema } from '@/helpers/formValidators';
import { GlobalInput } from '@/styles/global';
import colors from '@/theme/colors';
import { TSiteForm } from '@/types';
import { isValidLatLng } from '@/utils';
import { Button, MenuItem, Stack, Typography } from '@mui/material';
import { Form, Formik } from 'formik';
import dynamic from 'next/dynamic';

const SiteMapComponent = dynamic(
  () => import('@/app/(main)/console/sites/[id]/_components/SiteMapComponent'),
  { loading: () => <p>Site map is loading</p>, ssr: false },
);

interface SiteDetailsStepProps {
  initialValues: TSiteForm;
  networks: NetworkDto[];
  address: string;
  currentAddress: string;
  isAddressLoading: boolean;
  addressError: string | null;
  addSiteLoading: boolean;
  onBack: () => void;
  onCancel: () => void;
  onComplete: (values: TSiteForm) => void;
  onFetchAddress: (lat: number, lng: number) => void;
  onFieldChange: (field: string, value: unknown) => void;
}

/**
 * Step 1 of ConfigureSiteDialog.
 * Collects site name, network assignment, and lat/lng with a live map preview.
 */
const SiteDetailsStep: React.FC<SiteDetailsStepProps> = ({
  initialValues,
  networks,
  address,
  currentAddress,
  isAddressLoading,
  addressError,
  addSiteLoading,
  onBack,
  onCancel,
  onComplete,
  onFetchAddress,
  onFieldChange,
}) => (
  <Formik
    initialValues={initialValues}
    onSubmit={onComplete}
    validationSchema={AddSiteValidationSchema[1]}
  >
    {({ values, errors, touched, isValid, setFieldValue, validateField }) => {
      const syncField = (field: string, value: unknown) => {
        setFieldValue(field, value);
        onFieldChange(field, value);
      };

      const maybeRefetchAddress = () => {
        if (
          !errors.latitude &&
          !errors.longitude &&
          isValidLatLng([values.latitude, values.longitude])
        ) {
          onFetchAddress(values.latitude, values.longitude);
          syncField('address', address);
        }
      };

      return (
        <Form>
          <Stack spacing={2}>
            {currentAddress && (
              <SiteMapComponent
                posix={[
                  values.latitude.toString(),
                  values.longitude.toString(),
                ]}
                address={currentAddress}
                height="200px"
                id="configure-site-map-dialog"
              />
            )}

            <Typography variant="body2" mt={1}>
              {Boolean(errors.latitude) || Boolean(errors.longitude) ? (
                <span style={{ color: colors.redMatt }}>
                  Please enter valid coordinates
                </span>
              ) : isAddressLoading ? (
                'Loading address...'
              ) : addressError ? (
                <span style={{ color: colors.redMatt }}>
                  Error fetching address. Please try again.
                </span>
              ) : (
                currentAddress
              )}
            </Typography>

            <GlobalInput
              fullWidth
              required
              margin="normal"
              name="siteName"
              label="Site name"
              placeholder="site-name"
              error={touched.siteName && Boolean(errors.siteName)}
              helperText={touched.siteName && errors.siteName}
              slotProps={{ inputLabel: { shrink: true } }}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                syncField('siteName', e.target.value)
              }
            />

            <GlobalInput
              fullWidth
              select
              required
              name="network"
              label="Network"
              margin="normal"
              slotProps={{ inputLabel: { shrink: true } }}
              error={touched.network && Boolean(errors.network)}
              helperText={touched.network && errors.network}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                syncField('network', e.target.value)
              }
            >
              <MenuItem value="" disabled>
                Choose a network to add your site to
              </MenuItem>
              {networks.map((network) => (
                <MenuItem key={network.id} value={network.id}>
                  {network.name}
                </MenuItem>
              ))}
            </GlobalInput>

            <GlobalInput
              required
              fullWidth
              type="number"
              label="Latitude"
              name="latitude"
              value={values.latitude}
              onBlur={() => {
                validateField('latitude');
                maybeRefetchAddress();
              }}
              slotProps={{ inputLabel: { shrink: true } }}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                syncField('latitude', parseFloat(e.target.value))
              }
              error={touched.latitude && Boolean(errors.latitude)}
              helperText={touched.latitude && errors.latitude}
            />

            <GlobalInput
              required
              fullWidth
              type="number"
              label="Longitude"
              name="longitude"
              value={values.longitude}
              onBlur={() => {
                validateField('longitude');
                maybeRefetchAddress();
              }}
              slotProps={{ inputLabel: { shrink: true } }}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                syncField('longitude', parseFloat(e.target.value))
              }
              error={touched.longitude && Boolean(errors.longitude)}
              helperText={touched.longitude && errors.longitude}
            />
          </Stack>

          <Stack
            direction="row"
            spacing={1}
            justifyContent="flex-end"
            sx={{ mt: 2 }}
          >
            <Button onClick={onCancel}>Cancel</Button>
            <Button onClick={onBack}>Back</Button>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={
                !isValid || addSiteLoading || isAddressLoading || !address
              }
            >
              Submit
            </Button>
          </Stack>
        </Form>
      );
    }}
  </Formik>
);

export default SiteDetailsStep;
