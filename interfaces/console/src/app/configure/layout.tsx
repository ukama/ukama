/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { useGetNetworksQuery } from '@/client/graphql/generated';
import AppSnackbar from '@/components/AppSnackbar/page';
import { ONBOARDING_FLOW } from '@/constants';
import { useAppContext } from '@/context';
import { CenterContainer, GradiantBarNoRadius } from '@/styles/global';
import colors from '@/theme/colors';
import { ConfigureStep, isValidLatLng } from '@/utils';
import { Box, Typography } from '@mui/material';
import Grid from '@mui/material/Grid2';
import {
  useParams,
  usePathname,
  useRouter,
  useSearchParams,
} from 'next/navigation';
import CustomNetworkInfo from '../../../public/svg/CustomNetworkInfo';
import { Logo } from '../../../public/svg/Logo';

const ConfigureLayout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  const router = useRouter();
  const params = useParams<{ id: string; name: string }>();
  const path = usePathname();
  const searchParams = useSearchParams();
  const qpLat = searchParams.get('lat') ?? '';
  const qpLng = searchParams.get('lng') ?? '';
  const pstep = parseInt(searchParams.get('step') ?? '1');
  const flow = searchParams.get('flow') ?? ONBOARDING_FLOW;
  const { currentStep, totalStep } = ConfigureStep(path, flow, pstep);
  const { network, setNetwork, setSnackbarMessage } = useAppContext();

  useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (data.getNetworks.networks.length > 0) {
        setNetwork({
          id: data.getNetworks.networks[0].id,
          name: data.getNetworks.networks[0].name,
        });
        // router.push(`/configure/check?flow=${CHECK_SITE_FLOW}`);
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

  // if (flow !== ONBOARDING_FLOW && flow !== NETWORK_FLOW)
  //   return (
  //     <Box width="100%" height="100%" overflow="hidden">
  //       <Stack height="100%">
  //         <GradiantBarNoRadius />
  //         <Container maxWidth={'sm'} sx={{ height: '100%' }}>
  //           <CenterContainer>
  //             <Stack spacing={2} p={{ xs: 4, md: 8 }}>
  //               <Grid container size={12} rowSpacing={6} height={'fit-content'}>
  //                 <Grid size={12}>
  //                   {Logo({
  //                     color: colors.primaryMain,
  //                     width: 120,
  //                     height: 37,
  //                   })}
  //                 </Grid>
  //                 <Grid size={12}>
  //                   <Typography
  //                     fontWeight={600}
  //                     variant="caption"
  //                     lineHeight={'18px'}
  //                     letterSpacing={'1.5px'}
  //                     color={colors.tertiary}
  //                   >
  //                     STEP {`${currentStep}/${totalStep}`}
  //                   </Typography>
  //                 </Grid>
  //                 <Grid size={12} height={'100%'}>
  //                   {children}
  //                 </Grid>
  //               </Grid>
  //             </Stack>
  //           </CenterContainer>
  //         </Container>
  //       </Stack>
  //       <AppSnackbar />
  //     </Box>
  //   );

  return (
    <Box width="100%" height="100%" overflow="hidden">
      <Grid container height={'100%'}>
        <Grid size={12}>
          <GradiantBarNoRadius />
        </Grid>
        <Grid container height={'100%'} size={12}>
          <Grid
            container
            spacing={1}
            height={'fit-content'}
            size={{ xs: 12, md: 6 }}
            p={{ xs: 4, md: 12 }}
          >
            <Grid container size={12} rowSpacing={6} height={'fit-content'}>
              <Grid size={12}>
                {Logo({ color: colors.primaryMain, width: 120, height: 37 })}
              </Grid>
              <Grid size={12}>
                <Typography
                  fontWeight={600}
                  variant="caption"
                  lineHeight={'18px'}
                  letterSpacing={'1.5px'}
                  color={colors.tertiary}
                >
                  STEP {`${currentStep}/${totalStep}`}
                </Typography>
              </Grid>
            </Grid>
            <Grid size={12} height={'100%'}>
              {children}
            </Grid>
          </Grid>
          <Grid
            size={{ xs: 0, md: 6 }}
            bgcolor={colors.solitude}
            display={{ xs: 'none', md: 'flex' }}
          >
            <CenterContainer>
              {CustomNetworkInfo({
                networkName: network.name ? network.name : 'NETWORK',
                networkColor: network.name ? colors.primaryMain : '#333333',
                networkIconColor: network.name ? colors.primaryMain : '#6F7979',
                siteOneIconColor:
                  qpLat &&
                  qpLng &&
                  isValidLatLng([parseFloat(qpLat), parseFloat(qpLng)])
                    ? colors.primaryMain
                    : '#6F7979',
                siteOneName: params.name ? params.name : 'SITE 1',
                siteOneColor: params.name ? colors.primaryMain : '#333333',
              })}
            </CenterContainer>
          </Grid>
        </Grid>
      </Grid>
      <AppSnackbar />
    </Box>
  );
};

export default ConfigureLayout;
