/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  useAddNetworkMutation,
  useGetNetworksQuery,
  useSetDefaultNetworkMutation,
  useUpdateNotificationMutation,
} from '@/client/graphql/generated';
import {
  NotificationsResDto,
  Role_Type,
  useGetNotificationsLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import AddNetworkDialog from '@/components/AddNetworkDialog';
import AppSnackbar from '@/components/AppSnackbar/page';
import AppLayout from '@/components/Layout';
import { useAppContext } from '@/context';
import { getMetaInfo } from '@/lib/MetaInfo';
import ServerNotificationSubscription from '@/lib/NotificationSubscription';
import '@/styles/console.css';
import { TNotificationResDto } from '@/types';
import ErrorBoundary from '@/wrappers/errorBoundary';
import { Box } from '@mui/material';
import PubSub from 'pubsub-js';
import { useEffect, useState } from 'react';

export default function ConosleLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const {
    user,
    network,
    metaInfo,
    pageName,
    isDarkMode,
    setNetwork,
    setMetaInfo,
    setSnackbarMessage,
    subscriptionClient,
  } = useAppContext();
  const [notifications, setNotifications] = useState<
    NotificationsResDto[] | []
  >([]);
  const [startTimeStamp] = useState<string>(new Date().getTime().toString());
  const [showAddNetwork, setShowAddNetwork] = useState<boolean>(false);
  const {
    data: networksData,
    loading: networksLoading,
    refetch: refetchNetworks,
  } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (data.getNetworks.networks.length >= 1 && network.id === '') {
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

  const [addNetwork, { loading: addNetworkLoading }] = useAddNetworkMutation({
    onCompleted: () => {
      refetchNetworks();
      setSnackbarMessage({
        id: 'add-networks-success',
        message: 'Network added successfully',
        type: 'success',
        show: true,
      });
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'add-networks-error',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [setDefaultNetwork] = useSetDefaultNetworkMutation({
    fetchPolicy: 'network-only',
  });

  const [updateNotificationCall] = useUpdateNotificationMutation({
    onCompleted: () => {
      refetchNotifications().then((res) => {
        setNotifications(res.data?.getNotifications.notifications);
      });
    },
  });

  const [
    getNotifications,
    { refetch: refetchNotifications, loading: notificationsLoading },
  ] = useGetNotificationsLazyQuery({
    onCompleted: (data) => {
      if (data.getNotifications.notifications.length > 0) {
        setNotifications(data.getNotifications.notifications);
      }
      ServerNotificationSubscription(
        `notification-${user.orgId}-${user.id}-${user.role}-${network.id}`,
        user.role as Role_Type,
        user.orgId,
        user.id,
        user.orgName,
        network.id,
        startTimeStamp,
      );
    },
  });

  useEffect(() => {
    if (metaInfo.ip === '') {
      fetchInfo();
    }
  }, [metaInfo]);

  useEffect(() => {
    if (user.id && network.id && user.orgId && user.orgName) {
      getNotifications({
        client: subscriptionClient,
        variables: {
          data: {
            userId: user.id,
            subscriberId: '',
            orgId: user.orgId,
            orgName: user.orgName,
            networkId: network.id,
            role: user.role as Role_Type,
            startTimestamp: startTimeStamp,
          },
        },
      });

      PubSub.subscribe(
        `notification-${user.orgId}-${user.id}-${user.role}-${network.id}`,
        handleNotification,
      );
    }
  }, [user.id, network.id]);

  const fetchInfo = async () => {
    const res = await getMetaInfo();
    setMetaInfo({
      ip: res.ip,
      city: res.city,
      lat: res.lat,
      lng: res.lng,
      languages: res.languages,
      currency: res.currency,
      timezone: res.timezone,
      region_code: res.region_code,
      country_code: res.country_code,
      country_name: res.country_name,
      country_calling_code: res.country_calling_code,
    });
    typeof window !== 'undefined' &&
      localStorage.setItem('metaInfo', JSON.stringify(res));
  };

  const handleNotification = (_: any, data: string) => {
    const parsedData: TNotificationResDto = JSON.parse(data);
    const { id, type, scope, title, isRead, description, createdAt } =
      parsedData.data.notificationSubscription;
    setNotifications((prev: any) => {
      if (!prev) return prev;
      return [
        {
          id,
          type,
          scope,
          title,
          isRead,
          createdAt,
          description,
        },
        ...prev,
      ];
    });
  };

  const handleNotificationRead = (id: string) => {
    if (id) {
      updateNotificationCall({
        variables: {
          isRead: true,
          updateNotificationId: id,
        },
      });
    }
  };

  const handleAddNetworkAction = () => setShowAddNetwork(true);

  const handleAddNetwork = (values: any) => {
    addNetwork({
      variables: {
        data: {
          name: values.name,
          budget: values.budget,
          networks: values.networks,
          countries: values.countries,
        },
      },
    })
      .then((res) => {
        if (values.isDefault && res.data?.addNetwork.id) {
          setDefaultNetwork({
            variables: {
              data: {
                id: res.data?.addNetwork.id,
              },
            },
          });
        }
      })
      .finally(() => {
        setShowAddNetwork(false);
      });
  };

  const handleNetworkChange = (id: string) => {
    if (id) {
      setNetwork({
        id: id,
        name:
          networksData?.getNetworks.networks.filter((n) => n.id === id)[0]
            .name ?? '',
      });
    }
  };
  const handleConfigureSite = (nodeId: string) => {
    //console.log(nodeId)
    //Here we will call the nodeState service to get lat,lng and state for node
    console.log('NODE ID :', nodeId);
    // router.push(
    //   `/configure/node/${nodeState.nodeId}?lat=${nodeState.latitude}&lng=${nodeState.longitude}`,
    // );
  };

  return (
    <ErrorBoundary>
      <Box
        sx={{
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <AppLayout
          isDarkMode={isDarkMode}
          isLoading={networksLoading}
          placeholder={'Select Network'}
          handleNotificationRead={handleNotificationRead}
          handleAddNetwork={handleAddNetworkAction}
          handleNetworkChange={handleNetworkChange}
          networks={networksData?.getNetworks.networks ?? []}
          notifications={notifications}
          onConfigureSite={handleConfigureSite}
        >
          {children}
        </AppLayout>
        <AppSnackbar />
        <AddNetworkDialog
          title={'Add Network'}
          isOpen={showAddNetwork}
          labelSuccessBtn={'Submit'}
          labelNegativeBtn={'Cancel'}
          loading={addNetworkLoading}
          handleSuccessAction={handleAddNetwork}
          description={'Add network in organization'}
          handleCloseAction={() => setShowAddNetwork(false)}
        />
      </Box>
    </ErrorBoundary>
  );
}
