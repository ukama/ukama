import { isDarkmode } from '@/app-recoil';
import { DRAWER_WIDTH, SIDEBAR_MENU1 } from '@/constants';
import { colors } from '@/styles/theme';
import { MenuItemType } from '@/types';
import { LoadingWrapper } from '@/ui/components';
import {
  Box,
  Drawer,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Paper,
  Stack,
  Toolbar,
  Typography,
} from '@mui/material';
import { makeStyles } from '@mui/styles';
import { useRouter } from 'next/router';
import React from 'react';
import { useRecoilValue } from 'recoil';

const Logo = React.lazy(() =>
  import('../../../public/svg/Logo').then((module) => ({
    default: module.Logo,
  })),
);

const useStyles = makeStyles(() => ({
  listItem: {
    opacity: 1,
    height: '40px',
    marginTop: '6px',
    padding: '8px 12px',
    borderRadius: '4px',
  },
  listItemText: {},
  listItemDone: {
    opacity: 1,
    height: '40px',
    marginTop: '8px',
    padding: '8px 12px',
    borderRadius: '4px',
  },
}));

type SidebarProps = {
  page: string;
  isOpen: boolean;
  isLoading: boolean;
  handlePageChange: Function;
  handleDrawerToggle: Function;
};

const Sidebar = (
  {
    page,
    isOpen,
    isLoading,
    handlePageChange,
    handleDrawerToggle,
  }: SidebarProps,
  props: any,
) => {
  const { window } = props;
  const classes = useStyles();
  const router = useRouter();
  const _isDarkMod = useRecoilValue(isDarkmode);

  const MenuList = (list: any) => (
    <List sx={{ padding: '8px 20px' }}>
      {list.map(({ id, title, Icon, route }: MenuItemType) => (
        <ListItem
          key={id}
          onClick={() => {
            handlePageChange(title);
            router.push(route);
          }}
          className={title === page ? classes.listItemDone : classes.listItem}
        >
          <ListItemIcon sx={{ minWidth: '44px' }}>
            <Icon
              fontSize="medium"
              sx={{
                fill: _isDarkMod ? colors.solitude : colors.vulcan,
              }}
            />
          </ListItemIcon>
          <ListItemText>
            <Stack
              direction="row"
              justifyContent="space-between"
              alignItems="center"
              spacing={1}
            >
              <Typography
                variant="body1"
                color={_isDarkMod ? colors.solitude : 'textColorPrimary'}
                fontWeight={title === page ? '600' : 'normal'}
                className={title !== page ? classes.listItemText : ''}
              >
                {title}
              </Typography>
            </Stack>
          </ListItemText>
        </ListItem>
      ))}
    </List>
  );

  const drawer = (
    <Paper
      sx={{
        height: '100%',
        overflowY: 'auto',
      }}
    >
      <Toolbar sx={{ padding: '33px 0px 12px 0px' }}>
        <Logo
          width={'100%'}
          height={'36px'}
          color={_isDarkMod ? colors.white : colors.primaryMain}
        />
      </Toolbar>
      <Stack
        sx={{
          display: 'flex',
          minHeight: '200px',
          height: `calc(100% - 400px)`,
        }}
      >
        {MenuList(SIDEBAR_MENU1)}
      </Stack>
    </Paper>
  );

  const container = undefined;
  // window !== undefined ? () => window().document.body : undefined;

  return (
    <Box
      component="nav"
      sx={{
        flexShrink: { sm: 0 },
        width: { xs: 0, sm: DRAWER_WIDTH },
        boxShadow: '6px 0px 18px rgba(0, 0, 0, 0.06)',
      }}
      aria-label="mailbox folders"
    >
      <LoadingWrapper isLoading={isLoading}>
        <Drawer
          open={isOpen}
          variant="temporary"
          container={container}
          onClose={() => handleDrawerToggle()}
          ModalProps={{
            keepMounted: true,
          }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': {
              boxSizing: 'border-box',
              width: DRAWER_WIDTH,
            },
            borderRight: '0px',
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          open
          variant="permanent"
          sx={{
            display: { xs: 'none', sm: 'block' },
            '& .MuiDrawer-paper': {
              boxSizing: 'border-box',
              width: DRAWER_WIDTH,
              boxShadow: 'none',
            },
          }}
        >
          {drawer}
        </Drawer>
      </LoadingWrapper>
    </Box>
  );
};

export default Sidebar;
