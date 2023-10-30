import { commonData, snackbarMessage } from '@/app-recoil';
import { SUBSCRIBER_TABLE_COLUMNS, SUBSCRIBER_TABLE_MENU } from '@/constants';
import {
  SubscribersResDto,
  useGetSubscribersByNetworkQuery,
  useAddSubscriberMutation,
  useGetSubscriberLazyQuery,
  useGetPackagesQuery,
  PackagesResDto,
  SimsResDto,
  // useAllocateSimMutation,
  // useGetSimLazyQuery,
  // useSetActivePackageForSimMutation,
  // useToggleSimStatusMutation,
  useGetSimpoolStatsQuery,
  useGetSimsQuery,
  // useGetPackagesForSimLazyQuery,
  useDeleteSubscriberMutation,
  // useDeleteSimMutation,
  useUpdateSubscriberMutation,
  // useAddPackageToSimMutation,
  // useGetSimBySubscriberLazyQuery,
  // SubscriberToSimDto,
} from '@/generated';
import {
  ContainerMax,
  PageContainer,
  VerticalContainer,
} from '@/styles/global';
import { colors } from '@/styles/theme';
import { TCommonData, TSnackMessage } from '@/types';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import DataTableWithOptions from '@/ui/molecules/DataTableWithOptions';
import AddSubscriberDialog from '@/ui/molecules/AddSubscriber';
import DeleteConfirmation from '@/ui/molecules/DeleteDialog';
import TopUpData from '@/ui/molecules/TopUpData';
import SubscriberDetails from '@/ui/molecules/SubscriberDetails';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import { AlertColor } from '@mui/material';
import { useCallback, useEffect, useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';

const Page = () => {
  const [search, setSearch] = useState<string>('');
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const [simId, setSimId] = useState<string>('');
  const [qrCode, setQrCode] = useState<string>('');
  const [simPlan, setSimPlan] = useState<string>('');
  const [iccid, setIccid] = useState<string>('');
  const [openAddSubscriber, setOpenAddSubscriber] = useState<boolean>(false);
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);
  const [deletedSubscriber, setDeletedSubscriber] = useState<string>('');
  const [activePackageId, setActivePackageId] = useState<string>('');
  const [subscriberSuccess, setSubscriberSuccess] = useState<boolean>(false);
  const [isToPupData, setIsToPupData] = useState<boolean>(false);
  const [selectedSubscriber, setSelectedSubscriber] = useState<any>();
  const [subcriberInfo, setSubscriberInfo] = useState<any>();
  const [subscriberSimList, setSubscriberSimList] =
    useState<SubscriberToSimDto[]>();
  const [isSubscriberDetailsOpen, setIsSubscriberDetailsOpen] =
    useState<boolean>(false);
  const [topUpDetails, setTopUpDetails] = useState<any>({
    simId: '',
    subscriberId: '',
  });
  const [subscriber, setSubscriber] = useState<SubscribersResDto>({
    subscribers: [],
  });
  const [packages, setPackages] = useState<PackagesResDto>({
    packages: [],
  });
  const [simList, setSimList] = useState<SimsResDto>({
    sims: [],
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

  const [getSim] = useGetSimLazyQuery({
    onCompleted: (res) => {
      if (res.getSim) {
        console.log('SIM', res.getSim);
      }
    },
  });

  const [getSimBySubscriber] = useGetSimBySubscriberLazyQuery({
    onCompleted: (res) => {
      if (res.getSimBySubscriber) {
        setSubscriberSimList(res.getSimBySubscriber.sims);
      }
    },
  });

  const [getPackagesForSim] = useGetPackagesForSimLazyQuery({
    onCompleted: (res) => {
      if (res.getPackagesForSim.packages) {
      }
    },
  });

  const [getSubscriber] = useGetSubscriberLazyQuery({
    onCompleted: (res) => {
      if (res.getSubscriber) {
        setSubscriberInfo(res.getSubscriber);
        getPackagesForSim({
          variables: {
            data: {
              simId: res.getSubscriber.sim[0].id,
            },
          },
        });
      }
    },
  });

  useEffect(() => {
    if (search.length > 3) {
      const subscribers = data?.getSubscribersByNetwork.subscribers.filter(
        (subscriber) => {
          const s = search.toLowerCase();
          if (
            subscriber.firstName?.toLowerCase().includes(s) ||
            subscriber.lastName?.toLowerCase().includes(s)
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

  const onTableMenuItem = (id: string, type: string) => {
    console.log('TYPE', type, 'ID', id);
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
          subscriberId: id,
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

  const structureData = useCallback(
    (data: SubscribersResDto) =>
      data.subscribers.map((subscriber) => ({
        id: subscriber.uuid,
        email: subscriber.email,
        name: `${subscriber.firstName} ${subscriber.lastName}`,
        dataUsage: '',
        dataPlan: subscriber.sim.length > 0 ? subscriber.sim[0].package : '',
        actions: '',
      })),
    [],
  );
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
  const [updateSubscriber, { loading: updateSubscriberLoading }] =
    useUpdateSubscriberMutation({
      onCompleted: (res) => {
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
              network_id: _commonData.networkId,
              package_id: simPlan ?? '',
              subscriber_id: res.addSubscriber.uuid,
              sim_type: 'test',
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
  const [addPackageToSim, { loading: addPackageToSimLoading }] =
    useAddPackageToSimMutation({
      onCompleted: () => {
        setSnackbarMessage({
          id: 'sim-package-add-success',
          message: 'Package added successfully',
          type: 'success' as AlertColor,
          show: true,
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'sim-package-add-error',
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
        setQrCode(res.allocateSim.sim.iccid);
        setSimId(res.allocateSim.sim.id);
        getSim({
          variables: {
            data: {
              simId: res.allocateSim.sim.id,
            },
          },
        });

        setActivePackageId(res.allocateSim.packageId);
        toggleSimStatus({
          variables: {
            data: {
              sim_id: res.allocateSim.sim.id ?? '',
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
              sim_id: simId,
              package_id: activePackageId,
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

  const OnCloseAddSubcriber = () => {
    setOpenAddSubscriber(false);
  };
  const handleAddSubscriberModal = () => {
    setOpenAddSubscriber(true);
  };

  const { data: _sims } = useGetSimsQuery({
    variables: { type: 'test' },

    fetchPolicy: 'cache-first',
    onCompleted: (_sims) => {
      if (
        _sims &&
        _sims.getSims &&
        _sims.getSims.sims &&
        _sims.getSims.sims.length > 0
      ) {
        const simsArray = _sims.getSims.sims || [];

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
  const handleDeleteSubscriber = () => {
    deleteSubscriber({
      variables: {
        subscriberId: deletedSubscriber,
      },
    });
  };

  const handleCancel = () => {
    setIsConfirmationOpen(false);
  };
  const handleRoamingInstallation = async (values: any) => {
    const { plan, simIccid, email, name } = values;
    setSimPlan(plan);
    setIccid(simIccid);

    await addSubscriber({
      variables: {
        data: {
          email: email as string,
          first_name: name as string,
          last_name: '',
          network_id: _commonData.networkId,
          org_id: _commonData.orgId,
          address: 'test',
          phone: '',
          dob: 'Fri, 01 Jan 1990 00:00:00 GMT',
          gender: '',
          id_serial: '',
          proof_of_identification: '',
        },
      },
    });
  };
  const handleCloseTopUp = () => {
    setIsToPupData(false);
  };
  const handleTopUp = async (topUpplan: string, selectedSim: string) => {
    const data = {
      sim_id: selectedSim,
      package_id: topUpplan,
    };
    await addPackageToSim({ variables: { data } });
    await activatePackageSim({ variables: { data } });
  };

  const handleCloseSubscriberDetails = () => {
    setIsSubscriberDetailsOpen(false);
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
    }
  };
  const handleUpdateSubscriber = (subscriberId: string, email: string) => {
    updateSubscriber({
      variables: {
        subscriberId: subscriberId,
        data: {
          email: email,
        },
      },
    });
  };
  return (
    <LoadingWrapper
      radius="small"
      width={'100%'}
      isLoading={loading}
      cstyle={{
        backgroundColor: loading ? colors.white : 'transparent',
      }}
    >
      <PageContainer>
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
              columns={SUBSCRIBER_TABLE_COLUMNS}
              dataset={structureData(subscriber)}
              menuOptions={SUBSCRIBER_TABLE_MENU}
              onMenuItemClick={onTableMenuItem}
              emptyViewLabel={'No subscribers yet!'}
            />
          </ContainerMax>
        </VerticalContainer>
        <AddSubscriberDialog
          onSuccess={subscriberSuccess}
          open={openAddSubscriber}
          handleRoamingInstallation={handleRoamingInstallation}
          onClose={OnCloseAddSubcriber}
          qrCode={qrCode}
          submitButtonState={
            addSubscriberLoading || allocateSimLoading || packagesLoading
            // simsLoading
          }
          sims={simList?.sims?.filter((sim) => sim.isPhysical === 'true') || []}
          pkgList={packages.packages}
          loading={
            addSubscriberLoading ||
            allocateSimLoading ||
            packagesLoading ||
            toggleSimStatusLoading ||
            activatePackageSimLoading
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
        <TopUpData
          isToPup={isToPupData}
          onCancel={handleCloseTopUp}
          subscriberId={topUpDetails.subscriberId}
          handleTopUp={handleTopUp}
          loadingTopUp={packagesLoading || addPackageToSimLoading}
          packages={packages.packages}
          sims={subscriberSimList || []}
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
        />
      </PageContainer>
    </LoadingWrapper>
  );
};

export default Page;
