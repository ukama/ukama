import {
  MANAGE_MENU_LIST,
  MANAGE_NODE_POOL_COLUMN,
  MANAGE_SIM_POOL_COLUMN,
  MANAGE_TABLE_COLUMN,
} from '@/constants';
import { colors } from '@/styles/theme';
import { SimpleDataTable } from '@/ui/components';
import PageContainerHeader from '@/ui/components/PageContainerHeader';
import PeopleAlt from '@mui/icons-material/PeopleAlt';
import {
  Grid,
  ListItemIcon,
  ListItemText,
  MenuItem,
  MenuList,
  Paper,
  Stack,
  Typography,
} from '@mui/material';
import { useEffect, useState } from 'react';

const MEMBERS_DATA = [
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

const SIMPOOL_DATA = [
  {
    id: '123',
    iccid: '8910300000003540855',
    type: 'Physical sim',
    qrcode: '1231412414',
  },
  {
    id: '123',
    iccid: '8910300000003540833',
    type: 'E sim',
    qrcode: '123120412414',
  },
];

const NODE_POOL_DATA = [
  {
    type: 'Tower Node',
    dateClaimed: '1231412414',
    id: '8910-3333-0000-3540-833',
  },
  {
    type: 'Amplifier Node',
    dateClaimed: '123120412414',
    id: '8910-3000-0000-3540-833',
  },
];

const DATA_PLAN_DATA = [
  {
    name: 'Monthly Data Plan',
    usage: '$0 / 2 GB / Month',
    users: '4',
  },
  {
    name: 'Weekly Data Plan',
    usage: '$0 / 1 GB / Weekly',
    users: '8',
  },
];

interface IMemberContainer {
  data: any;
  search: string;
  setSearch: (value: string) => void;
}
interface ISimPoolContainer {
  data: any;
}

interface INodePoolContainer {
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
      height: 'calc(100vh - 200px)',
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
    <SimpleDataTable dataset={data} columns={MANAGE_TABLE_COLUMN} />
  </Paper>
);

const SimPoolContainer = ({ data }: ISimPoolContainer) => (
  <Paper
    sx={{
      py: 3,
      px: 4,
      width: '100%',
      borderRadius: '10px',
      height: 'calc(100vh - 200px)',
    }}
  >
    <PageContainerHeader
      subtitle={'2'}
      showSearch={false}
      title={'My SIM pool'}
      buttonTitle={'IMPORT SIMS'}
      handleButtonAction={() => console.log('IMPORT SIMS')}
    />
    <br />
    <SimpleDataTable dataset={data} columns={MANAGE_SIM_POOL_COLUMN} />
  </Paper>
);

const NodePoolContainer = ({ data, search, setSearch }: INodePoolContainer) => (
  <Paper
    sx={{
      py: 3,
      px: 4,
      width: '100%',
      borderRadius: '10px',
      height: 'calc(100vh - 200px)',
    }}
  >
    <PageContainerHeader
      subtitle={'2'}
      search={search}
      title={'My node pool'}
      buttonTitle={'CLAIM NODE'}
      onSearchChange={(e: string) => setSearch(e)}
      handleButtonAction={() => console.log('CLAIM NODE')}
    />
    <br />
    <SimpleDataTable dataset={data} columns={MANAGE_NODE_POOL_COLUMN} />
  </Paper>
);

const DataPlanContainer = ({ data }: ISimPoolContainer) => (
  <Paper
    sx={{
      py: 3,
      px: 4,
      width: '100%',
      borderRadius: '10px',
      height: 'calc(100vh - 200px)',
    }}
  >
    <PageContainerHeader
      showSearch={false}
      title={'Data plans'}
      buttonTitle={'CREATE DATA PLAN'}
      handleButtonAction={() => console.log('IMPORT SIMS')}
    />
    <br />
    <Grid container rowSpacing={2} columnSpacing={2}>
      {data.map(({ name, usage, users }: any) => (
        <Grid item xs={12} sm={6} md={4} key={name}>
          <Paper
            variant="outlined"
            sx={{
              px: 3,
              py: 2,
              display: 'flex',
              boxShadow: 'none',
              borderRadius: '4px',
              textAlign: 'center',
              justifyContent: 'center',
              borderTop: `4px solid ${colors.primaryMain}`,
            }}
          >
            <Stack spacing={1}>
              <Typography variant="h5" sx={{ fontWeight: 400 }}>
                {name}
              </Typography>
              <Typography variant="body2" fontWeight={400}>
                {usage}
              </Typography>
              <Stack
                spacing={0.6}
                direction={'row'}
                alignItems={'flex-end'}
                justifyContent={'center'}
              >
                <PeopleAlt htmlColor={colors.black54} />
                <Typography variant="body2" fontWeight={400}>
                  {users}
                </Typography>
              </Stack>
            </Stack>
          </Paper>
        </Grid>
      ))}
    </Grid>
  </Paper>
);

const Manage = () => {
  const [menu, setMenu] = useState<string>('manage-members');
  const [memberSearch, setMemberSearch] = useState<string>('');
  const [nodeSearch, setNodeSearch] = useState<string>('');
  const [data, setData] = useState<any>({
    members: MEMBERS_DATA,
    simPool: SIMPOOL_DATA,
    node: NODE_POOL_DATA,
    dataPlan: DATA_PLAN_DATA,
  });

  useEffect(() => {
    if (memberSearch.length > 3) {
      const members = MEMBERS_DATA.filter((member) => {
        const s = memberSearch.toLowerCase();
        if (member.name.toLowerCase().includes(s)) return member;
      });
      setData((prev: any) => ({ ...prev, members: members ?? [] }));
    } else if (memberSearch.length === 0) {
      setData((prev: any) => ({ ...prev, members: MEMBERS_DATA }));
    }
  }, [memberSearch]);

  useEffect(() => {
    if (nodeSearch.length > 3) {
      const nodes = NODE_POOL_DATA.filter((node) => {
        if (node.id.includes(nodeSearch)) return node;
      });
      setData((prev: any) => ({ ...prev, node: nodes ?? [] }));
    } else if (nodeSearch.length === 0) {
      setData((prev: any) => ({ ...prev, node: NODE_POOL_DATA }));
    }
  }, [nodeSearch]);

  const onMenuItemClick = (id: string) => setMenu(id);

  return (
    <Stack mt={3} direction={{ xs: 'column', md: 'row' }} spacing={3}>
      <ManageMenu selectedId={menu} onMenuItemClick={onMenuItemClick} />
      {menu === 'manage-members' && (
        <MemberContainer
          data={data.members}
          search={memberSearch}
          setSearch={setMemberSearch}
        />
      )}
      {menu === 'manage-sim' && <SimPoolContainer data={data.simPool} />}
      {menu === 'manage-node' && (
        <NodePoolContainer
          data={data.node}
          search={nodeSearch}
          setSearch={setNodeSearch}
        />
      )}
      {menu === 'manage-data-plan' && (
        <DataPlanContainer data={data.dataPlan} />
      )}
    </Stack>
  );
};
export default Manage;
