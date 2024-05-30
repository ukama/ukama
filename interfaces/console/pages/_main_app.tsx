/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

'use client';
import { useAppContext } from '@/context';
import {
  useAddNetworkMutation,
  useGetNetworksQuery,
  useSetDefaultNetworkMutation,
} from '@/generated';
import { MyAppProps } from '@/types';
import AddNetworkDialog from '@/ui/molecules/AddNetworkDialog';
import { getTitleFromPath } from '@/utils';
import dynamic from 'next/dynamic';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

const Layout = dynamic(() => import('@/ui/layout'), {
  ssr: false,
});

const MainApp = ({ Component, pageProps }: MyAppProps) => {
  const route = useRouter();
  const [isFullScreen, setIsFullScreen] = useState<boolean>(false);
  const [showAddNetwork, setShowAddNetwork] = useState<boolean>(false);
  const {
    token,
    user,
    network,
    setNetwork,
    pageName,
    isDarkMode,
    setPageName,
    skeltonLoading,
    isValidSession,
    snackbarMessage,
    setSnackbarMessage,
    setSkeltonLoading,
  } = useAppContext();

  const {
    data: networksData,
    loading: networksLoading,
    refetch: refetchNetworks,
  } = useGetNetworksQuery({
    skip: user?.orgId === '',
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

  useEffect(() => {
    if (route.pathname) {
      setIsFullScreen(
        route.pathname === '/manage' ||
          route.pathname === '/settings' ||
          route.pathname === '/unauthorized' ||
          route.pathname === '/onboarding' ||
          getTitleFromPath(route.pathname, route.query['id'] as string) ===
            '404',
      );
      setPageName(
        getTitleFromPath(route.pathname, route.query['id'] as string),
      );
      if (
        getTitleFromPath(route.pathname, '') === '404' &&
        route.pathname !== '/404'
      ) {
        route.replace('/404');
      }
    }
  }, [route]);

  useEffect(() => {
    if (networksLoading && !skeltonLoading) setSkeltonLoading(true);
    else if (!networksLoading && skeltonLoading) setSkeltonLoading(false);
  }, [networksLoading, skeltonLoading, setSkeltonLoading]);

  const handlePageChange = (page: string) => setPageName(page);

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

  const handleAddNetworkAction = () => setShowAddNetwork(true);

  const handleAddNetwork = (values: any) => {
    addNetwork({
      variables: {
        data: {
          name: values.name,
          budget: values.budget,
          networks: values.networks,
          org: user.orgName,
          countries: values.countries,
        },
      },
    })
      // .then((res) => {
      //   if (values.isDefault && res.data?.addNetwork.id) {
      //     setDefaultNetwork({
      //       variables: {
      //         data: {
      //           id: res.data?.addNetwork.id,
      //         },
      //       },
      //     });
      //   }
      // })
      .finally(() => {
        setShowAddNetwork(false);
      });
  };

  return (
    <Layout
      page={pageName}
      isDarkMode={isDarkMode}
      isFullScreen={isFullScreen}
      placeholder={'Select Network'}
      handlePageChange={handlePageChange}
      handleNetworkChange={handleNetworkChange}
      handleAddNetwork={handleAddNetworkAction}
      networks={networksData?.getNetworks.networks || []}
      isLoading={networksLoading || skeltonLoading}
    >
      <Component {...pageProps} />
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
    </Layout>
  );
};

export default MainApp;
