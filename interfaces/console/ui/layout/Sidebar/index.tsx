/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { useAppContext } from '@/context';
import { NetworkDto } from '@/generated';
import { NavList } from '@/router/config';
import colors from '@/styles/theme/colors';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import {
  Box,
  Divider,
  Drawer,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Stack,
  styled,
} from '@mui/material';
import dynamic from 'next/dynamic';

const BasicDropdown = dynamic(() => import('@/ui/molecules/BasicDropdown'), {
  ssr: false,
});

const drawerWidth: number = 218;

const UkamaDrawer = styled(Drawer, {
  shouldForwardProp: (prop) => prop !== 'open',
})(({ theme, open }) => ({
  '& .MuiDrawer-paper': {
    paddingTop: 60,
    whiteSpace: 'nowrap',
    width: drawerWidth,
    boxSizing: 'border-box',
    ...(!open && {
      overflowX: 'hidden',
      width: theme.spacing(7.8),
    }),
    [theme.breakpoints.down('md')]: {
      paddingTop: 44,
    },
  },
}));

interface ISidebarProps {
  page: string;
  isOpen: boolean;
  isLoading: boolean;
  isDarkMode: boolean;
  placeholder: string;
  onNavigate: Function;
  networks: NetworkDto[];
  handleAddNetwork: Function;
  handleNetworkChange: Function;
}

const Sidebar = ({
  page,
  isOpen,
  isLoading,
  onNavigate,
  isDarkMode,
  placeholder,
  networks = [],
  handleAddNetwork,
  handleNetworkChange,
}: ISidebarProps) => {
  const { user, network } = useAppContext();
  const getDropDownData = () =>
    networks?.map((network) => ({
      id: network.id,
      label: network.name,
      value: network.id,
    }));

  return (
    <UkamaDrawer
      open={isOpen}
      variant="permanent"
      style={{ marginTop: 46, height: '100%', backgroundColor: 'white' }}
    >
      <LoadingWrapper
        isLoading={isLoading}
        radius="none"
        cstyle={{ backgroundColor: 'white' }}
      >
        <Stack direction={'column'}>
          <Box mx={{ xs: '18px', md: '28px' }} my={{ xs: 1, md: 1.7 }}>
            <BasicDropdown
              value={network.id}
              isLoading={isLoading}
              list={getDropDownData()}
              placeholder={placeholder}
              handleOnChange={handleNetworkChange}
              handleAddNetwork={handleAddNetwork}
            />
          </Box>
          <Divider sx={{ m: 0 }} />
          <List component="nav">
            {NavList.map(({ name, path, icon: Icon }) => (
              <ListItemButton
                key={name}
                onClick={() => onNavigate(name, path)}
                sx={{
                  backgroundColor:
                    page === name
                      ? isDarkMode
                        ? colors.primaryMain02
                        : colors.solitude
                      : 'transparent',
                  '.MuiListItemText-root .MuiTypography-root': {
                    fontWeight: page === name ? 600 : 400,
                  },
                  ':hover': {
                    backgroundColor: isDarkMode
                      ? colors.darkGreen05
                      : colors.solitude,
                  },
                }}
              >
                <ListItemIcon>
                  <Icon />
                </ListItemIcon>
                <ListItemText primary={name} />
              </ListItemButton>
            ))}
          </List>
        </Stack>
      </LoadingWrapper>
    </UkamaDrawer>
  );
};

export default Sidebar;
