import { snackbarMessage } from '@/app-recoil';
import {
  MANAGE_MENU_LIST,
  MANAGE_NODE_POOL_COLUMN,
  MANAGE_SIM_POOL_COLUMN,
  MANAGE_TABLE_COLUMN,
} from '@/constants';
import {
  OrgMembersResDto,
  useAddMemberMutation,
  useGetOrgMemberQuery,
  useGetPackagesLazyQuery,
  useGetSimsLazyQuery,
} from '@/generated';
import { colors } from '@/styles/theme';
import { TObject, TSnackMessage } from '@/types';
import {
  InviteMemberDialog,
  LoadingWrapper,
  SimpleDataTable,
} from '@/ui/components';
import PageContainerHeader from '@/ui/components/PageContainerHeader';
import { getDataPlanUsage } from '@/utils';
import PeopleAlt from '@mui/icons-material/PeopleAlt';
import {
  AlertColor,
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
import { useSetRecoilState } from 'recoil';

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

interface IMemberContainer {
  data: any;
  search: string;
  setSearch: (value: string) => void;
  handleButtonAction: () => void;
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

const MemberContainer = ({
  data,
  search,
  setSearch,
  handleButtonAction,
}: IMemberContainer) => (
  <Paper
    sx={{
      py: 3,
      px: 4,
      width: '100%',
      overflow: 'hidden',
      borderRadius: '5px',
      height: 'calc(100vh - 200px)',
    }}
  >
    <PageContainerHeader
      search={search}
      title={'My members'}
      buttonTitle={'Invite member'}
      onSearchChange={(e: string) => setSearch(e)}
      handleButtonAction={handleButtonAction}
    />
    <br />
    <SimpleDataTable
      dataKey="uuid"
      dataset={data}
      columns={MANAGE_TABLE_COLUMN}
    />
  </Paper>
);

const SimPoolContainer = ({ data }: ISimPoolContainer) => (
  <Paper
    sx={{
      py: 3,
      px: 4,
      width: '100%',
      overflow: 'hidden',
      borderRadius: '5px',
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
      borderRadius: '5px',
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
      borderRadius: '5px',
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
      {data.map(
        ({
          uuid,
          name,
          duration,
          users,
          currency,
          dataVolume,
          dataUnit,
          amount,
        }: any) => (
          <Grid item xs={12} sm={6} md={4} key={uuid}>
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
                  {getDataPlanUsage(
                    duration,
                    currency,
                    amount,
                    dataVolume,
                    dataUnit,
                  )}
                </Typography>
                {false && (
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
                )}
              </Stack>
            </Paper>
          </Grid>
        ),
      )}
    </Grid>
  </Paper>
);

const Manage = () => {
  const [isInviteMember, setIsInviteMember] = useState<boolean>(false);
  const [menu, setMenu] = useState<string>('manage-members');
  const [memberSearch, setMemberSearch] = useState<string>('');
  const [nodeSearch, setNodeSearch] = useState<string>('');
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const [data, setData] = useState<any>({
    members: [],
    simPool: [],
    dataPlan: [],
    node: NODE_POOL_DATA,
  });

  const {
    data: members,
    loading: membersLoading,
    refetch: refetchMembers,
  } = useGetOrgMemberQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setData((prev: any) => ({ ...prev, members: data?.getOrgMembers ?? [] }));
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'org-members',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [getSims, { loading: simsLoading }] = useGetSimsLazyQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setData((prev: any) => ({
        ...prev,
        simPool: data?.getSims.sim ?? [],
      }));
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'sim-pool',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [getPackages, { loading: packagesLoading }] = useGetPackagesLazyQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setData((prev: any) => ({
        ...prev,
        dataPlan: data?.getPackages.packages ?? [],
      }));
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'packages',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [addMember, { loading: addMemberLoading }] = useAddMemberMutation({
    onCompleted: () => {
      refetchMembers();
      setSnackbarMessage({
        id: 'add-member',
        message: 'Invitation sent successfully',
        type: 'success' as AlertColor,
        show: true,
      });
      setIsInviteMember(false);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'add-member-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  useEffect(() => {
    if (memberSearch.length > 3) {
      const _members = members?.getOrgMembers.members.filter((member) => {
        const s = memberSearch.toLowerCase();
        if (member.uuid.includes(s)) return member;
      });
      setData((prev: any) => ({ ...prev, members: _members ?? [] }));
    } else if (memberSearch.length === 0) {
      setData((prev: any) => ({
        ...prev,
        members: members?.getOrgMembers.members ?? [],
      }));
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

  const onMenuItemClick = (id: string) => {
    if (id === 'manage-sim')
      getSims({
        variables: {
          type: 'unknown',
        },
      });
    else if (id === 'manage-data-plan') getPackages();
    setMenu(id);
  };

  const structureData = (data: OrgMembersResDto) =>
    data.members?.map((member) => ({
      name: member.user.name,
      email: member.user.email,
      role: 'member',
      uuid: member.uuid,
    }));

  const handleAddMemberAction = (member: TObject) => {
    addMember({
      variables: {
        data: {
          email: member.email as string,
          role: member.role as string,
        },
      },
    });
  };

  const isLoading =
    packagesLoading || simsLoading || membersLoading || addMemberLoading;

  return (
    <Stack mt={3} direction={{ xs: 'column', md: 'row' }} spacing={3}>
      <ManageMenu selectedId={menu} onMenuItemClick={onMenuItemClick} />
      <LoadingWrapper
        width="100%"
        radius="small"
        isLoading={isLoading}
        cstyle={{ backgroundColor: isLoading ? colors.white : 'transparent' }}
      >
        <>
          {menu === 'manage-members' && (
            <MemberContainer
              search={memberSearch}
              setSearch={setMemberSearch}
              data={structureData(data.members)}
              handleButtonAction={() => setIsInviteMember(true)}
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
        </>
      </LoadingWrapper>
      {isInviteMember && (
        <InviteMemberDialog
          title={'Invite member'}
          isOpen={isInviteMember}
          labelNegativeBtn={'Cancel'}
          labelSuccessBtn={'Invite member'}
          handleSuccessAction={handleAddMemberAction}
          handleCloseAction={() => setIsInviteMember(false)}
        />
      )}
    </Stack>
  );
};
export default Manage;
