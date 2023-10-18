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
  useGetNetworksLazyQuery,
  useGetOrgsLazyQuery,
  useGetUserLazyQuery,
} from '@/generated';
import { MyAppProps, TCommonData, TSnackMessage, TUser } from '@/types';
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
  const [_commonData, setCommonData] = useRecoilState<TCommonData>(commonData);
  const resetData = useResetRecoilState(user);
  const resetPageName = useResetRecoilState(pageName);

  const [getUser] = useGetUserLazyQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      _setUser({
        role: '',
        isFirstVisit: false,
        id: data.getUser.uuid,
        name: data.getUser.name,
        email: data.getUser.email,
      });
    },
  });

  const [
    getNetworks,
    { data: networksData, error: networksError, loading: networksLoading },
  ] = useGetNetworksLazyQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (data.getNetworks.networks.length > 1) {
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
      setSkeltonLoading(false);
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
          getTitleFromPath(route.pathname, route.query['id'] as string) ===
            '404',
      );
      setPage(getTitleFromPath(route.pathname, route.query['id'] as string));
      if (getTitleFromPath(route.pathname, '') === '404') route.replace('/404');
    }
  }, [route]);

  useEffect(() => {
    if (_commonData.userId) {
      getUser({
        variables: {
          userId: _commonData.userId,
        },
      });
    }
    if (_commonData.orgId) {
      setSkeltonLoading(true);
      getNetworks();
    }
  }, [commonData]);

  useEffect(() => {
    const { id, name, email } = _user;
    const pathname =
      typeof window !== 'undefined' && window.location.pathname
        ? window.location.pathname
        : '';
    setPage(getTitleFromPath(pathname, (route.query['id'] as string) || ''));
    if (id && name && email) {
      if (!doesHttpOnlyCookieExist('ukama_session')) handleGoToLogin();
    }
  }, [_user]);

  const handleGoToLogin = () => {
    resetData();
    resetPageName();
    typeof window !== 'undefined' &&
      window.location.replace(process.env.NEXT_PUBLIC_AUTH_APP_URL || '');
  };

  const handlePageChange = (page: string) => setPage(page);
  const handleNetworkChange = (id: string) => {
    setCommonData({
      ..._commonData,
      networkId: id,
      networkName:
        networksData?.getNetworks.networks.filter((n) => n.id === id)[0].name ??
        '',
    });
  };

  return (
    <Layout
      page={page}
      isFullScreen={isFullScreen}
      isDarkMode={_isDarkMod}
      isLoading={false}
      placeholder={'Select Network'}
      handlePageChange={handlePageChange}
      handleNetworkChange={handleNetworkChange}
      networks={networksData?.getNetworks.networks}
    >
      <Component {...pageProps} />
    </Layout>
  );
};

export default MainApp;
