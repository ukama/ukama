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
import PubSub from 'pubsub-js';
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
import {
  DialogStates,
  INITIAL_OPERATION_DATA,
  INITIAL_SUBSCRIBER_DATA,
  INITIAL_SUBSCRIBER_DIALOG_STATES,
  INITIAL_SUBSCRIBER_UI_STATE,
  OperationData,
  SubscriberData,
  UIState,
} from '@/types';

const Page = () => {
  const query = useSearchParams();
  const { setSnackbarMessage, network, env, user } = useAppContext();
  const scrollContainerRef = useRef<HTMLDivElement | null>(null);

  const [dialogStates, setDialogStates] = useState<DialogStates>(
    INITIAL_SUBSCRIBER_DIALOG_STATES,
  );
  const [subscriberData, setSubscriberData] = useState<SubscriberData>(
    INITIAL_SUBSCRIBER_DATA,
  );
  const [operationData, setOperationData] = useState<OperationData>(
    INITIAL_OPERATION_DATA,
  );
  const [uiState, setUIState] = useState<UIState>(INITIAL_SUBSCRIBER_UI_STATE);

  const updateDialogState = useCallback((updates: Partial<DialogStates>) => {
    setDialogStates((prev) => ({ ...prev, ...updates }));
  }, []);

  const updateSubscriberData = useCallback(
    (updates: Partial<SubscriberData>) => {
      setSubscriberData((prev) => ({ ...prev, ...updates }));
    },
    [],
  );

  const updateOperationData = useCallback((updates: Partial<OperationData>) => {
    setOperationData((prev) => ({ ...prev, ...updates }));
  }, []);

  const updateUIState = useCallback((updates: Partial<UIState>) => {
    setUIState((prev) => ({ ...prev, ...updates }));
  }, []);

  const closeDialog = useCallback((dialogName: keyof DialogStates) => {
    setDialogStates((prev) => ({ ...prev, [dialogName]: false }));
  }, []);

  const openDialog = useCallback(
    (dialogName: keyof DialogStates, data?: any) => {
      setDialogStates((prev) => ({ ...prev, [dialogName]: true }));
      if (data && dialogName === 'subscriberDetails') {
        updateSubscriberData({ details: data });
      }
    },
    [updateSubscriberData],
  );

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
        updateSubscriberData({ simList: res.getSimsBySubscriber.sims });
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
      updateUIState({
        subscribers: {
          subscribers: [...data.getSubscribersByNetwork.subscribers],
        },
      });

      if (query.size > 0) {
        const iccid = query.get('iccid');
        updateUIState({ search: iccid ?? '' });
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
      if (subscriberData.details) {
        getSimBySubscriber({
          variables: {
            data: {
              subscriberId: subscriberData.details.uuid,
            },
          },
        });
      }

      refetchSubscribers().then((response) => {
        if (response?.data?.getSubscribersByNetwork) {
          updateUIState({
            subscribers: {
              subscribers: [
                ...response.data.getSubscribersByNetwork.subscribers,
              ],
            },
          });
        }
      });

      setSnackbarMessage({
        id: 'delete-sim-success',
        message: `SIM ${operationData.simToDelete.iccid} deletion initiated. ${
          operationData.simToDelete.isLastSim
            ? 'This will also delete the subscriber.'
            : ''
        }`,
        type: 'info' as AlertColor,
        show: true,
      });

      closeDialog('simDeleteConfirmation');
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
          updateSubscriberData({
            packageHistories: [
              ...subscriberData.packageHistories,
              ...res.getPackagesForSim.packages.map((pkg) => ({
                ...pkg,
                simId: res.getPackagesForSim.sim_id,
              })),
            ],
          });
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
      refetchSubscribers().then(() => {});

      setSnackbarMessage({
        id: 'sim-toggle-success',
        message: 'SIM status updated successfully',
        type: 'success',
        show: true,
      });

      refetchSubscribers().then((response) => {
        if (response?.data?.getSubscribersByNetwork) {
          updateUIState({
            subscribers: {
              subscribers: [
                ...response.data.getSubscribersByNetwork.subscribers,
              ],
            },
          });

          if (subscriberData.details) {
            getSimBySubscriber({
              variables: {
                data: {
                  subscriberId: subscriberData.details.uuid,
                },
              },
              onCompleted: (res) => {
                if (res.getSimsBySubscriber && res.getSimsBySubscriber.sims) {
                  updateSubscriberData({
                    simList: res.getSimsBySubscriber.sims,
                  });

                  const updatedSubscriberInfo =
                    response.data.getSubscribersByNetwork.subscribers.find(
                      (sub) => sub.uuid === subscriberData.details.uuid,
                    );

                  if (updatedSubscriberInfo) {
                    updateSubscriberData({
                      details: {
                        ...updatedSubscriberInfo,
                        packageId:
                          updatedSubscriberInfo.sim?.[0]?.package?.package_id,
                        dataUsage: subscriberData.details.dataUsage,
                        dataPlan: subscriberData.details.dataPlan,
                        simIccid: updatedSubscriberInfo.sim?.[0]?.iccid,
                      },
                    });
                  }
                }
              },
            });
          }
        }
      });
      closeDialog('subscriberDetails');
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

        if (subscriberData.simList && subscriberData.simList.length > 0) {
          updateSubscriberData({ packageHistories: [] });
          getPackagesForSim({
            variables: { data: { sim_id: subscriberData.simList[0].id } },
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
          updateUIState({
            subscribers: {
              subscribers: [...data.data.getSubscribersByNetwork.subscribers],
            },
          });
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
            updateUIState({
              subscribers: {
                subscribers: [
                  ...response.data.getSubscribersByNetwork.subscribers,
                ],
              },
            });
          }
        });
        setSnackbarMessage({
          id: 'delete-subscriber-success',
          message: `Deletion started for "${operationData.deletedSubscriber.name}". SIMs are being processed.`,
          type: 'info' as AlertColor,
          show: true,
        });
        closeDialog('confirmation');
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
              updateUIState({
                subscribers: {
                  subscribers: [
                    ...response.data.getSubscribersByNetwork.subscribers,
                  ],
                },
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
  }, [
    user.id,
    network.id,
    user.orgId,
    user.role,
    refetchSubscribers,
    updateUIState,
  ]);

  const [updateSubscriber, { loading: updateSubscriberLoading }] =
    useUpdateSubscriberMutation({
      onCompleted: (res) => {
        refetchSubscribers().then((data) => {
          updateUIState({
            subscribers: {
              subscribers: [...data.data.getSubscribersByNetwork.subscribers],
            },
          });
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
        subscriberId: operationData.deletedSubscriber.id,
      },
    });
  };

  const handleTopUpDataPreparation = (id: string) => {
    const subscriberInfo = data?.getSubscribersByNetwork.subscribers.find(
      (subscriber) => subscriber.uuid === id,
    );

    updateDialogState({
      topupData: true,
      subscriberDetails: false,
    });

    getSimBySubscriber({
      variables: {
        data: {
          subscriberId: id,
        },
      },
    });

    if (subscriberInfo) {
      updateSubscriberData({ topUpSubscriberName: subscriberInfo.name });
    }
  };

  const onTableMenuItem = async (id: string, type: string) => {
    switch (type) {
      case 'delete-sub':
        const subscriberToDelete =
          data?.getSubscribersByNetwork.subscribers.find(
            (subscriber) => subscriber.uuid === id,
          );

        updateDialogState({ confirmation: true });
        updateOperationData({
          deletedSubscriber: {
            id: id,
            name: subscriberToDelete?.name ?? 'This subscriber',
          },
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
      updateSubscriberData({ packageHistories: [] });

      const subscriberInfo = data?.getSubscribersByNetwork.subscribers.find(
        (subscriber) => subscriber.uuid === id,
      );

      openDialog('subscriberDetails');

      getSimBySubscriber({
        variables: {
          data: {
            subscriberId: id,
          },
        },
        onCompleted: (res) => {
          if (res.getSimsBySubscriber && res.getSimsBySubscriber.sims) {
            updateSubscriberData({ simList: res.getSimsBySubscriber.sims });

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

      const usageData = dataUsageData?.getDataUsages.usages.find(
        (usage) => usage.simId === subscriberInfo?.sim?.[0]?.id,
      );

      updateSubscriberData({ dataUsageForSim: usageData?.usage ?? '' });

      if (subscriberInfo) {
        const plan = packagesData?.getPackages.packages.find(
          (pkg) => pkg.uuid === subscriberInfo.sim?.[0]?.package?.package_id,
        );

        updateSubscriberData({
          details: {
            ...subscriberInfo,
            packageId: subscriberInfo.sim?.[0]?.package?.package_id,
            dataUsage: `${formatBytesToGB(Number(usageData?.usage)) || 0} GB`,
            dataPlan: plan?.name ?? 'No active plan',
            simIccid: subscriberInfo.sim?.[0]?.iccid,
          },
        });
      }
    },
    [
      data?.getSubscribersByNetwork.subscribers,
      dataUsageData?.getDataUsages.usages,
      packagesData?.getPackages.packages,
      getSimBySubscriber,
      getPackagesForSim,
      openDialog,
      updateSubscriberData,
    ],
  );

  const handleAddSubscriberModal = () => {
    openDialog('addSubscriber');
    refetchSims();
  };
  const handleSimAction = (
    action: string,
    simId: string,
    additionalData?: any,
  ) => {
    switch (action) {
      case 'deactivateSim':
        toggleSimStatus({
          variables: {
            data: {
              sim_id: simId,
              status: 'inactive',
            },
          },
        });
        break;
      case 'activateSim':
        toggleSimStatus({
          variables: {
            data: {
              sim_id: simId,
              status: 'active',
            },
          },
        });
        break;

      case 'deleteSim':
        closeDialog('subscriberDetails');
        updateOperationData({
          simToDelete: {
            id: simId,
            iccid: additionalData.iccid,
            isLastSim: additionalData.isLastSim,
          },
        });
        openDialog('simDeleteConfirmation');
        break;

      default:
        break;
    }
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

      closeDialog('topupData');
    } catch (error) {
      console.error('Error handling top up:', error);
      closeDialog('topupData');
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

      updateDialogState({
        confirmation: true,
        subscriberDetails: false,
      });
      updateOperationData({
        deletedSubscriber: {
          id: subscriberId,
          name: subscriberToDelete?.name ?? 'This subscriber',
        },
      });
    }
  };

  const handleAddSubscriber = async (
    subscriber: any,
  ): Promise<AllocateSimApiDto> => {
    try {
      updateSubscriberData({ details: subscriber });
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

  const handleMenuItemClick = (id: string, type: string) => {
    const subscriber = data?.getSubscribersByNetwork.subscribers.find(
      (subscriber) => subscriber.uuid === id,
    );

    if (subscriber && subscriber.subscriberStatus === 'pending_deletion') {
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
            search={uiState.search}
            title={'My subscribers'}
            buttonTitle={'Add Subscriber'}
            handleButtonAction={handleAddSubscriberModal}
            onSearchChange={(e: string) => updateUIState({ search: e })}
            subtitle={`${data?.getSubscribersByNetwork.subscribers.length}`}
          />
          <br />

          <DataTableWithOptions
            icon={SubscriberIcon}
            isRowClickable={false}
            columns={SUBSCRIBER_TABLE_COLUMNS}
            dataset={structureData(uiState.subscribers)}
            menuOptions={SUBSCRIBER_TABLE_MENU}
            onMenuItemClick={handleMenuItemClick}
            emptyViewLabel={'No subscribers yet!'}
          />
        </Paper>
      )}

      <AddSubscriberStepperDialog
        isOpen={dialogStates.addSubscriber}
        currencySymbol={currencyData?.getCurrencySymbol.symbol ?? ''}
        handleCloseAction={() => closeDialog('addSubscriber')}
        handleAddSubscriber={handleAddSubscriber}
        sims={simPoolData?.getSimsFromPool.sims ?? []}
        packages={packagesData?.getPackages.packages ?? []}
        isLoading={addSubscriberLoading || allocateSimLoading}
      />

      <DeleteConfirmation
        open={dialogStates.confirmation}
        onDelete={handleDeleteSubscriber}
        onCancel={() => closeDialog('confirmation')}
        itemName={operationData.deletedSubscriber.name}
        itemType="subscriber"
        loading={deleteSubscriberLoading}
      />

      <SubscriberDetailsDialog
        open={dialogStates.subscriberDetails}
        onClose={() => closeDialog('subscriberDetails')}
        subscriber={{
          id: subscriberData.details?.uuid || '',
          firstName: subscriberData.details?.name || '',
          email: subscriberData.details?.email || '',
        }}
        onUpdateSubscriber={(updates) =>
          handleUpdateSubscriber(subscriberData.details?.uuid, updates)
        }
        onDeleteSubscriber={() =>
          handleSubscriberMenuAction(
            'deleteSubscriber',
            subscriberData.details?.uuid,
          )
        }
        onTopUpPlan={handleTopUpDataPreparation}
        sims={subscriberData.simList ?? []}
        onSimAction={handleSimAction}
        packageHistories={subscriberData.packageHistories}
        packagesData={packagesData?.getPackages}
        loadingPackageHistories={packagesForSimLoading}
        dataUsage={subscriberData.dataUsageForSim}
        currencySymbol={currencyData?.getCurrencySymbol.symbol ?? ''}
      />

      <TopUpData
        isToPup={dialogStates.topupData}
        onCancel={() => closeDialog('topupData')}
        handleTopUp={handleTopUp}
        loadingTopUp={packagesLoading || addPackagesToSimLoading}
        packages={packagesData?.getPackages.packages ?? []}
        sims={subscriberData.simList ?? []}
        subscriberName={subscriberData.topUpSubscriberName}
      />

      <DeleteConfirmation
        open={dialogStates.simDeleteConfirmation}
        onDelete={() => {
          deleteSim({
            variables: {
              data: {
                simId: operationData.simToDelete.id,
              },
            },
          });
        }}
        onCancel={() => closeDialog('simDeleteConfirmation')}
        itemName={operationData.simToDelete.iccid}
        itemType="sim"
        loading={deleteSimLoading}
      />
    </Stack>
  );
};

export default Page;
