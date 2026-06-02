/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import {
  AllocateSimApiDto,
  GetSimsBySubscriberQuery,
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
  useGetSimsFromPoolQuery,
  useGetSubscribersByNetworkQuery,
  useToggleSimStatusMutation,
  useUpdateSubscriberMutation,
} from '@/client/graphql/generated';
import { useEnvContext, useUserContext, useNetworkContext, useUIContext } from '@/context';
import {
  SubscriberDetailsType,
  TSubscriberDetails,
  TSubscriberTableRow,
} from '@/types';
import { formatBytesToGB } from '@/utils';
import { AlertColor } from '@mui/material';
import { useSearchParams } from 'next/navigation';
import { useCallback, useRef, useState } from 'react';

export function useSubscribersPage() {
  const query = useSearchParams();
  const [search, setSearch] = useState<string>('');
  const { env } = useEnvContext();
  const { user } = useUserContext();
  const { network } = useNetworkContext();
  const { setSnackbarMessage } = useUIContext();
  const [openAddSubscriber, setOpenAddSubscriber] = useState(false);
  const [isTopupData, setIsTopupData] = useState<boolean>(false);
  const [subscriberDetails, setSubscriberDetails] =
    useState<TSubscriberDetails | null>(null);
  const [isSubscriberDetailsOpen, setIsSubscriberDetailsOpen] =
    useState<boolean>(false);
  type QuerySim =
    GetSimsBySubscriberQuery['getSimsBySubscriber']['sims'][number];
  const [subscriberSimList, setSubscriberSimList] = useState<QuerySim[]>([]);
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);
  const [deletedSubscriber, setDeletedSubscriber] = useState<string>('');
  const scrollContainerRef = useRef<HTMLDivElement | null>(null);
  const [topUpSubscriberName, setTopUpSubscriberName] = useState('');
  const [subscriber, setSubscriber] = useState<SubscribersResDto>({
    subscribers: [],
  });

  const notify = (id: string, message: string, type: AlertColor | 'error') =>
    setSnackbarMessage({ id, message, type, show: true });

  const { data: packagesData, loading: packagesLoading } = useGetPackagesQuery({
    fetchPolicy: 'cache-and-network',
    onError: (e) => notify('packages', e.message, 'error'),
  });

  const { data: simPoolData, refetch: refetchSims } = useGetSimsFromPoolQuery({
    variables: {
      data: {
        status: Sim_Status.Unassigned,
        type: env.SIM_TYPE as Sim_Types,
      },
    },
    fetchPolicy: 'network-only',
    onError: (e) => notify('sims-error-msg', e.message, 'error'),
  });

  const [getSimBySubscriber] = useGetSimsBySubscriberLazyQuery({
    onCompleted: (res) => {
      if (res.getSimsBySubscriber) {
        setSubscriberSimList(res.getSimsBySubscriber.sims);
      }
    },
  });

  const {
    data,
    loading: subscribersLoading,
    refetch: refetchSubscribers,
  } = useGetSubscribersByNetworkQuery({
    skip: !network.id,
    variables: { networkId: network.id },
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      setSubscriber({
        subscribers: [...data.getSubscribersByNetwork.subscribers],
      });
      if (query.size > 0) setSearch(query.get('iccid') ?? '');
      getDataUsages({
        variables: {
          data: { type: Sim_Types.UkamaData, networkId: network.id },
        },
      });
    },
    onError: (e) => notify('subscriber-msg', e.message, 'error'),
  });

  const [toggleSimStatus, { loading: toggleSimStatusLoading }] =
    useToggleSimStatusMutation({
      onCompleted: () => {
        notify(
          'sim-activated-success',
          'Sim state updated successfully!',
          'success',
        );
        refetchSubscribers();
      },
      onError: (e) => notify('sim-activated-error', e.message, 'error'),
    });

  const [addPackagesToSim, { loading: addPackagesToSimLoading }] =
    useAddPackagesToSimMutation({
      onCompleted: () =>
        notify(
          'packages-added-success',
          'Packages added successfully!',
          'success',
        ),
      onError: (e) => notify('packages-added-error', e.message, 'error'),
    });

  const [allocateSim, { loading: allocateSimLoading }] = useAllocateSimMutation(
    {
      onCompleted: () =>
        notify(
          'allocate-sim-success',
          'SIM allocated successfully!',
          'success',
        ),
      onError: (e) => notify('allocate-sim-error', e.message, 'error'),
    },
  );

  const [addSubscriber, { loading: addSubscriberLoading }] =
    useAddSubscriberMutation({
      onCompleted: () => {
        refetchSubscribers().then((res) =>
          setSubscriber({
            subscribers: [...res.data.getSubscribersByNetwork.subscribers],
          }),
        );
        notify(
          'add-subscriber-success',
          'Subscriber added successfully!',
          'success',
        );
        refetchSims();
      },
      onError: (e) => notify('add-subscriber-error', e.message, 'error'),
    });

  const [deleteSubscriber, { loading: deleteSubscriberLoading }] =
    useDeleteSubscriberMutation({
      onCompleted: () => {
        refetchSubscribers();
        notify(
          'delete-subscriber-success',
          'Subscriber deleted successfully!',
          'success',
        );
        setIsConfirmationOpen(false);
      },
      onError: (e) => notify('delete-subscriber-error', e.message, 'error'),
    });

  const [updateSubscriber, { loading: updateSubscriberLoading }] =
    useUpdateSubscriberMutation({
      onCompleted: () => {
        refetchSubscribers().then((res) =>
          setSubscriber({
            subscribers: [...res.data.getSubscribersByNetwork.subscribers],
          }),
        );
        notify(
          'update-subscriber-success',
          'Subscriber updated successfully!',
          'success',
        );
      },
      onError: (e) => notify('update-subscriber-error', e.message, 'error'),
    });

  const { data: currencyData } = useGetCurrencySymbolQuery({
    skip: !user.currency,
    fetchPolicy: 'cache-first',
    variables: { code: user.currency },
    onError: (e) => notify('currency-info-error', e.message, 'error'),
  });

  const [getDataUsages, { data: dataUsageData, loading: dataUsageLoading }] =
    useGetDataUsagesLazyQuery({
      pollInterval: 30000,
      fetchPolicy: 'network-only',
      variables: {
        data: { type: Sim_Types.UkamaData, networkId: network.id },
      },
    });

  const buildTableRows = useCallback(
    (data: SubscribersResDto): TSubscriberTableRow[] => {
      if (!(packagesData?.getPackages.packages?.length ?? 0) || !network)
        return [];
      return data.subscribers.map((sub) => {
        const sim = sub?.sim?.length ? sub.sim[0] : null;
        const pkg = packagesData?.getPackages.packages.find(
          (p) => p.uuid === sim?.package?.package_id,
        );
        const usage = dataUsageData?.getDataUsages.usages.find(
          (u) => u.simId === sim?.id,
        );
        return {
          id: sub.uuid,
          name: sub.name,
          email: sub.email,
          packageId: sim?.package?.package_id,
          dataPlan: pkg?.name ?? 'No active plan',
          dataUsage: `${formatBytesToGB(Number(usage?.usage)) || 0} GB`,
          actions: '',
        };
      });
    },
    [
      packagesData?.getPackages.packages,
      dataUsageData?.getDataUsages.usages,
      network,
    ],
  );

  const handleOpenSubscriberDetails = useCallback(
    (id: string) => {
      const sub = data?.getSubscribersByNetwork.subscribers.find(
        (s) => s.uuid === id,
      );
      setIsSubscriberDetailsOpen(true);
      if (sub) {
        const usage = dataUsageData?.getDataUsages.usages.find(
          (u) => u.simId === sub.sim?.[0]?.id,
        );
        const plan = packagesData?.getPackages.packages.find(
          (p) => p.uuid === sub.sim?.[0]?.package?.package_id,
        );
        setSubscriberDetails({
          ...sub,
          packageId: sub.sim?.[0]?.package?.package_id,
          dataUsage: `${formatBytesToGB(Number(usage?.usage)) || 0} GB`,
          dataPlan: plan?.name ?? 'No active plan',
          simIccid: sub.sim?.[0]?.iccid,
        });
      }
    },
    [
      data?.getSubscribersByNetwork.subscribers,
      dataUsageData?.getDataUsages.usages,
      packagesData?.getPackages.packages,
    ],
  );

  const handleTopUpDataPreparation = (id: string) => {
    const sub = data?.getSubscribersByNetwork.subscribers.find(
      (s) => s.uuid === id,
    );
    setIsTopupData(true);
    getSimBySubscriber({ variables: { data: { subscriberId: id } } });
    if (sub) setTopUpSubscriberName(sub.name);
  };

  const handleTableMenuAction = (id: string, type: string) => {
    switch (type) {
      case 'delete-sub':
        setIsConfirmationOpen(true);
        setDeletedSubscriber(id);
        break;
      case 'top-up-data':
        handleTopUpDataPreparation(id);
        break;
      case 'edit-sub':
        handleOpenSubscriberDetails(id);
        break;
    }
  };

  const handleDeleteSubscriber = () => {
    deleteSubscriber({ variables: { subscriberId: deletedSubscriber } });
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
        setIsTopupData(true);
        break;
    }
  };

  const handleTopUp = async (simId: string, planIds: string[]) => {
    const packages = planIds.map((planId) => ({
      package_id: planId,
      start_date: new Date(Date.now() + 60000).toISOString(),
    }));
    await addPackagesToSim({
      variables: { data: { sim_id: simId, packages } },
    });
    setIsTopupData(false);
  };

  const handleUpdateSubscriber = (
    subscriberId: string,
    updates: { name?: string; phone?: string },
  ) => {
    updateSubscriber({ variables: { subscriberId, data: updates } });
    refetchSubscribers();
  };

  const handleSubscriberMenuAction = (action: string, subscriberId: string) => {
    if (action === 'deleteSubscriber') {
      deleteSubscriber({ variables: { subscriberId } });
    }
  };

  const handleAddSubscriber = async (
    sub: SubscriberDetailsType,
  ): Promise<AllocateSimApiDto> => {
    const subscriberResponse = await addSubscriber({
      variables: {
        data: {
          name: sub.name,
          network_id: network.id,
          email: sub.email,
          phone: '',
        },
      },
    });
    if (!subscriberResponse.data) throw new Error('Failed to add subscriber');

    const simResponse = await allocateSim({
      variables: {
        data: {
          network_id: subscriberResponse.data.addSubscriber.networkId,
          package_id: sub.plan,
          subscriber_id: subscriberResponse.data.addSubscriber.uuid,
          sim_type: env.SIM_TYPE,
          iccid: sub.simIccid,
          traffic_policy: 0,
        },
      },
    });
    if (!simResponse.data) throw new Error('Failed to allocate SIM');
    return simResponse.data.allocateSim;
  };

  const scroll = (direction: 'left' | 'right'): void => {
    if (scrollContainerRef.current) {
      const amount = scrollContainerRef.current.clientWidth / 2;
      scrollContainerRef.current.scrollLeft +=
        direction === 'left' ? -amount : amount;
    }
  };

  return {
    search,
    setSearch,
    openAddSubscriber,
    isTopupData,
    subscriberDetails,
    isSubscriberDetailsOpen,
    subscriberSimList,
    isConfirmationOpen,
    deletedSubscriber,
    scrollContainerRef,
    topUpSubscriberName,
    subscriber,
    packagesData,
    packagesLoading,
    simPoolData,
    subscribersLoading,
    dataUsageLoading,
    currencyData,
    toggleSimStatusLoading,
    addSubscriberLoading,
    allocateSimLoading,
    deleteSubscriberLoading,
    updateSubscriberLoading,
    addPackagesToSimLoading,
    subscriberCount: data?.getSubscribersByNetwork.subscribers.length ?? 0,
    buildTableRows,
    handleAddSubscriberModal: () => {
      setOpenAddSubscriber(true);
      refetchSims();
    },
    handleCloseAddSubscriber: () => setOpenAddSubscriber(false),
    handleCloseSubscriberDetails: () => setIsSubscriberDetailsOpen(false),
    handleCancel: () => setIsConfirmationOpen(false),
    handleDeleteSubscriber,
    handleSimAction,
    handleTopUp,
    handleCloseTopUp: () => setIsTopupData(false),
    handleUpdateSubscriber,
    handleSubscriberMenuAction,
    handleAddSubscriber,
    handleTableMenuAction,
    scroll,
  };
}
