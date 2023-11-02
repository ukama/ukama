/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { commonData, snackbarMessage } from '@/app-recoil';
import { SUBSCRIBER_TABLE_COLUMNS, SUBSCRIBER_TABLE_MENU } from '@/constants';
import {
  SubscribersResDto,
  useGetPackagesQuery,
  useGetSubscribersByNetworkQuery,
  useAddSubscriberMutation,
  useGetSimLazyQuery,
  useGetSimsQuery,
  useSetActivePackageForSimMutation,
  useToggleSimStatusMutation,
  PackagesResDto,
  useAllocateSimMutation,
} from '@/generated';
import {
  ContainerMax,
  PageContainer,
  VerticalContainer,
} from '@/styles/global';
import { colors } from '@/styles/theme';
import { TCommonData, TSnackMessage } from '@/types';
import AddSubscriberDialog from '@/ui/molecules/AddSubscriber';
import DataTableWithOptions from '@/ui/molecules/DataTableWithOptions';
import EmptyView from '@/ui/molecules/EmptyView';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import PlanCard from '@/ui/molecules/PlanCard';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import { AlertColor, Grid, Stack, Typography } from '@mui/material';
import { useCallback, useEffect, useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';

const Page = () => {
  const [search, setSearch] = useState<string>('');
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const [openAddSubscriber, setOpenAddSubscriber] = useState<boolean>(false);
  const [simId, setSimId] = useState<string>('');
  const [subscriberSuccess, setSubscriberSuccess] = useState<boolean>(false);
  const [qrCode, setQrCode] = useState<string>('');
  const [activePackageId, setActivePackageId] = useState<string>('');
  const [simPlan, setSimPlan] = useState<string>('');
  const [iccid, setIccid] = useState<string>('');
  const [simList, setSimList] = useState<any>({
    sim: [],
  });
  const [packages, setPackages] = useState<PackagesResDto>({
    packages: [],
  });
  const [subscriber, setSubscriber] = useState<SubscribersResDto>({
    subscribers: [],
  });
  const { loading: packagesLoading, data: _packages } = useGetPackagesQuery({
    fetchPolicy: 'cache-first',
    onCompleted: (_packages) => {
      if (_packages.getPackages.packages.length > 0) {
        setPackages(() => ({
          packages: [..._packages.getPackages.packages],
        }));
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'packages-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  const { data: dataPlanData, loading: dataPlanLoading } = useGetPackagesQuery({
    fetchPolicy: 'cache-and-network',
    onError: (error) => {
      setSnackbarMessage({
        id: 'data-plan-err-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  const { data: _sims } = useGetSimsQuery({
    variables: { type: 'test' },

    fetchPolicy: 'cache-first',
    onCompleted: (_sims) => {
      if (
        _sims &&
        _sims.getSims &&
        _sims.getSims.sim &&
        _sims.getSims.sim.length > 0
      ) {
        const simsArray = _sims.getSims.sim || [];

        setSimList(() => ({
          sims: [...simsArray],
        }));
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'sims-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  useEffect(() => {
    if (search.length > 3) {
      const subscribers = data?.getSubscribersByNetwork.subscribers.filter(
        (subscriber) => {
          const s = search.toLowerCase();
          if (
            subscriber.firstName.toLowerCase().includes(s) ||
            subscriber.lastName.toLowerCase().includes(s)
          )
            return subscriber;
        },
      );
      setSubscriber(() => ({ subscribers: subscribers ?? [] }));
    } else if (search.length === 0) {
      setSubscriber(() => ({
        subscribers: data?.getSubscribersByNetwork.subscribers ?? [],
      }));
    }
  }, [search]);

  const onTableMenuItem = (id: string, type: string) => {};

  const [activatePackageSim, { loading: activatePackageSimLoading }] =
    useSetActivePackageForSimMutation({
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
  const {
    loading,
    data,
    refetch: refetchSubscribers,
  } = useGetSubscribersByNetworkQuery({
    variables: { networkId: _commonData.networkId },
    fetchPolicy: 'cache-first',
    onCompleted: (data) => {
      if (data.getSubscribersByNetwork.subscribers.length > 0) {
        console.log('SUBSCRIBER :', data);
        setSubscriber(() => ({
          subscribers: [...data.getSubscribersByNetwork.subscribers],
        }));
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'subscriber-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const structureData = useCallback(
    (data: SubscribersResDto) =>
      data.subscribers.map((subscriber) => ({
        id: subscriber.uuid,
        email: subscriber.email,
        name: `${subscriber.firstName} ${subscriber.lastName}`,
        dataUsage: '',
        dataPlan: '',
        actions: '',
      })),
    [],
  );
  const handleAddSubscriberModal = () => {
    setOpenAddSubscriber(true);
  };
  const OnCloseAddSubcriber = () => {
    setOpenAddSubscriber(false);
  };
  const [getSim] = useGetSimLazyQuery({
    onCompleted: (res) => {
      if (res.getSim) {
        console.log('SIM', res.getSim);
      }
    },
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
        activatePackageSim({
          variables: {
            data: {
              simId: simId,
              packageId: activePackageId,
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

  const [allocateSim, { loading: allocateSimLoading }] = useAllocateSimMutation(
    {
      onCompleted: (res) => {
        setSnackbarMessage({
          id: 'sim-allocated-success',
          message: 'Sim allocated successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
        setQrCode(res.allocateSim.iccid);
        setSimId(res.allocateSim.id);
        getSim({
          variables: {
            data: {
              simId: res.allocateSim.id,
            },
          },
        });

        // setActivePackageId(res.allocateSim.packageId);
        // setActivePackageId(res.allocateSim.packageId);

        toggleSimStatus({
          variables: {
            data: {
              simId: res.allocateSim.id ?? '',
              status: 'active',
            },
          },
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'sim-allocated-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    },
  );
  const [addSubscriber, { loading: addSubscriberLoading }] =
    useAddSubscriberMutation({
      onCompleted: (res) => {
        console.log('SUBSCRIBER ADDED', res)
        refetchSubscribers();
        setSnackbarMessage({
          id: 'add-subscriber-success',
          message: 'Subscriber added successfully!',
          type: 'success' as AlertColor,
          show: true,
        });

        allocateSim({
          variables: {
            data: {
              networkId: _commonData.networkId,
              packageId: simPlan ?? '',
              subscriberId: res.addSubscriber.uuid,
              simType: 'test',
              iccid: iccid,
            },
          },
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'add-subscriber-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });
  const handleRoamingInstallation = async (values: any) => {
    const { plan, simIccid, email, name } = values;
    setSimPlan(plan);
    setIccid(simIccid);

    await addSubscriber({
      variables: {
        data: {
          email: email as string,
          first_name: name as string,
          last_name: 'pacifique',
          network_id: '9fa532e5-cdb5-4ba9-be7b-4375f42e4ac3',
          org_id: '8c6c2bec-5f90-4fee-8ffd-ee6456abf4fc',
          address: 'test',
          phone: '+1 555-123-4567',
          dob: '1990-01-01T00:00:00Z',
          gender: 'make',
          id_serial: '9802830928',
          proof_of_identification: 'passwport',
        },
      },
    });
  };
  return (
    <Stack direction={'column'}>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={loading}
        cstyle={{
          backgroundColor: loading ? colors.white : 'transparent',
        }}
      >
        <PageContainer
          sx={{ height: 'fit-content', maxHeight: 'calc(100vh - 400px)' }}
        >
          <PageContainerHeader
            title={'My subscribers'}
            subtitle={`${subscriber.subscribers.length}`}
            buttonTitle={'Add Subscriber'}
            handleButtonAction={handleAddSubscriberModal}
            onSearchChange={(e: string) => setSearch(e)}
            search={search}
          />
          <VerticalContainer>
            <ContainerMax mt={4.5}>
              <DataTableWithOptions
                icon={SubscriberIcon}
                isRowClickable={false}
                columns={SUBSCRIBER_TABLE_COLUMNS}
                dataset={structureData(subscriber)}
                menuOptions={SUBSCRIBER_TABLE_MENU}
                onMenuItemClick={onTableMenuItem}
                emptyViewLabel={'No subscribers yet!'}
              />
            </ContainerMax>
          </VerticalContainer>
          <AddSubscriberDialog
            onSuccess={false}
            open={openAddSubscriber}
            handleRoamingInstallation={handleRoamingInstallation}
            onClose={OnCloseAddSubcriber}
            qrCode={qrCode}
            submitButtonState={
              addSubscriberLoading || allocateSimLoading || packagesLoading
              // simsLoading
            }
            sims={
              simList?.sims?.filter(
                (sim: { isPhysical: string }) => sim.isPhysical === 'true',
              ) || []
            }
            pkgList={packages.packages}
            loading={
              addSubscriberLoading ||
              allocateSimLoading ||
              packagesLoading ||
              toggleSimStatusLoading
              // activatePackageSimLoading
            }
            pSimCount={1}
            eSimCount={0}
          />
        </PageContainer>
      </LoadingWrapper>
      {/* <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={dataPlanLoading}
        cstyle={{
          backgroundColor: dataPlanLoading ? colors.white : 'transparent',
        }}
      >
        <PageContainer
          sx={{ height: 'fit-content', maxHeight: 'calc(100vh - 550px)' }}
        >
          <Stack direction={'row'} alignItems={'center'}>
            <Typography variant="h6" mr={1}>
              Data plans
            </Typography>
            <Typography variant="subtitle2">
              <i>(view only)</i>
            </Typography>
          </Stack>
          <Stack my={4}>
            {dataPlanData?.getPackages &&
            dataPlanData?.getPackages?.packages?.length > 0 ? (
              <Grid container rowSpacing={2} columnSpacing={2}>
                {dataPlanData?.getPackages?.packages.map(
                  ({
                    uuid,
                    name,
                    duration,
                    users,
                    currency,
                    dataVolume,
                    dataUnit,
                    amount,
                  }: any) => (
                    <Grid item xs={12} sm={6} md={4} key={uuid}>
                      <PlanCard
                        uuid={uuid}
                        name={name}
                        users={users}
                        amount={amount}
                        dataUnit={dataUnit}
                        duration={duration}
                        currency={currency}
                        dataVolume={dataVolume}
                        isOptions={false}
                      />
                    </Grid>
                  ),
                )}
              </Grid>
            ) : (
              <EmptyView
                size="medium"
                title={
                  'No data plans yet! Go to “Manage data plans” in your organization settings to add one'
                }
                icon={SubscriberIcon}
              />
            )}
          </Stack>
          <AddSubscriberDialog
            onSuccess={false}
            open={openAddSubscriber}
            handleRoamingInstallation={handleRoamingInstallation}
            onClose={OnCloseAddSubcriber}
            qrCode={''}
            submitButtonState={true}
            sims={[]}
            pkgList={[]}
            loading={false}
            pSimCount={1}
            eSimCount={0}
          />
        </PageContainer>
      </LoadingWrapper> */}
    </Stack>
  );
};

export default Page;
