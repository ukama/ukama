/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  AllocateSimApiDto,
  Sim_Status,
  Sim_Types,
  SubscribersResDto,
  useAddSubscriberMutation,
  useAllocateSimMutation,
  useDeleteSubscriberMutation,
  useGetNetworksQuery,
  useGetPackagesForSimLazyQuery,
  useGetPackagesQuery,
  useGetSimLazyQuery,
  useGetSimPoolStatsQuery,
  useGetSimsBySubscriberLazyQuery,
  useGetSimsQuery,
  useAddPackagesToSimMutation,
  useGetSubscriberLazyQuery,
  useGetSubscribersByNetworkQuery,
  useToggleSimStatusMutation,
  useUpdateSubscriberMutation,
} from '@/client/graphql/generated';
import AddSubscriberStepperDialog from '@/components/AddSubscriber';
import DataTableWithOptions from '@/components/DataTableWithOptions';
import DeleteConfirmation from '@/components/DeleteDialog';
import LoadingWrapper from '@/components/LoadingWrapper';
import PageContainerHeader from '@/components/PageContainerHeader';
import PlanCard from '@/components/PlanCard';
import SubscriberDetails from '@/components/SubscriberDetails';
import TopUpData from '@/components/TopUpData';
import { SUBSCRIBER_TABLE_COLUMNS, SUBSCRIBER_TABLE_MENU } from '@/constants';
import { useAppContext } from '@/context';
import {
  CardWrapper,
  DataPlanEmptyView,
  NavigationButton,
  NavigationWrapper,
  ScrollableContent,
  ScrollContainer,
} from '@/styles/global';
import colors from '@/theme/colors';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';
import { AlertColor, Box, Paper, Stack, Typography } from '@mui/material';
import { useCallback, useEffect, useRef, useState } from 'react';

