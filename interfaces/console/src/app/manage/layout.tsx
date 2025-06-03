/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { useGetNetworksQuery } from '@/client/graphql/generated';
import AppSnackbar from '@/components/AppSnackbar/page';
import BackButton from '@/components/BackButton';
import { useAppContext } from '@/context';
import { MANAGE_MENU_LIST } from '@/routes';
import '@/styles/console.css';
import colors from '@/theme/colors';
import {
  Container,
  Divider,
  Paper,
  Stack,
  SvgIconTypeMap,
  Typography,
  useMediaQuery,
  useTheme,
} from '@mui/material';
import { OverridableComponent } from '@mui/material/OverridableComponent';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import React from 'react';

interface MenuItemProps {
  id: string;
  icon: OverridableComponent<SvgIconTypeMap>;
  name: string;
  path: string;
  isActive: boolean;
  isDarkMode: boolean;
  isCompactView: boolean;
}

interface MenuSectionProps {
  isCompactView: boolean;
  pathname: string;
  isDarkMode: boolean;
}

interface ManageLayoutProps {
  children: React.ReactNode;
}

const MenuItem: React.FC<MenuItemProps> = ({
  id,
  icon: Icon,
  name,
  path,
  isActive,
  isDarkMode,
  isCompactView,
}) => (
  <Link
    id={id}
    href={path}
    data-testid={id}
    prefetch={id === 'manage-members'}
    style={{
      borderRadius: 4,
      textDecoration: 'none',
      backgroundColor: isActive ? colors.white38 : 'transparent',
      minWidth: isCompactView ? 120 : 'auto',
      flex: isCompactView ? '0 0 auto' : undefined,
    }}
  >
    <Stack
      pl={{ xs: 1, md: 2 }}
      pr={{ xs: 2, md: 4 }}
      py={1}
      spacing={{ xs: 1, md: 2 }}
      alignItems={{ xs: 'center', md: 'flex-start' }}
      direction={{ xs: 'column', md: 'row' }}
      sx={{
        '& svg': {
          color: isActive
            ? colors.black
            : isDarkMode
              ? colors.black10
              : colors.black54,
        },
      }}
    >
      <Icon
        sx={{
          color: isActive
            ? isDarkMode
              ? colors.white
              : colors.vulcan100
            : isDarkMode
              ? colors.white70
              : colors.vulcan100,
        }}
      />
      <Typography
        variant="body1"
        sx={{ color: isActive ? colors.black : colors.vulcan }}
      >
        {name}
      </Typography>
    </Stack>
  </Link>
);

const MenuSection: React.FC<MenuSectionProps> = ({
  isCompactView,
  pathname,
  isDarkMode,
}) => (
  <Stack
    spacing={1}
    sx={{
      width: '100%',
      flexDirection: isCompactView ? 'row' : 'column',
      overflowX: isCompactView ? 'auto' : 'visible',
      WebkitOverflowScrolling: 'touch',
      scrollbarWidth: 'none',
      '&::-webkit-scrollbar': { display: 'none' },
      gap: isCompactView ? 2 : 1,
    }}
  >
    {MANAGE_MENU_LIST.map(({ id, icon, name, path }) => (
      <MenuItem
        key={id}
        id={id}
        icon={icon}
        name={name}
        path={path}
        isActive={pathname === path}
        isDarkMode={isDarkMode}
        isCompactView={isCompactView}
      />
    ))}
  </Stack>
);

const ManageLayout: React.FC<ManageLayoutProps> = ({ children }) => {
  const theme = useTheme();
  const pathname = usePathname();
  const isCompactView = useMediaQuery(theme.breakpoints.down('md'));
  const { isDarkMode, network, setNetwork } = useAppContext();

  useGetNetworksQuery({
    fetchPolicy: 'network-only',
    skip: !!network.id,
    onCompleted: (data) => {
      if (data.getNetworks.networks.length > 0)
        setNetwork({
          id: data.getNetworks.networks[0].id,
          name: data.getNetworks.networks[0].name,
        });
    },
  });

  return (
    <Container maxWidth={'xl'} sx={{ my: { xs: 2, md: 8 } }}>
      <Stack direction={'column'} spacing={{ xs: 2, md: 4 }}>
        <Stack direction="row" spacing={14} alignItems="center">
          <BackButton title="BACK TO CONSOLE" />
          <Typography
            sx={{
              fontSize: { xs: '1.2rem', md: '1.5rem' },
            }}
          >
            Settings
          </Typography>
        </Stack>

        <Divider />

        <Stack
          direction={{ xs: 'column', md: 'row' }}
          spacing={{ xs: 2, md: 4 }}
        >
          <Paper
            sx={{
              display: 'flex',
              borderRadius: '10px',
              width: { xs: '100%', md: '300px' },
              height: 'fit-content',
              overflowX: isCompactView ? 'auto' : 'visible',
            }}
          >
            <Stack
              px={{ xs: 1.4, md: 1 }}
              py={{ xs: 1.5, md: 3 }}
              spacing={{ xs: 0.5, md: 1.5 }}
              direction={{ xs: 'row', md: 'column' }}
              width="100%"
              sx={{
                overflowX: isCompactView ? 'auto' : 'visible',
                WebkitOverflowScrolling: 'touch',
                scrollbarWidth: 'none',
                '&::-webkit-scrollbar': { display: 'none' },
              }}
            >
              <MenuSection
                isCompactView={isCompactView}
                pathname={pathname}
                isDarkMode={isDarkMode}
              />
            </Stack>
          </Paper>
          {children}
          <AppSnackbar />
        </Stack>
      </Stack>
    </Container>
  );
};

export default ManageLayout;
