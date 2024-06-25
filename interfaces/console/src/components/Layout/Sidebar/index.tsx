/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NetworkDto } from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import { NavList } from '@/routes';
import colors from '@/theme/colors';
import {
  Box,
  Divider,
  Drawer,
  Skeleton,
  Stack,
  Typography,
  styled,
} from '@mui/material';
import dynamic from 'next/dynamic';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const BasicDropdown = dynamic(() => import('@/components/BasicDropdown'), {
  ssr: false,
  loading: () => <Skeleton variant="rectangular" width={'100%'} height={44} />,
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
  const pathname = usePathname();
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
      style={{ height: '100%', backgroundColor: 'white' }}
    >
      <Stack direction={'column'}>
        <Box mx={2} my={1.5}>
          <BasicDropdown
            value={network.id}
            isLoading={isLoading}
            list={getDropDownData()}
            placeholder={placeholder}
            handleOnChange={handleNetworkChange}
            handleAddNetwork={handleAddNetwork}
          />
        </Box>
        <Divider sx={{ mx: 2, my: 0 }} />
        <Stack direction="column" spacing={1.5} px={2} py={2}>
          {NavList.map(({ name, path, icon: Icon }) => (
            <Link
              key={path}
              href={path}
              prefetch={true}
              style={{
                borderRadius: 4,
                textDecoration: 'none',
                backgroundColor:
                  pathname === path ? colors.white38 : 'transparent',
              }}
            >
              <Stack
                px={2}
                py={1}
                spacing={2}
                direction={'row'}
                alignItems={'flex-start'}
              >
                <Icon
                  sx={{
                    color: isDarkMode ? colors.white70 : colors.vulcan100,
                  }}
                />
                <Typography
                  variant="body1"
                  fontWeight={400}
                  color={colors.vulcan100}
                >
                  {name}
                </Typography>
              </Stack>
            </Link>
          ))}
        </Stack>
      </Stack>
    </UkamaDrawer>
  );
};

export default Sidebar;
