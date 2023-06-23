import { SETTING_MENU } from '@/constants';
import colors from '@/styles/theme/colors';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import {
  Divider,
  ListItemText,
  MenuItem,
  MenuList,
  Paper,
  Stack,
} from '@mui/material';
import dynamic from 'next/dynamic';
import { useState } from 'react';

const PersonalSettings = dynamic(() => import('./_personalSetting'));
const Billing = dynamic(() => import('./_billing'));
const NetworkSettings = dynamic(() => import('./_networkSetting'));
const Alerts = dynamic(() => import('./_alerts'));
const ConsoleSettings = dynamic(() => import('./_consoleSetting'));

interface IManageMenu {
  selectedId: string;
  onMenuItemClick: (id: string) => void;
  onLogoutClick: () => void;
}

const ManageMenu = ({
  selectedId,
  onMenuItemClick,
  onLogoutClick,
}: IManageMenu) => (
  <Paper
    sx={{
      py: 2,
      px: 2,
      width: 258,
      maxHeight: 320,
      overflow: 'auto',
      height: 'inderit',
      borderRadius: '4px',
    }}
  >
    <MenuList sx={{ p: 0, width: '100%' }}>
      {SETTING_MENU.map(({ id, name }: any) => (
        <MenuItem
          key={id}
          sx={{
            py: 0.8,
            px: 1.8,
            mb: 1.5,
            borderRadius: '4px',
            backgroundColor:
              selectedId === id ? colors.solitude : 'transparent',
            '.MuiListItemText-root .MuiTypography-root': {
              fontWeight: selectedId === id ? 600 : 400,
            },
          }}
          onClick={() => onMenuItemClick(id)}
        >
          <ListItemText>{name}</ListItemText>
        </MenuItem>
      ))}
      <Divider />
      <MenuItem
        sx={{
          '.MuiListItemText-root .MuiTypography-root': {
            color: colors.error,
          },
        }}
        onClick={() => onLogoutClick()}
      >
        <ListItemText>Logout</ListItemText>
      </MenuItem>
    </MenuList>
  </Paper>
);

export default function Page() {
  const [menu, setMenu] = useState('alerts');
  const [isLoading, setIsLoading] = useState(false);
  const onMenuItemClick = (id: string) => setMenu(id);

  return (
    <Stack mt={3} direction={{ xs: 'column', md: 'row' }} spacing={3}>
      <ManageMenu
        onLogoutClick={() => {}}
        selectedId={menu}
        onMenuItemClick={onMenuItemClick}
      />
      <LoadingWrapper
        width="100%"
        radius="small"
        isLoading={isLoading}
        cstyle={{ backgroundColor: isLoading ? colors.white : 'transparent' }}
      >
        <>
          {menu === 'personal-settings' && <PersonalSettings />}
          {menu === 'billing' && <Billing />}
          {menu === 'network-settings' && <NetworkSettings />}
          {menu === 'alerts' && <Alerts />}
          {menu === 'console-settings' && <ConsoleSettings />}
        </>
      </LoadingWrapper>
    </Stack>
  );
}
