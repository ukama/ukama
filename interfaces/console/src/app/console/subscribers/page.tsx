/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  Sim_Status,
  Sim_Types,
  SubscribersResDto,
  useAddPackageToSimMutation,
  useAddSubscriberMutation,
  useAllocateSimMutation,
  useDeleteSimMutation,
  useDeleteSubscriberMutation,
  useGetNetworksQuery,
  useGetPackagesForSimLazyQuery,
  useGetPackagesQuery,
  useGetSimLazyQuery,
  useGetSimPoolStatsQuery,
  useGetSimsBySubscriberLazyQuery,
  useGetSimsQuery,
  useGetSubscriberLazyQuery,
  useGetSubscribersByNetworkQuery,
  useSetActivePackageForSimMutation,
  useToggleSimStatusMutation,
  useUpdateSubscriberMutation,
} from '@/client/graphql/generated';
import AddSubscriberDialog from '@/components/AddSubscriber';
import DataTableWithOptions from '@/components/DataTableWithOptions';
import DeleteConfirmation from '@/components/DeleteDialog';
import EmptyView from '@/components/EmptyView';
import LoadingWrapper from '@/components/LoadingWrapper';
import PageContainerHeader from '@/components/PageContainerHeader';
import PlanCard from '@/components/PlanCard';
import SubscriberDetails from '@/components/SubscriberDetails';
import TopUpData from '@/components/TopUpData';
import { SUBSCRIBER_TABLE_COLUMNS, SUBSCRIBER_TABLE_MENU } from '@/constants';
import { useAppContext } from '@/context';
import { TAddSubscriberData } from '@/types';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';
import { AlertColor, Box, Grid, Paper, Stack, Typography } from '@mui/material';
import { useCallback, useEffect, useState } from 'react';

