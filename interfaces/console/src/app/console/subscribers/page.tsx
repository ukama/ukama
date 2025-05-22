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
  useAddPackagesToSimMutation,
  useAddSubscriberMutation,
  useAllocateSimMutation,
  useDeleteSubscriberMutation,
  useGetCurrencySymbolQuery,
  useGetDataUsagesLazyQuery,
  useGetPackagesQuery,
  useGetSimsBySubscriberLazyQuery,
  useGetPackagesForSimLazyQuery,
  useGetSimsFromPoolQuery,
  useGetSubscribersByNetworkQuery,
  useToggleSimStatusMutation,
  useUpdateSubscriberMutation,
  useDeleteSimMutation,
} from '@/client/graphql/generated';
import AddSubscriberStepperDialog from '@/components/AddSubscriber';
import DataTableWithOptions from '@/components/DataTableWithOptions';
import DeleteConfirmation from '@/components/DeleteDialog';
import LoadingWrapper from '@/components/LoadingWrapper';
import PageContainerHeader from '@/components/PageContainerHeader';
import PlanCard from '@/components/PlanCard';
import SubscriberDetailsDialog from '@/components/SubscriberDetailsDialog';
import TopUpData from '@/components/TopUpData';
import {
  SUBSCRIBER_ERROR_MESSAGES,
  SUBSCRIBER_TABLE_COLUMNS,
  SUBSCRIBER_TABLE_MENU,
} from '@/constants';
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
import { formatBytesToGB, getDisplayStatus } from '@/utils';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';
import { AlertColor, Box, Paper, Stack, Typography } from '@mui/material';
import { useSearchParams } from 'next/navigation';
import { useCallback, useEffect, useRef, useState } from 'react';

