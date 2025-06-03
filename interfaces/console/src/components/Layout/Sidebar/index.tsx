/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NetworkDto, Role_Type } from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import { NavList } from '@/routes';
import colors from '@/theme/colors';
import { Box, Divider, Skeleton, Stack, Typography } from '@mui/material';
import dynamic from 'next/dynamic';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const BasicDropdown = dynamic(() => import('@/components/BasicDropdown'), {
  ssr: false,
  loading: () => <Skeleton variant="rectangular" width={'100%'} height={29} />,
});

const drawerWidth: number = 232;

interface ISidebarProps {
  isOpen: boolean;
  isDarkMode: boolean;
  placeholder: string;
  networks: NetworkDto[];
  handleAddNetwork: () => void;
  handleNetworkChange: (value: string) => void;
}

const Sidebar = ({
  isDarkMode,
  placeholder,
  networks = [],
  handleAddNetwork,
  handleNetworkChange,
}: ISidebarProps) => {
  const pathname = usePathname();
  const { network, user } = useAppContext();
  const isOwner =
    user.role === Role_Type.RoleOwner || user.role === Role_Type.RoleAdmin;
  const getDropDownData = () =>
    networks?.map((network) => ({
      id: network.id,
      label: network.name,
      value: network.id,
    }));

  return (
    <Box
      sx={{
        height: '100%',
        width: drawerWidth,
        whiteSpace: 'nowrap',
        boxSizing: 'border-box',
        backgroundColor: 'white',
      }}
    >
      <Stack direction={'column'}>
        <Box mx={2} my={2}>
          <BasicDropdown
            value={network.id}
            id={'create-network'}
            list={getDropDownData()}
            isShowAddOption={isOwner}
            placeholder={placeholder}
            handleOnChange={handleNetworkChange}
            handleAddNetwork={handleAddNetwork}
            isDisableAddOption={networks.length >= 3}
          />
        </Box>
        <Divider sx={{ mx: 2, my: 0 }} />
        <Stack direction="column" spacing={1.5} px={2} py={2}>
          {NavList.map(
            ({ name, path, icon: Icon, forRoles }) =>
              forRoles.includes(user.role as Role_Type) && (
                <Link
                  href={path}
                  key={path}
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
                    sx={{
                      ':hover': {
                        borderRadius: '4px',
                        backgroundColor: colors.white38,
                      },
                    }}
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
              ),
          )}
        </Stack>
      </Stack>
    </Box>
  );
};

export default Sidebar;
