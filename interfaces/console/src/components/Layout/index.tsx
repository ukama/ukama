/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NetworkDto } from '@/client/graphql/generated';
import { NotificationsRes } from '@/client/graphql/generated/subscriptions';
import { useAppContext } from '@/context';
import { getTitleFromPath } from '@/utils';
import { Divider, Stack, Typography, useMediaQuery } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import { usePathname, useRouter } from 'next/navigation';
import React, { useEffect } from 'react';
import UDrawer from './Drawer';
import Header from './Header';
import Sidebar from './Sidebar';

interface ILayoutProps {
  isLoading: boolean;
  placeholder: string;
  isDarkMode: boolean;
  networks: NetworkDto[];
  children: React.ReactNode;
  handleAddNetwork: () => void;
  notifications: NotificationsRes;
  handleNetworkChange: (value: string) => void;
  handleAction: (action: string, id: string) => void;
}

const isHaveId = (pathname: string) => {
  const parts = pathname.split('/');
  return pathname.startsWith('/console/') && parts.length > 3;
};

const AppLayout = ({
  children,
  networks,
  isLoading,
  isDarkMode,
  placeholder,
  handleAction,
  notifications,
  handleAddNetwork,
  handleNetworkChange,
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
    router.push(path);
  };
  const dynamicId = pathname.startsWith('/console/sites/')
    ? selectedDefaultSite
    : id;
  return (
    <Stack direction={'column'} height={'100%'}>
      <Header
        isOpen={open}
        isLoading={isLoading}
        onNavigate={onNavigate}
        notifications={notifications}
        handleAction={handleAction}
      />
      <Stack height={'100%'} direction={'row'} spacing={2}>
        {!matches && (
          <Sidebar
            isOpen={open}
            isDarkMode={isDarkMode}
            placeholder={placeholder}
            networks={networks ?? []}
            handleAddNetwork={handleAddNetwork}
            handleNetworkChange={handleNetworkChange}
          />
        )}
        <Stack
          width={'100%'}
          height={'100%'}
          overflow={'hidden'}
          direction={'column'}
          pt={{ xs: 1, md: 2 }}
          px={{ xs: 2, md: 3 }}
        >
          <Stack direction={'row'} spacing={{ xs: 2, md: 0 }}>
            {matches && (
              <UDrawer
                placeholder={placeholder}
                networks={networks ?? []}
                handleAddNetwork={handleAddNetwork}
                handleNetworkChange={handleNetworkChange}
              />
            )}
            <Typography variant="h5" fontWeight={400}>
              {getTitleFromPath(pathname, dynamicId)}
            </Typography>
          </Stack>
          <Divider sx={{ mb: 1 }} />
          {children}
        </Stack>
      </Stack>
    </Stack>
  );
};

export default AppLayout;
