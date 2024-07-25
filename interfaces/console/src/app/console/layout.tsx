/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import { getMetricsClient } from '@/client/client';
import {
  useAddNetworkMutation,
  useGetNetworksQuery,
  useSetDefaultNetworkMutation,
  useUpdateNotificationMutation,
} from '@/client/graphql/generated';
import {
  NotificationsResDto,
  useGetNotificationsQuery,
  useNotificationSubscriptionSubscription,
} from '@/client/graphql/generated/metrics';
import AddNetworkDialog from '@/components/AddNetworkDialog';
import AppSnackbar from '@/components/AppSnackbar/page';
import AppLayout from '@/components/Layout';
import { useAppContext } from '@/context';
import '@/styles/console.css';
import { getRoleType, getScopesByRole } from '@/utils';
import ErrorBoundary from '@/wrappers/errorBoundary';
import { Box } from '@mui/material';
import { useState } from 'react';

export default function ConosleLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const {
    env,
    user,
    network,
    pageName,
    isDarkMode,
    setNetwork,
    setSnackbarMessage,
  } = useAppContext();
  const [alerts, setAlerts] = useState<NotificationsResDto[] | undefined>(
    undefined,
  );
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

  const [updateNotificationCall] = useUpdateNotificationMutation();

  useGetNotificationsQuery({
    client: getMetricsClient(env.METRIC_URL, env.METRIC_WEBSOCKET_URL),
    fetchPolicy: 'cache-first',

    variables: {
      data: {
        orgId: user.orgId,
        userId: user.id,
        networkId: network.id,
        forRole: getRoleType(user.role),
        scopes: getScopesByRole(user.role),
      },
    },
    onCompleted: (data) => {
      const alerts = data.getNotifications.notifications;
      setAlerts(alerts);
    },
  });

  useNotificationSubscriptionSubscription({
    client: getMetricsClient(env.METRIC_URL, env.METRIC_WEBSOCKET_URL),
    variables: {
      networkId: network.id,
      orgId: user.orgId,
      userId: user.id,
      forRole: getRoleType(user.role),
      scopes: getScopesByRole(user.role),
    },
    onData: ({ data }) => {
      const newAlert = data.data?.notificationSubscription;
      if (newAlert) {
        setAlerts((prev) => (prev ? [...prev, newAlert] : [newAlert]));
      }
    },
  });

  const handleAlertRead = (index: number) => {
    if (alerts) {
      let alertId = alerts[index].id;
      updateNotificationCall({
        variables: {
          updateNotificationId: alertId,
          isRead: true,
        },
      });
      setAlerts((prev: any) => {
        if (!prev) return prev;
        const newAlerts = [...prev];
        newAlerts[index] = { ...newAlerts[index], isRead: true };
        return newAlerts;
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
          alerts={[]}
          page={pageName}
          isDarkMode={isDarkMode}
          handlePageChange={() => {}}
          isLoading={networksLoading}
          placeholder={'Select Network'}
          handleAlertRead={handleAlertRead}
          handleAddNetwork={handleAddNetworkAction}
          handleNetworkChange={handleNetworkChange}
          networks={networksData?.getNetworks.networks ?? []}
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
