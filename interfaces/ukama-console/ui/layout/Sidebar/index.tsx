import { NavList } from '@/router/config';
import { colors } from '@/styles/theme';
import { SelectItemType } from '@/types';
import { LoadingWrapper } from '@/ui/components';
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

const BasicDropdown = dynamic(() => import('@/ui/components/BasicDropdown'), {
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
  networkId: string;
  isLoading: boolean;
  isDarkMode: boolean;
  onNavigate: Function;
  networks: SelectItemType[];
  handleNetworkChange: Function;
}

const Sidebar = ({
  page,
  isOpen,
  networks,
  isLoading,
  networkId,
  onNavigate,
  isDarkMode,
  handleNetworkChange,
}: ISidebarProps) => {
  return (
    <UkamaDrawer
      open={isOpen}
      variant="permanent"
      style={{ marginTop: 60, height: '100%' }}
    >
      <LoadingWrapper isLoading={isLoading} radius="none">
        <Stack direction={'column'}>
          <Box mx={{ xs: '18px', md: '28px' }} my={{ xs: 1, md: 1.7 }}>
            <BasicDropdown
              value={networkId}
              networkList={networks}
              isLoading={isLoading}
              handleOnChange={handleNetworkChange}
            />
          </Box>
          <Divider sx={{ m: 0 }} />
          <List component="nav">
            {NavList.map(({ name, path, icon: Icon }) => {
              return (
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
              );
            })}
          </List>
        </Stack>
      </LoadingWrapper>
    </UkamaDrawer>
  );
};

export default Sidebar;