const Page = () => {
  const query = useSearchParams();
  const [search, setSearch] = useState<string>('');
  const { setSnackbarMessage, network, env, user } = useAppContext();
  const [simPackageHistories, setSimPackageHistories] = useState<any[]>([]);
  const [openAddSubscriber, setOpenAddSubscriber] = useState(false);
  const [isTopupData, setIsTopupData] = useState<boolean>(false);
  const [subscriberDetails, setSubscriberDetails] = useState<any>();
  const [isSubscriberDetailsOpen, setIsSubscriberDetailsOpen] =
    useState<boolean>(false);
  const [subscriberSimList, setSubscriberSimList] = useState<any[]>();
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);
  const [dataUsageForSim, setDataUsageForSim] = useState<string>('');
  const [deletedSubscriber, setDeletedSubscriber] = useState<{
    id: string;
    name: string;
  }>({ id: '', name: '' });
  const [isSimDeleteConfirmationOpen, setIsSimDeleteConfirmationOpen] =
    useState(false);
  const [simToDelete, setSimToDelete] = useState<{
    id: string;
    iccid: string;
    isLastSim: boolean;
  }>({ id: '', iccid: '', isLastSim: false });
  const scrollContainerRef = useRef<HTMLDivElement | null>(null);
  const [topUpSubscriberName, setTopUpSubscriberName] = useState('');
  const [subscriber, setSubscriber] = useState<SubscribersResDto>({
    subscribers: [],
  });
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

  const { data: simPoolData, refetch: refetchSims } = useGetSimsFromPoolQuery({
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

  const [getSimBySubscriber] = useGetSimsBySubscriberLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (res) => {
      if (res.getSimsBySubscriber) {
        setSubscriberSimList(res.getSimsBySubscriber.sims);
      }
    },
  });

  const {
    data,
    loading: getSubscriberByNetworkLoading,
    refetch: refetchSubscribers,
  } = useGetSubscribersByNetworkQuery({
    skip: !network.id,
    variables: {
      networkId: network.id,
    },
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      setSubscriber({
        subscribers: [...data.getSubscribersByNetwork.subscribers],
      });
      if (query.size > 0) {
        const iccid = query.get('iccid');
        setSearch(iccid ?? '');
      }

      getDataUsages({
        variables: {
          data: {
            type: Sim_Types.UkamaData,
            networkId: network.id,
          },
        },
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

  const [deleteSim, { loading: deleteSimLoading }] = useDeleteSimMutation({
    onCompleted: () => {
      if (subscriberDetails) {
        getSimBySubscriber({
          variables: {
            data: {
              subscriberId: subscriberDetails.uuid,
            },
          },
        });
      }

      refetchSubscribers().then((response) => {
        if (response?.data?.getSubscribersByNetwork) {
          setSubscriber({
            subscribers: [...response.data.getSubscribersByNetwork.subscribers],
          });
        }
      });

      setSnackbarMessage({
        id: 'delete-sim-success',
        message: `SIM ${simToDelete.iccid} deletion initiated. ${
          simToDelete.isLastSim ? 'This will also delete the subscriber.' : ''
        }`,
        type: 'info' as AlertColor,
        show: true,
      });

      setIsSimDeleteConfirmationOpen(false);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'delete-sim-error',
        message: `Error deleting SIM: ${error.message}`,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  const [getPackagesForSim, { loading: packagesForSimLoading }] =
    useGetPackagesForSimLazyQuery({
      fetchPolicy: 'network-only',
      onCompleted: (res) => {
        if (res.getPackagesForSim && res.getPackagesForSim.packages) {
          setSimPackageHistories((prevHistories) => [
            ...prevHistories,
            ...res.getPackagesForSim.packages.map((pkg) => ({
              ...pkg,
              simId: res.getPackagesForSim.sim_id,
            })),
          ]);
        }
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'packages-history-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const [toggleSimStatus] = useToggleSimStatusMutation({
    onCompleted: () => {
      setSnackbarMessage({
        id: 'sim-activated-success',
        message: 'Sim state updated successfully!',
        type: 'success' as AlertColor,
        show: true,
      });

      refetchSubscribers().then((response) => {
        if (response?.data?.getSubscribersByNetwork) {
          setSubscriber({
            subscribers: [...response.data.getSubscribersByNetwork.subscribers],
          });
          if (subscriberDetails) {
            getSimBySubscriber({
              variables: {
                data: {
                  subscriberId: subscriberDetails.uuid,
                },
              },
              onCompleted: (res) => {
                if (res.getSimsBySubscriber && res.getSimsBySubscriber.sims) {
                  setSubscriberSimList(res.getSimsBySubscriber.sims);
                  const updatedSubscriberInfo =
                    response.data.getSubscribersByNetwork.subscribers.find(
                      (sub) => sub.uuid === subscriberDetails.uuid,
                    );

                  if (updatedSubscriberInfo) {
                    setSubscriberDetails({
                      ...updatedSubscriberInfo,
                      packageId:
                        updatedSubscriberInfo.sim?.[0]?.package?.package_id,
                      dataUsage: subscriberDetails.dataUsage,
                      dataPlan: subscriberDetails.dataPlan,
                      simIccid: updatedSubscriberInfo.sim?.[0]?.iccid,
                    });
                  }
                }
              },
            });
          }
        }
      });
      setIsSubscriberDetailsOpen(false);
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

  const [addPackagesToSim, { loading: addPackagesToSimLoading }] =
    useAddPackagesToSimMutation({
      onCompleted: () => {
        setSnackbarMessage({
          id: 'packages-added-success',
          message: 'Packages added successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
        if (subscriberSimList && subscriberSimList.length > 0) {
          setSimPackageHistories([]);
          getPackagesForSim({
            variables: { data: { sim_id: subscriberSimList[0].id } },
          });
        }
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
        const errorMsg =
          (error.graphQLErrors?.[0]?.extensions?.response as any)?.body
            ?.error || '';
        const userFriendlyMessage = errorMsg.includes(
          'idx_subscribers_active_email',
        )
          ? SUBSCRIBER_ERROR_MESSAGES.DUPLICATE_EMAIL
          : error.message;

        setSnackbarMessage({
          id: 'add-subscriber-error',
          message: userFriendlyMessage,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const [deleteSubscriber, { loading: deleteSubscriberLoading }] =
    useDeleteSubscriberMutation({
      onCompleted: () => {
        refetchSubscribers().then((response) => {
          if (response?.data?.getSubscribersByNetwork) {
            setSubscriber({
              subscribers: [
                ...response.data.getSubscribersByNetwork.subscribers,
              ],
            });
          }
        });
        setSnackbarMessage({
          id: 'delete-subscriber-success',
          message: `Deletion started for "${deletedSubscriber.name}". SIMs are being processed.`,
          type: 'info' as AlertColor,
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
  useEffect(() => {
    const handleSubscriberNotification = (_: any, data: string) => {
      try {
        const parsedData = JSON.parse(data);
        const { eventKey } = parsedData.data.notificationSubscription;

        if (eventKey === 'EventSubscriberDelete') {
          refetchSubscribers().then((response) => {
            if (response?.data?.getSubscribersByNetwork) {
              setSubscriber({
                subscribers: [
                  ...response.data.getSubscribersByNetwork.subscribers,
                ],
              });
            }
          });
        }
      } catch (error) {
        console.error('Error processing notification:', error);
      }
    };

    if (user.id && network.id && user.orgId) {
      const topic = `notification-${user.orgId}-${user.id}-${user.role}-${network.id}`;
      PubSub.subscribe(topic, handleSubscriberNotification);

      return () => {
        PubSub.unsubscribe(topic);
      };
    }
  }, [user.id, network.id, user.orgId, user.role, refetchSubscribers]);

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

  const { data: currencyData } = useGetCurrencySymbolQuery({
    skip: !user.currency,
    fetchPolicy: 'cache-first',
    variables: {
      code: user.currency,
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'currency-info-error',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [getDataUsages, { data: dataUsageData, loading: dataUsageLoading }] =
    useGetDataUsagesLazyQuery({
      pollInterval: 120000,
      fetchPolicy: 'network-only',
      variables: {
        data: {
          type: Sim_Types.UkamaData,
          networkId: network.id,
        },
      },
    });

  const handleDeleteSubscriber = () => {
    deleteSubscriber({
      variables: {
        subscriberId: deletedSubscriber.id,
      },
    });
  };

  const handleTopUpDataPreparation = (id: string) => {
    const subscriberInfo = data?.getSubscribersByNetwork.subscribers.find(
      (subscriber) => subscriber.uuid === id,
    );

    setIsTopupData(true);
    setIsSubscriberDetailsOpen(false);

    getSimBySubscriber({
      variables: {
        data: {
          subscriberId: id,
        },
      },
    });

    if (subscriberInfo) {
      setTopUpSubscriberName(subscriberInfo.name);
    }
  };

  const onTableMenuItem = async (id: string, type: string) => {
    switch (type) {
      case 'delete-sub':
        const subscriberToDelete =
          data?.getSubscribersByNetwork.subscribers.find(
            (subscriber) => subscriber.uuid === id,
          );
        setIsConfirmationOpen(true);
        setDeletedSubscriber({
          id: id,
          name: subscriberToDelete?.name ?? 'This subscriber',
        });
        break;

      case 'top-up-data':
        handleTopUpDataPreparation(id);
        break;

      case 'edit-sub':
        handleOpenSubscriberDetails(id);
        break;
    }
  };

  const structureData = useCallback(
    (data: SubscribersResDto) => {
      if ((packagesData?.getPackages.packages?.length ?? 0) > 0 && network) {
        return data.subscribers.map((subscriber) => {
          const sim =
            subscriber?.sim && subscriber.sim?.length > 0
              ? subscriber?.sim[0]
              : null;
          const pkg = packagesData?.getPackages.packages.find(
            (pkg) => pkg.uuid === sim?.package?.package_id,
          );
          const dataUsage = dataUsageData?.getDataUsages.usages.find(
            (usage) => usage.simId === sim?.id,
          );

          return {
            id: subscriber.uuid,
            name: subscriber.name,
            email: subscriber.email,
            packageId: sim?.package?.package_id,
            dataPlan: pkg?.name ?? 'No active plan',
            dataUsage: `${isNaN(Number(dataUsage?.usage)) ? 0 : formatBytesToGB(Number(dataUsage?.usage))} GB`,
            subscriberStatus: getDisplayStatus(subscriber.subscriberStatus),
            actions: '',
          };
        });
      }
    },
    [
      packagesData?.getPackages.packages,
      dataUsageData?.getDataUsages.usages,
      network,
    ],
  );

  const handleOpenSubscriberDetails = useCallback(
    (id: string) => {
      setSimPackageHistories([]);

      const subscriberInfo = data?.getSubscribersByNetwork.subscribers.find(
        (subscriber) => subscriber.uuid === id,
      );
      setIsSubscriberDetailsOpen(true);

      getSimBySubscriber({
        variables: {
          data: {
            subscriberId: id,
          },
        },
        onCompleted: (res) => {
          if (res.getSimsBySubscriber && res.getSimsBySubscriber.sims) {
            setSubscriberSimList(res.getSimsBySubscriber.sims);

            if (res.getSimsBySubscriber.sims.length > 0) {
              const firstSim = res.getSimsBySubscriber.sims[0];
              if (firstSim?.id) {
                getPackagesForSim({
                  variables: { data: { sim_id: firstSim.id } },
                });
              }
            }
          }
        },
      });
      getDataUsages({
        variables: {
          data: {
            type: Sim_Types.UkamaData,
            networkId: network.id,
          },
        },
      });
      setDataUsageForSim(
        dataUsageData?.getDataUsages.usages.find(
          (usage) => usage.simId === subscriberInfo?.sim?.[0]?.id,
        )?.usage ?? '',
      );
      if (subscriberInfo) {
        const usageData = dataUsageData?.getDataUsages.usages.find(
          (usage) => usage.simId === subscriberInfo.sim?.[0]?.id,
        );
        const plan = packagesData?.getPackages.packages.find(
          (pkg) => pkg.uuid === subscriberInfo.sim?.[0]?.package?.package_id,
        );

        setSubscriberDetails({
          ...subscriberInfo,
          packageId: subscriberInfo.sim?.[0]?.package?.package_id,
          dataUsage: `${formatBytesToGB(Number(usageData?.usage)) || 0} GB`,
          dataPlan: plan?.name ?? 'No active plan',
          simIccid: subscriberInfo.sim?.[0]?.iccid,
        });
      }
    },
    [
      data?.getSubscribersByNetwork.subscribers,
      dataUsageData?.getDataUsages.usages,
      packagesData?.getPackages.packages,
      getSimBySubscriber,
      getPackagesForSim,
    ],
  );

  const handleAddSubscriberModal = () => {
    setOpenAddSubscriber(true);
    refetchSims();
  };

  const handleCloseSubscriberDetails = () => {
    setIsSubscriberDetailsOpen(false);
  };

  const handleCancel = () => {
    setIsConfirmationOpen(false);
  };

  const handleSimAction = (
    action: string,
    simId: string,
    additionalData?: any,
  ) => {
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
      case 'deleteSim':
        setIsSubscriberDetailsOpen(false);
        setSimToDelete({
          id: simId,
          iccid: additionalData.iccid,
          isLastSim: additionalData.isLastSim,
        });
        setIsSimDeleteConfirmationOpen(true);
        break;
      case 'topUp':
        setIsSubscriberDetailsOpen(false);
        setIsTopupData(true);
        break;
      default:
        break;
    }
  };
  const handleCloseTopUp = () => {
    setIsTopupData(false);
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

      setIsTopupData(false);
    } catch (error) {
      console.error('Error handling top up:', error);
      setIsTopupData(false);
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
    refetchSubscribers();
  };

  const handleSubscriberMenuAction = async (
    action: string,
    subscriberId: string,
  ) => {
    if (action === 'deleteSubscriber') {
      const subscriberToDelete = data?.getSubscribersByNetwork.subscribers.find(
        (subscriber) => subscriber.uuid === subscriberId,
      );
      setIsConfirmationOpen(true);
      setDeletedSubscriber({
        id: subscriberId,
        name: subscriberToDelete?.name ?? 'This subscriber',
      });
      setIsSubscriberDetailsOpen(false);
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
  console.log('subscriber');

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
                          currency={
                            currencyData?.getCurrencySymbol.symbol ?? ''
                          }
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
      {getSubscriberByNetworkLoading || dataUsageLoading ? (
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
            overflow: 'hidden',
            borderRadius: '10px',
            px: { xs: 2, md: 3 },
            py: { xs: 2, md: 4 },
          }}
        >
          <PageContainerHeader
            search={search}
            title={'My subscribers'}
            buttonTitle={'Add Subscriber'}
            handleButtonAction={handleAddSubscriberModal}
            onSearchChange={(e: string) => setSearch(e)}
            subtitle={`${data?.getSubscribersByNetwork.subscribers.length}`}
          />
          <br />

          <DataTableWithOptions
            icon={SubscriberIcon}
            isRowClickable={false}
            columns={SUBSCRIBER_TABLE_COLUMNS}
            dataset={structureData(subscriber)}
            menuOptions={SUBSCRIBER_TABLE_MENU}
            onMenuItemClick={(id: string, type: string) => {
              const rows = structureData(subscriber) || [];
              const row = rows.find((r) => r.id === id);
              if (row && row.subscriberStatus === 'Deleting...') {
                if (type === 'delete-sub') {
                  setSnackbarMessage({
                    id: 'retry-deletion-info',
                    message: SUBSCRIBER_ERROR_MESSAGES.RETRY_DELETION,
                    type: 'info' as AlertColor,
                    show: true,
                  });
                  deleteSubscriber({
                    variables: { subscriberId: id },
                  });
                  return;
                } else {
                  setSnackbarMessage({
                    id: 'action-blocked',
                    message: SUBSCRIBER_ERROR_MESSAGES.ACTION_BLOCKED_DELETING,
                    type: 'warning' as AlertColor,
                    show: true,
                  });
                  return;
                }
              }
              onTableMenuItem(id, type);
            }}
            emptyViewLabel={'No subscribers yet!'}
          />
        </Paper>
      )}
      <AddSubscriberStepperDialog
        isOpen={openAddSubscriber}
        currencySymbol={currencyData?.getCurrencySymbol.symbol ?? ''}
        handleCloseAction={() => setOpenAddSubscriber(false)}
        handleAddSubscriber={handleAddSubscriber}
        sims={simPoolData?.getSimsFromPool.sims ?? []}
        packages={packagesData?.getPackages.packages ?? []}
        isLoading={addSubscriberLoading || allocateSimLoading}
      />
      <DeleteConfirmation
        open={isConfirmationOpen}
        onDelete={handleDeleteSubscriber}
        onCancel={handleCancel}
        itemName={deletedSubscriber.name}
        itemType="subscriber"
        loading={deleteSubscriberLoading}
      />
      <SubscriberDetailsDialog
        open={isSubscriberDetailsOpen}
        onClose={handleCloseSubscriberDetails}
        subscriber={{
          id: subscriberDetails?.uuid || '',
          firstName: subscriberDetails?.name || '',
          email: subscriberDetails?.email || '',
        }}
        onUpdateSubscriber={(updates) =>
          handleUpdateSubscriber(subscriberDetails?.uuid, updates)
        }
        onDeleteSubscriber={() =>
          handleSubscriberMenuAction(
            'deleteSubscriber',
            subscriberDetails?.uuid,
          )
        }
        onTopUpPlan={handleTopUpDataPreparation}
        sims={subscriberSimList ?? []}
        onSimAction={handleSimAction}
        packageHistories={simPackageHistories}
        packagesData={packagesData?.getPackages}
        loadingPackageHistories={packagesForSimLoading}
        dataUsage={dataUsageForSim}
        currencySymbol={currencyData?.getCurrencySymbol.symbol ?? ''}
      />
      <TopUpData
        isToPup={isTopupData}
        onCancel={handleCloseTopUp}
        handleTopUp={handleTopUp}
        loadingTopUp={packagesLoading || addPackagesToSimLoading}
        packages={packagesData?.getPackages.packages ?? []}
        sims={subscriberSimList ?? []}
        subscriberName={topUpSubscriberName}
      />
      <DeleteConfirmation
        open={isSimDeleteConfirmationOpen}
        onDelete={() => {
          deleteSim({
            variables: {
              data: {
                simId: simToDelete.id,
              },
            },
          });
        }}
        onCancel={() => setIsSimDeleteConfirmationOpen(false)}
        itemName={simToDelete.iccid}
        itemType="sim"
        loading={deleteSimLoading}
      />
    </Stack>
  );
};

export default Page;
