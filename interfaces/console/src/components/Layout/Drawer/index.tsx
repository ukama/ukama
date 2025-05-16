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
import { colors } from '@/theme';
import MenuIcon from '@mui/icons-material/Menu';
import {
  Divider,
  IconButton,
  Link,
  Skeleton,
  Stack,
  Typography,
} from '@mui/material';
import Box from '@mui/material/Box';
import SwipeableDrawer from '@mui/material/SwipeableDrawer';
import dynamic from 'next/dynamic';
import { usePathname } from 'next/navigation';
import * as React from 'react';
const BasicDropdown = dynamic(() => import('@/components/BasicDropdown'), {
  ssr: false,
  loading: () => <Skeleton variant="rectangular" width={'100%'} height={29} />,
});

interface IUDrawer {
  placeholder: string;
  networks: NetworkDto[];
  handleAddNetwork: () => void;
  handleNetworkChange: (value: string) => void;
}

export default function UDrawer({
  networks,
  placeholder,
  handleAddNetwork,
  handleNetworkChange,
}: IUDrawer) {
  const pathname = usePathname();
  const [anchor, setAnchor] = React.useState(false);
  const { user, isDarkMode, network } = useAppContext();
  const isOwner =
    user.role === Role_Type.RoleOwner || user.role === Role_Type.RoleAdmin;

  const getDropDownData = () =>
    networks?.map((network) => ({
      id: network.id,
      label: network.name,
      value: network.id,
    }));

  const toggleDrawer =
    (open: boolean) => (event: React.KeyboardEvent | React.MouseEvent) => {
      if (
        event &&
        event.type === 'keydown' &&
        ((event as React.KeyboardEvent).key === 'Tab' ||
          (event as React.KeyboardEvent).key === 'Shift')
      ) {
        return;
      }

      setAnchor(open);
    };

  const list = () => (
    <Box
      sx={{ width: 250, marginTop: 8 }}
      onClick={toggleDrawer(false)}
      onKeyDown={toggleDrawer(false)}
    >
      <Box mx={2} my={2}>
        <BasicDropdown
          id={'network-dropdown'}
          value={network.id}
          list={getDropDownData()}
          isShowAddOption={isOwner}
          placeholder={placeholder}
          handleOnChange={handleNetworkChange}
          handleAddNetwork={handleAddNetwork}
        />
      </Box>
      <Divider sx={{ mx: 2, my: 0 }} />
      {NavList.map(
        ({ name, path, icon: Icon, forRoles }) =>
          forRoles.includes(user.role as Role_Type) && (
            <Link
              href={path}
              key={path}
              style={{
                borderRadius: 4,
                textDecoration: 'none',
              }}
            >
              <Stack
                px={2}
                py={2}
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
                    color:
                      pathname === path
                        ? colors.primaryMain
                        : isDarkMode
                          ? colors.white70
                          : colors.vulcan100,
                  }}
                />
                <Typography
                  variant="body1"
                  fontWeight={400}
                  color={
                    pathname === path
                      ? colors.primaryMain
                      : isDarkMode
                        ? colors.white70
                        : colors.vulcan100
                  }
                >
                  {name}
                </Typography>
              </Stack>
            </Link>
          ),
      )}
    </Box>
  );

  return (
    <React.Fragment>
      <IconButton sx={{ p: 0 }} onClick={toggleDrawer(true)}>
        <MenuIcon />
      </IconButton>
      <SwipeableDrawer
        anchor={'left'}
        open={anchor}
        onClose={toggleDrawer(false)}
        onOpen={toggleDrawer(true)}
      >
        {list()}
      </SwipeableDrawer>
    </React.Fragment>
  );
}
