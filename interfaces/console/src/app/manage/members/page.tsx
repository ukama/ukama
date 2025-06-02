/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  Invitation_Status,
  Role_Type,
  useCreateInvitationMutation,
  useDeleteInvitationMutation,
  useGetInvitationsQuery,
  useGetMembersQuery,
  useRemoveMemberMutation,
  useUpdateMemberMutation,
} from '@/client/graphql/generated';
import DataTableWithOptions from '@/components/DataTableWithOptions';
import InviteMemberDialog from '@/components/InviteMemberDialog';
import LoadingWrapper from '@/components/LoadingWrapper';
import SimpleDataTable from '@/components/SimpleDataTable';
import {
  INVITATION_TABLE_COLUMN,
  MEMBER_TABLE_COLUMN,
  MEMBER_TABLE_MENU,
} from '@/constants';
import { useAppContext } from '@/context';
import colors from '@/theme/colors';
import { TObject } from '@/types';
import { Search } from '@mui/icons-material';
import PeopleAltIcon from '@mui/icons-material/PeopleAlt';
import {
  AlertColor,
  Box,
  Button,
  Grid,
  Paper,
  Stack,
  Tab,
  Tabs,
  TextField,
  Typography,
} from '@mui/material';
import React, { useEffect, useState } from 'react';

const Page = () => {
  const [tabIndex, setTabIndex] = useState(0);
  const { setSnackbarMessage } = useAppContext();
  const [search, setSearch] = useState<string>('');
  const [isInviteMember, setIsInviteMember] = useState<boolean>(false);
  const [data, setData] = useState({ members: [], invites: [] });

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

  const {
    data: invitationsData,
    loading: invitationsLoading,
    refetch: refetchInvitations,
  } = useGetInvitationsQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setData((prev: any) => ({
        ...prev,
        invites:
          data?.getInvitations.invitations.filter(
            (i) => i.status != Invitation_Status.InviteAccepted,
          ) ?? [],
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

  const [deleteInvite, { loading: deleteInviteLoading }] =
    useDeleteInvitationMutation({
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

  useEffect(() => {
    if (data.members.length > 2) {
      const _members = membersData?.getMembers.members.filter((member) => {
        const s = search.toLowerCase();
        if (member.name.toLowerCase().includes(s)) return member;
      });
      const _invitations = invitationsData?.getInvitations.invitations.filter(
        (invite) => {
          const s = search.toLowerCase();
          if (invite.name.toLowerCase().includes(s)) return invite;
        },
      );
      setData((prev: any) => ({
        ...prev,
        members: _members,
        invitations: _invitations,
      }));
    } else if (
      data.members.length === 0 &&
      data.members.length !== membersData?.getMembers.members.length &&
      data.invites.length !== invitationsData?.getInvitations.invitations.length
    ) {
      setData((prev: any) => ({
        ...prev,
        members: membersData?.getMembers.members ?? [],
        invites: invitationsData?.getInvitations.invitations ?? [],
      }));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [search]);

  const handleAddMemberAction = (member: TObject) => {
    sendInvitation({
      variables: {
        data: {
          email: (member.email as string).toLowerCase(),
          role: member.role as Role_Type,
          name: member.name as string,
        },
      },
    });
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabIndex(newValue);
  };

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

  const renderMemberDataTable = () => (
    <Stack direction={'column'}>
      <Typography variant="h6" fontWeight={500}>
        Members
      </Typography>
      <DataTableWithOptions
        dataset={data.members ?? []}
        icon={PeopleAltIcon}
        isRowClickable={false}
        columns={MEMBER_TABLE_COLUMN}
        menuOptions={MEMBER_TABLE_MENU}
        emptyViewLabel={'No members yet!'}
        onMenuItemClick={handleMemberAction}
      />
    </Stack>
  );

  return (
    <LoadingWrapper
      width={'100%'}
      radius="medium"
      height={'calc(100vh - 244px)'}
      isLoading={membersLoading ?? invitationsLoading ?? deleteInviteLoading}
    >
      <Paper
        sx={{
          py: { xs: 1.5, md: 3 },
          px: { xs: 2, md: 4 },
          overflow: 'scroll',
          borderRadius: '10px',
          height: '100%',
        }}
      >
        <Box sx={{ width: '100%', height: '100%' }}>
          <Tabs value={tabIndex} onChange={handleTabChange}>
            <Tab label="team members" />
          </Tabs>
          {tabIndex === 0 && (
            <Box sx={{ width: '100%', mt: { xs: 2, md: 4 } }}>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <TextField
                    id="subscriber-search"
                    variant="outlined"
                    size="small"
                    placeholder="Search"
                    defaultValue={search}
                    fullWidth
                    onChange={(e) => setSearch(e.target.value)}
                    InputLabelProps={{
                      shrink: false,
                    }}
                    InputProps={{
                      endAdornment: <Search htmlColor={colors.black54} />,
                    }}
                  />
                </Grid>
                <Grid item xs={6} container justifyContent="flex-end">
                  <Button
                    variant="contained"
                    color="primary"
                    fullWidth
                    sx={{ width: { xs: '100%', md: 'fit-content' } }}
                    onClick={() => setIsInviteMember(true)}
                  >
                    INVITE MEMBER
                  </Button>
                </Grid>
              </Grid>

              <br />
              {renderMemberDataTable()}
              <br />
              <br />
              {data.invites.length > 0 && (
                <Stack direction={'column'}>
                  <Typography variant="h6" fontWeight={500}>
                    Pending/Declined Invitations
                  </Typography>
                  <SimpleDataTable
                    dataKey="id"
                    dataset={data.invites ?? []}
                    columns={INVITATION_TABLE_COLUMN}
                    handleDeleteElement={(id: string) =>
                      handleDeleteInviteAction(id)
                    }
                  />
                </Stack>
              )}
            </Box>
          )}
        </Box>
        {isInviteMember && (
          <InviteMemberDialog
            title={'Invite member'}
            isOpen={isInviteMember}
            labelNegativeBtn={'Cancel'}
            labelSuccessBtn={'Invite member'}
            invitationLoading={sendInvitationLoading}
            handleSuccessAction={handleAddMemberAction}
            handleCloseAction={() => setIsInviteMember(false)}
          />
        )}
      </Paper>
    </LoadingWrapper>
  );
};

export default Page;
