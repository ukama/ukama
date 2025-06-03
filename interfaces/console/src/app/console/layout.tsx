/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  GetNodesDocument,
  GetNodesQuery,
  NodeConnectivityEnum,
  Role_Type,
  useAddNetworkMutation,
  useGetNetworksQuery,
  useSetDefaultNetworkMutation,
  useUpdateNotificationMutation,
} from '@/client/graphql/generated';
import {
  NotificationsRes,
  useGetNotificationsLazyQuery,
} from '@/client/graphql/generated/subscriptions';
import AddNetworkDialog from '@/components/AddNetworkDialog';
import AppSnackbar from '@/components/AppSnackbar/page';
import AppLayout from '@/components/Layout';
import { useAppContext } from '@/context';
import ServerNotificationSubscription from '@/lib/NotificationSubscription';
import '@/styles/console.css';
import { TNotificationResDto } from '@/types';
import ErrorBoundary from '@/wrappers/errorBoundary';
import { ApolloClient, useApolloClient } from '@apollo/client';
import { Box } from '@mui/material';
import PubSub from 'pubsub-js';
import { useEffect, useState } from 'react';

export default function ConosleLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const {
    env,
    user,
    network,
    isDarkMode,
    setNetwork,
    setSnackbarMessage,
    subscriptionClient,
  } = useAppContext();
  const client = useApolloClient();
  const [notifications, setNotifications] = useState<NotificationsRes>({
    notifications: [],
  });
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
        setNotifications(res.data?.getNotifications);
      });
    },
  });

  const [getNotifications, { refetch: refetchNotifications }] =
    useGetNotificationsLazyQuery({
      fetchPolicy: 'network-only',
      onCompleted: (data) => {
        if (data.getNotifications.notifications.length > 0) {
          setNotifications(data.getNotifications);
        }
      },
      onError: () => {},
    });

  useEffect(() => {
    if (user.role === Role_Type.RoleInvalid) {
      window.location.reload();
    }
  }, []);

  useEffect(() => {
    if (user.id && network.id && user.orgId && user.orgName) {
      const startTimeStamp = new Date().getTime().toString();
      getNotifications({
        client: subscriptionClient,
        variables: {
          userId: user.id,
          subscriberId: '',
          orgId: user.orgId,
          orgName: user.orgName,
          networkId: network.id,
          role: user.role,
          startTimestamp: startTimeStamp,
        },
      }).then(() => {
        ServerNotificationSubscription(
          env.METRIC_URL,
          `notification-${user.orgId}-${user.id}-${user.role}-${network.id}`,
          user.role,
          user.orgId,
          user.id,
          user.orgName,
          network.id,
          startTimeStamp,
        );
      });

      PubSub.subscribe(
        `notification-${user.orgId}-${user.id}-${user.role}-${network.id}`,
        handleNotification,
      );
    }
  }, [user.id, network.id]);

  const handleNotification = (_: any, data: string) => {
    const parsedData: TNotificationResDto = JSON.parse(data);
    const {
      id,
      type,
      scope,
      title,
      isRead,
      eventKey,
      createdAt,
      resourceId,
      description,
      redirect: { action, title: redirectTitle },
    } = parsedData.data.notificationSubscription;

    const newNotification = {
      id,
      type,
      scope,
      title,
      isRead,
      eventKey,
      createdAt,
      resourceId,
      description,
      redirect: {
        action,
        title: redirectTitle,
      },
    };

    setNotifications((prev) => {
      return {
        notifications: [newNotification, ...prev.notifications].filter(
          (v, i, a) => a.findIndex((t) => t.id === v.id) === i,
        ),
      };
    });

    if (eventKey === 'EventNodeOnline' || eventKey === 'EventNodeOffline') {
      (client as ApolloClient<any>).cache.updateQuery<GetNodesQuery>(
        { query: GetNodesDocument },
        (data: GetNodesQuery | null) => {
          if (!data?.getNodes?.nodes) return data;

          const updatedNodes = data.getNodes.nodes.map(
            (node: GetNodesQuery['getNodes']['nodes'][0]) => {
              if (node.id === resourceId) {
                return {
                  ...node,
                  status: {
                    ...node.status,
                    connectivity:
                      eventKey === 'EventNodeOnline'
                        ? NodeConnectivityEnum.Online
                        : NodeConnectivityEnum.Offline,
                  },
                };
              }
              return node;
            },
          );

          return {
            ...data,
            getNodes: {
              ...data.getNodes,
              nodes: updatedNodes,
            },
          };
        },
      );
    }
  };

  const handleNotificationAction = (action: string, id: string) => {
    switch (action) {
      case 'mark-read':
        updateNotificationCall({
          variables: {
            isRead: true,
            updateNotificationId: id,
          },
        });
        break;
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
      const filterNetwork = networksData?.getNetworks.networks.find(
        (n) => n.id === id,
      );
      setNetwork({
        id: id,
        name: filterNetwork?.name ?? '',
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
          overflow: 'hidden',
          flexDirection: 'column',
        }}
      >
        <AppLayout
          isDarkMode={isDarkMode}
          isLoading={networksLoading}
          placeholder={'Select Network'}
          handleAction={handleNotificationAction}
          handleAddNetwork={handleAddNetworkAction}
          handleNetworkChange={handleNetworkChange}
          networks={networksData?.getNetworks.networks ?? []}
          notifications={notifications}
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
