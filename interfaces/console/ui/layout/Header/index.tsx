/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { useAppContext } from '@/context';
import { HorizontalContainerJustify } from '@/styles/global';
import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import ManageAccountsIcon from '@mui/icons-material/ManageAccounts';
import NotificationsIcon from '@mui/icons-material/Notifications';
import SettingsIcon from '@mui/icons-material/Settings';
import { Badge, IconButton, Stack, Toolbar, styled } from '@mui/material';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import dynamic from 'next/dynamic';
import Alert from './Alert';
import { NotificationsResDto } from '@/generated/metrics';

const Logo = dynamic(() =>
  import('../../../public/svg/Logo').then((module) => ({
    default: module.Logo,
  })),
);

interface IHeaderProps {
  isOpen: boolean;
  isLoading: boolean;
  isDarkMode: boolean;
  onNavigate: Function;
  alerts:NotificationsResDto[] | undefined;
  setAlerts: Function
  handleAlertRead: (index: number) => void
}

interface AppBarProps extends MuiAppBarProps {
  open: boolean;
  isloading: string;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open, isloading }) => ({
  zIndex: theme.zIndex.drawer + 1,
  boxShadow: 'none',
  ...(theme.palette.mode === 'dark' && {
    backgroundImage: 'none',
    backgroundColor: 'none',
    background: isloading === 'true' ? colors.nightGrey5 : 'none',
  }),
  ...(theme.palette.mode === 'light' && {
    backgroundColor: 'none',
    backgroundImage: 'none',
    background: isloading === 'true' ? colors.white : colors.darkBlueGradiant,
  }),
  ...(open && {
    width: '100%',
    height: 60,
  }),
}));

const IconStyle = {
  '.MuiSvgIcon-root': {
    width: '24px',
    height: '24px',
    fill: colors.white,
  },
  '.MuiBadge-root': {
    '.MuiSvgIcon-root': {
      width: '24px',
      height: '24px',
      fill: colors.white,
    },
  },
};

const Header = ({ onNavigate, isLoading, isOpen, alerts, setAlerts, handleAlertRead }: IHeaderProps) => {
  const { user, setUser } = useAppContext();
  const isManager =
    user.role === 'ADMIN' || user.role === 'OWNER' ? true : false;
  return (
    <AppBar
      open={isOpen}
      isloading={`${isLoading}`}
      sx={{
        py: isLoading ? 0 : 1,
        height: 'fit-content',
        justifyContent: 'center',
        boxShadow: '2px 2px 6px rgba(0, 0, 0, 0.05)',
      }}
    >
      <LoadingWrapper radius="none" isLoading={isLoading} height={48}>
        <Toolbar sx={{ alignSelf: 'center', width: '100%' }}>
          <HorizontalContainerJustify>
            <IconButton onClick={() => onNavigate('Home', '/home')}>
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
              <IconButton
                onClick={() => onNavigate('Settings', '/settings')}
                sx={{ ...IconStyle }}
              >
                <SettingsIcon />
              </IconButton>
              <Alert alerts={alerts} setAlerts={setAlerts} handleAlertRead={handleAlertRead}/>
              <IconButton
                sx={{
                  ...IconStyle,
                }}
              >
                <AccountCircleIcon />
              </IconButton>
            </Stack>
          </HorizontalContainerJustify>
        </Toolbar>
      </LoadingWrapper>
    </AppBar>
  );
};

export default Header;
