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
  useGetMemberLazyQuery,
  useGetNetworksQuery,
  useGetOrgsLazyQuery,
  useGetUserLazyQuery,
} from '@/generated';
import { MyAppProps, TCommonData, TSnackMessage, TUser } from '@/types';
import AddNetworkDialog from '@/ui/molecules/AddNetworkDialog';
import { doesHttpOnlyCookieExist, getTitleFromPath } from '@/utils';
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

  const [getMember] = useGetMemberLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      _setUser({
        ..._user,
        role: data.getMember.role,
      });
    },
  });

  const [getUser, { data: userData, loading: userLoading, error: userError }] =
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
            memberId: data.getUser.uuid,
          },
        });
      },
    });

  const {
    data: networksData,
    error: networksError,
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

  const [getOrgs, { data: orgsData, error: orgsError, loading: orgsLoading }] =
    useGetOrgsLazyQuery({
      onCompleted: (data) => {},
    });

  useEffect(() => {
    const orgId = route.query['org-id'] as string;
    const orgName = route.query['org-name'] as string;
    const userId = route.query['uid'] as string;

    if (orgId && orgName) {
      setCommonData({
        metaData: {},
        orgId,
        orgName,
        userId,
        networkId: '',
        networkName: '',
      });
      route.replace(route.pathname, undefined, { shallow: true });
    } else {
      getOrgs();
      // TODO: NO ORG FOUND AGAINST USER, Navigate to Org selection screen/ Invitation screen
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
    if (!_commonData.userId || (!_commonData.orgId && !_commonData.orgName)) {
      route.push('/unauthorized');
    }

    if (_commonData.userId) {
      getUser({
        variables: {
          userId: _commonData.userId,
        },
      });
    }
  }, [_commonData]);

  useEffect(() => {
    if (!userLoading && userData && userData.getUser && _user.id) {
      const { id, name, email } = _user;
      if (id && name && email) {
        if (!doesHttpOnlyCookieExist('ukama_session')) handleGoToLogin();
      }
    } else if (userError) {
      handleGoToLogin();
    }
  }, [_user]);

  useEffect(() => {
    if (networksLoading && orgsLoading && userLoading && !skeltonLoading)
      setSkeltonLoading(true);
    else if (!networksLoading && !orgsLoading && !userLoading && skeltonLoading)
      setSkeltonLoading(false);
  }, [
    networksLoading,
    orgsLoading,
    userLoading,
    skeltonLoading,
    setSkeltonLoading,
  ]);

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
    // const countriesName =
    //   values.countries.length > 0
    //     ? values.countries.map((item: any) => item.name)
    //     : [];
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
    }).finally(() => {
      setShowAddNetwork(false);
    });
  };

  return (
    <Layout
      page={page}
      isFullScreen={isFullScreen}
      isDarkMode={_isDarkMod}
      isLoading={
        networksLoading || orgsLoading || userLoading || skeltonLoading
      }
      placeholder={'Select Network'}
      handlePageChange={handlePageChange}
      handleNetworkChange={handleNetworkChange}
      networks={networksData?.getNetworks.networks || []}
      handleAddNetwork={handleAddNetworkAction}
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
        // networks={networksData?.getNetworks.networks || []}
        handleCloseAction={() => setShowAddNetwork(false)}
      />
    </Layout>
  );
};

export default MainApp;
