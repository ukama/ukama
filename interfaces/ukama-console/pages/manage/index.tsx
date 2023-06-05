import { snackbarMessage } from '@/app-recoil';
import { MANAGE_MENU_LIST } from '@/constants';
import {
  MemberObj,
  PackageDto,
  useAddMemberMutation,
  useAddPackageMutation,
  useDeletePacakgeMutation,
  useGetOrgMemberQuery,
  useGetPackagesLazyQuery,
  useGetSimsLazyQuery,
  useUpdatePacakgeMutation,
  useUploadSimsMutation,
} from '@/generated';
import { colors } from '@/styles/theme';
import { TObject, TSnackMessage } from '@/types';
import {
  DataPlanDialog,
  FileDropBoxDialog,
  InviteMemberDialog,
  LoadingWrapper,
} from '@/ui/components';
import {
  AlertColor,
  ListItemIcon,
  ListItemText,
  MenuItem,
  MenuList,
  Paper,
  Stack,
} from '@mui/material';
import dynamic from 'next/dynamic';
import { useEffect, useState } from 'react';
import { useSetRecoilState } from 'recoil';

const SimPool = dynamic(() => import('./_simpool'));
const NodePool = dynamic(() => import('./_nodepool'));
const Member = dynamic(() => import('./_member'));
const DataPlan = dynamic(() => import('./_dataplan'));

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

const structureData = (data: any) =>
  data && data.length > 0
    ? data.map((member: MemberObj) => ({
        name: member.user.name,
        email: member.user.email,
        role: 'member',
        uuid: member.uuid,
      }))
    : [];

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

