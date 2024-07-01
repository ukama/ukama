/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { useState } from 'react';
import dynamic from 'next/dynamic';
import NetworkIcon from '@mui/icons-material/Hub';
import { AlertColor, Box, Paper, Skeleton, Stack } from '@mui/material';
import Grid from '@mui/material/Unstable_Grid2';
import { LabelOverlayUI, SitesTree } from '@/components/NetworkMap/OverlayUI';
import { DataVolume, Throughput, UsersWithBG } from '@/../public/svg';
import { MONTH_FILTER, SIM_TYPE_OPERATOR, TIME_FILTER } from '@/constants';
import EmptyView from '@/components/EmptyView';
import LoadingWrapper from '@/components/LoadingWrapper';
import { useAppContext } from '@/context';
import NetworkStatus from '@/components/NetworkStatus';
import OnboardingCard from '@/components/OnboardingCard';
import SiteConfigurationStepperDialog from '@/components/SiteConfigurationStepperDialog';
import AddSubscriberDialog from '@/components/AddSubscriber';
import { TAddSubscriberData } from '@/types';
import {
  useAddSubscriberMutation,
  useAllocateSimMutation,
  useGetPackagesForSimLazyQuery,
  useGetPackagesQuery,
  useGetSimLazyQuery,
  useGetSimpoolStatsQuery,
  useGetSimsQuery,
  useToggleSimStatusMutation,
  useSetActivePackageForSimMutation,
} from '@/client/graphql/generated';

const NetworkMap = dynamic(() => import('@/components/NetworkMap'), {
  ssr: false,
  loading: () => (
    <Skeleton
      variant="rectangular"
      sx={{
        borderRadius: '10px',
        height: 'calc(100vh - 332px)',
        width: '100%',
      }}
    />
  ),
});

const StatusCard = dynamic(() => import('@/components/StatusCard'), {
  ssr: false,
  loading: () => (
    <Skeleton
      variant="rectangular"
      sx={{ borderRadius: '10px', height: '90px', width: '100%', m: 0 }}
    />
  ),
});

