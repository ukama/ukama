import { NavList } from '@/router/config';
import { colors } from '@/styles/theme';
import { LoadingWrapper } from '@/ui/components';
import {
  Drawer,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  styled,
} from '@mui/material';

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
  onNavigate: Function;
}

const Sidebar = ({
  page,
  isOpen,
  isLoading,
  onNavigate,
  isDarkMode,
}: ISidebarProps) => {
  return (
    <UkamaDrawer
      open={isOpen}
      variant="permanent"
      style={{ marginTop: 64, height: '100%' }}
    >
      <LoadingWrapper isLoading={isLoading} radius="none">
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
      </LoadingWrapper>
    </UkamaDrawer>
  );
};

export default Sidebar;