const Manage = () => {
  const [isInviteMember, setIsInviteMember] = useState<boolean>(false);
  const [isUploadSims, setIsUploadSims] = useState<boolean>(false);
  const [isDataPlan, setIsDataPlan] = useState<boolean>(false);
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
  const [dataplan, setDataplan] = useState({
    id: '',
    name: '',
    dataVolume: 0,
    dataUnit: '',
    amount: 0,
    duration: 0,
  });

  const {
    data: members,
    loading: membersLoading,
    refetch: refetchMembers,
  } = useGetOrgMemberQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setData((prev: any) => ({
        ...prev,
        members: structureData(members?.getOrgMembers.members),
      }));
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

  const [getSims, { loading: simsLoading, refetch: refetchSims }] =
    useGetSimsLazyQuery({
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

  const [getPackages, { loading: packagesLoading, refetch: getDataPlans }] =
    useGetPackagesLazyQuery({
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

  const [uploadSimPool, { loading: uploadSimsLoading }] = useUploadSimsMutation(
    {
      onCompleted: () => {
        refetchSims();
        setSnackbarMessage({
          id: 'sim-pool-uploaded',
          message: 'Sims uploaded successfully',
          type: 'success' as AlertColor,
          show: true,
        });
        setIsUploadSims(false);
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'sim-pool-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    },
  );

  const [addDataPlan, { loading: dataPlanLoading }] = useAddPackageMutation({
    onCompleted: () => {
      refetchSims();
      setSnackbarMessage({
        id: 'add-data-plan',
        message: 'Data plan added successfully',
        type: 'success' as AlertColor,
        show: true,
      });
      setIsDataPlan(false);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'data-plan-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [deletePackage, { loading: deletePkgLoading }] =
    useDeletePacakgeMutation({
      onCompleted: () => {
        getDataPlans();
        setSnackbarMessage({
          id: 'delete-data-plan',
          message: 'Data plan deleted successfully',
          type: 'success' as AlertColor,
          show: true,
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'data-plan-delete-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const [updatePackage, { loading: updatePkgLoading }] =
    useUpdatePacakgeMutation({
      onCompleted: () => {
        getDataPlans();
        setSnackbarMessage({
          id: 'update-data-plan',
          message: 'Data plan updated successfully',
          type: 'success' as AlertColor,
          show: true,
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'data-plan-update-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  useEffect(() => {
    if (memberSearch.length > 2) {
      const _members = members?.getOrgMembers.members.filter((member) => {
        const s = memberSearch.toLowerCase();
        if (member.user.name.toLowerCase().includes(s)) return member;
      });
      setData((prev: any) => ({
        ...prev,
        members: structureData(_members),
      }));
    } else if (
      memberSearch.length === 0 &&
      data.members.length !== members?.getOrgMembers.members.length
    ) {
      setData((prev: any) => ({
        ...prev,
        members: structureData(members?.getOrgMembers.members),
      }));
    }
  }, [memberSearch]);

  useEffect(() => {
    if (nodeSearch.length > 3) {
      const nodes = NODE_POOL_DATA.filter((node) => {
        if (node.id.includes(nodeSearch)) return node;
      });
      setData((prev: any) => ({ ...prev, node: nodes || [] }));
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

  const handleUploadSimsAction = (
    action: string,
    value: string,
    type: string,
  ) => {
    if (action === 'error') {
      setSnackbarMessage({
        id: 'sim-pool-parsing-error',
        message: value,
        type: 'error' as AlertColor,
        show: true,
      });
    } else if (action === 'success') {
      uploadSimPool({
        variables: {
          data: {
            data: value,
            simType: type,
          },
        },
      });
    }
  };

  const handleDataPlanAction = (action: string) => {
    if (action === 'add') {
      addDataPlan({
        variables: {
          data: {
            name: dataplan.name,
            amount: dataplan.amount,
            dataUnit: dataplan.dataUnit,
            dataVolume: dataplan.dataVolume,
            duration: dataplan.duration,
          },
        },
      });
    } else if (action === 'update') {
      updatePackage({
        variables: {
          packageId: dataplan.id,
          data: {
            name: dataplan.name,
            active: true,
          },
        },
      });
    }
  };

  const handleOptionMenuItemAction = (id: string, action: string) => {
    if (action === 'delete') {
      deletePackage({
        variables: {
          packageId: id,
        },
      });
    } else if (action === 'edit') {
      const d: PackageDto = data.dataPlan.find(
        (pkg: PackageDto) => pkg.uuid === id,
      );
      setDataplan({
        id: id,
        amount: d.rate.amount,
        dataUnit: d.dataUnit,
        dataVolume: parseInt(parseInt(d.dataVolume).toFixed(2)),
        duration: parseInt(d.duration),
        name: d.name,
      });
      setIsDataPlan(true);
    }
  };

  const isLoading =
    packagesLoading ||
    simsLoading ||
    membersLoading ||
    addMemberLoading ||
    uploadSimsLoading ||
    dataPlanLoading ||
    deletePkgLoading ||
    updatePkgLoading;

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
            <Member
              search={memberSearch}
              setSearch={setMemberSearch}
              data={data.members}
              handleButtonAction={() => setIsInviteMember(true)}
            />
          )}
          {menu === 'manage-sim' && (
            <SimPool
              data={data.simPool}
              handleActionButon={() => setIsUploadSims(true)}
            />
          )}
          {menu === 'manage-node' && (
            <NodePool
              data={data.node}
              search={nodeSearch}
              setSearch={setNodeSearch}
            />
          )}
          {menu === 'manage-data-plan' && (
            <DataPlan
              data={data.dataPlan}
              handleActionButon={() => setIsDataPlan(true)}
              handleOptionMenuItemAction={handleOptionMenuItemAction}
            />
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
      {isUploadSims && (
        <FileDropBoxDialog
          title={'Upload Sims in Sim Pool'}
          isOpen={isUploadSims}
          labelNegativeBtn={'Cancel'}
          labelSuccessBtn={'Upload'}
          handleSuccessAction={handleUploadSimsAction}
          handleCloseAction={() => setIsUploadSims(false)}
        />
      )}
      {isDataPlan && (
        <DataPlanDialog
          data={dataplan}
          action={dataplan.id ? 'update' : 'add'}
          isOpen={isDataPlan}
          setData={setDataplan}
          title={'Create data plan'}
          labelNegativeBtn={'Cancel'}
          labelSuccessBtn={dataplan.id ? 'Update Data Plan' : 'Save Data Plan'}
          handleSuccessAction={handleDataPlanAction}
          handleCloseAction={() => setIsDataPlan(false)}
        />
      )}
    </Stack>
  );
};
export default Manage;
