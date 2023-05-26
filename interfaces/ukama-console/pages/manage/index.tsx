import { MANAGE_MENU_LIST, MANAGE_TABLE_COLUMN } from '@/constants';
import { colors } from '@/styles/theme';
import { SimpleDataTable } from '@/ui/components';
import PageContainerHeader from '@/ui/components/PageContainerHeader';
import {
  ListItemIcon,
  ListItemText,
  MenuItem,
  MenuList,
  Paper,
  Stack,
} from '@mui/material';
import { useEffect, useState } from 'react';

const DATA = [
  {
    id: 'testuuid',
    name: 'Salman',
    role: 'admin',
    email: 'salman@ukama.com',
  },
  {
    id: 'testuuid1',
    name: 'Brackly',
    role: 'admin',
    email: 'brackley@ukama.com',
  },
];

interface IMemberContainer {
  data: any;
  search: string;
  setSearch: (value: string) => void;
}

interface IManageMenu {
  selectedId: string;
  onMenuItemClick: (id: string) => void;
}

const ManageMenu = ({ selectedId, onMenuItemClick }: IManageMenu) => (
  <Paper
    sx={{
      py: 2,
      px: 2,
      width: 248,
      height: 230,
      borderRadius: '10px',
      maxWidth: 'fit-content',
    }}
  >
    <MenuList sx={{ p: 0 }}>
      {MANAGE_MENU_LIST.map(({ id, icon: Icon, name, path }) => (
        <MenuItem
          key={id}
          href={path}
          sx={{
            py: 1,
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
          <ListItemIcon sx={{ mr: 1 }}>
            <Icon />
          </ListItemIcon>
          <ListItemText>{name}</ListItemText>
        </MenuItem>
      ))}
    </MenuList>
  </Paper>
);

const MemberContainer = ({ data, search, setSearch }: IMemberContainer) => (
  <Paper
    sx={{
      py: 3,
      px: 4,
      width: '100%',
      borderRadius: '10px',
    }}
  >
    <PageContainerHeader
      search={search}
      title={'My members'}
      buttonTitle={'Invite member'}
      onSearchChange={(e: string) => setSearch(e)}
      handleButtonAction={() => console.log('Invite member')}
    />
    <br />
    <SimpleDataTable
      dataset={data.members}
      height={`calc(100vh - 350px)`}
      columns={MANAGE_TABLE_COLUMN}
    />
  </Paper>
);

const Manage = () => {
  const [menu, setMenu] = useState<string>('');
  const [search, setSearch] = useState<string>('');
  const [data, setData] = useState<any>({
    members: DATA,
  });

  useEffect(() => {
    if (search.length > 3) {
      const members = DATA.filter((member) => {
        const s = search.toLowerCase();
        if (member.name.toLowerCase().includes(s)) return member;
      });
      setData({ members: members ?? [] });
    } else if (search.length === 0) {
      setData({ members: DATA });
    }
  }, [search]);

  const onMenuItemClick = (id: string) => setMenu(id);

  return (
    <Stack mt={3} direction={{ xs: 'column', md: 'row' }} spacing={3}>
      <ManageMenu selectedId={menu} onMenuItemClick={onMenuItemClick} />
      {menu === 'manage-members' && (
        <MemberContainer data={data} search={search} setSearch={setSearch} />
      )}
    </Stack>
  );
};
export default Manage;
