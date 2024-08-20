/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import {
  NotificationsResDto,
  Role_Type,
} from '@/client/graphql/generated/subscriptions';
import { useAppContext } from '@/context';
import { HorizontalContainerJustify, IconStyle } from '@/styles/global';
import colors from '@/theme/colors';
import ManageAccountsIcon from '@mui/icons-material/ManageAccounts';
import { IconButton, Stack, Toolbar, styled } from '@mui/material';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import dynamic from 'next/dynamic';
import AccountPopover from './AccountPopover';
import Alert from './Alert';

const Logo = dynamic(() =>
  import('../../../../public/svg/Logo').then((module) => ({
    default: module.Logo,
  })),
);

interface IHeaderProps {
  isOpen: boolean;
  isLoading: boolean;
  onNavigate: Function;
  notifications: NotificationsResDto[];
  handleNotificationRead: (id: string) => void;
}

interface AppBarProps extends MuiAppBarProps {
  open: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
  zIndex: theme.zIndex.drawer + 1,
  boxShadow: 'none',
  ...(theme.palette.mode === 'dark' && {
    backgroundImage: 'none',
    backgroundColor: 'none',
    background: 'none',
  }),
  ...(theme.palette.mode === 'light' && {
    backgroundColor: 'none',
    backgroundImage: 'none',
    background: colors.darkBlueGradiant,
  }),
  ...(open && {
    width: '100%',
    height: 60,
  }),
}));

const Header = ({
  isOpen,
  isLoading,
  onNavigate,
  notifications,
  handleNotificationRead,
}: IHeaderProps) => {
  const { user } = useAppContext();
  const isManager =
    user.role === Role_Type.RoleOwner || user.role === Role_Type.RoleAdmin;

  return (
    <AppBar
      open={isOpen}
      sx={{
        py: 1,
        px: 3,
        height: 'fit-content',
        justifyContent: 'center',
        boxShadow: '2px 2px 6px rgba(0, 0, 0, 0.05)',
      }}
    >
      <Toolbar sx={{ alignSelf: 'center', width: '100%' }}>
        <HorizontalContainerJustify>
          <IconButton onClick={() => onNavigate('Root', '/')}>
            <Logo width={'100%'} height={'28px'} color={colors.white} />
          </IconButton>
          <Stack direction={'row'} alignItems={'center'} spacing={1.75}>
            {isManager && (
              <IconButton
                onClick={() => onNavigate('Manage', '/manage')}
                sx={{
                  ...IconStyle,
                  '.MuiSvgIcon-root': {
                    width: '28px',
                    height: '28px',
                    fill: colors.white,
                  },
                }}
              >
                <ManageAccountsIcon />
              </IconButton>
            )}
            {/* <NotificationSubscription /> */}
            <Alert
              alerts={notifications}
              handleNotificationRead={handleNotificationRead}
            />
            <AccountPopover />
          </Stack>
        </HorizontalContainerJustify>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