const Page = () => {
  const [search, setSearch] = useState<string>('');
  const { setSnackbarMessage, network, env } = useAppContext();
  const [openAddSubscriber, setOpenAddSubscriber] = useState(false);
  const [isToPupData, setIsToPupData] = useState<boolean>(false);
  const [isPackageActivationNeeded, setIsPackageActivationNeeded] =
    useState(false);
  const [subscriberDetails, setSubscriberDetails] = useState<any>();
  const [isSubscriberDetailsOpen, setIsSubscriberDetailsOpen] =
    useState<boolean>(false);
  const [subscriberSimList, setSubscriberSimList] = useState<any[]>();
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);
  const [deletedSubscriber, setDeletedSubscriber] = useState<string>('');
  const [selectedNetwork, setSelectedNetwork] = useState<string | null>(null);
  const scrollContainerRef = useRef<HTMLDivElement | null>(null);
  const [subscriber, setSubscriber] = useState<SubscribersResDto>({
    subscribers: [],
  });

  useEffect(() => {
    setSelectedNetwork(null);
    setSubscriber({ subscribers: [] });
  }, [network.id]);

  const getActiveNetworkId = (): string => {
    if (selectedNetwork) {
      return selectedNetwork;
    }
    return network.id;
  };

  const { data: packagesData, loading: packagesLoading } = useGetPackagesQuery({
    fetchPolicy: 'cache-and-network',
    onError: (error) => {
      setSnackbarMessage({
        id: 'packages',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: simPoolData, refetch: refetchSims } = useGetSimsQuery({
    variables: {
      data: {
        status: Sim_Status.Unassigned,
        type: env.SIM_TYPE as Sim_Types,
      },
    },
    fetchPolicy: 'network-only',
    onError: (error) => {
      setSnackbarMessage({
        id: 'sims-error-msg',
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
          if (subscriber.name.toLowerCase().includes(s)) return subscriber;
        },
      );
      setSubscriber(() => ({ subscribers: subscribers ?? [] }));
    } else if (search.length === 0) {
      setSubscriber(() => ({
        subscribers: data?.getSubscribersByNetwork.subscribers ?? [],
      }));
    }
  }, [search]);

  const [getSimBySubscriber] = useGetSimsBySubscriberLazyQuery({
    onCompleted: (res) => {
      if (res.getSimsBySubscriber) {
        setSubscriberSimList(res.getSimsBySubscriber.sims);
      }
    },
  });
  const [getPackagesForSim] = useGetPackagesForSimLazyQuery({
    onCompleted: (res) => {
      if (
        res.getPackagesForSim.packages &&
        res.getPackagesForSim.packages.length > 0 &&
        isPackageActivationNeeded
      ) {
      }
    },
  });

  const fetchPackagesForSim = useCallback(
    (simId: string) => {
      getPackagesForSim({
        variables: {
          data: {
            sim_id: simId,
          },
        },
      });
    },
    [getPackagesForSim],
  );

  const [getSubscriber] = useGetSubscriberLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (res) => {
      if (res?.getSubscriber?.sim && res.getSubscriber?.sim.length > 0) {
        fetchPackagesForSim(res.getSubscriber.sim[0].id);
      }
    },
  });

  const onTableMenuItem = async (id: string, type: string) => {
    if (type === 'delete-sub') {
      setIsConfirmationOpen(true);
      setDeletedSubscriber(id);
    }
    if (type === 'top-up-data') {
      setIsToPupData(true);
      getSubscriber({
        variables: {
          subscriberId: id,
        },
      });
      getSimBySubscriber({
        variables: {
          data: {
            subscriberId: id,
          },
        },
      });
    }
    if (type === 'edit-sub') {
      handleOpenSubscriberDetails(id);
    }
  };

  const {
    data,
    loading: getSubscriberByNetworkLoading,
    refetch: refetchSubscribers,
  } = useGetSubscribersByNetworkQuery({
    skip: !getActiveNetworkId(),
    variables: {
      networkId: getActiveNetworkId(),
    },
    fetchPolicy: 'network-only',
    nextFetchPolicy: 'network-only',
    onCompleted: (data) => {
      setSubscriber({
        subscribers: [...data.getSubscribersByNetwork.subscribers],
      });
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
  const { data: networkList } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',

    onError: (error) => {
      setSnackbarMessage({
        id: 'networks-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  // const { data: sitesData, loading: sitesLoading } = useGetSitesQuery({
  //   variables: {
  //     networkId: network.id,
  //   },
  //   fetchPolicy: 'cache-and-network',

  //   onError: (error) => {
  //     setSnackbarMessage({
  //       id: 'sites-msg',
  //       message: error.message,
  //       type: 'error' as AlertColor,
  //       show: true,
  //     });
  //   },
  // });

  const structureData = useCallback(
    (data: SubscribersResDto) => {
      if (
        (packagesData?.getPackages.packages?.length ?? 0) > 0 &&
        (networkList?.getNetworks?.networks?.length ?? 0) > 0
      ) {
        return data.subscribers.map((subscriber) => {
          const networkName =
            networkList?.getNetworks?.networks.find(
              (net) => net.id === subscriber.networkId,
            )?.name ?? '';
          const sim =
            subscriber?.sim && subscriber.sim?.length > 0
              ? subscriber?.sim[0]
              : null;
          const pkg = packagesData?.getPackages.packages.find(
            (pkg) => pkg.uuid === sim?.package?.package_id,
          );

          return {
            id: subscriber.uuid,
            network: networkName,
            name: subscriber.name,
            dataPlan: pkg?.name ?? '',
            email: subscriber.email,
            dataUsage: '',
            actions: '',
          };
        });
      }
    },
    [packagesData?.getPackages.packages, networkList?.getNetworks?.networks],
  );
  const handleOpenSubscriberDetails = useCallback(
    (id: string) => {
      const subscriberInfo = data?.getSubscribersByNetwork.subscribers.find(
        (subscriber) => subscriber.uuid === id,
      );
      setIsSubscriberDetailsOpen(true);
      if (subscriberInfo) {
        setSubscriberDetails(subscriberInfo);
      }
    },
    [data?.getSubscribersByNetwork.subscribers],
  );
  const [getSim] = useGetSimLazyQuery({
    onCompleted: (res) => {},
  });

  const [toggleSimStatus, { loading: toggleSimStatusLoading }] =
    useToggleSimStatusMutation({
      onCompleted: () => {
        setSnackbarMessage({
          id: 'sim-activated-success',
          message: 'Sim state updated successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
        refetchSubscribers();
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

  const handleDeleteSubscriber = () => {
    deleteSubscriber({
      variables: {
        subscriberId: deletedSubscriber,
      },
    });
  };

  const [addPackagesToSim, { loading: addPackagesToSimLoading }] =
    useAddPackagesToSimMutation({
      onCompleted: () => {
        setSnackbarMessage({
          id: 'packages-added-success',
          message: 'Package added successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'packages-added-error',
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
          id: 'allocate-sim-success',
          message: 'SIM allocated successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'allocate-sim-error',
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
        refetchSubscribers().then((data) => {
          setSubscriber(() => ({
            subscribers: [...data.data.getSubscribersByNetwork.subscribers],
          }));
        });
        setSnackbarMessage({
          id: 'add-subscriber-success',
          message: 'Subscriber added successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
        refetchSims();
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

  const { data: simStatData, refetch: refetchSimPoolStats } =
    useGetSimPoolStatsQuery({
      variables: {
        data: {
          type: env.SIM_TYPE as Sim_Types,
        },
      },
      fetchPolicy: 'cache-and-network',
      onError: (error) => {
        setSnackbarMessage({
          id: 'sim-stats-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const [deleteSubscriber, { loading: deleteSubscriberLoading }] =
    useDeleteSubscriberMutation({
      onCompleted: () => {
        refetchSubscribers();
        setSnackbarMessage({
          id: 'delete-subscriber-success',
          message: 'Subscriber deleted successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
        setIsConfirmationOpen(false);
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'delete-subscriber-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const [updateSubscriber, { loading: updateSubscriberLoading }] =
    useUpdateSubscriberMutation({
      onCompleted: (res) => {
        refetchSubscribers().then((data) => {
          setSubscriber(() => ({
            subscribers: [...data.data.getSubscribersByNetwork.subscribers],
          }));
        });
        setSnackbarMessage({
          id: 'update-subscriber-success',
          message: 'Subscriber updated successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'update-subscriber-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const handleAddSubscriberModal = () => {
    setOpenAddSubscriber(true);
    refetchSims();
  };

  const getSelectedNetwork = (network: string) => {
    setSelectedNetwork(network);
  };

  const handleCloseSubscriberDetails = () => {
    setIsSubscriberDetailsOpen(false);
  };

  const handleCancel = () => {
    setIsConfirmationOpen(false);
  };

  const handleSimAction = (action: string, simId: string) => {
    switch (action) {
      case 'deactivateSim':
      case 'activateSim':
        toggleSimStatus({
          variables: {
            data: {
              sim_id: simId,
              status: action === 'deactivateSim' ? 'inactive' : 'active',
            },
          },
        });
        break;
      case 'topUp':
        setIsToPupData(true);
        break;
      default:
        break;
    }
  };
  const handleCloseTopUp = () => {
    setIsToPupData(false);
  };

  const handleTopUp = async (simId: string, planIds: string[]) => {
    try {
      const packages = planIds.map((planId) => ({
        package_id: planId,
        start_date: new Date(Date.now() + 1 * 60000).toISOString(),
      }));

      await addPackagesToSim({
        variables: {
          data: {
            sim_id: simId,
            packages: packages,
          },
        },
      });

      setIsToPupData(false);
    } catch (error) {
      console.error('Error handling top up:', error);
      setIsToPupData(false);
      throw error;
    }
  };
  const handleUpdateSubscriber = (
    subscriberId: string,
    updates: { name?: string; phone?: string },
  ) => {
    updateSubscriber({
      variables: {
        subscriberId: subscriberId,
        data: updates,
      },
    });
    getSubscriber({
      variables: {
        subscriberId: subscriberId,
      },
    });
    refetchSubscribers();
  };

  const handleSubscriberMenuAction = async (
    action: string,
    subscriberId: string,
  ) => {
    if (action === 'deleteSubscriber') {
      deleteSubscriber({
        variables: {
          subscriberId: subscriberId,
        },
      });
    }
  };

  const handleAddSubscriber = async (
    subscriber: any,
  ): Promise<AllocateSimApiDto> => {
    try {
      setSubscriberDetails(subscriber);

      const subscriberResponse = await addSubscriber({
        variables: {
          data: {
            name: subscriber.name,
            network_id: network.id,
            email: subscriber.email,
            phone: subscriber.phone,
          },
        },
      });

      if (!subscriberResponse.data) {
        throw new Error('Failed to add subscriber');
      }

      const simResponse = await allocateSim({
        variables: {
          data: {
            network_id: subscriberResponse.data.addSubscriber.networkId,
            package_id: subscriber.plan,
            subscriber_id: subscriberResponse.data.addSubscriber.uuid,
            sim_type: env.SIM_TYPE,
            iccid: subscriber.simIccid,
            traffic_policy: 0,
          },
        },
      });

      if (!simResponse.data) {
        throw new Error('Failed to allocate SIM');
      }

      return simResponse.data.allocateSim;
    } catch (error) {
      throw error;
    }
  };

  const scroll = (direction: 'left' | 'right'): void => {
    if (scrollContainerRef.current) {
      const scrollAmount = scrollContainerRef.current.clientWidth / 2;
      scrollContainerRef.current.scrollLeft +=
        direction === 'left' ? -scrollAmount : scrollAmount;
    }
  };

  return (
    <Stack
      mt={2}
      spacing={2}
      direction={'column'}
      sx={{ height: { xs: 'calc(100vh - 158px)', md: 'calc(100vh - 172px)' } }}
    >
      {getSubscriberByNetworkLoading ? (
        <LoadingWrapper
          radius="small"
          width={'100%'}
          isLoading={true}
          cstyle={{
            height: '240px',
          }}
        >
          <br />
        </LoadingWrapper>
      ) : (
        <Paper
          elevation={1}
          sx={{
            borderRadius: '10px',
            p: { xs: 2, md: 4 },
          }}
        >
          <Stack direction="column" spacing={{ xs: 0.5, md: 1.5 }}>
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              <Stack direction={'row'} spacing={1} alignItems={'center'}>
                <Typography variant="h6">Data plans</Typography>
                <Typography variant="subtitle2" sx={{ color: colors.black38 }}>
                  ({packagesData?.getPackages.packages?.length ?? 0})
                </Typography>
              </Stack>

              {packagesData &&
                packagesData?.getPackages.packages?.length > 4 && (
                  <NavigationWrapper>
                    <NavigationButton
                      onClick={() => scroll('left')}
                      disabled={!packagesData?.getPackages.packages?.length}
                    >
                      <KeyboardArrowLeftIcon fontSize="small" />
                    </NavigationButton>

                    <NavigationButton
                      onClick={() => scroll('right')}
                      disabled={!packagesData?.getPackages.packages?.length}
                    >
                      <KeyboardArrowRightIcon fontSize="small" />
                    </NavigationButton>
                  </NavigationWrapper>
                )}
            </Box>

            <ScrollContainer>
              <ScrollableContent ref={scrollContainerRef}>
                {!packagesData?.getPackages.packages?.length ? (
                  <DataPlanEmptyView>
                    <UpdateIcon sx={{ fontSize: 40, mb: 1 }} />
                    <Typography variant="body1">
                      No data plan created yet!
                    </Typography>
                  </DataPlanEmptyView>
                ) : (
                  packagesData.getPackages.packages.map(
                    ({
                      uuid,
                      name,
                      duration,
                      currency,
                      dataVolume,
                      dataUnit,
                      amount,
                    }: any) => (
                      <CardWrapper key={uuid}>
                        <PlanCard
                          uuid={uuid}
                          name={name}
                          amount={amount}
                          isOptions={false}
                          dataUnit={dataUnit}
                          duration={duration}
                          currency={currency}
                          dataVolume={dataVolume}
                        />
                      </CardWrapper>
                    ),
                  )
                )}
              </ScrollableContent>
            </ScrollContainer>
          </Stack>
        </Paper>
      )}
      {getSubscriberByNetworkLoading ? (
        <LoadingWrapper
          radius="small"
          width={'100%'}
          isLoading={true}
          cstyle={{
            height: '100%',
          }}
        >
          <br />
        </LoadingWrapper>
      ) : (
        <Paper
          sx={{
            height: '100%',
            borderRadius: '10px',
            px: { xs: 2, md: 3 },
            py: { xs: 2, md: 4 },
          }}
        >
          <PageContainerHeader
            title={'My subscribers'}
            subtitle={`${subscriber.subscribers.length}`}
            buttonTitle={'Add Subscriber'}
            handleButtonAction={handleAddSubscriberModal}
            onSearchChange={(e: string) => setSearch(e)}
            search={search}
          />
          <br />
          <DataTableWithOptions
            icon={SubscriberIcon}
            isRowClickable={false}
            columns={SUBSCRIBER_TABLE_COLUMNS}
            dataset={structureData(subscriber)}
            menuOptions={SUBSCRIBER_TABLE_MENU}
            onMenuItemClick={onTableMenuItem}
            emptyViewLabel={'No subscribers yet!'}
            getSelectedNetwork={getSelectedNetwork}
            networkList={networkList?.getNetworks?.networks ?? []}
          />
        </Paper>
      )}
      <AddSubscriberStepperDialog
        isOpen={openAddSubscriber}
        handleCloseAction={() => setOpenAddSubscriber(false)}
        handleAddSubscriber={handleAddSubscriber}
        sims={simPoolData?.getSims.sim ?? []}
        packages={packagesData?.getPackages.packages ?? []}
        isLoading={addSubscriberLoading || allocateSimLoading}
      />

      <DeleteConfirmation
        open={isConfirmationOpen}
        onDelete={handleDeleteSubscriber}
        onCancel={handleCancel}
        itemName={deletedSubscriber}
        loading={deleteSubscriberLoading}
      />
      <SubscriberDetails
        ishowSubscriberDetails={isSubscriberDetailsOpen}
        handleClose={handleCloseSubscriberDetails}
        subscriberInfo={subscriberDetails}
        handleSimActionOption={handleSimAction}
        handleUpdateSubscriber={handleUpdateSubscriber}
        loading={updateSubscriberLoading}
        handleDeleteSubscriber={handleSubscriberMenuAction}
        simStatusLoading={toggleSimStatusLoading}
        currentSite={'-'}
      />
      <TopUpData
        isToPup={isToPupData}
        onCancel={handleCloseTopUp}
        handleTopUp={handleTopUp}
        loadingTopUp={packagesLoading || addPackagesToSimLoading}
        packages={packagesData?.getPackages.packages ?? []}
        sims={subscriberSimList ?? []}
        subscriberName={subscriberDetails?.name}
      />
    </Stack>
  );
};

export default Page;
