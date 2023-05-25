'use client';
import { commonData } from '@/app-recoil';
import client from '@/client/ApolloClient';
import { MyAppProps, TCommonData } from '@/types';
import { ApolloProvider, HttpLink } from '@apollo/client';
import dynamic from 'next/dynamic';
import { RecoilRoot, useRecoilValue } from 'recoil';
import '../styles/global.css';
import ErrorBoundary from '@/ui/wrappers/errorBoundary';

const MainApp = dynamic(() => import('@/pages/_main_app'));

const ClientWrapper = (appProps: MyAppProps) => {
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const httpLink = new HttpLink({
    uri: process.env.NEXT_PUBLIC_REACT_APP_API,
    credentials: 'include',
    headers: {
      'org-id': _commonData.orgId,
      'user-id': _commonData.userId,
      'org-name': _commonData.orgName,
    },
  });

  const getClient = (): any => {
    client.setLink(httpLink);
    return client;
  };

  return (
    <ApolloProvider client={getClient()}>
      <MainApp {...appProps} />
    </ApolloProvider>
  );
};

const RootWrapper = (appProps: MyAppProps) => {
  return (
    <ErrorBoundary>
      <RecoilRoot>
        <ClientWrapper {...appProps} />
      </RecoilRoot>
    </ErrorBoundary>
  );
};

export default RootWrapper;