export default function Page() {
  const { network, setSnackbarMessage } = useAppContext();
  const [isOnboardingOpen, setOnboardingOpen] = useState(true);
  const [qrCode, setQrCode] = useState<string>('');
  const [simId, setSimId] = useState<string>('');
  const [subscriberSuccess, setSubscriberSuccess] = useState<boolean>(false);
  const [addSubscriberData, setAddSubscriberData] =
    useState<TAddSubscriberData>({
      email: '',
      iccid: '',
      name: '',
      phone: '',
      plan: '',
      simType: 'eSim',
      roamingStatus: false,
    });
  const [activeDialogs, setActiveDialogs] = useState<{
    [key: number]: boolean;
  }>({
    1: false,
    2: false,
    3: false,
  });
  const [activatePackageSim] = useSetActivePackageForSimMutation({
    onCompleted: () => {
      setSnackbarMessage({
        id: 'package-activated-success',
        message: 'Package activated successfully!',
        type: 'success' as AlertColor,
        show: true,
      });
      setSubscriberSuccess(true);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'package-activated-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: packagesData, loading: packagesLoading } = useGetPackagesQuery({
    fetchPolicy: 'cache-and-network',
    onError: (error) =>
      setSnackbarMessage({
        id: 'packages',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      }),
  });

  const [toggleSimStatus, { loading: toggleSimStatusLoading }] =
    useToggleSimStatusMutation({
      onCompleted: () => {
        setSnackbarMessage({
          id: 'sim-activated-success',
          message: 'Sim activated successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
        getSim({
          variables: {
            data: {
              simId: simId,
            },
          },
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'sim-activated-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const [addSubscriber, { loading: addSubscriberLoading }] =
    useAddSubscriberMutation({
      onCompleted: (res) => {
        setSnackbarMessage({
          id: 'add-subscriber-success',
          message: 'Subscriber added successfully!',
          type: 'success' as AlertColor,
          show: true,
        });

        allocateSim({
          variables: {
            data: {
              network_id: res.addSubscriber.networkId,
              package_id: addSubscriberData.plan ?? '',
              subscriber_id: res.addSubscriber.uuid,
              sim_type: SIM_TYPE_OPERATOR,
              iccid: addSubscriberData.iccid,
              traffic_policy: 10,
            },
          },
        });
      },
      onError: (error) =>
        setSnackbarMessage({
          id: 'add-subscriber-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        }),
    });

  const { data: simStatData } = useGetSimpoolStatsQuery({
    variables: { type: SIM_TYPE_OPERATOR },
    fetchPolicy: 'cache-and-network',
    onError: (error) =>
      setSnackbarMessage({
        id: 'sims-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      }),
  });

  const handleOnboardingClose = () => setOnboardingOpen(false);

  const handleStepClick = (step: number) => {
    console.log('Clicked step:', step);
    setActiveDialogs({ ...activeDialogs, [step]: true });
  };

  const handleCloseDialog = (step: number) =>
    setActiveDialogs({ ...activeDialogs, [step]: false });

  const handleSiteConfig = (formData: any) =>
    console.log('Form data submitted:', formData);

  const [getSim] = useGetSimLazyQuery({
    onCompleted: (res) => {},
  });

  const [getPackagesForSim] = useGetPackagesForSimLazyQuery({
    onCompleted: (res) => {
      if (res.getPackagesForSim.packages) {
        activatePackageSim({
          variables: {
            data: {
              sim_id: res.getPackagesForSim.sim_id,
              package_id: res.getPackagesForSim.packages[0].package_id,
            },
          },
        });
      }
    },
  });

  const [allocateSim, { loading: allocateSimLoading }] = useAllocateSimMutation(
    {
      onCompleted: (res) => {
        setSimId(res.allocateSim.id);
        setSnackbarMessage({
          id: 'sim-allocated-success',
          message: 'Sim allocated successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
        setQrCode(res.allocateSim.iccid);
        getPackagesForSim({
          variables: { data: { sim_id: res.allocateSim.id } },
        });
        getSim({
          variables: { data: { simId: res.allocateSim.id } },
        });
      },
      onError: (error) =>
        setSnackbarMessage({
          id: 'sim-allocated-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        }),
    },
  );

  const { data: simPoolData } = useGetSimsQuery({
    variables: { type: SIM_TYPE_OPERATOR },
    fetchPolicy: 'network-only',
    onError: (error) =>
      setSnackbarMessage({
        id: 'sims-error-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      }),
  });

  const handleAddSubscriber = async (values: TAddSubscriberData) => {
    setAddSubscriberData(values);
    await addSubscriber({
      variables: {
        data: {
          email: values.email,
          phone: values.phone,
          first_name: values.name,
          last_name: 'name',
          network_id: network.id,
        },
      },
    });
  };

  const renderInstallationDialog = () => (
    <>
      {activeDialogs[1] && (
        <SiteConfigurationStepperDialog
          open={activeDialogs[1]}
          handleClose={() => handleCloseDialog(1)}
          handleFormDataSubmit={handleSiteConfig}
        />
      )}
      {activeDialogs[2] && (
        <SiteConfigurationStepperDialog
          open={activeDialogs[2]}
          handleClose={() => handleCloseDialog(1)}
          handleFormDataSubmit={handleSiteConfig}
        />
      )}
      {activeDialogs[3] && (
        <AddSubscriberDialog
          qrCode={qrCode}
          pkgList={packagesData?.getPackages.packages ?? []}
          onSuccess={subscriberSuccess}
          onClose={() => handleCloseDialog(3)}
          onSubmit={handleAddSubscriber}
          open={activeDialogs[3]}
          sims={simPoolData?.getSims.sim ?? []}
          pSimCount={simStatData?.getSimPoolStats.physical}
          eSimCount={simStatData?.getSimPoolStats.physical}
          submitButtonState={
            addSubscriberLoading ?? allocateSimLoading ?? packagesLoading
          }
          loading={
            addSubscriberLoading ?? allocateSimLoading ?? packagesLoading
          }
        />
      )}
    </>
  );

  return (
    <Grid container rowSpacing={2} columnSpacing={2}>
      <Grid xs={12}>
        <NetworkStatus
          title={
            network.name
              ? `${network.name} is created.`
              : `No network selected.`
          }
          subtitle={network.name ? 'No node attached to this network.' : ''}
          loading={false}
          availableNodes={undefined}
          statusType="ONLINE"
          tooltipInfo="Network is online"
        />
      </Grid>
      <Grid xs={12}>
        <Stack direction={'row'}>
          <StatusCard
            option={''}
            subtitle2={''}
            Icon={UsersWithBG}
            subtitle1={`${0}`}
            options={TIME_FILTER}
            loading={false}
            title={'Active subscribers'}
            handleSelect={(value: string) => {}}
          />
          <Box p={1} />
          <StatusCard
            Icon={DataVolume}
            option={'usage'}
            subtitle2={`dBM`}
            subtitle1={`${0}`}
            options={TIME_FILTER}
            loading={false}
            title={'Data Volume'}
            handleSelect={(value: string) => {}}
          />
          <Box p={1} />
          <StatusCard
            option={'bill'}
            subtitle2={`bps`}
            subtitle1={`${0}`}
            Icon={Throughput}
            options={MONTH_FILTER}
            loading={false}
            title={'Average throughput'}
            handleSelect={(value: string) => {}}
          />
        </Stack>
      </Grid>
      <Grid xs={12}>
        <Paper sx={{ borderRadius: '10px', height: 'calc(100vh - 332px)' }}>
          {network.id ? (
            <LoadingWrapper radius="small" width={'100%'} isLoading={false}>
              <NetworkMap
                id="network-map"
                zoom={10}
                className="network-map"
                markersData={{ nodes: [], networkId: '' }}
              >
                {() => (
                  <>
                    <LabelOverlayUI name={network.name} />
                    <SitesTree sites={[]} />
                  </>
                )}
              </NetworkMap>
            </LoadingWrapper>
          ) : (
            <EmptyView
              title="No network selected"
              icon={NetworkIcon}
              size="medium"
            />
          )}
        </Paper>
        <OnboardingCard
          open={isOnboardingOpen}
          onClose={handleOnboardingClose}
          onStepClick={handleStepClick}
        />
        {renderInstallationDialog()}
      </Grid>
    </Grid>
  );
}