const Page = () => {
  const [search, setSearch] = useState<string>('');
  const { setSnackbarMessage, network, env } = useAppContext();
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
  const [isAddSubscriberDialogOpen, setIsAddSubscriberDialogOpen] =
    useState(false);
  const [isToPupData, setIsToPupData] = useState<boolean>(false);
  const [topUpDetails, setTopUpDetails] = useState<any>({
    simId: '',
    subscriberId: '',
  });
  const [simId, setSimId] = useState<string>('');
  const [qrCode, setQrCode] = useState<string>('');
  const [subscriberSuccess, setSubscriberSuccess] = useState<boolean>(false);
  const [selectedSubscriber, setSelectedSubscriber] = useState<any>();
  const [isSubscriberDetailsOpen, setIsSubscriberDetailsOpen] =
    useState<boolean>(false);
  const [subscriberSimList, setSubscriberSimList] = useState<any[]>();
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);
  const [isPackageActivationNeeded, setIsPackageActivationNeeded] =
    useState(false);
  const [deletedSubscriber, setDeletedSubscriber] = useState<string>('');
  const [selectedNetwork, setSelectedNetwork] = useState<string>('');
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
          if (
            subscriber.firstName.toLowerCase().includes(s) ??
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
        activatePackageSim({
          variables: {
            data: {
              sim_id: res.getPackagesForSim.sim_id,
              package_id: res.getPackagesForSim.packages[0].package_id,
            },
          },
        });
        setIsPackageActivationNeeded(false);
      }
    },
  });

  const [getSubscriber, { data: subcriberInfo }] = useGetSubscriberLazyQuery({
    onCompleted: (res) => {
      if (res?.getSubscriber?.sim && res.getSubscriber?.sim.length > 0) {
        fetchPackagesForSim(res.getSubscriber.sim[0].id);
      }
    },
  });
  const handleOpenSubscriberDetails = useCallback(
    (id: string, shouldActivatePackage: boolean = false) => {
      setIsSubscriberDetailsOpen(true);
      setSelectedSubscriber(id);
      setIsPackageActivationNeeded(shouldActivatePackage);

      getSubscriber({
        variables: {
          subscriberId: id,
        },
      });
    },
    [getSubscriber],
  );

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

  const onTableMenuItem = (id: string, type: string) => {
    if (type === 'delete-sub') {
      setIsConfirmationOpen(true);
      setDeletedSubscriber(id);
    }
    if (type === 'top-up-data') {
      setIsToPupData(true);
      setTopUpDetails({
        simId: id,
        subscriberId: id,
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
      handleOpenSubscriberDetails(id, false);
    }
    if (type === 'pause-service') {
      toggleSimStatus({
        variables: {
          data: {
            sim_id: id,
            status: 'terminated',
          },
        },
      });
    }
  };
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
  const {
    data,
    loading: getSubscriberByNetworkLoading,
    refetch: refetchSubscribers,
  } = useGetSubscribersByNetworkQuery({
    skip: !network.id,
    variables: {
      networkId: network.id,
    },
    fetchPolicy: 'cache-first',
    onCompleted: (data) => {
      if (data.getSubscribersByNetwork.subscribers.length > 0) {
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

  const structureData = useCallback(
    (data: SubscribersResDto) =>
      data.subscribers.map((subscriber) => {
        const networkName =
          networkList?.getNetworks?.networks.find(
            (net) => net.id === subscriber.networkId,
          )?.name ?? '';

        return {
          id: subscriber.uuid,
          email: subscriber.email,
          name: `${subscriber.firstName}`,
          dataUsage: '',
          dataPlan: '',
          actions: '',
          network: networkName,
        };
      }),
    [networkList],
  );

  const [getSim] = useGetSimLazyQuery({
    onCompleted: (res) => {},
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

  const handleDeleteSubscriber = () => {
    deleteSubscriber({
      variables: {
        subscriberId: deletedSubscriber,
      },
    });
  };
  const [addPackageToSim, { loading: addPackageToSimLoading }] =
    useAddPackageToSimMutation({
      onCompleted: () => {
        setSnackbarMessage({
          id: 'package-added-success',
          message: 'Package added successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'package-added-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });
  const [allocateSim, { loading: allocateSimLoading }] = useAllocateSimMutation(
    {
      onCompleted: (res) => {
        setSimId(res.allocateSim.id);
        refetchSubscribers();
        setSnackbarMessage({
          id: 'sim-allocated-success',
          message: 'Sim allocated successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
        setQrCode(res.allocateSim.iccid);
        // getPackagesForSim({
        //   variables: {
        //     data: {
        //       sim_id: res.allocateSim.id,
        //     },
        //   },
        // });
        // getSim({
        //   variables: {
        //     data: {
        //       simId: res.allocateSim.id,
        //     },
        //   },
        // });
        refetchSimPoolStats();
        refetchSims();
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
        allocateSim({
          variables: {
            data: {
              network_id: res.addSubscriber.networkId,
              package_id: addSubscriberData.plan ?? '',
              subscriber_id: res.addSubscriber.uuid,
              sim_type: env.SIM_TYPE,
              iccid: addSubscriberData.iccid,
              traffic_policy: 0,
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

  const [deleteSim, { loading: deleteSimLoading }] = useDeleteSimMutation({
    onCompleted: () => {
      setSnackbarMessage({
        id: 'sim-delete-success',
        message: 'Sim deleted successfully',
        type: 'success' as AlertColor,
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'sim-delete-error',
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
    setIsAddSubscriberDialogOpen(true);
  };

  const OnCloseAddSubcriber = () => {
    setIsAddSubscriberDialogOpen(false);
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
    if (action === 'deleteSim') {
      deleteSim({
        variables: {
          data: {
            simId: simId,
          },
        },
      });
    } else if (action === 'deactivateSim') {
      toggleSimStatus({
        variables: {
          data: {
            sim_id: simId,
            status: 'inactive',
          },
        },
      });
    } else if (action === 'topUp') {
      setIsToPupData(true);
      setTopUpDetails({
        simId: simId,
        subscriberId: selectedSubscriber,
      });
    }
  };

  const handleCloseTopUp = () => {
    setIsToPupData(false);
  };

  const handleTopUp = async (topUpplan: string, selectedSim: string) => {
    const data = {
      sim_id: selectedSim,
      package_id: topUpplan,
    };

    await addPackageToSim({
      variables: {
        data: {
          sim_id: selectedSim,
          package_id: topUpplan,
          start_date: new Date(Date.now() + 5 * 60000).toISOString(),
        },
      },
    });
    await activatePackageSim({ variables: { data } });
  };

  const handleUpdateSubscriber = (
    subscriberId: string,
    firstName: string,
    phone: string,
  ) => {
    updateSubscriber({
      variables: {
        subscriberId: subscriberId,
        data: {
          phone: phone,
          first_name: firstName,
        },
      },
    });
    refetchSubscribers();
  };

  const handleDeleteSubscriberModal = (
    action: string,
    subscriberId: string,
  ) => {
    if (action === 'deleteSubscriber') {
      deleteSubscriber({
        variables: {
          subscriberId: subscriberId,
        },
      });
    } else if (action === 'pauseService') {
      toggleSimStatus({
        variables: {
          data: {
            sim_id: subscriberId,
            status: 'terminated',
          },
        },
      });
    }
  };

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

  return (
    <Stack
      direction={'column'}
      spacing={2}
      mt={2}
      sx={{ height: 'calc(100vh - 200px)' }}
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
          sx={{
            height: '218px',
            borderRadius: '10px',
            padding: '24px 32px',
          }}
        >
          <Stack direction={'column'} spacing={1.5}>
            <Typography variant="h6" mr={1}>
              Data plans
            </Typography>
            <Box>
              {packagesData?.getPackages.packages.length === 0 ? (
                <EmptyView
                  icon={UpdateIcon}
                  title="No data plan created yet!"
                />
              ) : (
                <Grid
                  container
                  rowSpacing={2}
                  columnSpacing={2}
                  overflow={'scroll'}
                >
                  {packagesData?.getPackages.packages.map(
                    ({
                      uuid,
                      name,
                      duration,
                      currency,
                      dataVolume,
                      dataUnit,
                      amount,
                    }: any) => (
                      <Grid item xs={12} sm={6} md={3} key={uuid}>
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
                      </Grid>
                    ),
                  )}
                </Grid>
              )}
            </Box>
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
            padding: '24px 32px',
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

      <AddSubscriberDialog
        qrCode={qrCode}
        pkgList={packagesData?.getPackages.packages ?? []}
        onSuccess={subscriberSuccess}
        onClose={OnCloseAddSubcriber}
        onSubmit={handleAddSubscriber}
        open={isAddSubscriberDialogOpen}
        sims={simPoolData?.getSims.sim ?? []}
        pSimCount={simStatData?.getSimPoolStats.physical}
        eSimCount={simStatData?.getSimPoolStats.esim}
        submitButtonState={
          addSubscriberLoading ?? allocateSimLoading ?? packagesLoading
        }
        loading={addSubscriberLoading ?? allocateSimLoading ?? packagesLoading}
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
        subscriberId={selectedSubscriber}
        subscriberInfo={subcriberInfo?.getSubscriber}
        handleSimActionOption={handleSimAction}
        handleUpdateSubscriber={handleUpdateSubscriber}
        loading={updateSubscriberLoading ?? deleteSimLoading}
        handleDeleteSubscriber={handleDeleteSubscriberModal}
        simStatusLoading={toggleSimStatusLoading}
        currentSite={'-'}
      />

      <TopUpData
        isToPup={isToPupData}
        onCancel={handleCloseTopUp}
        subscriberId={topUpDetails.subscriberId}
        handleTopUp={handleTopUp}
        loadingTopUp={packagesLoading ?? addPackageToSimLoading}
        packages={packagesData?.getPackages.packages ?? []}
        sims={subscriberSimList ?? []}
      />
    </Stack>
  );
};

export default Page;
