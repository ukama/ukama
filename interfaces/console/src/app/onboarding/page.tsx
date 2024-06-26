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
  useGetInvitationsQuery,
  useUpdateInvitationMutation,
} from '@/client/graphql/generated';
import DataTableWithOptions from '@/components/DataTableWithOptions';
import { INVITATION_TABLE_COLUMN, INVITATION_TABLE_MENU } from '@/constants';
import { useAppContext } from '@/context';
import '@/styles/console.css';
import { CenterContainer } from '@/styles/global';
import { colors } from '@/theme';
import PeopleAltIcon from '@mui/icons-material/PeopleAlt';
import { Box, Container, Paper, Stack, Typography } from '@mui/material';

const Page = () => {
  const { user } = useAppContext();

  const { data: invitationsData, refetch: refetchInvitations } =
    useGetInvitationsQuery({
      fetchPolicy: 'network-only',
      variables: {
        email: user.email,
      },
      onCompleted: (data) => {
        if (data.getInvitations.status === Invitation_Status.InviteAccepted) {
          // TODO: ON ACCEPT INVITE REDIRECT TO ROOT SO THAT TOKEN CAN BE REFRESHED
        }
      },
    });

  const [updateInvitation] = useUpdateInvitationMutation({
    fetchPolicy: 'network-only',
    onCompleted: () => {
      refetchInvitations();
    },
  });

  const handleInviteAction = (id: string, type: string) => {
    if (type === 'accept-invite') {
      updateInvitation({
        variables: {
          data: {
            id,
            status: Invitation_Status.InviteAccepted,
          },
        },
      });
    } else if (type === 'reject-invite') {
      updateInvitation({
        variables: {
          data: {
            id,
            status: Invitation_Status.InviteDeclined,
          },
        },
      });
    }
  };

  return (
    <CenterContainer>
      <Container maxWidth={'md'}>
        <Paper sx={{ p: 3, bgcolor: colors.white, borderRadius: '10px' }}>
          <Stack spacing={2}>
            <Typography variant={'h6'}>Welcome to Ukama!</Typography>
            <Typography variant={'body1'}>
              Ukama is currently in beta access. If you do not currently belong
              to an organization partaking in a pilot program with us, but would
              like to, please contact us at hello@ukama.com.
              <br />
              <br />
              If you do belong to an organization, please request that your
              network owner or admin sends you an invite with the following
              credentials, so that you can access your organization’s Console.
            </Typography>
            <Box
              px={2}
              py={2.5}
              width={'fit-content'}
              sx={{
                borderRadius: '12px',
                border: `1px solid ${colors.primaryLight}`,
              }}
            >
              <Stack
                spacing={0.5}
                direction={'row'}
                width={'fit-content'}
                alignItems={'center'}
              >
                <Typography variant={'body1'} fontWeight={600}>
                  {user.name}
                </Typography>
                <span>•</span>
                <Typography variant={'body1'}>{user.email}</Typography>
              </Stack>
            </Box>
            <br />
            <br />
            <DataTableWithOptions
              icon={PeopleAltIcon}
              isRowClickable={false}
              withStatusColumn={true}
              columns={INVITATION_TABLE_COLUMN}
              menuOptions={INVITATION_TABLE_MENU}
              emptyViewLabel={'No invitation yet!'}
              onMenuItemClick={handleInviteAction}
              dataset={
                invitationsData?.getInvitations
                  ? [invitationsData?.getInvitations]
                  : []
              }
            />
          </Stack>
        </Paper>
      </Container>
    </CenterContainer>
  );
};

export default Page;
