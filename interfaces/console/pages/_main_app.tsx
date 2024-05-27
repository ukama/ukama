/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

'use client';
import {
  commonData,
  isDarkmode,
  isSkeltonLoading,
  pageName,
  snackbarMessage,
  user,
} from '@/app-recoil';
import {
  useAddNetworkMutation,
  useGetMemberByUserIdLazyQuery,
  useGetNetworksQuery,
  useGetUserLazyQuery,
  useSetDefaultNetworkMutation,
} from '@/generated';
import { MyAppProps, TCommonData, TSnackMessage, TUser } from '@/types';
import AddNetworkDialog from '@/ui/molecules/AddNetworkDialog';
import { getTitleFromPath } from '@/utils';
import dynamic from 'next/dynamic';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';
import { useRecoilState, useRecoilValue, useResetRecoilState } from 'recoil';

const Layout = dynamic(() => import('@/ui/layout'), {
  ssr: false,
});

const MainApp = ({ Component, pageProps }: MyAppProps) => {
  const route = useRouter();
  const [isFullScreen, setIsFullScreen] = useState<boolean>(false);
  const [_user, _setUser] = useRecoilState<TUser>(user);
  const [page, setPage] = useRecoilState(pageName);
  const _isDarkMod = useRecoilValue<boolean>(isDarkmode);
  const [_snackbarMessage, setSnackbarMessage] =
    useRecoilState<TSnackMessage>(snackbarMessage);
  const [skeltonLoading, setSkeltonLoading] =
    useRecoilState<boolean>(isSkeltonLoading);
  const [showAddNetwork, setShowAddNetwork] = useState<boolean>(false);
  const [_commonData, setCommonData] = useRecoilState<TCommonData>(commonData);
  const resetData = useResetRecoilState(user);
  const resetPageName = useResetRecoilState(pageName);

  const [getMember] = useGetMemberByUserIdLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      _setUser({
        ..._user,
        role: data.getMemberByUserId.role,
      });
    },
  });

  const [getUser, { data: userData, loading: userLoading }] =
    useGetUserLazyQuery({
      fetchPolicy: 'cache-and-network',
      onCompleted: (data) => {
        _setUser({
          role: '',
          isFirstVisit: false,
          id: data.getUser.uuid,
          name: data.getUser.name,
          email: data.getUser.email,
        });
        const pathname =
          typeof window !== 'undefined' && window.location.pathname
            ? window.location.pathname
            : '';
        setPage(
          getTitleFromPath(pathname, (route.query['id'] as string) || ''),
        );
        getMember({
          variables: {
            userId: data.getUser.uuid,
          },
        });
      },
    });

  const {
    data: networksData,
    loading: networksLoading,
    refetch: refetchNetworks,
  } = useGetNetworksQuery({
    skip: _commonData?.orgId === '',
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (
        data.getNetworks.networks.length >= 1 &&
        _commonData.networkId === ''
      ) {
        setCommonData({
          ..._commonData,
          networkId: data.getNetworks.networks[0].id,
          networkName: data.getNetworks.networks[0].name,
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
    const userId = route.query['uid'] as string;
    const orgId = route.query['org-id'] as string;
    const orgName = route.query['org-name'] as string;
    if (userId) {
      getUser({
        variables: {
          userId: userId,
        },
      });
    }
    setCommonData({
      metaData: {},
      orgId,
      orgName,
      userId,
      networkId: '',
      networkName: '',
    });
    if (!orgId && !orgName) {
      route.replace('/onboarding', undefined, { shallow: true });
    } else {
      route.replace(route.pathname, undefined, { shallow: true });
    }

    if (route.pathname) {
      setIsFullScreen(
        route.pathname === '/manage' ||
          route.pathname === '/settings' ||
          route.pathname === '/unauthorized' ||
          route.pathname === '/onboarding' ||
          getTitleFromPath(route.pathname, route.query['id'] as string) ===
            '404',
      );
      setPage(getTitleFromPath(route.pathname, route.query['id'] as string));
      if (
        getTitleFromPath(route.pathname, '') === '404' &&
        route.pathname !== '/404'
      ) {
        route.replace('/404');
      }
    }
  }, [route]);

  useEffect(() => {
    if (networksLoading && userLoading && !skeltonLoading)
      setSkeltonLoading(true);
    else if (!networksLoading && !userLoading && skeltonLoading)
      setSkeltonLoading(false);
  }, [networksLoading, userLoading, skeltonLoading, setSkeltonLoading]);

  const handleGoToLogin = () => {
    resetData();
    resetPageName();
    typeof window !== 'undefined' &&
      window.location.replace(process.env.NEXT_PUBLIC_AUTH_APP_URL || '');
  };

  const handlePageChange = (page: string) => setPage(page);

  const handleNetworkChange = (id: string) => {
    if (id) {
      setCommonData({
        ..._commonData,
        networkId: id,
        networkName:
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
          org: _commonData.orgName,
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
      page={page}
      isDarkMode={_isDarkMod}
      isFullScreen={isFullScreen}
      placeholder={'Select Network'}
      handlePageChange={handlePageChange}
      handleNetworkChange={handleNetworkChange}
      handleAddNetwork={handleAddNetworkAction}
      networks={networksData?.getNetworks.networks || []}
      isLoading={networksLoading || userLoading || skeltonLoading}
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
