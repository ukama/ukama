/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Role_Type } from '@/client/graphql/generated';
import { NotificationsRes } from '@/client/graphql/generated/subscriptions';
import { useAppContext } from '@/context';
import { HorizontalContainerJustify, IconStyle } from '@/styles/global';
import colors from '@/theme/colors';
import ManageAccountsIcon from '@mui/icons-material/ManageAccounts';
import {
  Box,
  Divider,
  IconButton,
  Stack,
  Toolbar,
  Typography,
  styled,
  useMediaQuery,
  useTheme,
} from '@mui/material';
import { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import dynamic from 'next/dynamic';
import AccountPopover from './AccountPopover';
import Alert from './Alert';

const Logo = dynamic(() =>
  import('../../../../public/svg/Logo').then((module) => ({
    default: module.Logo,
  })),
);
const ULogo = dynamic(() =>
  import('../../../../public/svg/ULogo').then((module) => ({
    default: module.ULogo,
  })),
);

interface IHeaderProps {
  isOpen: boolean;
  isLoading: boolean;
  onNavigate: (page: string, path: string) => void;
  notifications: NotificationsRes;
  handleAction: (action: string, id: string) => void;
}

interface AppBarProps extends MuiAppBarProps {
  open: boolean;
}

const AppBar = styled(Box, {
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
    height: 60,
  }),
}));

const Header = ({
  isOpen,
  onNavigate,
  notifications,
  handleAction,
}: IHeaderProps) => {
  const theme = useTheme();
  const { user } = useAppContext();
  const isManager = user.role === Role_Type.RoleOwner;
  const matches = useMediaQuery(theme.breakpoints.up('sm'));

  return (
    <AppBar
      open={isOpen}
      sx={{
        py: 1,
        px: { xs: 1.5, md: 3 },
        height: 'fit-content',
        justifyContent: 'center',
        boxShadow: '2px 2px 6px rgba(0, 0, 0, 0.05)',
      }}
    >
      <Toolbar sx={{ alignSelf: 'center', padding: '0px !important' }}>
        <HorizontalContainerJustify>
          <IconButton onClick={() => onNavigate('Root', '/')}>
            {matches ? (
              <Logo width={'100%'} height={'28px'} color={colors.white} />
            ) : (
              <ULogo width={'100%'} height={'28px'} color={colors.white} />
            )}
          </IconButton>
          <Stack
            direction={'row'}
            alignItems={'center'}
            spacing={{ xs: 1, md: 1.75 }}
          >
            <Typography variant="body1" fontWeight={600} color={colors.white}>
              {user.orgName}
            </Typography>
            <Divider
              orientation="vertical"
              sx={{ width: '0.5px', height: '24px', bgcolor: colors.darkGray }}
            />
            {isManager && (
              <IconButton
                data-testid="manage-btn"
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
            <Alert notifications={notifications} handleAction={handleAction} />
            <AccountPopover />
          </Stack>
        </HorizontalContainerJustify>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
