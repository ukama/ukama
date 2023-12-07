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
  PackagesResDto,
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
  useGetSimpoolStatsQuery,
  useGetSimsBySubscriberLazyQuery,
  useGetSimsQuery,
  useGetSitesQuery,
  useGetSubscriberLazyQuery,
  useGetSubscribersByNetworkQuery,
  useSetActivePackageForSimMutation,
  useToggleSimStatusMutation,
  useUpdateSubscriberMutation,
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
import DeleteConfirmation from '@/ui/molecules/DeleteDialog';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import SubscriberDetails from '@/ui/molecules/SubscriberDetails';
import TopUpData from '@/ui/molecules/TopUpData';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import { AlertColor, Stack } from '@mui/material';
import { useCallback, useEffect, useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';

const Page = () => {
  const [search, setSearch] = useState<string>('');
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const [isAddSubscriberDialogOpen, setIsAddSubscriberDialogOpen] =
    useState(false);
  const [isToPupData, setIsToPupData] = useState<boolean>(false);
  const [topUpDetails, setTopUpDetails] = useState<any>({
    simId: '',
    subscriberId: '',
  });
  const [simId, setSimId] = useState<string>('');
  const [qrCode, setQrCode] = useState<string>('');
  const [simPlan, setSimPlan] = useState<string>('');
  const [iccid, setIccid] = useState<string>('');
  const [subscriberSuccess, setSubscriberSuccess] = useState<boolean>(false);
  const [selectedSubscriber, setSelectedSubscriber] = useState<any>();
  const [isSubscriberDetailsOpen, setIsSubscriberDetailsOpen] =
    useState<boolean>(false);
  const [subcriberInfo, setSubscriberInfo] = useState<any>();
  const [subscriberSimList, setSubscriberSimList] = useState<any[]>();
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);
  const [deletedSubscriber, setDeletedSubscriber] = useState<string>('');
  const [selectedNetwork, setSelectedNetwork] = useState<string>('');
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
  const [getSimBySubscriber] = useGetSimsBySubscriberLazyQuery({
    onCompleted: (res) => {
      if (res.getSimsBySubscriber) {
        setSubscriberSimList(res.getSimsBySubscriber.sims);
      }
    },
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
  const [getSubscriber] = useGetSubscriberLazyQuery({
    onCompleted: (res) => {
      if (
        res.getSubscriber &&
        res.getSubscriber.sim &&
        res.getSubscriber?.sim.length > 0
      ) {
        setSubscriberInfo(res.getSubscriber);
        getPackagesForSim({
          variables: {
            data: {
              sim_id: res?.getSubscriber && res.getSubscriber.sim[0].id,
            },
          },
        });
      }
    },
  });

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
      setIsSubscriberDetailsOpen(true);
      getSubscriber({
        variables: {
          subscriberId: id,
        },
      });

      setSelectedSubscriber(id);
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
    loading: getSubscriberByNetworkLoading,
    data,
    refetch: refetchSubscribers,
  } = useGetSubscribersByNetworkQuery({
    variables: {
      networkId: selectedNetwork || _commonData.networkId,
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

  const { data: sitesData, loading: sitesLoading } = useGetSitesQuery({
    variables: {
      networkId: _commonData.networkId,
    },
    fetchPolicy: 'cache-and-network',

    onError: (error) => {
      setSnackbarMessage({
        id: 'sites-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  const { data: networkList, loading: netLoading } = useGetNetworksQuery({
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
          )?.name || '';

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

  const handleAddSubscriberModal = () => {
    setIsAddSubscriberDialogOpen(true);
  };
  const OnCloseAddSubcriber = () => {
    setIsAddSubscriberDialogOpen(false);
  };
  const [getSim] = useGetSimLazyQuery({
    onCompleted: (res) => {
      if (res.getSim) {
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
        setSnackbarMessage({
          id: 'sim-allocated-success',
          message: 'Sim allocated successfully!',
          type: 'success' as AlertColor,
          show: true,
        });
        setQrCode(res.allocateSim.iccid);
        getPackagesForSim({
          variables: {
            data: {
              sim_id: res.allocateSim.id,
            },
          },
        });
        getSim({
          variables: {
            data: {
              simId: res.allocateSim.id,
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
              network_id: res.addSubscriber.networkId,
              package_id: simPlan ?? '',
              subscriber_id: res.addSubscriber.uuid,
              sim_type: 'test',
              iccid: iccid,
              traffic_policy: 10,
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
  const handleRoamingInstallation = async (values: any) => {
    const { plan, simIccid, email, name, phone } = values;
    setSimPlan(plan);
    setIccid(simIccid);

    await addSubscriber({
      variables: {
        data: {
          email: email as string,
          first_name: name as string,
          network_id: _commonData.networkId,
          org_id: _commonData.orgId,
          phone: phone,
          last_name: '',
          proof_of_identification: '',
          dob: '',
          address: '',
        },
      },
    });
  };
  const { data: simStatData } = useGetSimpoolStatsQuery({
    variables: { type: 'test' },

    fetchPolicy: 'cache-first',
    onCompleted: (stats) => {
      if (stats.getSimPoolStats.available) {
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
  const getSelectedNetwork = (network: string) => {
    setSelectedNetwork(network);
  };

  const handleCloseSubscriberDetails = () => {
    setIsSubscriberDetailsOpen(false);
  };
  const handleCancel = () => {
    setIsConfirmationOpen(false);
  };
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
  const [updateSubscriber, { loading: updateSubscriberLoading }] =
    useUpdateSubscriberMutation({
      onCompleted: () => {
        refetchSubscribers();
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
  const handleUpdateSubscriber = (
    subscriberId: string,
    email: string,
    firstName: string,
  ) => {
    updateSubscriber({
      variables: {
        subscriberId: subscriberId,
        data: {
          email: email,
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
  return (
    <Stack direction={'column'}>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={sitesLoading || getSubscriberByNetworkLoading}
        cstyle={{
          backgroundColor: getSubscriberByNetworkLoading
            ? colors.white
            : 'transparent',
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
                getSelectedNetwork={getSelectedNetwork}
                networkList={networkList?.getNetworks?.networks ?? []}
              />
            </ContainerMax>
          </VerticalContainer>
          <AddSubscriberDialog
            onSuccess={subscriberSuccess}
            open={isAddSubscriberDialogOpen}
            handleRoamingInstallation={handleRoamingInstallation}
            onClose={OnCloseAddSubcriber}
            qrCode={qrCode}
            submitButtonState={
              addSubscriberLoading || allocateSimLoading || packagesLoading
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
            pSimCount={simStatData?.getSimPoolStats.physical}
            eSimCount={simStatData?.getSimPoolStats.physical}
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
            onCancel={handleCloseSubscriberDetails}
            subscriberInfo={subcriberInfo}
            handleSimActionOption={handleSimAction}
            handleUpdateSubscriber={handleUpdateSubscriber}
            loading={updateSubscriberLoading || deleteSimLoading}
            handleDeleteSubscriber={handleDeleteSubscriberModal}
            simStatusLoading={toggleSimStatusLoading}
            currentSite={
              sitesData?.getSites?.sites.length > 0
                ? sitesData?.getSites?.sites[0].name
                : '-'
            }
          />
          <TopUpData
            isToPup={isToPupData}
            onCancel={handleCloseTopUp}
            subscriberId={topUpDetails.subscriberId}
            handleTopUp={handleTopUp}
            loadingTopUp={packagesLoading || addPackageToSimLoading}
            packages={packages.packages}
            sims={subscriberSimList || []}
          />
        </PageContainer>
      </LoadingWrapper>
    </Stack>
  );
};

export default Page;
