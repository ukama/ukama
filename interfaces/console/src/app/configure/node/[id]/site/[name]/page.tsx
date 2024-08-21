'use client';

import {
  Component_Type,
  useAddNetworkMutation,
  useAddSiteMutation,
  useGetComponentsByUserIdQuery,
  useGetNetworksQuery,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import { globalUseStyles } from '@/styles/global';
import colors from '@/theme/colors';
import {
  AlertColor,
  Box,
  Button,
  CircularProgress,
  MenuItem,
  Paper,
  Skeleton,
  Stack,
  SvgIcon,
  TextField,
  Typography,
} from '@mui/material';
import { formatISO } from 'date-fns';
import { Field, FormikProvider, FormikValues, useFormik } from 'formik';
import { useRouter, useSearchParams } from 'next/navigation';
import { useState } from 'react';
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

const SiteLoadingState = ({ msg }: { msg: string }) => {
  return (
    <Stack direction={'column'} alignItems={'center'} my={3}>
      <CircularProgress />
      <Typography variant="body1">{msg}</Typography>
    </Stack>
  );
};

interface IPage {
  params: {
    id: string;
    name: string;
  };
}

const Page = ({ params }: IPage) => {
  const { id, name } = params;
  const router = useRouter();
  const gclasses = globalUseStyles();
  const { setSnackbarMessage } = useAppContext();
  const searchParams = useSearchParams();
  const qpLat = searchParams.get('lat') ?? '';
  const qpLng = searchParams.get('lng') ?? '';
  const qpPower = searchParams.get('power') || '';
  const qpSwitch = searchParams.get('switch') || '';
  const qpAddress = searchParams.get('address') ?? '';
  const qpbackhaul = searchParams.get('backhaul') || '';
  const [loadingMessage, setLoadingMessage] = useState<string>('');
  const [isCreateNetwork, setIsCreateNetwork] = useState<boolean>(false);
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

  const { data: networksData, loading: networksLoading } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (data.getNetworks.networks.length === 0) {
        setIsCreateNetwork(true);
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

  const [addNetwork, { loading: addNetworkLoading }] = useAddNetworkMutation({
    onCompleted: (data) => {
      setSnackbarMessage({
        id: 'add-networks-success',
        message: 'Network added successfully',
        type: 'success',
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'add-networks-error',
        message: error.message,
        type: 'error',
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
      router.push('/configure/success');
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

  const handleBack = async () => router.back();

  const handleSubmit = (values: FormikValues) => {
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

    if (isCreateNetwork) {
      setLoadingMessage('Creating network...');
      addNetwork({
        variables: {
          data: {
            isDefault: false,
            name: values.name,
            budget: values.budget,
            networks: values.networks,
            countries: values.countries,
          },
        },
      }).then((res) => {
        setLoadingMessage('Creating site...');
        addSiteCall(accessId, spectrumId, res.data?.addNetwork.id ?? '');
      });
    } else {
      setLoadingMessage('Creating site...');
      addSiteCall(accessId, spectrumId, values.name);
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

  return (
    <Paper elevation={0} sx={{ px: 4, py: 2 }}>
      {networksLoading || addSiteLoading ? (
        <SiteLoadingState msg={''} />
      ) : (
        <Box>
          <Typography variant="h6" fontWeight={500}>
            Name network -{' '}
            <span style={{ color: colors.black70, fontWeight: 400 }}>
              <i>optional</i> (6/6)
            </span>
          </Typography>
          <FormikProvider value={formik}>
            <form onSubmit={formik.handleSubmit}>
              <Stack
                my={3}
                spacing={3}
                direction="column"
                alignItems={'center'}
              >
                <Typography variant="body1">
                  You have successfully created your first network, and can
                  always add more sites to it later! Please name it for your
                  ease of reference.
                </Typography>

                <SvgIcon sx={{ width: 240, height: 176, mt: 2, mb: 4 }}>
                  {NetworkInfo}
                </SvgIcon>

                {isCreateNetwork ? (
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
                ) : networksLoading ? (
                  <Skeleton variant="rounded" width={'100%'} height={42} />
                ) : (
                  <Field
                    as={TextField}
                    fullWidth
                    select
                    required
                    name="network"
                    label="Network"
                    margin="normal"
                    InputLabelProps={{
                      shrink: true,
                    }}
                    InputProps={{
                      classes: {
                        input: gclasses.inputFieldStyle,
                      },
                    }}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                      formik.setFieldValue('name', e.target.value);
                    }}
                    helperText={formik.touched.name && formik.errors.name}
                    error={formik.touched.name && Boolean(formik.errors.name)}
                  >
                    <MenuItem value="" disabled>
                      Choose a network to add your site to
                    </MenuItem>
                    {networksData?.getNetworks.networks.map((network) => (
                      <MenuItem key={network.id} value={network.id}>
                        {network.name}
                      </MenuItem>
                    ))}
                  </Field>
                )}
                {networksData?.getNetworks.networks &&
                  networksData?.getNetworks.networks.length > 0 && (
                    <Button
                      variant="text"
                      sx={{
                        textTransform: 'none',
                        px: 0,
                        alignSelf: 'end',
                        mt: '0px !important',
                      }}
                      onClick={() => {
                        formik.setFieldValue('name', '');
                        setIsCreateNetwork(!isCreateNetwork);
                      }}
                    >
                      {isCreateNetwork
                        ? 'Choose existing network'
                        : 'Create new network'}
                    </Button>
                  )}
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
        </Box>
      )}
    </Paper>
  );
};

export default Page;
