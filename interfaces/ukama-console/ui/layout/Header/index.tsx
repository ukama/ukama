import { HorizontalContainerJustify } from '@/styles/global';
import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import ManageAccountsIcon from '@mui/icons-material/ManageAccountsSharp';
import NotificationsIcon from '@mui/icons-material/Notifications';
import SettingsIcon from '@mui/icons-material/Settings';
import { Badge, IconButton, Stack, Toolbar, styled } from '@mui/material';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import dynamic from 'next/dynamic';

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

const Header = ({ onNavigate, isLoading, isOpen }: IHeaderProps) => {
  return (
    <AppBar
      open={isOpen}
      isloading={`${isLoading}`}
      sx={{ justifyContent: 'center' }}
    >
      <LoadingWrapper
        radius="none"
        isLoading={isLoading}
        cstyle={{ display: 'flex' }}
        height={isOpen ? '60px' : '44px'}
      >
        <Toolbar sx={{ alignSelf: 'center', width: '100%' }}>
          <HorizontalContainerJustify>
            <IconButton onClick={() => onNavigate('Home', '/home')}>
              <Logo width={'100%'} height={'28px'} color={colors.white} />
            </IconButton>
            <Stack direction={'row'} spacing={1.75}>
              <IconButton
                onClick={() => onNavigate('Manage', '/manage')}
                sx={{ ...IconStyle }}
              >
                <ManageAccountsIcon />
              </IconButton>
              <IconButton
                onClick={() => onNavigate('Setting', '/setting')}
                sx={{ ...IconStyle }}
              >
                <SettingsIcon />
              </IconButton>
              <IconButton sx={{ ...IconStyle }}>
                <Badge badgeContent={4} color="secondary">
                  <NotificationsIcon />
                </Badge>
              </IconButton>
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
