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
  ComponentsResDto,
  useGetComponentsByUserIdQuery,
  useGetNetworksQuery,
} from '@/client/graphql/generated';
import AppSnackbar from '@/components/AppSnackbar/page';
import { ONBOARDING_FLOW } from '@/constants';
import { useAppContext } from '@/context';
import { CenterContainer, GradiantBarNoRadius } from '@/styles/global';
import colors from '@/theme/colors';
import { ConfigureStep } from '@/utils';
import { AlertColor, Box, Typography } from '@mui/material';
import Grid from '@mui/material/Grid2';
import {
  useParams,
  usePathname,
  useRouter,
  useSearchParams,
} from 'next/navigation';
import { useEffect, useState } from 'react';
import { Logo } from '../../../public/svg/Logo';
import OnBoardingDynamic from '../../../public/svg/OnBoardingDynamic';

const ConfigureLayout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  const router = useRouter();
  const path = usePathname();
  const params = useParams<{ id: string; name: string }>();
  const searchParams = useSearchParams();
  const isSimsPath = path.includes('sims');
  const qpLat = searchParams.get('lat') ?? '';
  const qpLng = searchParams.get('lng') ?? '';
  const nid = searchParams.get('nid') ?? '';
  const accessId = searchParams.get('access') ?? '';
  const siteName = searchParams.get('name') ?? '';
  const flow = searchParams.get('flow') ?? ONBOARDING_FLOW;
  const { currentStep, totalStep } = ConfigureStep(path, flow);
  const { user, network, setNetwork, setSnackbarMessage } = useAppContext();
  const [parts, setParts] = useState({
    switchId: '',
    powerName: '',
    backhaulName: '',
  });

  useGetNetworksQuery({
    skip: path.includes('/configure/network'),
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (data.getNetworks.networks.length > 0) {
        setNetwork({
          id: data.getNetworks.networks[0].id,
          name: data.getNetworks.networks[0].name,
        });
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

  const { data: components } = useGetComponentsByUserIdQuery({
    fetchPolicy: 'cache-first',
    variables: {
      data: {
        category: Component_Type.All,
      },
    },
    onCompleted: (data) => {
      const { components } = data.getComponentsByUserId;

      if (components.length === 0) return;

      const nodeId = params.id ?? nid;
      if (!nodeId) return;

      const component = components.find((comp) => comp.partNumber === nodeId);
      if (!component?.id && flow !== ONBOARDING_FLOW) {
        setSnackbarMessage({
          id: 'components-msg',
          message: 'Access point not found in inventory.',
          type: 'warning',
          show: true,
        });
        router.push('/console/home');
        return;
      }

      mapComponents(data.getComponentsByUserId);
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
    if (components && components?.getComponentsByUserId.components.length > 0) {
      mapComponents(components.getComponentsByUserId);
    }
  }, [searchParams, components]);

  const mapComponents = (components: ComponentsResDto) => {
    const p = searchParams.get('power') ?? '';
    const s = searchParams.get('switch') ?? '';
    const b = searchParams.get('backhaul') ?? '';
    const switchRecords = components.components.find((comp) => comp.id === s);

    const powerRecords = components.components.find((comp) => comp.id === p);

    const backhaulRecords = components.components.find((comp) => comp.id === b);

    setParts({
      switchId: switchRecords?.description ?? '',
      powerName: powerRecords?.description ?? '',
      backhaulName: backhaulRecords?.description ?? '',
    });
  };

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
                  visibility={path.includes('complete') ? 'hidden' : 'visible'}
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
              <Box py={20} px={6} height={'100%'}>
                {OnBoardingDynamic({
                  textColor: colors.black,
                  isShowSimpool: isSimsPath,
                  selectedColor: colors.primaryMain,
                  isShowSite: params.id ? true : accessId ? true : false,
                  isShowComponents: params.name || siteName ? true : false,
                  siteName:
                    params.name || siteName
                      ? (params.name ?? siteName)
                      : 'Site 1',
                  networkName: network.name ? network.name : 'Network',
                  orgName: user.orgName ? user.orgName : 'Organization',
                  backhaulName: parts.backhaulName
                    ? parts.backhaulName
                    : 'Backhaul',
                  powerName: parts.powerName ? parts.powerName : 'Power',
                  switchName: parts.switchId ? parts.switchId : 'Switch',
                  nodeName: params.id ?? accessId,
                })}
              </Box>
            </CenterContainer>
          </Grid>
        </Grid>
      </Grid>
      <AppSnackbar />
    </Box>
  );
};

export default ConfigureLayout;
