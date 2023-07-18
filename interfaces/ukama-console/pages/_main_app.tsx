'use client';
import {
  commonData,
  isDarkmode,
  isSkeltonLoading,
  pageName,
  snackbarMessage,
  user,
} from '@/app-recoil';
import { useGetNetworksLazyQuery, useWhoamiLazyQuery } from '@/generated';
import { MyAppProps, TCommonData, TSnackMessage, TUser } from '@/types';
import { getTitleFromPath } from '@/utils';
import { AlertColor } from '@mui/material';
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
  const [getWhoami, { data, loading, error }] = useWhoamiLazyQuery({
    fetchPolicy: 'cache-first',
    onCompleted: (data) => {
      if (data.whoami) {
        if (!_user?.id)
          _setUser({
            id: data.whoami.id,
            name: data.whoami.name,
            email: data.whoami.email,
            role: data.whoami.role,
            isFirstVisit: data.whoami.isFirstVisit,
          });
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'whoami-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
      // resetData();
      // window.location.replace(`${process.env.NEXT_PUBLIC_REACT_AUTH_APP_URL}`);
    },
  });
  const [
    getNetworks,
    { data: networksData, error: networksError, loading: networksLoading },
  ] = useGetNetworksLazyQuery({
    fetchPolicy: 'cache-first',
    onCompleted: (data) => {
      if (data.getNetworks.networks.length === 1) {
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
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  useEffect(() => {
    setSkeltonLoading(true);
    getWhoami();
    getNetworks();
  }, []);

  useEffect(() => {
    if (
      data?.whoami.id &&
      networksData?.getNetworks?.networks &&
      networksData?.getNetworks?.networks.length > 0
    ) {
      setSkeltonLoading(false);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data, networksData]);

  // useEffect(() => {
  //   const { id, name, email } = _user;
  //   const pathname =
  //     typeof window !== 'undefined' && window.location.pathname
  //       ? window.location.pathname
  //       : '';
  //   setPage(getTitleFromPath(pathname));
  //   if (id && name && email) {
  //     if (
  //       !doesHttpOnlyCookieExist('id') &&
  //       doesHttpOnlyCookieExist('ukama_session')
  //     ) {
  //       resetData();
  //       resetPageName();
  //       window.location.replace(
  //         `${process.env.NEXT_PUBLIC_REACT_AUTH_APP_URL}/logout`,
  //       );
  //     } else if (
  //       doesHttpOnlyCookieExist('id') &&
  //       !doesHttpOnlyCookieExist('ukama_session')
  //     )
  //       handleGoToLogin();
  //   } else {
  //     if (process.env.NEXT_PUBLIC_NODE_ENV === 'test') return;
  //     handleGoToLogin();
  //   }
  // }, []);

  useEffect(() => {
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
  }, [route.pathname]);

  const handleGoToLogin = () => {
    setPage('Home');
    typeof window !== 'undefined' &&
      window.location.replace(process.env.NEXT_PUBLIC_REACT_AUTH_APP_URL || '');
  };

  const handlePageChange = (page: string) => setPage(page);
  const handleNetworkChange = (id: string) =>
    setCommonData({
      ..._commonData,
      networkId: id,
      networkName:
        networksData?.getNetworks.networks.filter((n) => n.id === id)[0].name ??
        '',
    });

  return (
    <Layout
      page={page}
      isFullScreen={isFullScreen}
      isDarkMode={_isDarkMod}
      isLoading={false}
      placeholder={'Select Network'}
      networkId={_commonData?.networkId}
      handlePageChange={handlePageChange}
      handleNetworkChange={handleNetworkChange}
      networks={networksData?.getNetworks.networks}
    >
      <Component {...pageProps} />
    </Layout>
  );
};

export default MainApp;
