/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { MANAGE_MENU_LIST } from '@/constants';
import { useAppContext } from '@/context';
import {
  PackageDto,
  useAddPackageMutation,
  useCreateInvitationMutation,
  useDeleteInvitationMutation,
  useDeletePackageMutation,
  useGetMembersQuery,
  // useGetInvitationsByOrgLazyQuery,
  useGetNetworksLazyQuery,
  useGetPackagesLazyQuery,
  useGetSimsLazyQuery,
  useInvitationsQuery,
  useRemoveMemberMutation,
  useUpdateMemberMutation,
  useUpdatePacakgeMutation,
  useUploadSimsMutation,
} from '@/generated';
import { colors } from '@/styles/theme';
import { TObject } from '@/types';
import DataPlanDialog from '@/ui/molecules/DataPlanDialog';
import FileDropBoxDialog from '@/ui/molecules/FileDropBoxDialog';
import InviteMemberDialog from '@/ui/molecules/InviteMemberDialog';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
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

const SimPool = dynamic(() => import('./_simpool'));
const NodePool = dynamic(() => import('./_nodepool'));
const Member = dynamic(() => import('./_member'));
const DataPlan = dynamic(() => import('./_dataplan'));

const INIT_DATAPLAN = {
  id: '',
  name: '',
  dataVolume: 0,
  dataUnit: '',
  amount: 0,
  duration: 0,
};

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
  const { setSnackbarMessage } = useAppContext();

  const [data, setData] = useState<any>({
    members: [],
    simPool: [],
    node: [],
    invitations: [],
    networkList: [],
  });
  const [dataplan, setDataplan] = useState(INIT_DATAPLAN);

  const {
    data: membersData,
    loading: membersLoading,
    refetch: refetchMembers,
  } = useGetMembersQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setData((prev: any) => ({
        ...prev,
        members: data?.getMembers.members ?? [],
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

  const [getNetworks, { loading: networkLoading }] = useGetNetworksLazyQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setData((prev: any) => ({
        ...prev,
        networkList: data?.getNetworks.networks ?? [],
      }));
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'network',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  // const [getNodes, { loading: getNodesLoading }] = useGetNodesLazyQuery({
  //   fetchPolicy: 'cache-and-network',
  //   onCompleted: (data) => {
  //     const filteredNodes = data?.getNodes.nodes;
  // .filter((node) => node.created_at)
  // .map((node) => ({
  //   ...node,
  //   created_at: format(parseISO(node.created_at), 'dd MMM yyyy'),
  // }));

  //     setData((prev: any) => ({
  //       ...prev,
  //       node: filteredNodes ?? [],
  //     }));
  //   },
  //   onError: (error) => {
  //     setSnackbarMessage({
  //       id: 'node',
  //       message: error.message,
  //       type: 'error' as AlertColor,
  //       show: true,
  //     });
  //   },
  // });

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

  const [
    getPackages,
    { data: packagesData, loading: packagesLoading, refetch: getDataPlans },
  ] = useGetPackagesLazyQuery({
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

  const {
    data: invitationsData,
    loading: invitationsLoading,
    refetch: refetchInvitations,
  } = useInvitationsQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setData((prev: any) => ({
        ...prev,
        invitations: data?.getInvitationsByOrg.invitations ?? [],
      }));
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'invitations',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [sendInvitation, { loading: sendInvitationLoading }] =
    useCreateInvitationMutation({
      onCompleted: () => {
        refetchMembers();
        refetchInvitations();

        setSnackbarMessage({
          id: 'invitation-success',
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
      getDataPlans();
      setSnackbarMessage({
        id: 'add-data-plan',
        message: 'Data plan added successfully',
        type: 'success' as AlertColor,
        show: true,
      });
      setIsDataPlan(false);
      setDataplan(INIT_DATAPLAN);
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
    useDeletePackageMutation({
      onCompleted: () => {
        getDataPlans().then((res) => {
          setData((prev: any) => ({
            ...prev,
            dataPlan: res?.data?.getPackages.packages ?? [],
          }));
        });
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

  const [deleteInvite] = useDeleteInvitationMutation({
    onCompleted: () => {
      refetchInvitations();
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'delete-invitation',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [deleteMember] = useRemoveMemberMutation({
    onCompleted: () => {
      refetchMembers();
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'delete-members',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [updateMember] = useUpdateMemberMutation({
    onCompleted: () => {
      refetchMembers();
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'update-members',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  useEffect(() => {
    if (memberSearch.length > 2) {
      const _members = membersData?.getMembers.members.filter((member) => {
        const s = memberSearch.toLowerCase();
        if (member.name.toLowerCase().includes(s)) return member;
      });
      const _invitations =
        invitationsData?.getInvitationsByOrg.invitations.filter((invite) => {
          const s = memberSearch.toLowerCase();
          if (invite.name.toLowerCase().includes(s)) return invite;
        });
      setData((prev: any) => ({
        ...prev,
        members: _members,
        invitations: _invitations,
      }));
    } else if (
      memberSearch.length === 0 &&
      data.members.length !== membersData?.getMembers.members.length &&
      data.invitations.length !==
        invitationsData?.getInvitationsByOrg.invitations.length
    ) {
      setData((prev: any) => ({
        ...prev,
        members: membersData?.getMembers.members ?? [],
        invitations: invitationsData?.getInvitationsByOrg.invitations ?? [],
      }));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [memberSearch]);

  // useEffect(() => {
  //   if (nodeSearch.length > 3) {
  //     const nodes = data.node.filter((node: any) => {
  //       if (node.id.includes(nodeSearch)) return node;
  //     });
  //     setData((prev: any) => ({ ...prev, node: nodes || [] }));
  //   } else if (nodeSearch.length === 0) {
  //     setData((prev: any) => ({ ...prev, node: data.node }));
  //   }
  // }, [nodeSearch]);

  const onMenuItemClick = (id: string) => {
    if (id === 'manage-sim')
      getSims({
        variables: {
          type: 'unknown',
        },
      });
    else if (id === 'manage-data-plan') getPackages();
    else if (id === 'manage-node') {
      getNetworks();
      // getNodes();
    }

    setMenu(id);
  };

  const handleAddMemberAction = (member: TObject) => {
    sendInvitation({
      variables: {
        data: {
          email: member.email as string,
          role: member.role as string,
          name: member.name as string,
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
      setIsDataPlan(false);
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
      const d: PackageDto | undefined = packagesData?.getPackages.packages.find(
        (pkg: PackageDto) => pkg.uuid === id,
      );
      setDataplan({
        id: id,
        name: d?.name ?? '',
        duration: d?.duration ?? 0,
        dataUnit: d?.dataUnit ?? '',
        dataVolume: d?.dataVolume ?? 0,
        amount: typeof d?.rate.amount === 'number' ? d.rate.amount : 0,
      });
      setIsDataPlan(true);
    }
  };

  const handleCreateNetwork = () => {};

  const handleMemberAction = (id: string, type: string) => {
    const m = membersData?.getMembers.members.find((mem) => mem.id === id);
    if (!m?.isDeactivated && type === 'remove-member') {
      setSnackbarMessage({
        id: 'deactivate-first-error',
        message: 'Please deactivate member first.',
        type: 'error' as AlertColor,
        show: true,
      });
      return;
    }

    if (type === 'member-status-update') {
      if (m)
        updateMember({
          variables: {
            memberId: id,
            data: {
              isDeactivated: !m.isDeactivated,
              role: m.role,
            },
          },
        });
    }
    if (type === 'remove-member') {
      deleteMember({
        variables: {
          memberId: id,
        },
      });
    }
  };

  const handleDeleteInviteAction = (uuid: string) => {
    if (uuid)
      deleteInvite({
        variables: {
          deleteInvitationId: uuid,
        },
      });
  };

  const isLoading =
    packagesLoading ||
    simsLoading ||
    // membersLoading ||
    // addMemberLoading ||
    // sendInvitationLoading ||
    uploadSimsLoading ||
    dataPlanLoading ||
    deletePkgLoading ||
    updatePkgLoading ||
    networkLoading;
  // invitationsLoading ||
  // getNodesLoading;

  return (
    <Stack mt={3} direction={{ xs: 'column', md: 'row' }} spacing={3}>
      <ManageMenu selectedId={menu} onMenuItemClick={onMenuItemClick} />
      <LoadingWrapper
        width="100%"
        radius="medium"
        isLoading={isLoading}
        cstyle={{
          overflow: 'scroll',
          height: isLoading ? '50vh' : '100%',
        }}
      >
        <Paper
          sx={{
            py: 3,
            px: 4,
            width: '100%',
            borderRadius: '10px',
            height: 'calc(100vh - 200px)',
          }}
        >
          {menu === 'manage-members' && (
            <Member
              search={memberSearch}
              memberData={data.members}
              setSearch={setMemberSearch}
              invitationsData={data.invitations}
              handleMemberAction={handleMemberAction}
              invitationTitle=" There is one pending invitation."
              handleButtonAction={() => setIsInviteMember(true)}
              handleDeleteInviteAction={handleDeleteInviteAction}
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
              networkList={data.networkList || []}
              handleCreateNetwork={handleCreateNetwork}
            />
          )}
          {menu === 'manage-data-plan' && (
            <DataPlan
              data={packagesData?.getPackages.packages ?? []}
              handleActionButon={() => {
                setDataplan(INIT_DATAPLAN);
                setIsDataPlan(true);
              }}
              handleOptionMenuItemAction={handleOptionMenuItemAction}
            />
          )}
        </Paper>
      </LoadingWrapper>
      {isInviteMember && (
        <InviteMemberDialog
          title={'Invite member'}
          isOpen={isInviteMember}
          labelNegativeBtn={'Cancel'}
          // invitationLoading={sendInvitationLoading}
          labelSuccessBtn={'Invite member'}
          handleSuccessAction={handleAddMemberAction}
          handleCloseAction={() => setIsInviteMember(false)}
        />
      )}
      {isUploadSims && (
        <FileDropBoxDialog
          isOpen={isUploadSims}
          labelSuccessBtn={'Upload'}
          labelNegativeBtn={'Cancel'}
          title={'Upload Sims in Sim Pool'}
          handleSuccessAction={handleUploadSimsAction}
          handleCloseAction={() => setIsUploadSims(false)}
        />
      )}
      {isDataPlan && (
        <DataPlanDialog
          data={dataplan}
          isOpen={isDataPlan}
          setData={setDataplan}
          title={'Create data plan'}
          labelNegativeBtn={'Cancel'}
          action={dataplan.id ? 'update' : 'add'}
          handleSuccessAction={handleDataPlanAction}
          handleCloseAction={() => setIsDataPlan(false)}
          labelSuccessBtn={dataplan.id ? 'Update Data Plan' : 'Save Data Plan'}
        />
      )}
    </Stack>
  );
};
export default Manage;
