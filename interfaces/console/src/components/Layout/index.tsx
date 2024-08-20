/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NetworkDto } from '@/client/graphql/generated';
import { NotificationsResDto } from '@/client/graphql/generated/subscriptions';
import { useAppContext } from '@/context';
import { getTitleFromPath } from '@/utils';
import { Divider, Stack, Typography, useMediaQuery } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import { usePathname, useRouter } from 'next/navigation';
import React, { useEffect } from 'react';
import Header from './Header';
import Sidebar from './Sidebar';

interface ILayoutProps {
  page: string;
  isLoading: boolean;
  placeholder: string;
  isDarkMode: boolean;
  handlePageChange: Function;
  networks: NetworkDto[];
  children: React.ReactNode;
  handleAddNetwork: Function;
  handleNetworkChange: Function;
  notifications: NotificationsResDto[];
  handleNotificationRead: (id: string) => void;
}

const isHaveId = (pathname: string) => {
  const parts = pathname.split('/');
  return pathname.startsWith('/console/') && parts.length > 3;
};

const AppLayout = ({
  page,
  children,
  networks,
  isLoading,
  isDarkMode,
  placeholder,
  handlePageChange,
  handleAddNetwork,
  handleNetworkChange,
  notifications,
  handleNotificationRead,
}: ILayoutProps) => {
  const pathname = usePathname();
  const id = isHaveId(pathname) ? pathname.split('/')[3] : '';
  const theme = useTheme();
  const router = useRouter();
  const { selectedDefaultSite } = useAppContext();

  const [open, setOpen] = React.useState(true);
  const matches = useMediaQuery(theme.breakpoints.down('md'));

  useEffect(() => {
    if (matches) {
      setOpen(false);
    } else {
      setOpen(true);
    }
  }, [matches]);

  const onNavigate = (name: string, path: string) => {
    handlePageChange(name);
    router.push(path);
  };
  const dynamicId = pathname.startsWith('/console/sites/')
    ? selectedDefaultSite
    : id;
  return (
    <Stack overflow={'hidden'}>
      <Header
        isOpen={open}
        isLoading={isLoading}
        onNavigate={onNavigate}
        notifications={notifications}
        handleNotificationRead={handleNotificationRead}
      />
      <Stack direction={'row'}>
        <Sidebar
          page={page}
          isOpen={open}
          isLoading={isLoading}
          onNavigate={onNavigate}
          isDarkMode={isDarkMode}
          placeholder={placeholder}
          networks={networks ?? []}
          handleAddNetwork={handleAddNetwork}
          handleNetworkChange={handleNetworkChange}
        />
        <Stack
          mr={2}
          ml={30}
          mt={8}
          p={2}
          width={'100%'}
          height={'100%'}
          direction={'column'}
        >
          <Typography variant="h5" fontWeight={400} mb={0.8}>
            {getTitleFromPath(pathname, dynamicId)}
          </Typography>
          <Divider sx={{ mb: 1 }} />
          {children}
        </Stack>
      </Stack>
    </Stack>
  );
};

export default AppLayout;
