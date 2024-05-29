import { commonData, user } from '@/app-recoil';
import { INVITATION_TABLE_COLUMN, INVITATION_TABLE_MENU } from '@/constants';
import {
  Invitation_Status,
  useGetInvitationsQuery,
  useUpdateInvitationMutation,
  useWhoamiLazyQuery,
} from '@/generated';
import colors from '@/styles/theme/colors';
import { TCommonData, TUser } from '@/types';
import DataTableWithOptions from '@/ui/molecules/DataTableWithOptions';
import PeopleAltIcon from '@mui/icons-material/PeopleAlt';
import { Box, Button, Stack, Typography } from '@mui/material';
import { useRecoilState } from 'recoil';

const OnBoarding = () => {
  const [_commonData, _setCommonData] = useRecoilState<TCommonData>(commonData);

  const [_user, _setUser] = useRecoilState<TUser>(user);

  const {
    data: invitationsData,
    loading: invitationsLoading,
    refetch: refetchInvitations,
  } = useGetInvitationsQuery({
    fetchPolicy: 'network-only',
    variables: {
      email: _user.email,
    },
    onCompleted: (data) => {
      if (data.getInvitations.status === Invitation_Status.Accepted) {
        _setUser({
          ..._user,
          role: data.getInvitations.role,
        });
        // _setCommonData({
        //   ..._commonData,
        //   userId: _user.id,
        // });

        whoami();
      }
    },
  });

  const [whoami] = useWhoamiLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      if (data.whoami) {
        if (data.whoami.memberOf.length > 0) {
          // _setCommonData({
          //   ..._commonData,
          //   orgId: data.whoami.memberOf[0].id,
          //   orgName: data.whoami.memberOf[0].name,
          // });
        }
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
            status: Invitation_Status.Accepted,
          },
        },
      }).then(() => {
        whoami();
      });
    } else if (type === 'reject-invite') {
      updateInvitation({
        variables: {
          data: {
            id,
            status: Invitation_Status.Declined,
          },
        },
      });
    }
  };

  return (
    <Box bgcolor={'white'} borderRadius={'8px'} p={4}>
      <Stack direction={'column'} spacing={3}>
        <Typography variant="h4" color={colors.primaryMain}>
          <b>Welcome to Ukama</b>
        </Typography>
        <Typography variant="body1" fontWeight={400}>
          Ukama offers an open access complete cellular network solution without
          the price tag and technical expertise required of traditional setups.
        </Typography>
        <Typography variant="body1" fontWeight={400}>
          Share below details with org owner/admin to get invited to the
          network.
        </Typography>

        <Stack
          px={2}
          pb={2}
          pt={1.6}
          spacing={0.8}
          borderRadius={2}
          width="300px"
          bgcolor={colors.black10}
        >
          <Typography variant="body1" fontWeight={500}>
            <b>Name: </b>
            {_user.name}
          </Typography>
          <Typography variant="body1" fontWeight={500}>
            <b>Email: </b>
            {_user.email}
          </Typography>
          <br />
          <Button
            variant="contained"
            sx={{ mb: 3 }}
            onClick={() =>
              navigator.clipboard.writeText(
                `Name: ${_user.name}, Email: ${_user.email}`,
              )
            }
          >
            Copy
          </Button>
        </Stack>
        <Typography variant="h6" fontWeight={600}>
          Invitations
        </Typography>
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
    </Box>
  );
};

export default OnBoarding;
